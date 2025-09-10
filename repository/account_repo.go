package repository

import (
	"BankingLedgerSystem/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

//go:generate sh -c "sh $(git rev-parse --show-toplevel)/scripts/mock_generator.sh $GOFILE"

// AccountRepository defines all DB operations for accounts
type AccountRepository interface {
	InsertAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error)
	FindAccountByID(ctx context.Context, id string) (*models.Account, error)
}

// PostgresAccountRepo implements AccountRepository using pgxpool
type PostgresAccountRepo struct {
	db *ConnectionProvider
}

func NewPostgresAccountRepo(db *ConnectionProvider) *PostgresAccountRepo {
	return &PostgresAccountRepo{db: db}
}

func (r *PostgresAccountRepo) InsertAccount(ctx context.Context, req models.CreateAccountRequest) (*models.Account, error) {
	query := `INSERT INTO accounts (account_number, owner_name, balance) 
	          VALUES ($1, $2, $3) 
	          RETURNING id, account_number, owner_name, balance, created_at`

	row := r.db.GetConnection(ctx).QueryRow(ctx, query, req.AccountNumber, req.Name, req.Balance)

	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.AccountNumber, &acc.Name, &acc.Balance, &acc.CreatedAt); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *PostgresAccountRepo) FindAccountByID(ctx context.Context, id string) (*models.Account, error) {
	query := `SELECT id, account_number, owner_name, balance, created_at 
	          FROM accounts WHERE account_number = $1`

	row := r.db.GetConnection(ctx).QueryRow(ctx, query, id)

	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.AccountNumber, &acc.Name, &acc.Balance, &acc.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &acc, nil
}
