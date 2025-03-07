package formatter

import (
	"fmt"
	"strings"

	"github.com/branila/restaurant-protocol/models"
)

// Formatta un ordine in una stringa leggibile
func FormatOrdine(ordine models.Ordine) string {
	var output strings.Builder

	// Intestazione dell'ordine
	output.WriteString(fmt.Sprintf("Ordine per il tavolo %d, data %s\n", ordine.Tavolo, ordine.Data))

	// Formatta ogni comanda
	for _, comanda := range ordine.Comande {
		output.WriteString(formatComanda(comanda))
	}

	return output.String()
}

// Formatta una comanda in una stringa leggibile
func formatComanda(comanda models.Comanda) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("  Comanda %d:\n", comanda.Numero))

	// Formatta il primo se presente
	if comanda.Primo != nil {
		output.WriteString(formatPiatto("Primo", comanda.Primo))
	}

	// Formatta il secondo se presente
	if comanda.Secondo != nil {
		output.WriteString(formatPiatto("Secondo", comanda.Secondo))
	}

	// Formatta il contorno se presente
	if comanda.Contorno != nil {
		output.WriteString(formatPiatto("Contorno", comanda.Contorno))
	}

	return output.String()
}

// Formatta un piatto in una stringa leggibile
func formatPiatto(categoria string, piatto *models.Piatto) string {
	return fmt.Sprintf("    %s: %s %s\n",
		categoria,
		piatto.Nome,
		formatModifiche(piatto.Modifiche))
}

// Formatta le modifiche in una stringa leggibile
func formatModifiche(modifiche []models.Modifica) string {
	if len(modifiche) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString("[")

	for i, mod := range modifiche {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(fmt.Sprintf("{%s %s}", mod.Tipo, mod.Voce))
	}

	result.WriteString("]")
	return result.String()
}
