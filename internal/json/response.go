package json

import (
	"encoding/json"
	"fmt"
)

type (
	Envelope map[string]any
)

// LogMessage display a pretty struct
func LogMessage(data any, message string) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(message, "\n", err.Error())
	}
	fmt.Println(message, "\n", string(dataJSON))
}
