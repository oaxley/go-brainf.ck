package main

// ----- structures
type stack []uint16

// ----- functions

// push a new value to the stack
func (s stack) Push(v uint16) stack {
	return append(s, v)
}

// pop a new value from the stack
func (s stack) Pop() (stack, uint16) {
	l := len(s)
	return s[:l-1], s[l-1]
}

// number of items in the stack
func (s stack) Count() int {
	return len(s)
}
