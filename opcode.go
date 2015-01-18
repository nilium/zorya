package zorya

// Opcode represents a single integer opcode in Zorya. The lower 16 bits are
// reserved for the Opcode integer itself (a 16-bit integer) and the remaining
// upper 15 bits, excluding sign bit, are reserved for attribute flags set on
// Opcodes to determine their usage (e.g., opDest signals that the first
// operand is a destination register).
//
// Attributes are not stored in bytecode.
type Opcode int

// Opcode attributes
const (
	// opDest marks an Opcode as having an output as its first operand.
	opDest = 1<<18 + iota
)

const (
	opArgShift = 28
	opArgMask  = 0x7
)

// These would be better as const structs if that were possible, but for now,
// this works.
const (
	// NAME             = iota | (ARGC)            | ATTRS
	OpAdd        Opcode = iota | (3 << opArgShift) | opDest
	OpSub        Opcode = iota | (3 << opArgShift) | opDest
	OpDiv        Opcode = iota | (3 << opArgShift) | opDest
	OpIDiv       Opcode = iota | (3 << opArgShift) | opDest
	OpMul        Opcode = iota | (3 << opArgShift) | opDest
	OpPow        Opcode = iota | (3 << opArgShift) | opDest
	OpMod        Opcode = iota | (3 << opArgShift) | opDest
	OpIMod       Opcode = iota | (3 << opArgShift) | opDest
	OpArithShift Opcode = iota | (3 << opArgShift) | opDest
	OpBitShift   Opcode = iota | (3 << opArgShift) | opDest
	OpOr         Opcode = iota | (3 << opArgShift) | opDest
	OpAnd        Opcode = iota | (3 << opArgShift) | opDest
	OpXor        Opcode = iota | (3 << opArgShift) | opDest
	OpNot        Opcode = iota | (2 << opArgShift) | opDest
	OpNeg        Opcode = iota | (2 << opArgShift) | opDest
	OpFloor      Opcode = iota | (2 << opArgShift) | opDest
	OpCeil       Opcode = iota | (2 << opArgShift) | opDest
	OpRound      Opcode = iota | (2 << opArgShift) | opDest
	OpTrunc      Opcode = iota | (2 << opArgShift) | opDest
	OpEq         Opcode = iota | (3 << opArgShift)
	OpLE         Opcode = iota | (3 << opArgShift)
	OpLT         Opcode = iota | (3 << opArgShift)
	OpJump       Opcode = iota | (1 << opArgShift)
	OpPush       Opcode = iota | (1 << opArgShift)
	OpPop        Opcode = iota | (1 << opArgShift) | opDest
	OpLoad       Opcode = iota | (2 << opArgShift) | opDest
	OpCall       Opcode = iota | (2 << opArgShift)
	OpReturn     Opcode = iota | (0 << opArgShift) | opDest
	OpRealloc    Opcode = iota | (3 << opArgShift) | opDest
	OpFree       Opcode = iota | (1 << opArgShift)
	OpMemmove    Opcode = iota | (5 << opArgShift) | opDest
	OpTrap       Opcode = iota | (0 << opArgShift) | opDest
	OpMemdup     Opcode = iota | (2 << opArgShift) | opDest
	OpMemlen     Opcode = iota | (2 << opArgShift) | opDest
	OpPeek       Opcode = iota | (4 << opArgShift) | opDest
	OpPoke       Opcode = iota | (4 << opArgShift)
	OpDefer      Opcode = iota | (1 << opArgShift) | opDest
	OpForce      Opcode = iota | (2 << opArgShift) | opDest
)

func (op Opcode) argc() int {
	return int(op>>opArgShift) & opArgMask
}

func (op Opcode) valid() bool {
	op = op.strict()
	return op >= OpAdd.strict() && op <= OpForce.strict()
}

func (op Opcode) strict() Opcode {
	return op & 0xFFFF
}

func (op Opcode) attrs() int {
	return int(op) & 0x7FFF0000
}

func (op Opcode) hasDst() bool {
	return op&opDest == opDest
}

func (op *Opcode) String() string {
	switch *op {
	case OpAdd:
		return "Add"
	case OpSub:
		return "Sub"
	case OpDiv:
		return "Div"
	case OpIDiv:
		return "IDiv"
	case OpMul:
		return "Mul"
	case OpPow:
		return "Pow"
	case OpMod:
		return "Mod"
	case OpIMod:
		return "IMod"
	case OpNeg:
		return "Neg"
	case OpNot:
		return "Not"
	case OpOr:
		return "Or"
	case OpAnd:
		return "And"
	case OpXor:
		return "Xor"
	case OpArithShift:
		return "ArithShift"
	case OpBitShift:
		return "BitShift"
	case OpFloor:
		return "Floor"
	case OpCeil:
		return "Ceil"
	case OpRound:
		return "Round"
	case OpTrunc:
		return "Trunc"
	case OpEq:
		return "Eq"
	case OpLE:
		return "LE"
	case OpLT:
		return "LT"
	case OpJump:
		return "Jump"
	case OpPush:
		return "Push"
	case OpPop:
		return "POp"
	case OpLoad:
		return "Load"
	case OpCall:
		return "Call"
	case OpReturn:
		return "Return"
	case OpRealloc:
		return "Realloc"
	case OpFree:
		return "Free"
	case OpMemmove:
		return "Memmove"
	case OpTrap:
		return "Trap"
	case OpMemdup:
		return "Memdup"
	case OpMemlen:
		return "Memlen"
	case OpPeek:
		return "Peek"
	case OpPoke:
		return "Poke"
	case OpDefer:
		return "Defer"
	case OpForce:
		return "Force"
	default:
		return "<UNKNOWN>"
	}
}
