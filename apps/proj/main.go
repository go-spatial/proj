package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {

	verbose := flag.Bool("verbose", false, "do lots of logging")

	flag.Parse()

	fmt.Printf("verbose: %t\n", *verbose)

	s := strings.Join(flag.Args(), " ")
	fmt.Printf("string: %s\n", s)
}
