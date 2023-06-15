package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/qeesung/image2ascii/convert"
)

func main() {
	// Create convert options
	convertOptions := convert.DefaultOptions
	convertOptions.Ratio = 0.12
	convertOptions.FixedWidth = 100
	convertOptions.FixedHeight = 40

	// Create the image converter
	filenameMd := "![](../imgs/gologo.png)"
	imagePath := strings.Trim(strings.Split(filenameMd, "![](")[1], ")")

	converter := convert.NewImageConverter()
	fmt.Print(converter.ImageFile2ASCIIString(imagePath, &convertOptions))
}
