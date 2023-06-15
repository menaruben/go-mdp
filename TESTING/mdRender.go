package main

import (
	"fmt"
	"io/ioutil"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func main() {
	path := "sample.md"
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	result := markdown.Render(string(source), 80, 6)

	fmt.Println(string(result))
}
