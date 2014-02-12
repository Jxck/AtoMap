AtoMap
======

atomic map


## Usage

```go
m := NewAtoMap()
m.Set(0, 1)
m.Get(0) // 1

m1 := m.Lock()
m1.Set(0, 1)

m2 := m.Lock()
m2.Set(0, 2)

m2.Unlock()
m1.Unlock()
m.Get(0) // 2
```

## License

The MIT License (MIT)
Copyright (c) 2014 Jxck
