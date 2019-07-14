package services

import (
	"account-sync/initialize/environment"
	"account-sync/interfaces"
	"account-sync/repository"
)

var BalanceProcessor interfaces.BalanceProcessor

func init() {
	if environment.GetMethod() == "serializable" {
		BalanceProcessor = repository.SerializableBalanceProcessor{}
	}
}