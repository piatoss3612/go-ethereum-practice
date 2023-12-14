package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type StandardJsonInput struct {
	Language interface{} `json:"language"`
	Sources  Sources     `json:"sources"`
	Settings interface{} `json:"settings"`
}

type Sources map[string]Source

type Source struct {
	Keccak256 string `json:"keccak256"`
	Content   string `json:"content"`
}

func main() {
	f, err := os.Open("build/MyToken_meta.json")
	handleErr(err)

	defer f.Close()

	var metadata map[string]interface{}

	err = json.NewDecoder(f).Decode(&metadata)
	handleErr(err)

	sources := make(Sources)

	metaSources := metadata["sources"].(map[string]interface{})

	for k, v := range metaSources {
		sources[k] = Source{
			Keccak256: v.(map[string]interface{})["keccak256"].(string),
			Content:   v.(map[string]interface{})["content"].(string),
		}
	}

	settings := metadata["settings"].(map[string]interface{})

	delete(settings, "compilationTarget")

	standardJsonInput := StandardJsonInput{
		Language: metadata["language"],
		Sources:  sources,
		Settings: settings,
	}

	standardJsonInputBytes, err := json.MarshalIndent(standardJsonInput, "", "  ")
	handleErr(err)

	// generate directory to current working directory if it doesn't exist
	path := filepath.Join(".", "verify")

	err = os.MkdirAll(path, os.ModePerm)
	handleErr(err)

	// create file in directory
	output, err := os.Create(filepath.Join(path, "MyToken_input.json"))
	handleErr(err)

	defer output.Close()

	// write to file
	_, err = output.Write(standardJsonInputBytes)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
