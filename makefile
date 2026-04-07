arch := "riscv"

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

machine: build
	./main -print=asm -arch=${arch} > program.s
	riscv64-elf-as -march=rv32i -mabi=ilp32 -mno-relax program.s -o program.o
	riscv64-elf-ld -m elf32lriscv -Ttext=0x00000000 -e 0 program.o -o program.elf
	riscv64-elf-objcopy -O binary program.elf program.bin
	riscv64-elf-objdump -d program.elf
	xxd -e -g 4 -c 4 program.bin | awk '{print $$2}' > program
	cat program

compile: build
	./main -input=./examples/cmd -output=main.s -arch=${arch}
	as -o main.o main.s
	ld -macos_version_min 15.1.0 -o main main.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e main -arch arm64
	cat main.s
	time ./main

