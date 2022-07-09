package server

import "sync"

type Operator struct {
	ID       int
	CallLink string
	Busy     bool
}

type CallCenter struct {
	Operators map[int]Operator
	m         sync.RWMutex
	HasFree   chan int
}

func NewCallCenter() *CallCenter {
	return &CallCenter{
		Operators: make(map[int]Operator),
		m:         sync.RWMutex{},
		HasFree:   make(chan int, 4),
	}
}

func (cc *CallCenter) AddOperator(ID int, link string) {
	cc.m.Lock()
	defer cc.m.Unlock()

	cc.Operators[ID] = Operator{
		ID:       ID,
		CallLink: link,
		Busy:     false,
	}

	go func() {
		cc.HasFree <- ID
	}()
}

func (cc *CallCenter) SetFree(ID int) {
	cc.m.Lock()
	defer cc.m.Unlock()

	operator := cc.Operators[ID]
	operator.Busy = false
	cc.Operators[ID] = operator

	go func() {
		cc.HasFree <- ID
	}()
}

func (cc *CallCenter) SetBusy(ID int) {
	cc.m.Lock()
	defer cc.m.Unlock()

	operator := cc.Operators[ID]
	operator.Busy = true
	cc.Operators[ID] = operator
}

func (cc *CallCenter) IsBusy(ID int) bool {
	cc.m.RLock()
	defer cc.m.RUnlock()

	return cc.Operators[ID].Busy
}

func (cc *CallCenter) GetLink(ID int) string {
	cc.m.RLock()
	defer cc.m.RUnlock()

	return cc.Operators[ID].CallLink
}
