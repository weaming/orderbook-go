package orderbook

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func TestOrderbook(t *testing.T) {
	ob := NewOrderBook("ETH")

	for line := range Lines("order-queue.txt") {
		o := parseOrder(line)
		// log.Println(line)
		log.Printf("order %+v", o)
		ob.Match(o)
		ob.Debug()
	}
}

func Lines(file string) <-chan string {
	f, e := os.Open(file)
	if e != nil {
		panic(e)
	}

	ch := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(ch)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	return ch
}

func Uint64(s string) uint64 {
	v, e := strconv.ParseUint(s, 10, 64)
	if e == nil {
		return v
	}
	panic(e)
}

func parseOrder(line string) *Order {
	parts := strings.Split(line, ",")
	// TODO: check index error
	action_, orderID, side_, qty_, price_ := parts[0], parts[1], parts[2], parts[3], parts[4]

	action := ADD
	if action_ == "X" {
		action = DEL
	}
	side := BUY
	if side_ == "S" {
		side = SELL
	}

	qty, e := decimal.NewFromString(qty_)
	if e != nil {
		panic(e)
	}

	price, e := decimal.NewFromString(price_)
	if e != nil {
		panic(e)
	}

	return &Order{
		Action: action,
		ID:     Uint64(orderID),
		Side:   side,
		Price:  price,
		Qty:    qty,
		// Timestamp: time.Now(),
	}
}
