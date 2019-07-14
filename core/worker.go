package core

import (
	"account-sync/initialize/postgers"
	"account-sync/initialize/rabbitmq"
	"account-sync/models"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
)

func RunWorker() {
	queue, err := rabbitmq.Listen()
	if err != nil {
		panic("Failed to establish listener")
	}
	for delivery := range queue {
		// deserialize to object
		inputMessage, err := prepareMessage(delivery.Body)
		if err != nil {
			continue
		}

		// работаем с одним пользователем
		if inputMessage.SingleUser != nil {
			var account models.Account
			postgers.AccountDatabase.First(&account)
		}

		// работаем с переводом со счета на счет
		if inputMessage.Transfer != nil {

		}
	}
}

func prepareMessage(body []byte) (*models.InputMessage, error) {
	var err error
	err = validateMessage(body)
	if err != nil {
		errorHandler(body, err.Error())
		return nil, err
	}

	var message models.InputMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		// пакет невозможно распарсить
		errorHandler(body, err.Error())
		return nil, err
	}
	return &message, nil
}

func validateMessage(body []byte) error {
	var err error
	var schemeBytes []byte
	contentLoader := gojsonschema.NewStringLoader(string(body))

	// генерируем схему валидации
	schemeBytes, err = json.Marshal(models.InputMessageValidator)
	schemeContent := string(schemeBytes)
	schemaLoader := gojsonschema.NewStringLoader(schemeContent)

	// валидируем json
	var validationResult *gojsonschema.Result
	validationResult, err = gojsonschema.Validate(schemaLoader, contentLoader)
	if err != nil {
		return err
	}

	if !validationResult.Valid() {
		var buffer bytes.Buffer
		for _, e := range validationResult.Errors() {
			buffer.WriteString(e.String() + "; ")
		}
		return errors.New(buffer.String())
	}
	return nil
}

func errorHandler(body []byte, reason string) {
	outputMessage := models.OutputMessage{
		Status:        "error",
		ErrorMessage:  reason,
		OriginalBytes: body,
	}
	err := rabbitmq.EnqueueMessage(&outputMessage)
	if err != nil {
		panic("Can't send output message: " + err.Error())
	}
}
