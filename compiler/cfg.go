package compiler

import (
	"sort"
)

type Cfg struct {
	name string

	gh map[int][]int
	mp map[int]*BasicBlock

	blocks []*BasicBlock

	pos int
	cur Instruction

	tac   []Instruction
	size  map[string]int64
	types map[string]*Struct
}

func NewCfg(tac *Tac) *Cfg {
	c := &Cfg{
		name: tac.name,

		gh: make(map[int][]int),
		mp: make(map[int]*BasicBlock),

		tac:   tac.code,
		size:  tac.size,
		types: tac.structs,
	}

	c.constructBlocks()
	c.placeReturn()

	// TODO refactor this щит
	for from, to := range c.gh {
		bb := c.mp[from]
		for i := 0; i < len(to); i++ {
			suc := c.mp[to[i]]
			if suc == nil {
				continue
			}

			bb.Succ = append(bb.Succ, suc)
			suc.Pred = append(suc.Pred, bb)
		}
	}

	sortBlocks := func(b []*BasicBlock) {
		sort.Slice(b, func(i, j int) bool {
			return b[i].ID < b[j].ID
		})
	}

	for i := 0; i < len(c.blocks); i++ {
		sortBlocks(c.blocks[i].Succ)
		sortBlocks(c.blocks[i].Pred)
	}

	return c
}

func (c *Cfg) constructBlocks() {
	c.consume()
	if c.cur == nil {
		return
	}

	l, ok := c.cur.(*Label)
	if !ok {
		return
	}

	bb := c.NewBasicBlock(l.id)
	c.mp[bb.ID] = bb
	c.blocks = append(c.blocks, bb)

	for c.peek() != nil {
		c.consume()
		switch instr := c.cur.(type) {
		case *Label:
			b := c.NewBasicBlock(instr.id)
			c.mp[b.ID] = b
			c.blocks = append(c.blocks, b)

			// TODO refactor:
			if len(bb.Instr) > 0 {
				switch bb.Instr[len(bb.Instr)-1].(type) {
				case *IfGoto, *Goto, *Return:
				default:
					bb.Instr = append(bb.Instr, &Goto{label: instr.id})
					c.gh[bb.ID] = append(c.gh[bb.ID], b.ID)
				}
			} else {
				bb.Instr = append(bb.Instr, &Goto{label: instr.id})
				c.gh[bb.ID] = append(c.gh[bb.ID], b.ID)
			}

			bb = b
		case *Goto:
			bb.Instr = append(bb.Instr, instr)

			c.gh[bb.ID] = append(c.gh[bb.ID], instr.label)

			for c.peek() != nil {
				if _, ok := c.peek().(*Label); ok {
					break
				}

				c.consume()
			}
		case *Return:
			bb.Instr = append(bb.Instr, instr)

			for c.peek() != nil {
				if _, ok := c.peek().(*Label); ok {
					break
				}

				c.consume()
			}
		case *IfGoto:
			bb.Instr = append(bb.Instr, instr)

			c.gh[bb.ID] = append(c.gh[bb.ID], instr.label)
			c.gh[bb.ID] = append(c.gh[bb.ID], instr.fall)
		default:
			bb.Instr = append(bb.Instr, instr)
		}
	}
}

func (c *Cfg) placeReturn() {
	for i := 0; i < len(c.blocks); i++ {
		bb := c.blocks[i]
		if len(bb.Instr) > 0 {
			switch bb.Instr[len(bb.Instr)-1].(type) {
			case *IfGoto, *Goto, *Return:
			default:
				bb.Instr = append(bb.Instr, &Return{})
			}
			continue
		}

		bb.Instr = append(bb.Instr, &Return{})
	}
}

func (c *Cfg) NewBasicBlock(id int) *BasicBlock {
	return &BasicBlock{
		ID: id,
	}
}

func (c *Cfg) consume() {
	if c.pos >= len(c.tac) {
		return
	}

	c.cur = c.tac[c.pos]
	c.pos++
}

func (c *Cfg) peek() Instruction {
	if c.pos >= len(c.tac) {
		return nil
	}

	return c.tac[c.pos]
}
