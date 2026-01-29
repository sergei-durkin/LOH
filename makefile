define TIMED
	@start=$$(date +%s%N); \
	$(1); \
	end=$$(date +%s%N); \
	elapsed=$$(echo "scale=3; ($$end - $$start)/1000000" | bc); \
	echo "Elapsed: $$elapsed ms"
endef

ast:
	$(call TIMED, make aast)

aast:
	go run main.go main.loh main.s print=ast

tac:
	$(call TIMED, make _tac)
	
_tac:
	go run main.go main.loh main.s print=tac

asm:
	$(call TIMED, make _asm)

_asm:
	go run main.go main.loh main.s print=asm

cfg:
	$(call TIMED, make _cfg)

_cfg:
	go run main.go main.loh main.s print=cfg

ssa:
	$(call TIMED, make _ssa)

_ssa:
	go run main.go main.loh main.s print=ssa

llir:
	$(call TIMED, make _llir)

_llir:
	go run main.go main.loh main.s print=lir

aarch:
	$(call TIMED, make _aarch)

_aarch:
	go run main.go main.loh main.s print=aarch

compile:
	go run main.go ./examples/cmd main.s
	as -o main.o main.s
	ld -macos_version_min 15.1.0 -o main main.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e main -arch arm64
	cat main.s
	./main
	
