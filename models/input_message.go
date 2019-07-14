package models

import "account-sync/validation"

type InputMessage struct {
	SingleUser *SingleUser `json:"single_user"`
	Transfer   *Transfer   `json:"transfer"`
}

type SingleUser struct {
	UserName string `json:"user_name"`
	Amount   int    `json:"amount"`
}

type Transfer struct {
	FromUser string `json:"from_user"`
	ToUser   string `json:"to_user"`
	Amount   int    `json:"amount"`
}

// объект валидации входного пакета
var InputMessageValidator = validation.Node{
	Type:               validation.ObjectType,
	AdditionalProperty: &validation.ValueFalse,
	OneOf: []validation.Node{
		{
			// в корне должен быть узел
			Required: []string{"single_user"},
		},
		{
			// или этот узел
			Required: []string{"transfer"},
		},
	},
	Properties: validation.Properties{
		"single_user": SingleUserValidator,
	},
}

var SingleUserValidator = validation.Node{
	Type:               validation.ObjectType,
	AdditionalProperty: &validation.ValueFalse,
	Required:           []string{"single_user", "amount"},
	Properties: validation.Properties{
		"user_name": UserNameValidator,
		"amount":    AmountValidator,
	},
}

var UserNameValidator = validation.Node{
	Type:      validation.StringType,
	MinLength: &validation.Value6,
	MaxLength: &validation.Value256,
}

var AmountValidator = validation.Node{
	Type: validation.IntegerType,
}

var TransferValidator = validation.Node{
	Type:               validation.ObjectType,
	AdditionalProperty: &validation.ValueFalse,
	Required:           []string{"from_user", "to_user", "amount"},
	Properties: validation.Properties{
		"from_user": UserNameValidator,
		"to_user":   UserNameValidator,
		"amount":    PositiveAmountValidator,
	},
}

var PositiveAmountValidator = validation.Node{
	Type:    validation.IntegerType,
	Minimum: &validation.Value1,
}
