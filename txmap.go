package txmap

import (
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type RequestType int

const (
	SET Type = iota
	GET
	BEGINTX
	ENDTX
)

type Request struct {
	Type   Type
	Key    int
	Value  int
	Result chan int
	Tx     chan Request
}

type TxMap struct {
	Tx     chan Request      // メインループ
	Parent chan chan Request // tx に入った時の退避用
}

func NewTxMap() *TxMap {
	txm := &TxMap{
		Tx:     make(chan Request),
		Parent: make(chan chan Request, 1),
	}

	go func() {
		m := make(map[int]int) // 実際の map
		HandleRequests(m, txm.Tx)
	}()

	return txm
}

func HandleRequests(m map[int]int, tx chan Request) {
	for {
		req := <-tx
		switch req.RequestType {
		case Get:
			req.Result <- m[req.key]
		case Set:
			m[req.key] = req.Value
		case BeginTx:
			HandleRequests(m, req.Tx)
		case EndTx:
			return
		}
	}
}

func (txm *TxMap) Set(key int, value int) {
	txm.Tx <- Request{
		RequestType: SET,
		Key:         key,
		Value:       value,
	}
}

func (txm *TxMap) Get(key int) int {
	result := make(chan int)
	txm.Tx <- Request{
		RequestType: Get,
		Key:         key,
		Result:      result,
	}
	return <-result
}

func (txm *TxMap) BeginTx() {
	tmp := txm.Tx
	txm.Parent <- tmp           // buffer chan に保存しておく
	txm.Tx = make(chan Request) // 小階層の Tx
	tmp <- Request{
		RequestType: BEGINTX,
		Tx:          txm.Tx,
	}
}

func (txm *TxMap) EndTx() {
	txm.req <- Request{
		requestType: EndTx,
	}
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
