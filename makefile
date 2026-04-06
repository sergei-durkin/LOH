arch := "aarch"

build:
	time go build main.go

ast: build
	time ./main -print=ast -arch=${arch}

tac: build
	time ./main -print=tac -arch=${arch}

asm: build
	time ./main -print=asm -arch=${arch}

cfg: build
	time ./main -print=cfg -arch=${arch}

ssa: build
	time ./main -print=ssa -arch=${arch}

lir: build
	time ./main -print=lir -arch=${arch}

compile: build
	./main -input=./examples/cmd -output=main.s -arch=${arch}
	as -o main.o main.s
	ld -macos_version_min 15.1.0 -o main main.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e main -arch arm64
	cat main.s
	time ./main

