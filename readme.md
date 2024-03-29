# Memo

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/ganglio/memo/tree/master.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/ganglio/memo/tree/master)
[![codecov](https://codecov.io/gh/ganglio/memo/branch/master/graph/badge.svg)](https://codecov.io/gh/ganglio/memo)
[![GoDoc](https://godoc.org/github.com/ganglio/memo?status.svg)](https://godoc.org/github.com/ganglio/memo)
[![Go Report Card](https://goreportcard.com/badge/github.com/ganglio/memo)](https://goreportcard.com/report/github.com/ganglio/memo)

Teeny-weeny cached variable library with auto refresh and anti stampede.

## Usage

```go
v := 0
counter := memo.Memo(func() interface{} {
  v = v + 1
  return v
}, time.Second)

for {
  fmt.Printf("Counter %s", counter())
}
```
