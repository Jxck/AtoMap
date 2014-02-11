package txmap

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type Type int

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
	log.Printf("new TxMap: %+v\n", txm)

	go func() {
		m := make(map[int]int) // 実際の map
		HandleRequests(m, txm.Tx)
	}()

	return txm
}

func HandleRequests(m map[int]int, tx chan Request) {
	log.Println("handle request")
	for {
		req := <-tx
		switch req.Type {
		case GET:
			log.Println("GET", m)
			req.Result <- m[req.Key]
		case SET:
			log.Println("SET", m)
			m[req.Key] = req.Value
		case BEGINTX:
			log.Println("BEGINTX")
			HandleRequests(m, req.Tx)
		case ENDTX:
			log.Println("ENDTX")
			return
		}
	}
}

func (txm *TxMap) String() string {
	return fmt.Sprintf("Tx: %v, Parent: %v", txm.Tx, len(txm.Parent))
}

func (txm *TxMap) Set(key int, value int) {
	request := Request{
		Type:  SET,
		Key:   key,
		Value: value,
	}
	log.Printf("SET request %+v", request)
	txm.Tx <- request
}

func (txm *TxMap) Get(key int) int {
	result := make(chan int)
	request := Request{
		Type:   GET,
		Key:    key,
		Result: result,
	}
	log.Printf("GET request %+v", request)
	txm.Tx <- request
	return <-result
}

func (txm *TxMap) BeginTx() {
	log.Printf("BeginTx %+v", txm)
	parent := txm.Tx
	log.Printf("parent %+v", parent)
	txm.Parent <- parent        // buffer chan に保存しておく
	txm.Tx = make(chan Request) // 小階層の Tx
	parent <- Request{
		Type: BEGINTX,
		Tx:   txm.Tx,
	}
}

func (txm *TxMap) EndTx() {
	log.Println("EndTx")
	txm.Tx <- Request{
		Type: ENDTX,
	}
	txm.Tx = <-txm.Parent
}

// func main() {
// 	r := make(chan Request)
// 	go runMap(r)
// 	tx := beginTx(r)
// 	set(tx, 1, "hoge")
// 	log.Println(get(tx, 1))
// 	endTx(tx)
// }
