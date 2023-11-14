package gopher_mart

type Market interface {
	// users
	LoginUser()
	RegisterUser()

	// orders
	SendOrder()
	GetAllOrders()
	GetConfirmedOrders()
	GetDoneOrders()
	ConfirmOrder()
	DoneOrder()
}
