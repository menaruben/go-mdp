package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/qeesung/image2ascii/convert"
)

type Slide struct {
	Title string
	Lines []string
	Image string
	Page  int
}

type Presentation struct {
	Author string
	Title  string
	Date   string
	Slides []Slide
	Pages  int
	Width  int
	Height int
}

func (p *Presentation) setInfo(info []string) {
	switch info[0] {
	case "?author":
		p.Author = info[1]

	case "?title":
		p.Title = info[1]

	case "?date":
		p.Date = info[1]
	}
}

func (p *Presentation) readMd(path string) {
	readFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var line string
	var numPage int = 0
	for fileScanner.Scan() {
		line = fileScanner.Text()

		if line != "" {
			switch ch := line[0]; ch {
			// information about markdown file
			case '?':
				info := strings.SplitN(line, ": ", 2)
				p.setInfo(info)

			// titles (for slide titles)
			case '#':
				titleLine := strings.Split(line, "# ")
				if len(titleLine) > 1 {
					numPage++
					slide := Slide{
						Title: titleLine[1],
						Page:  numPage,
					}
					p.Slides = append(p.Slides, slide)
				}

			// images
			case '!':
				imageUrl := strings.Trim(strings.Split(line, "![](")[1], ")")
				p.Slides[numPage-1].Image = imageUrl
			}
		}
		if numPage > 0 && !(strings.HasPrefix(line, "# ")) && !(strings.HasPrefix(line, "![](")) {
			p.Slides[numPage-1].Lines = append(p.Slides[numPage-1].Lines, line)
		}
	}
}

func clearScreen() {
	fmt.Printf("\x1bc")
}

func (s Slide) printAsciiImage() {
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 80
	convertOptions.FixedHeight = 20

	converter := convert.NewImageConverter()
	fmt.Print(converter.ImageFile2ASCIIString(s.Image, &convertOptions))
}

func displaySlide(s Slide) {
	slideContent := strings.Join(s.Lines, "\n")
	// slideDisplayText := fmt.Sprintf("%s\n\t\n%s", s.Title, slideContent)
	slideDisplayText := markdown.Render(slideContent, 80, 6)

	// fmt.Println(s.Title)
	color.New(color.FgCyan, color.Bold).Printf("%s\n", s.Title)
	fmt.Println(string(slideDisplayText))

	if s.Image != "" {
		s.printAsciiImage()
	}
}

func (p *Presentation) displayPresentation() {
	clearScreen()
	// fmt.Println("Press the left or right arrow key. Press 'q' to quit. ")

	presentationStartText := fmt.Sprintf("%s\n%s\n%s\n", p.Title, p.Author, p.Date)
	color.Cyan(presentationStartText)

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	var firstKeyPressed bool = false
	var numPage int = 0
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyArrowLeft {
			clearScreen()
			numPage--
			if numPage < 0 {
				// p.displayStart()
				color.Cyan(presentationStartText)
			} else {
				displaySlide(p.Slides[numPage])
			}

		} else if key == keyboard.KeyArrowRight {
			clearScreen()
			if !firstKeyPressed {
				displaySlide(p.Slides[0])
				firstKeyPressed = true
			} else {
				numPage++
				if numPage >= len(p.Slides) {
					break
				}

				displaySlide(p.Slides[numPage])
			}

		} else if char == 'q' || char == 'Q' {
			fmt.Println("Quitting...")
			os.Exit(0)
		}
	}
}

func main() {
	filePath := flag.String("path", "main.md", "specifies the path of the markdown file")
	height := flag.Int("height", 1080, "specifies the height of the window in pixel")
	width := flag.Int("width", 1920, "specifies the width of the window in pixel")
	flag.Parse()

	presentation := Presentation{
		Height: *height,
		Width:  *width,
	}
	presentation.readMd(*filePath)

	presentation.displayPresentation()
}
