package transaction

type Rekeying interface {
	// Rekey the account
	Rekey(from, to string) error

	// If the account is rekeyed, then it will rekey the account back to itself
	RekeyBack(address string) error
}
