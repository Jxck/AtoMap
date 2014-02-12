package txmap

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type RequestType int

const (
	SET RequestType = iota
	GET
	BEGINTX
	ENDTX
)

type Request struct {
	requestType RequestType
	key         int
	value       int
	result      chan int
	tx          chan Request
}

func NewTxMap() chan Request {
	txmap := make(chan Request)
	go func() {
		m := make(map[int]int)
		HandleRequests(m, txmap)
	}()
	return txmap
}

func HandleRequests(m map[int]int, r chan Request) {
	for {
		req := <-r
		switch req.requestType {
		case GET:
			req.result <- m[req.key]
		case SET:
			m[req.key] = req.value
		case BEGINTX:
			HandleRequests(m, req.tx)
		case ENDTX:
			return
		}
	}
}

func Get(m chan Request, key int) int {
	result := make(chan int)
	request := Request{
		requestType: GET,
		key:         key,
		result:      result,
	}
	m <- request
	return <-result
}

func Set(m chan Request, key int, value int) {
	request := Request{
		requestType: SET,
		key:         key,
		value:       value,
	}
	m <- request
}

func BeginTx(m chan Request) chan Request {
	tx := make(chan Request)
	request := Request{
		requestType: BEGINTX,
		tx:          tx,
	}
	m <- request
	return tx
}

func EndTx(m chan Request) {
	request := Request{
		requestType: ENDTX,
	}
	m <- request
}
