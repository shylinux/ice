# ice Web Framework
ice is a web framework written in Go (Golang). 

## Quick start

```sh
# assume the following codes in main.go file
$ cat main.go
```

```go
package main

import "github.com/shylinux/ice"

func main() {
	ice.RunServe("9090") // listen and serve on 0.0.0.0:9090 (for windows "localhost:9090")
}
```

```sh
# run main.go and visit 0.0.0.0:9090/ (for windows "localhost:9090/") on browser
$ go run main.go
```

## API Examples

```go
func main() {
	ice.App("web.demo", "/tool", func(p *ice.Page) {
		p.Cmd("web.code.inner", ice.Arg("./", "main.go"))
		p.Cmd("cli.system", ice.Arg("pwd"))
		p.Cmd("zone", &ice.Zone{})
		p.Cmd("hash", &ice.Hash{})
	})
	ice.RunServe("9090")
}
```

```sh
# run main.go and visit 0.0.0.0:9090/demo/tool (for windows "localhost:9090/demo/tool") on browser
$ go run main.go
```

