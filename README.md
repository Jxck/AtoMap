AtoMap
======

atomic map


## generator

main/generate.go is a file generator script.
you can genrate original atomap with key type and value type which you want.

```
$ go run main/generate.go -p packagename -k string -v MyType
```

this generates Atomic Map with  map[string]MyType.

## Usage

default atomap has map[int]int.

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


## Reference

this library in reference to [http://www.amazon.com/dp/0321817141]


## License

The MIT License (MIT)
Copyright (c) 2014 Jxck
