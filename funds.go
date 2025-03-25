package main

import (
	"fmt"
	"sync"
)

type Fund struct {
	mu    sync.Mutex
	total int
}

func (f *Fund) Add(amount int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.total += amount
}

func (f *Fund) Subtract(amount int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.total < amount {
		return fmt.Errorf("insufficient funds")
	}
	f.total -= amount
	return nil
}

func (f *Fund) Balance() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.total
}

func main() {
	fund := &Fund{}

	fund.Add(100)
	fmt.Println("Balance after adding 100:", fund.Balance())

	err := fund.Subtract(50)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Balance after subtracting 50:", fund.Balance())
	}

	err = fund.Subtract(100)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Balance after subtracting 100:", fund.Balance())
	}
}
