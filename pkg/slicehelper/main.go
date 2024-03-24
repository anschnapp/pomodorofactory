package slicehelper

import (
	"fmt"
	"math"
)

func Copy2DSlice[T any](src [][]T, dest [][]T) {

	if len(dest) < len(src) {
		panic("dest must have at least the same length as the src")
	}
	if MaxWidth(src) > MinWidth(dest) {
		panic("dest must have at least equal min space and then source has max space src max widht is " + fmt.Sprint(MaxWidth(src)) + " dest min width is " + fmt.Sprint(MinWidth(dest)) + "dest max width is" + fmt.Sprint(MaxWidth(dest)))
	}
	for i := range src {
		for j := range src[i] {
			dest[i][j] = src[i][j]
		}
	}
}

func MaxWidth[T any](slice [][]T) int {
	maxWidth := 0
	for i := range slice {
		width := len(slice[i])
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func MinWidth([T any]slice [][]T) int {
	minWidth := math.MaxInt
	for i := range slice {
		width := len(slice[i])
		if width < minWidth {
			minWidth = width
		}
	}
	return minWidth
}
