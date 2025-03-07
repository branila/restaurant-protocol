package errors

import "fmt"

// Codici di errore
const (
	// Errori di validazione
	ErrCodePiattiMultipli = 1001 // Tentativo di ordinare più piatti dello stesso tipo
	ErrCodePiattoEsaurito = 1002 // Piatto non disponibile nell'inventario
	ErrCodeComandaVuota   = 1003 // Comanda senza piatti
	ErrCodeOrdineVuoto    = 1004 // Ordine senza comande

	// Errori sintattici
	ErrCodeSintassiGenerale = 2001 // Errore generico di sintassi
	ErrCodeFormatoData      = 2002 // Formato data non valido
	ErrCodeNumeroNonValido  = 2003 // Numero non valido per tavolo/comanda

	// Errori relativi alle modifiche
	ErrCodeModificaNonValida = 3001 // Modifica non consentita
)

// OrderError è un tipo di errore personalizzato che contiene informazioni dettagliate
type OrderError struct {
	Code    int    // Codice numerico dell'errore
	Message string // Messaggio descrittivo dell'errore
	Details string // Dettagli aggiuntivi per risolvere l'errore
}

// Error implementa l'interfaccia error
func (e *OrderError) Error() string {
	return fmt.Sprintf("[Errore %d] %s. %s", e.Code, e.Message, e.Details)
}

// NewPiattiMultipliError crea un errore per piatti multipli dello stesso tipo
func NewPiattiMultipliError(tipoPiatto, numeroComanda string) *OrderError {
	return &OrderError{
		Code:    ErrCodePiattiMultipli,
		Message: fmt.Sprintf("Tentativo di ordinare più %s nella stessa comanda", tipoPiatto),
		Details: fmt.Sprintf("Nella comanda %s è già presente un %s. Rimuovere uno dei piatti o inserirlo in una comanda diversa", numeroComanda, tipoPiatto),
	}
}

// NewPiattoEsauritoError crea un errore per piatti non disponibili
func NewPiattoEsauritoError(nomePiatto string, disponibilita int) *OrderError {
	details := "Il piatto è esaurito per oggi"
	if disponibilita > 0 {
		details = fmt.Sprintf("Rimangono solo %d porzioni di questo piatto", disponibilita)
	}

	return &OrderError{
		Code:    ErrCodePiattoEsaurito,
		Message: fmt.Sprintf("Il piatto '%s' non è disponibile nella quantità richiesta", nomePiatto),
		Details: details,
	}
}

// NewComandaVuotaError crea un errore per comande vuote
func NewComandaVuotaError(numeroComanda string) *OrderError {
	return &OrderError{
		Code:    ErrCodeComandaVuota,
		Message: fmt.Sprintf("La comanda %s non contiene piatti", numeroComanda),
		Details: "Aggiungere almeno un piatto (primo, secondo o contorno) alla comanda",
	}
}

// NewOrdineVuotoError crea un errore per ordini vuoti
func NewOrdineVuotoError() *OrderError {
	return &OrderError{
		Code:    ErrCodeOrdineVuoto,
		Message: "L'ordine non contiene comande",
		Details: "Aggiungere almeno una comanda all'ordine",
	}
}

// NewSyntaxError crea un errore di sintassi
func NewSyntaxError(riga string, dettaglio string) *OrderError {
	return &OrderError{
		Code:    ErrCodeSintassiGenerale,
		Message: fmt.Sprintf("Errore di sintassi alla riga: '%s'", riga),
		Details: dettaglio,
	}
}

// NewFormatoDataError crea un errore di formato data
func NewFormatoDataError(data string) *OrderError {
	return &OrderError{
		Code:    ErrCodeFormatoData,
		Message: fmt.Sprintf("Formato data non valido: '%s'", data),
		Details: "Il formato corretto della data è DD/MM/YYYY",
	}
}

// NewNumeroNonValidoError crea un errore per numeri non validi
func NewNumeroNonValidoError(contesto, valore string) *OrderError {
	return &OrderError{
		Code:    ErrCodeNumeroNonValido,
		Message: fmt.Sprintf("Numero %s non valido: '%s'", contesto, valore),
		Details: fmt.Sprintf("Il valore per %s deve essere un numero intero positivo", contesto),
	}
}

// NewModificaNonValidaError crea un errore per modifiche non valide
func NewModificaNonValidaError(piatto string, modifica string, motivazione string) *OrderError {
	return &OrderError{
		Code:    ErrCodeModificaNonValida,
		Message: fmt.Sprintf("Modifica '%s' non consentita per il piatto '%s'", modifica, piatto),
		Details: motivazione,
	}
}

// Is permette di confrontare i tipi di errore in base al codice
func (e *OrderError) Is(target error) bool {
	if t, ok := target.(*OrderError); ok {
		return e.Code == t.Code
	}
	return false
}
