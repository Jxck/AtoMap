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

func Get(req chan Request, key int) int {
	result := make(chan int)
	request := Request{
		requestType: GET,
		key:         key,
		result:      result,
	}
	req <- request
	return <-result
}

func Set(req chan Request, key int, value int) {
	request := Request{
		requestType: SET,
		key:         key,
		value:       value,
	}
	req <- request
}

func BeginTx(req chan Request) chan Request {
	tx := make(chan Request)
	request := Request{
		requestType: BEGINTX,
		tx:          tx,
	}
	req <- request
	return tx
}

func EndTx(req chan Request) {
	request := Request{
		requestType: ENDTX,
	}
	req <- request
}
