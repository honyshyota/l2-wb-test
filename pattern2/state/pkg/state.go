package pkg

type State interface {
	AddItem(int) error
	RequestItem() error
	InsertMoney(int) error
	DispenseItem() error
}
