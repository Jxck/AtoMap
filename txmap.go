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

func (atomap *AtoMap) Get(key int) int {
	result := make(chan int)
	request := request{
		requestType: GET,
		key:         key,
		result:      result,
	}
	atomap.Tx <- request
	return <-result
}

func (atomap *AtoMap) Set(key int, value int) {
	request := request{
		requestType: SET,
		key:         key,
		value:       value,
	}
	atomap.Tx <- request
}

func (atomap *AtoMap) Lock() *AtoMap {
	tx := make(chan request)
	request := request{
		requestType: LOCK,
		tx:          tx,
	}
	atomap.Tx <- request
	return &AtoMap{tx}
}

func (atomap *AtoMap) Unlock() {
	request := request{
		requestType: UNLOCK,
	}
	atomap.Tx <- request
}
