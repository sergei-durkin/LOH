package risc

import "fmt"

var regsMap = map[string]string{
	"RA":   "x1",
	"SP":   "x2",
	"TMP0": "x5",
	"TMP1": "x6",
	"TMP2": "x7",
	"FP":   "x8",
	"A0":   "x10",
	"A1":   "x11",
	"A2":   "x12",
	"A3":   "x13",
	"A4":   "x14",
	"A5":   "x15",
	"A6":   "x16",
	"A7":   "x17",
	"S0":   "x9",
	"S1":   "x18",
	"S2":   "x19",
	"S3":   "x20",
	"S4":   "x21",
	"S5":   "x22",
	"S6":   "x23",
	"S7":   "x24",
	"S8":   "x25",
	"S9":   "x26",
	"S10":  "x27",
	"TMP3": "x28",
	"TMP4": "x29",
	"TMP5": "x30",
	"TMP6": "x31",
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
