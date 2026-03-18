package parser

import (
	"fmt"
	"slices"
)

func PrintSyntaxError(err error, buf []byte) {
	errs := []error{}
	s := []int{}

	for {
		if e, ok := err.(*ParserError); ok {
			e.Lexeme.Print()
			s = append(s, e.Lexeme.Pos())
			errs = append(errs, e)

			err = e.error
		} else {
			break
		}
	}

	rows := [][2]int{{0, 0}}
	for i, b := range buf {
		if b == '\n' || b == 0 {
			rows[len(rows)-1][1] = i
			rows = append(rows, [2]int{i, 0})
		}
	}

	cmp := func(row [2]int, t int) int {
		if row[0] > t {
			return 1
		}
		if row[1] < t {
			return -1
		}
		return 0
	}

	affected := [][2]int{}
	for j, pos := range s {
		i, ok := slices.BinarySearchFunc(rows, pos, cmp)
		if !ok {
			panic(1)
		}

		affected = append(affected, [2]int{j, i})
	}

	extractWithTrimLn := func(row [2]int) string {
		r := string(buf[rows[row[1]][0]:rows[row[1]][1]])
		if len(r) > 0 && (r[len(r)-1] == '\n' || r[len(r)-1] == 0) {
			r = r[:len(r)-1]
		}
		if len(r) > 0 && (r[0] == '\n' || r[0] == 0) {
			r = r[1:]
		}

		return r
	}

	printRow := func(row [2]int) {
		fmt.Println(extractWithTrimLn(row))
	}

	printTargetRow := func(row [2]int, msg string) {
		fmt.Printf("%s <-- %s\n", extractWithTrimLn(row), msg)
	}

	for _, row := range affected {
		err := errs[row[0]].Error()
		if row[1] > 0 {
			printRow([2]int{0, row[1] - 1})
		}
		printTargetRow(row, err)
		if row[1] < len(rows) {
			printRow([2]int{0, row[1] + 1})
		}
	}
}
