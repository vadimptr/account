package core

import (
	"account-sync/initialize/rabbitmq"
	"account-sync/initialize/services"
	"account-sync/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/xeipuuv/gojsonschema"
)

func RunWorker() {
	var err error
	var queue <-chan amqp.Delivery
	queue, err = rabbitmq.Listen()
	if err != nil {
		panic("Failed to establish listener")
	}
	for {
		delivery, valid := <- queue
		if !valid {
			// выходим если канал закрыт
			fmt.Println("   Channel closed.")
			break
		}

		fmt.Printf("\nRecieved new message.\n")
		// deserialize to object
		var inputMessage *models.InputMessage
		inputMessage, err = prepareMessage(delivery.Body)
		if err != nil {
			errorHandler(delivery.Body, err.Error())
			delivery.Ack(false)
			continue
		}

		fmt.Printf("   Message valid. ")

		// работаем с одним пользователем
		if inputMessage.SingleUser != nil {
			singleUser := inputMessage.SingleUser
			fmt.Printf("%v\n", *singleUser)
			err = services.BalanceProcessor.ProcessSingleUser(singleUser.User, singleUser.Amount)
			if err != nil {
				errorHandler(delivery.Body, err.Error())
				delivery.Ack(false)
				continue
			}
			fmt.Printf("   Message processed.\n")
		}

		// работаем с переводом со счета на счет
		if inputMessage.Transfer != nil {
			transfer := inputMessage.Transfer
			fmt.Printf("%v\n", *transfer)
		}

		successHandler(inputMessage)
		delivery.Ack(false)
	}
}

func prepareMessage(body []byte) (*models.InputMessage, error) {
	var err error
	err = validateMessage(body)
	if err != nil {
		return nil, err
	}

	var message models.InputMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
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
	fmt.Printf("   Pushed response: %s %s\n", outputMessage.Status, reason)
}

func successHandler(message *models.InputMessage) {
	outputMessage := models.OutputMessage{
		Status:          "success",
		OriginalMessage: message,
	}
	err := rabbitmq.EnqueueMessage(&outputMessage)
	if err != nil {
		panic("Can't send output message: " + err.Error())
	}
	fmt.Printf("   Pushed response: %s\n", outputMessage.Status)
}
