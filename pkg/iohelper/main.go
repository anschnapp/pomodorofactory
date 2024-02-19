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
	for scanner.Scan() {
		fmt.Println("print line")
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func SplitMultilineStringToArray(data string) [][]rune {
	lines := strings.Split(data, "\n")
	return lines
}
