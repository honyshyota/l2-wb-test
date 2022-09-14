package pkg

import "fmt"

type HasItemState struct {
	wending *Wending
}

func (h *HasItemState) AddItem(count int) error {
	fmt.Printf("%d item added\n", count)
	h.wending.IncrementItemCount(count)
	return nil
}

func (h *HasItemState) RequestItem() error {
	if h.wending.itemCount == 0 {
		h.wending.SetState(h.wending.noItem)
		return fmt.Errorf("no item present")
	}

	fmt.Println("Item requested")
	h.wending.SetState(h.wending.itemRequest)
	return nil
}

func (h *HasItemState) InsertMoney(money int) error {
	return fmt.Errorf("please select item first")
}

func (h *HasItemState) DispenseItem() error {
	return fmt.Errorf("please select item first")
}
