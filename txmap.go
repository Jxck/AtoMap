package txmap

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type RequestType int

const (
	Set RequestType = iota
	Get
	BeginTx
	EndTx
)

type Request struct {
	requestType RequestType
	key         int
	value       int
	ret         chan int
	tx          chan Request
}

type TxMap struct {
	req chan Request      // メインループ
	tx  chan chan Request // tx に入った時の退避用
}

func NewTxMap() *TxMap {
	txm := new(TxMap)
	txm.req = make(chan Request)
	txm.tx = make(chan chan Request, 1)
	go runMap(txm.req)
	return txm
}

func runMap(r chan Request) {
	m := make(map[int]int)
	HandleRequests(m, r)
}

func HandleRequests(m map[int]int, r chan Request) {
	for {
		req := <-r
		switch req.requestType {
		case Get:
			req.ret <- m[req.key]
		case Set:
			m[req.key] = req.value
		case BeginTx:
			HandleRequests(m, req.tx)
		case EndTx:
			return
		}
	}
}

func (txm *TxMap) Set(key int, value int) {
	txm.req <- Request{Set, key, value, nil, nil}
}

func (txm *TxMap) Get(key int) int {
	result := make(chan int)
	txm.req <- Request{Get, key, 0, result, nil}
	return <-result
}

func (txm *TxMap) BeginTx() {
	tmp := txm.req
	txm.tx <- tmp
	txm.req = make(chan Request)
	tmp <- Request{BeginTx, 0, 0, nil, txm.req}
}

func (txm *TxMap) EndTx() {
	txm.req <- Request{EndTx, 0, 0, nil, nil}
	txm.req = <-txm.tx
}

// func main() {
// 	r := make(chan Request)
// 	go runMap(r)
// 	tx := beginTx(r)
// 	set(tx, 1, "hoge")
// 	log.Println(get(tx, 1))
// 	endTx(tx)
// }
