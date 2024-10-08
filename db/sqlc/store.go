package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTransaction(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries:New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTransaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return  err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID 	int64 `json:"from_account_id"`
	ToAccountID 	int64 `json:"to_account_id"`
	Amount 	int64 `json:"amount"`
}

// TransferTxResult contains the result of the transfer transaction
type TransferTxResult struct {
	Transfer 	Transfer `json:"transfer"`
	FromAccount 	Account `json:"from_account"`
	ToAccount 	Account `json:"to_account"`
	FromEntry 	Entry `json:"from_entry"`
	ToEntry 	Entry `json:"to_entry"`
}


// TransferTransaction performs a money transfer from oe account to the other.
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *SQLStore) TransferTransaction(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTransaction(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return  err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return  err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return  err
		}

		// Get account -> update its balance

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
	} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
	}
		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context, q *Queries, account1ID int64, amount1 int64,
	account2ID int64, amount2 int64) (account1 Account, account2 Account, err error)  {

		account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: account1ID,
			Amount: amount1,
		})
		if err != nil {
			return
		}

		account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: account2ID,
			Amount: amount2,
		})
		if err != nil {
			return
		}

		return
}