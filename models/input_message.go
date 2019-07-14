package models

import "account-sync/validation"

type InputMessage struct {
	SingleUser *SingleUser `json:"single_user"`
	Transfer   *Transfer   `json:"transfer"`
}

type SingleUser struct {
	User   string `json:"user"`
	Amount int    `json:"amount"`
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
		"transfer":    TransferValidator,
	},
}

var SingleUserValidator = validation.Node{
	Type:               validation.ObjectType,
	AdditionalProperty: &validation.ValueFalse,
	Required:           []string{"user", "amount"},
	Properties: validation.Properties{
		"user":   UserValidator,
		"amount": AmountValidator,
	},
}

var UserValidator = validation.Node{
	Type:      validation.StringType,
	MinLength: &validation.Value1,
	MaxLength: &validation.Value256,
}

var AmountValidator = validation.Node{
	Type: validation.IntegerType,
	Not: &validation.Node{
		Enum: []interface{}{0},
	},
}

var TransferValidator = validation.Node{
	Type:               validation.ObjectType,
	AdditionalProperty: &validation.ValueFalse,
	Required:           []string{"from_user", "to_user", "amount"},
	Properties: validation.Properties{
		"from_user": UserValidator,
		"to_user":   UserValidator,
		"amount":    PositiveAmountValidator,
	},
}

var PositiveAmountValidator = validation.Node{
	Type:    validation.IntegerType,
	Minimum: &validation.Value1,
}
