package interfaces

type BalanceProcessor interface {
	ProcessSingleUser(string, int) error
	ProcessTransfer(string, string, int) error
}