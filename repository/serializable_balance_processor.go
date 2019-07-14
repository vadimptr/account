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
	errProcess := TransactionWrapper(func(transaction *gorm.DB) error {
		var err error

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
					return err
				}
			} else {
				// не удалось выполнить поиск аккаунта
				return err
			}
		}

		// если операция вычитания и на счету недостаточно средств
		if amount < 0 && account.Amount < -amount {
			// если недостаточно на счету
			return errors.New(fmt.Sprintf("not enough amount. exist %d try to subsctruct %d", account.Amount, amount))
		}

		// собственно операция по изменению баланса
		account.Amount += amount
		err = transaction.Save(&account).Error
		if err != nil {
			// не удалось обновить аккаунт
			return err
		}

		fmt.Printf("   Now ballance: %v\n", account)
		return nil
	})
	return errProcess
}

func (SerializableBalanceProcessor) ProcessTransfer(fromUser string, toUser string, amount int) error {
	errProcess := TransactionWrapper(func(transaction *gorm.DB) error {
		var err error

		// ищем аккаунт
		var accountTo models.Account
		err = transaction.First(&accountTo, "name = ?", toUser).Error
		if err != nil {
			// если аккаунт не найден
			if gorm.IsRecordNotFoundError(err) {
				accountTo = models.Account{
					Name:   toUser,
					Amount: amount,
				}
				err = transaction.Create(&accountTo).Error
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		var accountFrom models.Account
		err = transaction.First(&accountFrom, "name = ?", fromUser).Error
		if err != nil {
			// если аккаунта нет то просто бросаем ошибку
			return err
		}

		// если на счету недостаточно средств
		if accountFrom.Amount < amount {
			// если недостаточно на счету
			return errors.New(fmt.Sprintf("not enough amount. exist %d try to subsctruct %d", accountFrom.Amount, amount))
		}

		// списываем
		accountFrom.Amount -= amount
		err = transaction.Save(&accountFrom).Error
		if err != nil {
			return err
		}

		// добавляем
		accountTo.Amount += amount
		err = transaction.Save(&accountTo).Error
		if err != nil {
			return err
		}

		fmt.Printf("   Now ballance: %v %v\n", accountFrom, accountTo)
		return nil
	})
	return errProcess
}

func TransactionWrapper(some func(*gorm.DB) error) error {
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

	err = some(transaction)
	if err != nil {
		transaction.Rollback()
		return err
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
