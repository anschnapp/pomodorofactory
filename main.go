package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type marginBorder struct {
	top    int
	left   int
	right  int
	bottom int
}
type ui struct {
	width int
}

func main() {
	gomodoroAscii, err := readFileInArray("gomodoro-asci")
	// todo put all margins together
	margins := marginBorder{5, 5, 5, 5}
	ui := ui{80}
	// todo make pomodoro ascii object with validation if not empty and convenient witdh attribute etc...
	pomodoroWidth := len(gomodoroAscii[0])
	pomodoroHeight := len(gomodoroAscii)

	view := generateBlankView(margins, ui, pomodoroWidth, pomodoroHeight)

	if err != nil {
		panic("could not read pomodoro file")
	}

	// todo should be render funuction, view should be changed by tick and then render should be called after all have reacted on tick
	for _, value := range view {
		fmt.Println(value)
	}
}
func generateBlankView(margin marginBorder, ui ui, pomodoroWidth int, pomodorHeight int) []string {
	blankView := make([]string, margin.top+pomodorHeight+margin.bottom)
	width := margin.left + margin.right + pomodoroWidth + ui.width

	for i := range blankView {
		if i == 0 || i == len(blankView)-1 {
			blankView[i] = createStringFilledWith(width, 'x')
		} else {
			blankView[i] = "x" + createStringFilledWith(width-2, ' ') + "x"
		}
	}
	return blankView
}
func createStringFilledWith(size int, character rune) string {
	filledString := ""
	for i := 0; i < size; i++ {
		filledString = filledString + string(character)
	}
	return filledString
}

func syncExample() {
	var wg sync.WaitGroup
	wg.Add(1)

	ticker := time.NewTicker(time.Duration(1000 * time.Millisecond))
	go printTimes(ticker, &wg)

	wg.Wait()
}

func readFileInArray(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func printTimes(ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		<-ticker.C
		fmt.Printf("\r%s", time.Now())
	}
}
