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

	// ищем аккаунт
	var account models.Account
	err = transaction.First(&account, "name = ?", name).Error
	if err != nil {
		// если прибавляем и аккаунт не найден
		if amount > 0 && gorm.IsRecordNotFoundError(err) {
			// аккаунта нет, нужно создавать
			account = models.Account{
				Name:   name,
				Amount: amount,
			}
			err = transaction.Create(&account).Error
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
	}

	// если операция вычитания и на счету недостаточно средств
	if amount < 0 && account.Amount < -amount {
		// если недостаточно на счету
		transaction.Rollback()
		return errors.New(fmt.Sprintf("not enough amount. exist %d try to subsctruct %d", account.Amount, amount))
	}

	// собственно операция по изменению баланса
	account.Amount += amount
	err = transaction.Save(&account).Error
	if err != nil {
		// не удалось обновить аккаунт
		transaction.Rollback()
		return err
	}

	err = transaction.Commit().Error
	if err != nil {
		// не удалось применить транзакцию
		transaction.Rollback()
		return err
	}

	fmt.Printf("   Now ballance: %v\n", account)

	// успех
	return nil
}

func (SerializableBalanceProcessor) ProcessTransfer(string, string, int) error {
	return nil
}
