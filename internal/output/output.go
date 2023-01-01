package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func saveAsJSON(data map[string][]string, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func saveAsText(data map[string][]string, w io.Writer) error {
	for _, ipAddrs := range data {
		for _, ip := range ipAddrs {
			if _, err := fmt.Fprintf(w, "%s\n", ip); err != nil {
				return err
			}
		}
	}
	return nil
}

func saveAsCSV(data map[string][]string, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	for hostname, ipAddrs := range data {
		for _, ip := range ipAddrs {
			if err := writer.Write([]string{hostname, ip}); err != nil {
				return err
			}
		}
	}
	return nil
}

// SaveToFile saves the DNS resolution results to a file in the specified format

func SaveToFile(data map[string][]string, filename, format string) error {
	if format == "" {
		format = "txt"
	}
	filename += "." + format
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	switch format {
	case "json":
		return saveAsJSON(data, file)
	case "text", "txt":
		return saveAsText(data, file)
	case "csv":
		return saveAsCSV(data, file)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
