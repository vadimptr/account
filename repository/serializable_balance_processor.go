package repository

import (
	"account-sync/initialize/postgers"
	"account-sync/models"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
)

// Реализация механизма обновления аккаунта посредством транзакции с уровнем изоляции - serializable
type SerializableBalanceProcessor struct{}

func (SerializableBalanceProcessor) ProcessSingleUser(name string, amount int) error {
	transaction := postgers.AccountDatabase.Begin()
	err := transaction.Error
	if err != nil {
		// не удалось открыть транзакцию
		transaction.Rollback()
		return err
	}

	err = transaction.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE ").Error
	if err != nil {
		// не удалось установить транзакцию в режим serializable
		transaction.Rollback()
		return err
	}

	var account models.Account
	if amount < 0 {
		// ищем аккаунт
		err = transaction.First(&account, "name = ?", name).Error
		if err != nil {
			// не удалось выполнить поиск аккаунта. либо аккаунт не был найден
			transaction.Rollback()
			return err
		}
		if account.Amount < amount {
			// если недостаточно на счету
			transaction.Rollback()
			return errors.New(fmt.Sprintf("not enough amount. exist %d try to subsctruct %d", account.Amount, amount))
		}
	} else {
		// ищем аккаунт и аккаунта может не быть (ленивое создание аккаунта)
		err = transaction.First(&account, "name = ?", name).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				// аккаунта нет, нужно создавать
				account = models.Account{
					Name:   name,
					Amount: amount,
				}
				err = transaction.Create(account).Error
				if err != nil {
					// не удалось добавить новый аккаунт
					transaction.Rollback()
					return err
				}
			} else {
				// не удалось выполнить поиск аккаунта
				transaction.Rollback()
				return err
			}
		} else {
			account.Amount -= amount
			err = transaction.Update(&account).Error
			if err != nil {
				// не удалось обновить аккаунт
				transaction.Rollback()
				return err
			}
		}
	}

	err = transaction.Commit().Error
	if err != nil {
		// не удалось применить транзакцию
		transaction.Rollback()
		return err
	}

	// успех
	return nil
}

func (SerializableBalanceProcessor) ProcessTransfer(string, string, int) error {
	return nil
}
