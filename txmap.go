package atomap

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
	LOCK
	UNLOCK
)

type request struct {
	requestType RequestType
	key         int
	value       int
	result      chan int
	tx          chan request
}

type AtoMap struct {
	Tx chan request
}

func NewAtoMap() *AtoMap {
	tx := make(chan request)
	go func() {
		m := make(map[int]int)
		handleRequests(m, tx)
	}()
	atomap := &AtoMap{
		Tx: tx,
	}
	return atomap
}

func handleRequests(m map[int]int, r chan request) {
	for {
		req := <-r
		switch req.requestType {
		case GET:
			req.result <- m[req.key]
		case SET:
			m[req.key] = req.value
		case LOCK:
			handleRequests(m, req.tx)
		case UNLOCK:
			return
		}
	}
}

func Get(req chan request, key int) int {
	result := make(chan int)
	request := request{
		requestType: GET,
		key:         key,
		result:      result,
	}
	req <- request
	return <-result
}

func Set(req chan request, key int, value int) {
	request := request{
		requestType: SET,
		key:         key,
		value:       value,
	}
	req <- request
}

func Lock(req chan request) chan request {
	tx := make(chan request)
	request := request{
		requestType: LOCK,
		tx:          tx,
	}
	req <- request
	return tx
}

func Unlock(req chan request) {
	request := request{
		requestType: UNLOCK,
	}
	req <- request
}
