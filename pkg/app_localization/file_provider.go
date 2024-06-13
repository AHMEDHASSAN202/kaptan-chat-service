package app_localization

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadLocalizationFiles(appType string) map[string]interface{} {
	// Directory containing your JSON files based on App type
	dirPath := fmt.Sprintf("./pkg/app_localization/%s", appType)

	// Map to hold the parsed data
	data := make(map[string]interface{})

	// List and read all JSON files in the directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .json extension
		if !info.IsDir() && filepath.Ext(info.Name()) == ".json" {
			parsedData, err := readJSONFile(path)
			if err != nil {
				log.Printf("Failed to read file %s: %v", path, err)
			} else {
				key := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
				data[key] = parsedData
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", dirPath, err)
	}

	return data
}

func readJSONFile(filePath string) (map[string]interface{}, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file contents
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}
