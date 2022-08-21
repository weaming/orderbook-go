package orderbook

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type OrderAction uint

const (
	ADD OrderAction = iota
	DEL
)

type OrderSide uint

const (
	BUY OrderSide = iota
	SELL
)

type OrderType uint

const (
	LIMIT OrderType = iota
	MARKET
)

type OrderState uint

const (
	INIT OrderState = iota
	PARTIAL_MATCHED
	FINISHED
	CANCELLED
)

type Order struct {
	Action OrderAction
	ID     uint64
	Side   OrderSide
	Price  decimal.Decimal // use decimal.Zero to indicate market order
	Qty    decimal.Decimal
	// Timestamp time.Time
}

func (o *Order) Clone() *Order {
	return &Order{
		Action: o.Action,
		ID:     o.ID,
		Side:   o.Side,
		Price:  o.Price,
		Qty:    o.Qty,
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("<Order Action:%v, ID:%v, Side:%v, Price:%v, Qty:%v>", o.Action, o.ID, o.Side, o.Price, o.Qty)
}
