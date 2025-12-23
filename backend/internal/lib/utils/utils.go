package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	fmt.Println("JSON:", string(json))
}
