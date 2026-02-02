package json_tools

import (
	"encoding/json"
	"fmt"
	"os"
)

func Parse(fileLocation string, object any) {
	file, err := os.ReadFile(fileLocation)
	if os.IsNotExist(err) {
		fmt.Print("File ", fileLocation, " not found in directory")
		os.Exit(1)
	}

	err = json.Unmarshal(file, &object)
	if err != nil {
		fmt.Print("Error during parsing ", fileLocation)
		os.Exit(1)
	}
}
