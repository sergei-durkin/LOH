package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"loh/compiler"
	"loh/lir"
	"loh/machine"
	"loh/machine/aarch"
	"loh/machine/risc"
	"loh/parser"
	"os"
	"path/filepath"
)

type config struct {
	inputPath  string
	outputPath string
	printStage string
	arch       string
}

var cfg config

func init() {
	flag.StringVar(&cfg.inputPath, "input", "main.loh", "Path to input file or directory")
	flag.StringVar(&cfg.outputPath, "output", "main.s", "Path to output file")
	flag.StringVar(&cfg.printStage, "print", "", "Print intermediate stage and exit (ast|tac|cfg|ssa|lir|asm)")
	flag.StringVar(&cfg.arch, "arch", "aarch", "Target architecture (riscv|aarch)")

	flag.Parse()
}

func main() {
	buf := loadBuf()
	ast, err := parser.Parse(buf)
	if err != nil {
		parser.PrintSyntaxError(err, buf)
		os.Exit(1)
	}

	if cfg.printStage == "ast" {
		ast.Print()
		return
	}

	t := compiler.NewTac(ast.Unit)
	if cfg.printStage == "tac" {
		for _, t := range t {
			t.Print(os.Stdout)
		}
		return
	}

	c := []*compiler.Cfg{}
	for _, t := range t {
		c = append(c, compiler.NewCfg(t))
	}
	if cfg.printStage == "cfg" {
		for _, c := range c {
			c.Print(os.Stdout)
		}
		return
	}

	ssa := compiler.NewSSA(c)
	if cfg.printStage == "ssa" {
		ssa.Print(os.Stdout)
		return
	}

	lir := lir.NewLir(ssa)
	if cfg.printStage == "lir" {
		lir.Print(os.Stdout)
		return
	}

	src := bytes.NewBuffer(nil)
	var m machine.Machine
	switch cfg.arch {
	case "aarch":
		m = &aarch.AARCH{}
	case "riscv":
		m = &risc.RISCV{}
	}
	m.Emit(src, lir)

	if cfg.printStage == "asm" {
		fmt.Println(src)
		return
	}

	if err = os.Remove(cfg.outputPath); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	o, err := os.OpenFile(cfg.outputPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer o.Close()

	o.Write(src.Bytes())
}

func loadBuf() []byte {
	fileinfo, err := os.Lstat(cfg.inputPath)
	if err != nil {
		panic(err)
	}

	if !fileinfo.IsDir() {
		f, err := os.Open(cfg.inputPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buf, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		return buf
	}

	var buf []byte
	err = filepath.WalkDir(cfg.inputPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".loh" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}

		buf = append(buf, '\n')
		buf = append(buf, b...)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return buf
}
