package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/branila/restaurant-protocol/errors"
	"github.com/branila/restaurant-protocol/inventory"
	"github.com/branila/restaurant-protocol/models"
)

var (
	// Regexp per validare il formato della data
	dateRegex = regexp.MustCompile(`^\d{2}/\d{2}/\d{4}$`)

	// Inventario globale
	Inventario *inventory.Inventory
)

// Inizializza l'inventario con valori predefiniti
func Init() {
	Inventario = inventory.DefaultInventory()
}

func ParseOrdine(lines []string) (models.Ordine, error) {
	// Inizializza l'inventario se non è già stato fatto
	if Inventario == nil {
		Init()
	}

	var ordine models.Ordine
	var comandaCorrenteIndex int = -1

	if len(lines) == 0 {
		return ordine, errors.NewOrdineVuotoError()
	}

	// Analizza l'intestazione dell'ordine (prima riga)
	if err := parseIntestazione(lines[0], &ordine); err != nil {
		return ordine, err
	}

	// Analizza il resto delle righe
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Controlla se è l'inizio di una nuova comanda
		if strings.HasPrefix(line, "COMANDA") {
			// Valida la comanda precedente, se presente
			if comandaCorrenteIndex >= 0 {
				if err := validaComanda(ordine.Comande[comandaCorrenteIndex]); err != nil {
					return ordine, err
				}
			}

			// Crea una nuova comanda
			comanda, err := parseComanda(line)
			if err != nil {
				return ordine, err
			}
			ordine.Comande = append(ordine.Comande, *comanda)
			comandaCorrenteIndex = len(ordine.Comande) - 1
		} else if comandaCorrenteIndex >= 0 {
			// Analizza il piatto all'interno della comanda corrente
			if err := parsePiatto(line, &ordine.Comande[comandaCorrenteIndex]); err != nil {
				return ordine, err
			}
		} else {
			return ordine, errors.NewSyntaxError(line, "Trovato piatto senza comanda di riferimento")
		}
	}

	// Valida l'ultima comanda se presente
	if comandaCorrenteIndex >= 0 {
		if err := validaComanda(ordine.Comande[comandaCorrenteIndex]); err != nil {
			return ordine, err
		}
	}

	// Controlla che l'ordine abbia almeno una comanda
	if len(ordine.Comande) == 0 {
		return ordine, errors.NewOrdineVuotoError()
	}

	return ordine, nil
}

// Analizza la riga di intestazione dell'ordine
func parseIntestazione(line string, ordine *models.Ordine) error {
	parts := strings.Split(line, " ")
	if len(parts) < 3 || parts[0] != "ORDINE" {
		return errors.NewSyntaxError(line, "L'intestazione deve essere nel formato 'ORDINE [numero] [data]'")
	}

	// Analizza il numero del tavolo
	tavolo, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.NewNumeroNonValidoError("tavolo", parts[1])
	}
	ordine.Tavolo = tavolo

	// Analizza la data
	data := parts[2]
	if !dateRegex.MatchString(data) {
		return errors.NewFormatoDataError(data)
	}
	ordine.Data = data

	return nil
}

// Analizza una riga di comanda
func parseComanda(line string) (*models.Comanda, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 || parts[0] != "COMANDA" {
		return nil, errors.NewSyntaxError(line, "La comanda deve essere nel formato 'COMANDA [numero]'")
	}

	// Analizza il numero della comanda
	numero, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errors.NewNumeroNonValidoError("comanda", parts[1])
	}

	return &models.Comanda{Numero: numero}, nil
}

// Analizza una riga di piatto
func parsePiatto(line string, comanda *models.Comanda) error {
	// Separa il tipo di piatto dal resto della riga
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return errors.NewSyntaxError(line, "Il formato del piatto non è valido")
	}

	tipoPiatto := parts[0]
	restoDellaPiatto := parts[1]

	// Estrai il nome del piatto tra virgolette
	regexPiatto := regexp.MustCompile(`"([^"]+)"`)
	nomeMatch := regexPiatto.FindStringSubmatch(restoDellaPiatto)
	if len(nomeMatch) < 2 {
		return errors.NewSyntaxError(line, "Il nome del piatto deve essere tra virgolette")
	}

	nomePiatto := nomeMatch[1]

	// Verifica la disponibilità del piatto
	if err := Inventario.VerificaDisponibilita(nomePiatto); err != nil {
		return err
	}

	// Crea il piatto con le modifiche
	piatto := &models.Piatto{
		Nome:      nomePiatto,
		Modifiche: []models.Modifica{},
	}

	// Analizza le modifiche (+ e -)
	modificheRegex := regexp.MustCompile(`([+-])"([^"]+)"`)
	modificheMatch := modificheRegex.FindAllStringSubmatch(restoDellaPiatto, -1)

	for _, mod := range modificheMatch {
		if len(mod) < 3 {
			continue
		}

		tipoModifica := mod[1]
		voceModifica := mod[2]

		// Verifica che la modifica sia consentita
		if err := Inventario.VerificaModifica(nomePiatto, tipoModifica, voceModifica); err != nil {
			return err
		}

		// Aggiungi la modifica al piatto
		piatto.Modifiche = append(piatto.Modifiche, models.Modifica{
			Tipo: tipoModifica,
			Voce: voceModifica,
		})
	}

	// Assegna il piatto alla comanda in base al tipo
	switch tipoPiatto {
	case "PRIMO":
		if comanda.Primo != nil {
			return errors.NewPiattiMultipliError("PRIMO", strconv.Itoa(comanda.Numero))
		}
		comanda.Primo = piatto
	case "SECONDO":
		if comanda.Secondo != nil {
			return errors.NewPiattiMultipliError("SECONDO", strconv.Itoa(comanda.Numero))
		}
		comanda.Secondo = piatto
	case "CONTORNO":
		if comanda.Contorno != nil {
			return errors.NewPiattiMultipliError("CONTORNO", strconv.Itoa(comanda.Numero))
		}
		comanda.Contorno = piatto
	default:
		return errors.NewSyntaxError(line, fmt.Sprintf("Tipo di piatto non riconosciuto: %s", tipoPiatto))
	}

	// Decrementa la disponibilità del piatto nell'inventario
	if err := Inventario.DecrementaDisponibilita(nomePiatto); err != nil {
		return err
	}

	return nil
}

// validaComanda verifica che una comanda sia valida
func validaComanda(comanda models.Comanda) error {
	// Controlla che la comanda abbia almeno un piatto
	if comanda.Primo == nil && comanda.Secondo == nil && comanda.Contorno == nil {
		return errors.NewComandaVuotaError(strconv.Itoa(comanda.Numero))
	}

	return nil
}
