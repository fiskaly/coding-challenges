package persistence

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// PostgresStore implements persistence using PostgreSQL.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store.
func NewPostgresStore(connectionString string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

// Close closes the database connection.
func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// TODO: Implement device and transaction persistence methods
// Example methods to implement:
// - CreateDevice(device *domain.Device) error
// - GetDevices() ([]*domain.Device, error)
// - GetDeviceByID(id string) (*domain.Device, error)
// - UpdateDeviceCounter(id string, counter int) error
// - CreateTransaction(tx *domain.Transaction) error
// - GetTransactions(deviceID string) ([]*domain.Transaction, error)
