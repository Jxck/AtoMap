package atomap

import (
	"runtime"
	"sync"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestAtoMap(t *testing.T) {
	var wg sync.WaitGroup
	txMap := NewAtoMap()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			tx := txMap.Lock()
			j := tx.Get(0)
			j = j + 1
			tx.Set(0, j)

			actual := tx.Get(0)
			expected := j
			tx.Unlock()

			if actual != expected {
				t.Errorf("\ngot  %v\nwant %v", actual, expected)
			}
		}()
	}
	wg.Wait()
	t.Log(txMap.Get(0))
}
