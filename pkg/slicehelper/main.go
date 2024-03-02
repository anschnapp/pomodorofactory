package slicehelper

func Copy2DSlice(src [][]rune, dest [][]rune) {
	for i := range src {
		copy(src[i], dest[i])
	}
}
