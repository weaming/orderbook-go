package orderbook

import (
	"log"

	"github.com/shopspring/decimal"
)

type OrderBook struct {
	Symbol string
	Asks   *OrderSingleSide
	Bids   *OrderSingleSide
}

func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{symbol, NewOrderSingleSide(SELL), NewOrderSingleSide(BUY)}
}

func (ob *OrderBook) Match(o *Order) {
	oq, oq2 := ob.Bids, ob.Asks
	if o.Side == SELL {
		oq = ob.Asks
		oq2 = ob.Bids
	}

	switch o.Action {
	case ADD:
		if decimal.Zero.Equals(o.Price) {
			// TODO: market order
			return
		}

		orderPartial := oq2.Match(o)
		if !orderPartial.Qty.Equal(decimal.Zero) {
			oq.Append(orderPartial)
		}
	case DEL:
		oq.Remove(o)
	}
}

func (ob *OrderBook) Debug() {
	log.Println("==========================")
	ob.Asks.Debug("")
	log.Println("--------------------------")
	ob.Bids.Debug("")
	log.Println("===^^^^^^^^^^^^^^^^^^^^===")
}
