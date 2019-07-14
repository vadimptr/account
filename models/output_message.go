package models

type OutputMessage struct {
	Status          string        `json:"status"`
	OriginalMessage *InputMessage `json:"original_message"`
	OriginalBytes   []byte        `json:"original_bytes"`
	ErrorMessage    string        `json:"error_message"`
}
