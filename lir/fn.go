package lir

import (
	"fmt"
	"loh/compiler"
	"loh/machine"
	"loh/token"
)

type Fn struct {
	name string

	blocks []*BasicBlock
	po     []*BasicBlock

	frameIndex map[string]int
	frameSize  int

	regs []machine.PhysicalRegister
}

func NewFn(
	name string,

	blocks []*compiler.BasicBlock,
	postorder []*compiler.BasicBlock,

	frameIndex map[string]int,
	frameSize int,

	regs map[string]compiler.VirtualRegister,
) *Fn {
	f := buildFn(
		name,
		blocks,
		postorder,
		frameIndex,
		frameSize,
		regs,
	)

	f.regalloc()

	return f
}

func (f *Fn) Name() string {
	return f.name
}

func (f *Fn) Blocks() []machine.Block {
	blocks := []machine.Block{}
	for i := 0; i < len(f.blocks); i++ {
		blocks = append(blocks, f.blocks[i])
	}
	return blocks
}

func (f *Fn) FrameSize() int {
	size, align := f.frameSize, 16

	size = max(size, align)
	if size%align != 0 {
		size += align - size%align
	}

	return size
}

func (f *Fn) CalleeSavedRegisters() []machine.PhysicalRegister {
	return f.regs
}

func buildFn(
	name string,

	blocks []*compiler.BasicBlock,
	postorder []*compiler.BasicBlock,

	frameIndex map[string]int,
	frameSize int,

	regs map[string]compiler.VirtualRegister,
) *Fn {
	convertArg := func(v compiler.Value) machine.Value {
		switch v := v.(type) {
		default:
			panic(fmt.Sprintf("undefined value: %T %+v", v, v))
		case nil:
			return nil
		case compiler.Variable:
			return &machine.Reg{ID: int(regs[v.Label()])}
		case *compiler.FP:
			return &machine.FP{}
		case *compiler.BoolConst:
			return &machine.BoolConst{Bool: v.Bool()}
		case *compiler.IntConst:
			return &machine.IntConst{Int: v.Int()}
		case *compiler.Reg:
			return &machine.Reg{ID: v.ID}
		case *compiler.ArgReg:
			return &machine.ArgReg{ID: v.ID}
		case *compiler.AddressOf:
			switch addr := v.Target.(type) {
			default:
				panic(v.Target)
			case *compiler.Reg:
				return &machine.Reg{ID: int(addr.ID)}
			}
		case *compiler.Dereference:
			switch addr := v.Addr.(type) {
			default:
				panic(fmt.Sprintf("undefined value: %T %+v", v.Addr, v.Addr))
			case *compiler.Reg:
				return &machine.Reg{ID: int(addr.ID)}
			case *compiler.IntConst:
				return &machine.IntConst{Int: addr.Int()}
			}
		}
	}

	lirBlocks := make([]*BasicBlock, len(blocks))
	for i := 0; i < len(blocks); i++ {
		bb := blocks[i]

		newBB := &BasicBlock{
			ID: bb.ID,
		}
		lirBlocks[bb.ID] = newBB

		instrs := make([]Instruction, 0, len(bb.Instr))
		for j := 0; j < len(bb.Instr); j++ {
			switch instr := bb.Instr[j].(type) {
			default:
				panic(fmt.Sprintf("undefined instruction: %T %+v", instr, instr))
			case *compiler.Assign:
				dst := convertArg(instr.Target().(compiler.Value))
				arg1 := convertArg(instr.Arg1())
				arg2 := convertArg(instr.Arg2())

				switch instr.Operation() {
				default:
					panic(fmt.Sprintf("undefined operation: %d %s", instr.Operation(), instr.Operation().String()))
				case token.AND:
					instrs = append(instrs, &AND{destination: dst, left: arg1, right: arg2})
				case token.OR:
					instrs = append(instrs, &OR{destination: dst, left: arg1, right: arg2})
				case token.XOR:
					instrs = append(instrs, &XOR{destination: dst, left: arg1, right: arg2})
				case token.PERCENT:
					instrs = append(instrs, &MOD{destination: dst, left: arg1, right: arg2})
				case token.SLASH:
					instrs = append(instrs, &DIV{destination: dst, left: arg1, right: arg2})
				case token.STAR:
					instrs = append(instrs, &MUL{destination: dst, left: arg1, right: arg2})
				case token.PLUS:
					instrs = append(instrs, &SUM{destination: dst, left: arg1, right: arg2})
				case token.MINUS:
					instrs = append(instrs, &SUB{destination: dst, left: arg1, right: arg2})
				case token.GT:
					instrs = append(instrs, &CMP{size: instr.Size(), token: token.GT, destination: dst, left: arg1, right: arg2})
				case token.GTE:
					instrs = append(instrs, &CMP{size: instr.Size(), token: token.GTE, destination: dst, left: arg1, right: arg2})
				case token.LT:
					instrs = append(instrs, &CMP{size: instr.Size(), token: token.LT, destination: dst, left: arg1, right: arg2})
				case token.LTE:
					instrs = append(instrs, &CMP{size: instr.Size(), token: token.LTE, destination: dst, left: arg1, right: arg2})
				case token.EQ:
					instrs = append(instrs, &CMP{size: instr.Size(), token: token.EQ, destination: dst, left: arg1, right: arg2})
				case token.NE:
					instrs = append(instrs, &CMP{size: instr.Size(), token: token.NE, destination: dst, left: arg1, right: arg2})

				case token.ASSIGN:
					switch arg := instr.Arg1().(type) {
					default:
						instrs = append(instrs, &MOV{destination: dst, source: convertArg(arg)})
					case *compiler.Dereference:
						instrs = append(instrs, &LDR{size: instr.Size(), destination: dst, source: convertArg(arg)})
					}
				}
			case *compiler.Call:
				args := []machine.Value{}
				for _, arg := range instr.Args() {
					args = append(args, convertArg(arg))
				}

				instrs = append(instrs, &CALL{destination: convertArg(instr.Target().(compiler.Value)), callee: instr.Callee(), args: args})
			case *compiler.Alloca:
				instrs = append(instrs, &ALLOCA{destination: convertArg(instr.Pointer()), size: convertArg(instr.Size())})
			case *compiler.Goto:
				instrs = append(instrs, &JMP{label: instr.Label(), fn: name})
			case *compiler.IfGoto:
				instrs = append(instrs, &CBZ{left: convertArg(instr.Condition()), label: instr.Label(), fn: name})
				instrs = append(instrs, &JMP{label: instr.Fall(), fn: name})
			case *compiler.Return:
				instrs = append(instrs, &RET{Value: convertArg(instr.Value())})
			case *compiler.Store:
				instrs = append(instrs, &STR{size: instr.Size, destination: convertArg(instr.Destination), source: convertArg(instr.Value)})
			case *compiler.Load:
				panic("not implemented")
				//instrs = append(instrs, &LDR{destination: convertArg(instr.Destination), source: convertArg(instr.Source)})
			}
		}
		newBB.Instr = instrs
	}

	lirPostorder := make([]*BasicBlock, len(postorder))
	for i := 0; i < len(postorder); i++ {
		bb := postorder[i]
		lbb := lirBlocks[bb.ID]
		lirPostorder[i] = lbb

		succ := make([]*BasicBlock, 0, len(bb.Succ))
		pred := make([]*BasicBlock, 0, len(bb.Pred))
		for j := 0; j < len(bb.Succ); j++ {
			succ = append(succ, lirBlocks[bb.Succ[j].ID])
		}
		for j := 0; j < len(bb.Pred); j++ {
			pred = append(pred, lirBlocks[bb.Pred[j].ID])
		}

		lbb.Succ = succ
		lbb.Pred = pred
	}

	legal(lirBlocks, len(regs))

	return &Fn{
		name: name,

		blocks: lirBlocks,
		po:     lirPostorder,

		frameSize: frameSize,
	}
}
