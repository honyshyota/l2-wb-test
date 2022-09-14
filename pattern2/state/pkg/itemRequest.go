package pkg

import "fmt"

type ItemRequestState struct {
	wending *Wending
}

func (i *ItemRequestState) AddItem(count int) error {
	return fmt.Errorf("item dispense in progres")
}

func (i *ItemRequestState) RequestItem() error {
	return fmt.Errorf("item already requested")
}

func (i *ItemRequestState) InsertMoney(money int) error {
	if money < i.wending.itemPrice {
		return fmt.Errorf("inserted money is less. Please insert %d", i.wending.itemPrice)
	}
	fmt.Println("money entered is OK")
	i.wending.SetState(i.wending.hasMoney)
	return nil
}

func (i *ItemRequestState) DispenseItem() error {
	return fmt.Errorf("please insert miney first")
}
