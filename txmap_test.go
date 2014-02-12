package txmap

import (
	"runtime"
	"sync"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestTxMap(t *testing.T) {
	var wg sync.WaitGroup
	r := make(chan Request)
	go runMap(r)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			tx = beginTx(r)
			j := get(tx, 0)
			j = j + 1
			set(tx, 0, j)

			actual := get(0)
			expected := j
			txMap.end()
			endTx(tx)

			if actual != expected {
				t.Errorf("\ngot  %v\nwant %v", actual, expected)
			}
		}()
	}
	wg.Wait()
	t.Log(txMap.Get(0))
}
