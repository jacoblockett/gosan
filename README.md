# Gosan

Utility functions to sanitize strings in Go.

> ❗ There's only one method in here right now. Hardly a package worth its name, but hey, we all start somewhere.

### Installation

```bash
go get github.com/jacoblockett/gosan/v3
```

You can read the godoc [here](https://pkg.go.dev/github.com/jacoblockett/gosan/v3) for detailed documentation.

### Quickstart

```go
package main

import "github.com/jacoblockett/gosan"

func main() {
	// Assuming a Windows environment
	filename := "<>:\"/\\|?*abc.txt" // "<>:"/\|?*abc.txt" without escape chars
	opts := &gosan.FilenameOptions{Environment: gosan.Windows}
	sanitized := gosan.Filename(filename, opts)

	fmt.Println(sanitized) // "abc.txt"

	// Assuming a Linux/Unix environment
	filename := "/.."
	opts := &gosan.FilenameOptions{Environment: gosan.Linux, Replacement: "x"}
	sanitized := gosan.Filename(filename, opts)

	fmt.Println(sanitized) // "x.."
}
```
