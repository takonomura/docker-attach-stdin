package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
}
