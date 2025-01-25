# errifscope

errifscope is linter to find if block that can encapsulate the scope of error variable.

By confining the scope of error to if, it is possible to reduce reassignment to err and avoid creating variables such as err1 and err2.

## Installation

```bash
go install github.com/otakakot/errifscope/cmd/errifscope@latest
```

## Usage

```bash
go vet -vettool=$(which errifscope) ./...
```

## Examples

```go
package main

func f() error {
	return nil
}

func main() {
	err := f()
	if err != nil {
		return
	}
}
```

Running errifscope on the above code will produce the following output:

```bash
./main.go:9:2: can be scoped with if block
```

You can be rewritten as:

```go
package main

func f() error {
	return nil
}

func main() {
	if err := f(); err != nil {
		return
	}
}
```
