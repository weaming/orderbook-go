package orderbook

import (
	"container/list"

	"github.com/shopspring/decimal"
)

type OrderQueue struct {
	Price  decimal.Decimal
	Volume decimal.Decimal
	Orders *list.List
}

func NewOrderQueue(price decimal.Decimal) *OrderQueue {
	return &OrderQueue{
		Price:  price,
		Volume: decimal.Zero,
		Orders: list.New(),
	}
}

func (oq *OrderQueue) Len() int {
	return oq.Orders.Len()
}

func (oq *OrderQueue) Append(o *Order) *list.Element {
	oq.Volume = oq.Volume.Add(o.Qty)
	return oq.Orders.PushBack(o)
}

func (oq *OrderQueue) Update(e *list.Element, o *Order) *list.Element {
	oq.Volume = oq.Volume.Sub(e.Value.(*Order).Qty)
	oq.Volume = oq.Volume.Add(o.Qty)
	e.Value = o
	return e
}

func (oq *OrderQueue) Remove(e *list.Element) *Order {
	oq.Volume = oq.Volume.Sub(e.Value.(*Order).Qty)
	return oq.Orders.Remove(e).(*Order)
}
