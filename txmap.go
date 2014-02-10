package main

import (
	"log"
	"runtime"
)

var debug func(a ...interface{})

func init() {
	log.SetFlags(log.Lshortfile)
	debug = func(a ...interface{}) {
		log.Printf("%+v", a)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type RequestType int

const (
	SET RequestType = iota
	GET
	BEGINTX
	ENDTX
)

type request struct {
	requestType RequestType
	key         int
	value       int
	result      chan int
	tx          chan request
}

type TxMap struct {
	m       map[int]int
	request chan request
}

func NewTxMap() *TxMap {
	txm := &TxMap{
		m:       make(map[int]int),
		request: make(chan request),
	}
	go txm.Handle(txm.request)
	return txm
}

func (txm *TxMap) Handle(requestChan chan request) {
	for {
		req := <-requestChan
		switch req.requestType {
		case GET:
			req.result <- txm.m[req.key]
		case SET:
			old := txm.m[req.key]
			txm.m[req.key] = req.value
			req.result <- old
		case BEGINTX:
			txm.Handle(req.tx)
		case ENDTX:
			return
		}
	}
}

func (txm *TxMap) Set(key int, value int) int {
	result := make(chan int)
	request := request{
		requestType: SET,
		key:         key,
		value:       value,
		result:      result,
	}
	txm.request <- request // lock
	return <-result        // unlock
}

func (txm *TxMap) Get(key int) int {
	result := make(chan int)
	request := request{
		requestType: GET,
		key:         key,
		result:      result,
	}
	txm.request <- request // lock
	return <-result        // unlock
}

func (txm *TxMap) BeginTx() {
	tx := make(chan request)
	request := request{
		requestType: BEGINTX,
		tx:          tx,
	}
	txm.request <- request
}

func (txm *TxMap) EndTx() {
	request := request{
		requestType: ENDTX,
	}
	txm.request <- request
}

func main() {
	txm := NewTxMap()
	txm.BeginTx()
	log.Println(txm.Set(1, 1))
	log.Println(txm.Get(1))
	txm.EndTx()
}
