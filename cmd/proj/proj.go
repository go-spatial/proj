package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-spatial/proj/support"
)

func main() {

	verbose := flag.Bool("verbose", false, "do lots of logging")
	inverse := flag.Bool("inverse", false, "go backwards")

	flag.Parse()

	fmt.Printf("verbose: %t\n", *verbose)
	fmt.Printf("inverse: %t\n", *inverse)

	s := strings.Join(flag.Args(), " ")

	s = "+proj=utm +zone=32 +ellps=GRS80" // TODO

	fmt.Printf("string: %s\n", s)

	_, err := support.NewProjString(s)
	if err != nil {
		panic(err)
	}

	var a1, a2 float64
	fmt.Fscanf(os.Stdin, "%f %f\n", &a1, &a2)
	fmt.Printf("-> %f %f\n", a1, a2)
}
