package models

type ChangedStatus struct {
	OrderId int
	UserId  int
	Status  string
	Email   string
}

type NotifyInfo struct {
	UserEmail   string
	OrderStatus string
	OrderId     int
	Address     string
	UserName    string
	OrderData   Order
}
