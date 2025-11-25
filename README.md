# Gosan

Utility functions to sanitize strings in Go.

> ❗ There's only one method in here right now. Hardly a package worth its name, but hey, we all start somewhere.

### Installation

```bash
go get github.com/jacoblockett/gosan/v2
```

You can read the godoc [here](https://pkg.go.dev/github.com/jacoblockett/gosan/v2) for detailed documentation.

### Quickstart

```go
package main

import "github.com/jacoblockett/gosan"

func main() {
	// Assuming a Windows environment
	filename := "<>:\"/\\|?*abc.txt" // "<>:"/\|?*abc.txt" without escape chars
	sanitized := gosan.Filename(filename)

	fmt.Println(sanitized) // "abc.txt"

	// Assuming a Linux/Unix environment
	filename := "/.."
	sanitized := gosan.Filename(filename)

	fmt.Println(sanitized) // ""
}
```
