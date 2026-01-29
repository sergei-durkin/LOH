1. [+] AST
2. [+] Parser
3. [+] HIR
4. [-] SSA
	1. [+] CFG simplification
	2. [+] Dominator Tree
	3. [+] Phi simplification
	4. [ ] Memory SSA
	5. [-] SSA Optimizations:
		1. [+] Copy Propagation
		2. [+] Constant Folding
		3. [+] Constant Propagation
		4. [ ] Common Subexp Elimination
		5. [+] Dead Code Elimination
		6. [ ] Loop optimizations:
			1. [ ] Loop Simplify
			2. [ ] LCSSA
			3. [ ] Loop Invariant Code Motion (LICM)
			4. [ ] Strength Reduction
			5. [ ] Induction Variable Simplification
			6. [ ] Loop Unrolling
		7. [ ] Interprocedural optimizations:
			1. [ ] Inlining
			2. [ ] Function cloning
			3. [ ] Dead function elimination
			4. [ ] Argument specialization
			5. [ ] Constant argument propagation
			6. [ ] Escape analysis
	6. [+] Liveness analysis
	7. [+] Lowering
	8. [+] Destroy
	9. [+] Instruction Selection
5. [+] LIR
6. [+] Target Machine Model
	1. [+] Basic AArch64
	2. [ ] RISC-V like arch
7. [+] Regalloc
8. [+] Instruction Scheduling
9. [+] Assembly Emission
