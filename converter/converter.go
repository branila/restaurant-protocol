package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/branila/restaurant-protocol/models"
	"gopkg.in/yaml.v2"
)

func ToJSON(ordine models.Ordine) (string, error) {
	data, err := json.MarshalIndent(ordine, "", "  ")
	if err != nil {
		return "", fmt.Errorf("errore nella conversione in JSON: %w", err)
	}
	return string(data), nil
}

func ToYAML(ordine models.Ordine) (string, error) {
	data, err := yaml.Marshal(ordine)
	if err != nil {
		return "", fmt.Errorf("errore nella conversione in YAML: %w", err)
	}
	return string(data), nil
}

func ToXML(ordine models.Ordine) (string, error) {
	data, err := xml.MarshalIndent(ordine, "", "  ")
	if err != nil {
		return "", fmt.Errorf("errore nella conversione in XML: %w", err)
	}
	return xml.Header + string(data), nil
}

// Converte un ordine nel formato specificato
// formats supportati: "json", "yaml", "xml"
func ToPrettyFormat(ordine models.Ordine, format string) (string, error) {
	switch format {
	case "json":
		return ToJSON(ordine)
	case "yaml":
		return ToYAML(ordine)
	case "xml":
		return ToXML(ordine)
	default:
		return "", fmt.Errorf("formato non supportato: %s", format)
	}
}
