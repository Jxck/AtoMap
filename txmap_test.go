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
	atoMap := NewAtoMap()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			m := atoMap.Lock()
			j := m.Get(0)
			j = j + 1
			m.Set(0, j)

			actual := m.Get(0)
			expected := j
			m.Unlock()

			if actual != expected {
				t.Errorf("\ngot  %v\nwant %v", actual, expected)
			}
		}()
	}
	wg.Wait()
	t.Log(atoMap.Get(0))
}
