package pkg

import "fmt"

type HasMoneyState struct {
	wending *Wending
}

func (h *HasMoneyState) AddItem(count int) error {
	return fmt.Errorf("item dispense in progres")
}

func (h *HasMoneyState) RequestItem() error {
	return fmt.Errorf("item dispense in progres")
}

func (h *HasMoneyState) InsertMoney(money int) error {
	return fmt.Errorf("item out of stock")
}

func (h *HasMoneyState) DispenseItem() error {
	fmt.Println("Dispensing item")

	h.wending.itemCount = h.wending.itemCount - 1

	if h.wending.itemCount == 0 {
		h.wending.SetState(h.wending.noItem)
	} else {
		h.wending.SetState(h.wending.hasItem)
	}

	return nil
}
