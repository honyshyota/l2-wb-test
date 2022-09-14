package main

import (
	"fmt"
)

type State interface {
	AddItem(int)
	DispenseItem(int)
}

type Context struct {
	noIntem State
	hasItem State
	current State
	count   int
}

func NewContext() *Context {
	ctx := &Context{
		count: 1,
	}
	noItemState := &noItemState{ctx}
	hasItemState := &hasItemState{ctx}
	ctx.hasItem = hasItemState
	ctx.noIntem = noItemState
	ctx.current = hasItemState

	return ctx
}

func (ctx *Context) AddItem(count int) {
	ctx.current.AddItem(count)
}

func (ctx *Context) DispenseItem(count int) {
	ctx.current.DispenseItem(count)
}

func (ctx *Context) SetState(s State) {
	ctx.current = s
}

func (ctx *Context) IncrementCount(count int) {
	ctx.count += count
}

func (ctx *Context) PrintCurrentState() {
	if ctx.current == ctx.noIntem {
		fmt.Println("no item")
	} else if ctx.current == ctx.hasItem {
		fmt.Println("has item")
	}
}

type noItemState struct {
	ctx *Context
}

func (n *noItemState) AddItem(count int) {
	fmt.Println("Добавлено ", count, "предметов")
	n.ctx.SetState(n.ctx.hasItem)
	n.ctx.IncrementCount(count)
}

func (n *noItemState) DispenseItem(count int) {
	fmt.Println("Нет предмета для выдачи")
}

type hasItemState struct {
	ctx *Context
}

func (h *hasItemState) AddItem(count int) {
	fmt.Println("Предмет уже присутствует, невозможно добавить")
}

func (h *hasItemState) DispenseItem(count int) {
	if h.ctx.count < count {
		fmt.Println("недостаточто предметов для выдачи")
		if h.ctx.count == 0 {
			h.ctx.SetState(h.ctx.noIntem)
		}
	} else {
		fmt.Println("Предметы выданы")
		h.ctx.count = h.ctx.count - count
		if h.ctx.count == 0 {
			h.ctx.SetState(h.ctx.noIntem)
		}
	}
}

func main() {
	ctx := NewContext()
	ctx.PrintCurrentState()
	ctx.DispenseItem(1)
	ctx.PrintCurrentState()
	ctx.DispenseItem(1)
	ctx.PrintCurrentState()
	ctx.AddItem(2)
	ctx.PrintCurrentState()
	ctx.DispenseItem(2)
}
