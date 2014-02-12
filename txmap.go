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
	ret         chan int
	tx          chan Request
}

func HandleRequests(m map[int]int, r chan Request) {
	for {
		req := <-r
		switch req.requestType {
		case GET:
			req.ret <- m[req.key]
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
	m <- Request{GET, key, 0, result, nil}
	return <-result
}

func Set(m chan Request, key int, value int) {
	m <- Request{SET, key, value, nil, nil}
}

func BeginTx(m chan Request) chan Request {
	tx := make(chan Request)
	m <- Request{BEGINTX, 0, 0, nil, tx}
	return tx
}

func EndTx(m chan Request) {
	m <- Request{ENDTX, 0, 0, nil, nil}
}

func RunMap(r chan Request) {
	m := make(map[int]int)
	HandleRequests(m, r)
}

func main() {
	r := make(chan Request)
	go RunMap(r)
	r = BeginTx(r)
	Set(r, 0, 1)
	log.Println(Get(r, 0))
	EndTx(r)
}
