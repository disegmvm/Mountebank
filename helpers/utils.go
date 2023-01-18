package helpers

import "encoding/json"

func ToJSON(input any) string {
	inputData, _ := json.Marshal(input)

	return string(inputData)
}
