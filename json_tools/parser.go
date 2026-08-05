package json_tools

import (
	"encoding/json"
	"fmt"
	"os"
)

func Parse(fileLocation string, object any) {
	file, err := os.ReadFile(fileLocation)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", fileLocation, err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &object)
	if err != nil {
		fmt.Printf("Error during parsing %s: %v\n", fileLocation, err)
		os.Exit(1)
	}
}
