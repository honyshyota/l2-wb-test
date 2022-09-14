package pkg

import "fmt"

type NoItemState struct {
	wending *Wending
}

func (n *NoItemState) AddItem(count int) error {
	n.wending.IncrementItemCount(count)
	n.wending.SetState(n.wending.hasItem)
	return nil
}

func (n *NoItemState) RequestItem() error {
	return fmt.Errorf("item out of stock")
}

func (n *NoItemState) InsertMoney(money int) error {
	return fmt.Errorf("item out of stock")
}

func (n *NoItemState) DispenseItem() error {
	return fmt.Errorf("item out of stock")
}
