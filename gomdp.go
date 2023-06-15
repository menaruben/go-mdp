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
)

type Slide struct {
	Title string
	Lines []string
	Page  int
}

type Presentation struct {
	Author string
	Title  string
	Date   string
	Slides []Slide
	Pages  int
}

func (p *Presentation) setInfo(info []string) {
	switch info[0] {
	case "!author":
		p.Author = info[1]

	case "!title":
		p.Title = info[1]

	case "!date":
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
			case '!':
				info := strings.SplitN(line, ": ", 2)
				p.setInfo(info)

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
			}
		}
		if numPage > 0 && !(strings.HasPrefix(line, "# ")) {
			p.Slides[numPage-1].Lines = append(p.Slides[numPage-1].Lines, line)
		}
	}
}

func clearScreen() {
	fmt.Printf("\x1bc")
}

func displaySlide(s Slide) {
	slideContent := strings.Join(s.Lines, "\n")
	// slideDisplayText := fmt.Sprintf("%s\n\t\n%s", s.Title, slideContent)
	slideDisplayText := markdown.Render(slideContent, 80, 6)

	fmt.Println(s.Title)
	fmt.Println(string(slideDisplayText))
}

func (p *Presentation) displayPresentation() {
	clearScreen()
	// fmt.Println("Press the left or right arrow key. Press 'q' to quit. ")

	presentationStart := fmt.Sprintf("%s\n%s\n%s", p.Title, p.Author, p.Date)
	// fmt.Println(presentationStart)
	color.Cyan(presentationStart)

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
				color.Cyan(presentationStart)
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
	filePath := flag.String("path", "", "specifies the path of the markdown file")
	flag.Parse()

	presentation := Presentation{}
	presentation.readMd(*filePath)

	presentation.displayPresentation()
}
