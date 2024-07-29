package utils

import "gorm.io/gorm"

// WithTransaction executes the given function within a transaction.
// GORM will automatically handle beginning the transaction, committing it if the function succeeds, or rolling it back if an error occurs.
func WithTransaction(db *gorm.DB, txFunc func(tx *gorm.DB) error) error {
	return db.Transaction(txFunc)
}
