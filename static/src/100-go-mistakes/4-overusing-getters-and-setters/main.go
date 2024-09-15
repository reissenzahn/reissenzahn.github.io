package main

// it is not considered mandatory or idiomatic to use getters and setters

// getters and setters present some advantages: encapsulating behavior associated with getting or setting a field, hide internal representation, provide a debugging interception point

// for a field balance, the getter method should be Balance() and the setter method should be SetBalance()

// do not overwhelm your code with getter and setters on structs if they don't bring any value

type Account struct {
	balance int
}

func (a Account) Balance() int {
	return a.balance
}

func (a *Account) SetBalance(balance int) {
	a.balance = balance
}

func main() {
	var account Account

	b := account.Balance()
	if b < 0 {
		account.SetBalance(0)
	}
}
