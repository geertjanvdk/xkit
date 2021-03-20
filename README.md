# xkit

The xkit Go package is a collection of handy functions and types practically
put in themed packages mimicking the Go standard library.

xkit is for personal projects the creator is working on. They can be
handy for others, hence open sourced. Bugs can be reported, and even feature requested,
even appreciated, but don't expect a swift response.

Copyright (c) 2020 Geert JM Vanderkelen

## Installation

```go get github.com/geertjanvdk/xkit```

Dependencies are kept to a minimum (see `go.mod`). In fact, most of what is in xkit
is to prevent dependencies.

## Overview

* **xutil** - sometimes silly functions like finding an int in a slice, but also `OrderedMap`
  which is map keeping order
* **xsql** - wrapper around the [MySQL Go driver][2] (the only dependency in xkit); providing
  handy things like HaveTrigger, HaveTable, etc..
* **xpath** - handy tools like `IsRegularFile` or `IsDir`, or `FilesInDir`.  
* **xnet** - tools around networking like `GetTCPPort` and `IsEmailAddress` (which is very
  opinionated what an email address is)
* **xid** - the author does not like UUID because they are usually misused; this packages
  contains next to UUID implementation also [nanoid][3]
* **xgraphql** - writing a GraphQL client is easy; writing it over and over again is tedious
* **xansi** - basic ANSI things, mostly styling and coloring text
* **xvenv** - virtual environment tools; for example to figure out whether the application
  runs in a Docker container or not

### xhttp - HTTP tools

There are quite a few packages and frameworks like Gorilla, Gin, Martini, etc.. In the
end, if you sit down, you can write down what you need in a few hours and never think
of the dependency again.

The `xhttp` wrappers around the Go's `http` package offering a `xhttp.Client` which makes
it a bit easier to communicate with web APIs for example.

It also offers `xhttp.ServeReMux` which is a multiplexer (mux) that uses regular expressions
and can restrict HTTP methods. It also can capture values using Django's style. All that
in 300 lines of simple code (we don't care about performance too much).

Example:

```go
mux := xhttp.NewServeReMux()
mux.Handle("^/images", web.Static(), xhttp.MethodGet)
mux.Handle("^/person/<int:id>", PersonHandler(), xhttp.MethodGet)
mux.Handle("^/$", web.Home{}, xhttp.MethodGet)
```
In the `PersonHandler()` from above example, one would then access the captured value
like this:

```go
captures, ok := ctx.Value(xhttp.CapturesContextKey).(*xhttp.Captures)
fmt.Println("Person ID:", captures["id"].AsInt64())
```

### xlog - logging

Heavily inspired by [logrus][1], this package provides the minimum depending only on the
functionality coming with Go itself. Some formatters might be added, but the text
and JSON formatters are in most cases enough.

### xt - unit testing

Quite a few testing frameworks exists with fancy features, and pretty much all based on
the stock Go `testing` package.

With the `xkit.xt` package, however, we wanted to have a simple wrapper providing just
the basics (inspired by an online article which got lost in the bookmarks):

* `OK` & `KO` functions to test whether error is returned
* `Eq` check equality (no, there is no Not-Eq)
* `Assert` checking if condition is true (can be used for anything, and is actually
   a wrapper around `Eq`)
* `Panics` checks if a function that must panic actually panics
* `Match` is probably something too much: uses regular expression to match a string

There is little chance anything else will be added.

Example:

```go
// excerpt from xkit.xlog.entry_test.go
func TestEntry_UnmarshalJSON(t *testing.T) {
    t.Run("field as time.Time", func (t *testing.T) {
    r, ok := (e).Fields["someTime"].(time.Time)
    xt.Assert(t, ok, "field is not time.Time")
    xt.Eq(t, now, r)
    })
}
```

## Contributing

`xkit` is used for personal project, and meant to inspire others to start their own. Of
course, it can be used as dependency, features can be requests. Contributions are
not denied but don't expect it to get merged quickly.

Also, the code is not that good yet, and still lacks documentation.

## License

MIT licensed: see LICENSE.txt.

[1]: https://github.com/sirupsen/logrus
[2]: https://github.com/go-sql-driver/mysql
[3]: https://zelark.github.io/nano-id-cc/
