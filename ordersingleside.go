package orderbook

import (
	"container/list"
	"log"
	"math"

	"github.com/google/btree"
	"github.com/shopspring/decimal"
)

type OrderSingleSide struct {
	priceTree *btree.BTreeG[decimal.Decimal]
	priceMap  map[string]*OrderQueue   // price => order queue
	orderMap  map[uint64]*list.Element // order id => order queue element
	Side      OrderSide
}

func Less(a, b decimal.Decimal) bool {
	return a.LessThan(b)
}

func NewOrderSingleSide(side OrderSide) *OrderSingleSide {
	return &OrderSingleSide{
		btree.NewG(2, Less),
		map[string]*OrderQueue{},
		map[uint64]*list.Element{},
		side,
	}
}

func (os *OrderSingleSide) Append(o *Order) {
	_, exist := os.priceTree.ReplaceOrInsert(o.Price)
	if exist {
		log.Println("exist price:", o.Price)
	} else {
		os.priceMap[o.Price.String()] = NewOrderQueue(o.Price)
	}
	element := os.priceMap[o.Price.String()].Append(o)
	os.orderMap[o.ID] = element
}

func (os *OrderSingleSide) Remove(o *Order) bool {
	element, ok := os.orderMap[o.ID]
	if !ok {
		return false
	}

	if q, ok := os.priceMap[o.Price.String()]; ok {
		q.Remove(element)
		if q.Len() == 0 {
			os.priceTree.Delete(o.Price)
		}
		return true
	}

	return false
}

func (os *OrderSingleSide) Debug(prefix string) {
	mapFunc := func(item decimal.Decimal) bool {
		price := item
		oq := os.priceMap[price.String()]
		log.Printf("\t%v%v x %v", prefix, oq.Price, oq.Volume)
		return true
	}

	os.priceTree.DescendRange(decimal.NewFromInt(math.MaxInt64), decimal.Zero, mapFunc)
}

func (os *OrderSingleSide) Match(o *Order) *Order {
	ordersDone := []*Order{}
	orderPartial := o.Clone()

	mapFunc := func(price decimal.Decimal) bool {
		if os.Side == SELL {
			if price.GreaterThan(o.Price) {
				return false
			}
		} else {
			if price.LessThan(o.Price) {
				return false
			}
		}

		oq := os.priceMap[price.String()]

		for {
			o2e := oq.Orders.Front()
			if o2e == nil {
				break
			}

			o2 := o2e.Value.(*Order).Clone()
			orderPartial.Qty = orderPartial.Qty.Sub(o2.Qty)
			if orderPartial.Qty.GreaterThan(decimal.Zero) {
				oq.Remove(o2e)
				ordersDone = append(ordersDone, o2)
			} else if orderPartial.Qty.Equal(decimal.Zero) {
				oq.Remove(o2e)
				ordersDone = append(ordersDone, o2)
				ordersDone = append(ordersDone, orderPartial)
				break
			} else {
				o2.Qty = orderPartial.Qty.Abs()
				oq.Update(o2e, o2)
				orderPartial.Qty = decimal.Zero
				ordersDone = append(ordersDone, orderPartial)
				break
			}
		}
		os.Debug("* ")
		log.Println("ordersDone: ", ordersDone)
		return false
	}

	if os.Side == SELL {
		os.priceTree.AscendRange(decimal.Zero, decimal.NewFromInt(math.MaxInt64), mapFunc)
	} else {
		os.priceTree.DescendRange(decimal.NewFromInt(math.MaxInt64), decimal.Zero, mapFunc)
	}
	return orderPartial
}
