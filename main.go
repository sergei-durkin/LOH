package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"loh/compiler"
	"loh/lir"
	"loh/machine/aarch"
	"loh/parser"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println(os.Args[1])

	buf := loadBuf()
	ast, err := parser.Parse(buf)
	if err != nil {
		parser.PrintSyntaxError(err, buf)
		os.Exit(1)
	}

	if len(os.Args) > 3 && os.Args[3] == "print=ast" {
		ast.Print()
		return
	}

	t := compiler.NewTac(ast.Unit)
	if len(os.Args) > 3 && os.Args[3] == "print=tac" {
		for _, t := range t {
			t.Print(os.Stdout)
		}
		return
	}

	c := []*compiler.Cfg{}
	for _, t := range t {
		c = append(c, compiler.NewCfg(t))
	}
	if len(os.Args) > 3 && os.Args[3] == "print=cfg" {
		for _, c := range c {
			c.Print(os.Stdout)
		}
		return
	}

	ssa := compiler.NewSSA(c)
	if len(os.Args) > 3 && os.Args[3] == "print=ssa" {
		ssa.Print(os.Stdout)
		return
	}

	lir := lir.NewLir(ssa)
	if len(os.Args) > 3 && os.Args[3] == "print=lir" {
		lir.Print(os.Stdout)
		return
	}

	src := bytes.NewBuffer(nil)
	a := aarch.AARCH{}
	a.Emit(src, lir)
	if len(os.Args) > 3 && os.Args[3] == "print=aarch" {
		fmt.Println(src)
		return
	}

	if err = os.Remove(os.Args[2]); err != nil {
		panic(err)
	}

	o, err := os.OpenFile(os.Args[2], os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	o.Write(src.Bytes())

	o.Close()
}

func loadBuf() []byte {
	fileinfo, err := os.Lstat(os.Args[1])
	if err != nil {
		panic(err)
	}

	if !fileinfo.IsDir() {
		f, err := os.Open(os.Args[1])
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
	err = filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".loh" {
			return nil
		}

		if !d.IsDir() {
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
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	return buf
}
