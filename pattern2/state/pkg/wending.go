package pkg

type Wending struct {
	hasItem      State
	itemRequest  State
	hasMoney     State
	noItem       State
	currentState State
	itemCount    int
	itemPrice    int
}

func NewWending(count int, price int) *Wending {
	w := &Wending{
		itemCount: count,
		itemPrice: price,
	}
	hasItemState := &HasItemState{w}
	itemRequestState := &ItemRequestState{w}
	hasMoneyState := &HasMoneyState{w}
	noItemState := &NoItemState{w}
	w.SetState(hasItemState)
	w.hasItem = hasItemState
	w.itemRequest = itemRequestState
	w.hasMoney = hasMoneyState
	w.noItem = noItemState
	return w
}

func (w *Wending) AddItem(count int) error {
	return w.currentState.AddItem(count)
}

func (w *Wending) RequestItem() error {
	return w.currentState.RequestItem()
}

func (w *Wending) InsertMoney(money int) error {
	return w.currentState.InsertMoney(money)
}

func (w *Wending) DispenseItem() error {
	return w.currentState.DispenseItem()
}

func (w *Wending) SetState(s State) {
	w.currentState = s
}

func (w *Wending) IncrementItemCount(count int) {
	w.itemCount += count
}
