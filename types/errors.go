package types

import "fmt"

// ReceiptMissingError is returned when trying to fetch a receipt that does not exist
type ReceiptMissingError struct {
	Hash string
}

func (err *ReceiptMissingError) Error() string {
	return fmt.Sprintf("Receipt for Transaction %q missing", err.Hash)
}

// BlockMissingError is returned when trying to fetch a block that does not exist
type BlockMissingError struct {
	Number int64
}

func (err *BlockMissingError) Error() string {
	return fmt.Sprintf("Block %v missing", err.Number)
}
