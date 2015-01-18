package zorya

// Opcode represents a single integer opcode in Zorya.
type Opcode int

const (
	opAdd Opcode = iota
	opSub
	opDiv
	opIDiv
	opMul
	opPow
	opMod
	opIMod
	opNeg
	opNot
	opOr
	opAnd
	opXor
	opArithShift
	opBitShift
	opFloor
	opCeil
	opRound
	opRInt
	opEq
	opLE
	opLT
	opJump
	opPush
	opPop
	opLoad
	opCall
	opReturn
	opRealloc
	opFree
	opMemmove
	opTrap
	opMemdup
	opMemlen
	opPeek
	opPoke
	opDefer
	opForce
)

func (op *Opcode) String() string {
	switch *op {
	case opAdd:
		return "Add"
	case opSub:
		return "Sub"
	case opDiv:
		return "Div"
	case opIDiv:
		return "IDiv"
	case opMul:
		return "Mul"
	case opPow:
		return "Pow"
	case opMod:
		return "Mod"
	case opIMod:
		return "IMod"
	case opNeg:
		return "Neg"
	case opNot:
		return "Not"
	case opOr:
		return "Or"
	case opAnd:
		return "And"
	case opXor:
		return "Xor"
	case opArithShift:
		return "ArithShift"
	case opBitShift:
		return "BitShift"
	case opFloor:
		return "Floor"
	case opCeil:
		return "Ceil"
	case opRound:
		return "Round"
	case opRInt:
		return "RInt"
	case opEq:
		return "Eq"
	case opLE:
		return "LE"
	case opLT:
		return "LT"
	case opJump:
		return "Jump"
	case opPush:
		return "Push"
	case opPop:
		return "Pop"
	case opLoad:
		return "Load"
	case opCall:
		return "Call"
	case opReturn:
		return "Return"
	case opRealloc:
		return "Realloc"
	case opFree:
		return "Free"
	case opMemmove:
		return "Memmove"
	case opTrap:
		return "Trap"
	case opMemdup:
		return "Memdup"
	case opMemlen:
		return "Memlen"
	case opPeek:
		return "Peek"
	case opPoke:
		return "Poke"
	case opDefer:
		return "Defer"
	case opForce:
		return "Force"
	default:
		return "<UNKNOWN>"
	}
}
