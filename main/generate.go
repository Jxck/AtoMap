package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

var Package, Key, Value, File string

func init() {
	flag.StringVar(&Package, "p", "atomap", "package name")
	flag.StringVar(&Key, "k", "int", "key type")
	flag.StringVar(&Value, "v", "int", "value type")
	flag.StringVar(&File, "f", "atomap.go", "file name")
	flag.Parse()
}

func main() {
	// define params
	var param = struct {
		Package, Key, Value string
	}{
		Package: Package,
		Key:     Key,
		Value:   Value,
	}

	// open output file
	fd, err := os.Create(File)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	// execute template
	t := template.Must(template.New("atomap").Parse(str))
	err = t.Execute(fd, param)
	if err != nil {
		log.Fatal(err)
	}
}

var str = `package atomap

// generated from https://github.com/Jxck/AtoMap

type RequestType int

const (
	SET RequestType = iota
	GET
	LOCK
	UNLOCK
)

type request struct {
	requestType RequestType
	key         {{.Key}}
	value       {{.Value}}
	result      chan {{.Key}}
	tx          chan request
}

type AtoMap struct {
	Tx chan request
}

func NewAtoMap() *AtoMap {
	tx := make(chan request)
	go func() {
		m := make(map[{{.Key}}]{{.Value}})
		handleRequests(m, tx)
	}()
	atomap := &AtoMap{
		Tx: tx,
	}
	return atomap
}

func handleRequests(m map[{{.Key}}]{{.Value}}, r chan request) {
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

func (atomap *AtoMap) Get(key {{.Key}}) {{.Value}} {
	result := make(chan {{.Value}})
	request := request{
		requestType: GET,
		key:         key,
		result:      result,
	}
	atomap.Tx <- request
	return <-result
}

func (atomap *AtoMap) Set(key {{.Key}}, value {{.Value}}) {
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
`
