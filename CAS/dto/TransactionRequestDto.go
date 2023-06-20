package dto

type TransactionRequestDto struct {
	TransactionId int
	SenderId      int
	RecipientId   int
	AmountMoney   float64
	Status        string
}
