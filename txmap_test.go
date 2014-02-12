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

			tx := Lock(txMap)
			j := Get(tx, 0)
			j = j + 1
			Set(tx, 0, j)

			actual := Get(tx, 0)
			expected := j
			Unlock(tx)

			if actual != expected {
				t.Errorf("\ngot  %v\nwant %v", actual, expected)
			}
		}()
	}
	wg.Wait()
	t.Log(Get(txMap, 0))
}
