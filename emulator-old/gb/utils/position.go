package utils

func RowColtoPos(row, column, width int) int {
	return row*width + column
}
