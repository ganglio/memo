# Memo

[![GoDoc](https://godoc.org/github.com/ganglio/memo?status.svg)](https://godoc.org/github.com/ganglio/memo)
[![Go Report Card](https://goreportcard.com/badge/github.com/ganglio/memo)](https://goreportcard.com/report/github.com/ganglio/memo)

Teeny-weeny cached variable library with auto refresh and anti stampede.

## Usage

```go
v := 0
counter := M[int](func() int{
  v=v+1
  return v
}).Memo(time.Second)

for {
  fmt.Printf("Counter %d", counter())
}
```

If you function can return an error, use `GX` like following:

```go
v := 0
counter := MX[int](func() (int, error){
  if (v<0) {
    return v, errors.New("Cannot start with a negative counter, for some reason...")
  }
  v := v+1
  return v, nil
}).Memo(time.Second)

for {
  fmt.Printf("Counter %d", counter())
}
```

In this case it works like this:

  * If the generator returns an error, the Memo method will return the error upon first cache warmup.
  * If not it will return the wrapped function.
  * After the first initialisation, if the generator returns an error, it will be returned without updating the cached value.
    * This means that the errors must be handled inside the generator.