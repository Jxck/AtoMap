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
	for i := 0; i < 1000; i++ {
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

func TestConcurrent(t *testing.T) {
	m1 := NewAtoMap()
	m2 := NewAtoMap()

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mm1 := m1.Lock()
			mm2 := m2.Lock()

			j := mm1.Get(0)
			j = j + 1
			mm1.Set(0, j)

			k := mm2.Get(0)
			k = k + 1
			mm2.Set(0, k)

			actual := mm1.Get(0)
			expected := j

			actual = mm2.Get(0)
			expected = k

			mm1.Unlock()
			mm2.Unlock()

			if actual != expected {
				t.Errorf("\ngot  %v\nwant %v", actual, expected)
			}

			if actual != expected {
				t.Errorf("\ngot  %v\nwant %v", actual, expected)
			}
		}()
	}
}

func TestNest(t *testing.T) {
	m := NewAtoMap()
	m.Set(0, 0)

	m1 := m.Lock()
	m1.Set(0, 1)

	m2 := m1.Lock()
	m2.Set(0, 2)

	m3 := m2.Lock()

	m3.Set(0, 3)
	if m3.Get(0) != 3 {
		t.Errorf("\ngot  %v\nwant %v", m3.Get(0), 3)
	}
	m3.Unlock()

	if m2.Get(0) != 3 {
		t.Errorf("\ngot  %v\nwant %v", m2.Get(0), 3)
	}
	m2.Unlock()

	if m1.Get(0) != 3 {
		t.Errorf("\ngot  %v\nwant %v", m1.Get(0), 3)
	}
	m1.Unlock()

	if m.Get(0) != 3 {
		t.Errorf("\ngot  %v\nwant %v", m.Get(0), 3)
	}
}
