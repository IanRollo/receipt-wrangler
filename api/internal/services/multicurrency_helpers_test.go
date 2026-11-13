package services

import (
	"encoding/json"

	"receipt-wrangler/api/internal/constants"
)

func jsonUnmarshalForTest(text string, target any) error {
	return json.Unmarshal([]byte(text), target)
}

func promptTextForTest() string {
	return constants.DefaultPromptText
}
