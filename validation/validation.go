package validation

import (
	"strconv"

	"github.com/branila/restaurant-protocol/errors"
	"github.com/branila/restaurant-protocol/models"
)

// Esegue una validazione completa di un ordine
func ValidateOrdine(ordine models.Ordine) error {
	// Verifica che l'ordine abbia almeno una comanda
	if len(ordine.Comande) == 0 {
		return errors.NewOrdineVuotoError()
	}

	// Verifica che ogni comanda sia valida
	for _, comanda := range ordine.Comande {
		if err := ValidateComanda(comanda); err != nil {
			return err
		}
	}

	return nil
}

// Esegue una validazione completa di una comanda
func ValidateComanda(comanda models.Comanda) error {
	// Verifica che la comanda abbia almeno un piatto
	if comanda.Primo == nil && comanda.Secondo == nil && comanda.Contorno == nil {
		return errors.NewComandaVuotaError(strconv.Itoa(comanda.Numero))
	}

	return nil
}

// TODO: Esegue una validazione completa di un piatto
func ValidatePiatto(piatto *models.Piatto) error {
	return nil
}
