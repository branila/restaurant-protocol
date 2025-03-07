package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/branila/restaurant-protocol/converter"
	"github.com/branila/restaurant-protocol/errors"
	"github.com/branila/restaurant-protocol/formatter"
	"github.com/branila/restaurant-protocol/parser"
)

func main() {
	// Inizializza il parser
	parser.Init()

	// Legge il file di input
	lines, err := readInputFile("ordine.txt")
	if err != nil {
		fmt.Printf("Errore nella lettura del file: %v\n", err)
		return
	}

	// Stampa il contenuto del file (debug)
	fmt.Println("Contenuto del file ordine.txt:")
	for i, line := range lines {
		fmt.Printf("%d: %s\n", i, line)
	}
	fmt.Println()

	// Analizza l'ordine
	ordine, err := parser.ParseOrdine(lines)
	if err != nil {
		handleError(err)
		return
	}

	// Formatta e stampa l'ordine
	output := formatter.FormatOrdine(ordine)
	fmt.Println(output)

	// Converte l'ordine in altri formati
	jsonOutput, err := converter.ToJSON(ordine)
	if err != nil {
		fmt.Printf("Errore nella conversione in JSON: %v\n", err)
	} else {
		fmt.Println("\nFormato JSON:")
		fmt.Println(jsonOutput)
	}
}

// legge il file di input e restituisce le righe come slice di stringhe
func readInputFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("errore nell'apertura del file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Aggiungi solo righe non vuote
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("errore nella lettura del file: %w", err)
	}

	return lines, nil
}

// Gstisce gli errori in modo specifico in base al tipo
func handleError(err error) {
	// Controlla se è un errore personalizzato
	if orderErr, ok := err.(*errors.OrderError); ok {
		// Output formattato per gli errori personalizzati
		fmt.Printf("=== ERRORE [%d] ===\n", orderErr.Code)
		fmt.Printf("Messaggio: %s\n", orderErr.Message)
		fmt.Printf("Soluzione: %s\n", orderErr.Details)

		// Suggerimenti aggiuntivi in base al codice di errore
		switch orderErr.Code {
		case errors.ErrCodePiattoEsaurito:
			fmt.Println("Suggerimento: Controllare il menu per piatti alternativi disponibili.")
		case errors.ErrCodeModificaNonValida:
			fmt.Println("Suggerimento: Consultare il personale per le modifiche consentite.")
		case errors.ErrCodeSintassiGenerale:
			fmt.Println("Suggerimento: Verificare la sintassi del file di ordine.")
		}
	} else {
		// Output standard per gli errori non personalizzati
		fmt.Printf("Errore: %v\n", err)
	}
}
