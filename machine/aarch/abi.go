package aarch

import (
	"fmt"
)

var regsMap = map[string]string{
	"A0":   "x0",
	"A1":   "x1",
	"A2":   "x2",
	"A3":   "x3",
	"A4":   "x4",
	"A5":   "x5",
	"A6":   "x6",
	"A7":   "x7",
	"TMP0": "x9",
	"TMP1": "x10",
	"TMP2": "x11",
	"TMP3": "x12",
	"TMP4": "x13",
	"TMP5": "x14",
	"TMP6": "x15",
	"IP0":  "x16",
	"IP1":  "x17",
	"S0":   "x19",
	"S1":   "x20",
	"S2":   "x21",
	"S3":   "x22",
	"S4":   "x23",
	"S5":   "x24",
	"S6":   "x25",
	"S7":   "x26",
	"S8":   "x27",
	"S9":   "x28",
	"FP":   "x29",
	"LR":   "x30",
	"SP":   "SP",
}

func CalleeSavedRegister(x int) string {
	res, ok := regsMap[fmt.Sprintf("S%d", x)]
	if !ok {
		panic(fmt.Sprintf("undefined callee-saved register %d", x))
	}

	return res
}
func ArgRegister(x int) string {
	res, ok := regsMap[fmt.Sprintf("A%d", x)]
	if !ok {
		panic(fmt.Sprintf("undefined arg register %d", x))
	}

	return res
}
