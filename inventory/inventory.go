package inventory

import (
	"fmt"
	"sync"

	"github.com/branila/restaurant-protocol/errors"
)

type Piatto struct {
	Nome                string
	Disponibilita       int
	ModificheConsentite map[string]bool // true = aggiungere, false = rimuovere
}

type Inventory struct {
	piatti map[string]Piatto
	mu     sync.RWMutex
}

func New() *Inventory {
	return &Inventory{
		piatti: make(map[string]Piatto),
	}
}

func DefaultInventory() *Inventory {
	inv := New()

	// Primi piatti
	inv.AddPiatto("pasta al pomodoro", 10, map[string]bool{
		"formaggio": true,
		"basilico":  false,
		"pomodoro":  false,
	})

	inv.AddPiatto("risotto ai funghi", 5, map[string]bool{
		"parmigiano": true,
		"funghi":     false,
		"burro":      false,
	})

	// Secondi
	inv.AddPiatto("Bistecca", 8, map[string]bool{
		"Salsa barbecue": true,
		"pepe":           true,
		"sale":           false,
	})

	// Contorni
	inv.AddPiatto("insalata", 15, map[string]bool{
		"aceto balsamico": true,
		"pomodorini":      true,
		"olio":            false,
	})

	return inv
}

func (inv *Inventory) AddPiatto(nome string, disponibilita int, modifiche map[string]bool) {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	inv.piatti[nome] = Piatto{
		Nome:                nome,
		Disponibilita:       disponibilita,
		ModificheConsentite: modifiche,
	}
}

func (inv *Inventory) GetPiatto(nome string) (Piatto, bool) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	piatto, exists := inv.piatti[nome]
	return piatto, exists
}

func (inv *Inventory) VerificaDisponibilita(nome string) error {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	piatto, exists := inv.piatti[nome]
	if !exists {
		return &errors.OrderError{
			Code:    errors.ErrCodePiattoEsaurito,
			Message: fmt.Sprintf("Il piatto '%s' non esiste nel menu", nome),
			Details: "Consultare il menu aggiornato per verificare i piatti disponibili",
		}
	}

	if piatto.Disponibilita <= 0 {
		return errors.NewPiattoEsauritoError(nome, 0)
	}

	return nil
}

func (inv *Inventory) VerificaModifica(nomePiatto string, tipoModifica string, voceModifica string) error {
	inv.mu.RLock()
	defer inv.mu.RUnlock()

	piatto, exists := inv.piatti[nomePiatto]
	if !exists {
		return &errors.OrderError{
			Code:    errors.ErrCodePiattoEsaurito,
			Message: fmt.Sprintf("Il piatto '%s' non esiste nel menu", nomePiatto),
			Details: "Consultare il menu aggiornato per verificare i piatti disponibili",
		}
	}

	isAggiunta := tipoModifica == "+"
	consentita, exists := piatto.ModificheConsentite[voceModifica]

	if !exists {
		return errors.NewModificaNonValidaError(
			nomePiatto,
			fmt.Sprintf("%s%s", tipoModifica, voceModifica),
			fmt.Sprintf("La modifica di '%s' non è prevista per questo piatto", voceModifica),
		)
	}

	if isAggiunta != consentita {
		azione := "aggiungere"
		if isAggiunta {
			azione = "rimuovere"
		}
		return errors.NewModificaNonValidaError(
			nomePiatto,
			fmt.Sprintf("%s%s", tipoModifica, voceModifica),
			fmt.Sprintf("Non è possibile %s '%s' per questo piatto", azione, voceModifica),
		)
	}

	return nil
}

func (inv *Inventory) DecrementaDisponibilita(nome string) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	piatto, exists := inv.piatti[nome]
	if !exists {
		return &errors.OrderError{
			Code:    errors.ErrCodePiattoEsaurito,
			Message: fmt.Sprintf("Il piatto '%s' non esiste nel menu", nome),
			Details: "Consultare il menu aggiornato per verificare i piatti disponibili",
		}
	}

	if piatto.Disponibilita <= 0 {
		return errors.NewPiattoEsauritoError(nome, 0)
	}

	piatto.Disponibilita--
	inv.piatti[nome] = piatto

	// Avvisa se stiamo per esaurire il piatto
	if piatto.Disponibilita <= 3 {
		fmt.Printf("Attenzione: rimangono solo %d porzioni di %s\n", piatto.Disponibilita, nome)
	}

	return nil
}
