package repositories

import (
	"context"
	"database/sql"
)

func BeginTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	ctx := context.Background()
	transaction, err := resultsRepository.dbHandler.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	runnersRepository.SetTransaction(transaction)
	resultsRepository.SetTransaction(transaction)

	return nil
}

func RollbackTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	transaction := runnersRepository.GetTransaction()

	runnersRepository.ClearTransaction()
	resultsRepository.ClearTransaction()

	return transaction.Rollback()
}

func CommitTransaction(runnersRepository *RunnersRepository, resultsRepository *ResultsRepository) error {
	transaction := runnersRepository.GetTransaction()

	runnersRepository.ClearTransaction()
	resultsRepository.ClearTransaction()

	return transaction.Commit()
}
