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
	txMap := NewTxMap()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			txMap.BeginTx()
			j := txMap.Get(0)
			j = j + 1
			txMap.Set(0, j)
			actual := txMap.Get(0)
			expected := j
			txMap.EndTx()
			if actual != expected {
				t.Errorf("\ngot  %v\nwant %v", actual, expected)
			}
		}()
	}
	wg.Wait()
	t.Log(txMap.Get(0))
}
