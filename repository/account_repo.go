package repository

import (
	"BankingLedgerSystem/models"
	"context"
	"database/sql"
	"errors"
)

// AccountRepository defines all DB operations for accounts
type AccountRepository interface {
	InsertAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error)
	FindAccountByID(ctx context.Context, id string) (*models.Account, error)
}

type PostgresAccountRepo struct {
	db *sql.DB
}

func NewPostgresAccountRepo(db *sql.DB) *PostgresAccountRepo {
	return &PostgresAccountRepo{db: db}
}

func (r *PostgresAccountRepo) InsertAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error) {
	query := `INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING id, name, balance`
	row := r.db.QueryRowContext(ctx, query, req.Name, req.Balance)

	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.Name, &acc.Balance); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *PostgresAccountRepo) FindAccountByID(ctx context.Context, id string) (*models.Account, error) {
	query := `SELECT id, name, balance FROM accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.Name, &acc.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &acc, nil
}
