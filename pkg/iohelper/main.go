package iohelper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFileInArray(filename string) ([][]rune, error) {
	fmt.Println("%s", filename)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := [][]rune{}
	var i int = 0
	for scanner.Scan() {
		fmt.Println("print line")
		lines[i] = []rune(scanner.Text())
		i++
	}
	return lines, nil
}

func SplitMultilineStringToSlice(data string) [][]rune {
	slice := [][]rune{}
	lines := strings.Split(data, "\n")
	for i, line := range lines {
		slice[i] = []rune(line)
	}
	return slice
}
