package zorya

import "math"

const (
	regCounter int = iota
	regStackBase
	regStackTip
	regReturn
)

// operand flag masks
const (
	mask1 uint = 0x1 << iota
	mask2
	mask3
	mask4
	mask5
	mask6
)

// operand flag numbers
const (
	arg1 int = iota
	arg2
	arg3
	arg4
	arg5
	arg6
)

// common operand flag numbers
const (
	argDst, maskDst = arg1, mask1
	argSrc, maskSrc = arg2, mask2
	argLhs, maskLhs = arg2, mask2
	argRhs, maskRhs = arg3, mask3
)

type State struct {
}

type Thread struct {
	reg   []float64
	stack []float64
}

// round is used to round a floating point number to the nearest integer.
func round(f float64) float64 {
	if math.IsNaN(f) || math.IsInf(f, 1) || math.IsInf(f, 0) {
		return f
	}

	neg := math.Signbit(f)
	if neg {
		f = -f
	}

	r := math.Floor(f)
	d := f - r
	if d >= 0.5 {
		r += 1
	}

	if neg {
		r = -r
	}

	return float64(r)
}

// Storage returns a read-write pointer to a particular storage area in the
// Thread. Indices greater than zero refer to absolute register indices,
// negative indices refer to values on the stack below ESP (i.e., you cannot
// use Storage to modify values outside the current stack range).
func (th *Thread) Storage(index int) (p *float64, err error) {
	if index >= 0 {
		if index >= len(th.reg) {
			return nil, mkerrorf(CodeBadAccess, "Storage offset for register index %d exceeds register count %d.", index, len(th.reg))
		}
		return &th.reg[index], nil
	} else {
		orig := index
		index += int(th.reg[regStackTip])
		if index < 0 {
			return nil, mkerrorf(CodeBadAccess, "Storage offset for stack index %d exceeds stack size %d.", orig, len(th.stack))
		}
		return &th.stack[index], nil
	}
}

// deref, given a float64, determines if the value references a storage
// location if the mask matches flags and dereferences the appropriate
// storage location and returns its contents. If not, returns val.
func (th *Thread) deref(val float64, flags, mask uint) float64 {
	if flags&mask != mask {
		return val
	}
	p, err := th.Storage(int(val))
	if err != nil {
		// Doing anything with this would put the program in an
		// undefined state, so panic now before it gets worse.
		panic(err)
	}
	return *p
}

// push pushes a value onto the stack, growing it as needed. The only case
// where this fails is when the program is unable to allocate memory, in which
// case it's probably not recoverable.
func (th *Thread) push(val float64) {
	tip := int(th.reg[regStackTip])
	// Let append grow slice capacity on its own, but then make the whole
	// capacity available immediately.
	if len(th.stack) == tip {
		newStack := append(th.stack, val)
		th.stack = newStack[:cap(newStack)]
	} else {
		th.stack[tip] = val
	}
	th.reg[regStackTip] += 1
}

// pop pulls a value off the stack and returns it. If the stack is empty,
// returns the error ErrStackUnderflow.
func (th *Thread) pop() (float64, error) {
	tip := int(th.reg[regStackTip]) - 1
	if tip < int(th.reg[regStackBase]) {
		return 0, ErrStackUnderflow
	}
	th.reg[regStackTip] -= 1
	return th.stack[tip], nil
}

// peek returns the value at the tip of the stack and returns it. If the stack
// is empty, returns the error ErrStackUnderflow.
func (th *Thread) peek() (float64, error) {
	tip := int(th.reg[regStackTip]) - 1
	if tip < int(th.reg[regStackBase]) {
		return 0, ErrStackUnderflow
	}
	th.reg[regStackTip] -= 1
	return th.stack[tip], nil
}

// exec executes the given opcode with the provided operands, using the flag
// to determine which arguments are constants.
func (th *Thread) exec(op Opcode, operands []float64, flag uint) error {
	if !op.valid() {
		return mkerrorf(CodeBadOpcode, "Unrecognized opcode: %d", op)
	} else if len(operands) != op.argc() {
		return mkerrorf(CodeBadOpcode, "Expected %d operands for %v, received %d.", op.argc(), op, len(operands))
	}

	var out *float64 = nil
	var err error
	if op.hasDst() {
		out, err = th.Storage(int(operands[argDst]))
		if err != nil {
			return err
		}
	}

	switch op {
	case OpAdd:
		*out = th.deref(operands[argLhs], flag, maskLhs) + th.deref(operands[argRhs], flag, maskRhs)
	case OpSub:
		*out = th.deref(operands[argLhs], flag, maskLhs) - th.deref(operands[argRhs], flag, maskRhs)
	case OpDiv:
		*out = th.deref(operands[argLhs], flag, maskLhs) / th.deref(operands[argRhs], flag, maskRhs)
	case OpIDiv:
		*out = float64(int64(th.deref(operands[argLhs], flag, maskLhs)) / int64(th.deref(operands[argRhs], flag, maskRhs)))
	case OpMul:
		*out = th.deref(operands[argLhs], flag, maskLhs) * th.deref(operands[argRhs], flag, maskRhs)
	case OpPow:
		*out = math.Pow(th.deref(operands[argLhs], flag, maskLhs), th.deref(operands[argRhs], flag, maskRhs))
	case OpMod:
		*out = math.Mod(th.deref(operands[argLhs], flag, maskLhs), th.deref(operands[argRhs], flag, maskRhs))
	case OpIMod:
		*out = float64(int64(th.deref(operands[argLhs], flag, maskLhs)) % int64(th.deref(operands[argRhs], flag, maskRhs)))
	case OpNeg:
		*out = -th.deref(operands[argSrc], flag, maskSrc)
	case OpNot:
		*out = float64(^uint64(th.deref(operands[argSrc], flag, maskSrc)) & 0xFFFFFFFF)
	case OpOr:
		*out = float64((uint64(th.deref(operands[argLhs], flag, maskLhs)) | uint64(th.deref(operands[argRhs], flag, maskRhs))) & 0xFFFFFFFF)
	case OpAnd:
		*out = float64((uint64(th.deref(operands[argLhs], flag, maskLhs)) & uint64(th.deref(operands[argRhs], flag, maskRhs))) & 0xFFFFFFFF)
	case OpXor:
		*out = float64((uint64(th.deref(operands[argLhs], flag, maskLhs)) ^ uint64(th.deref(operands[argRhs], flag, maskRhs))) & 0xFFFFFFFF)
	case OpArithShift:
		value := int32(th.deref(operands[argLhs], flag, maskLhs))
		shift := int(th.deref(operands[argRhs], flag, maskRhs))
		switch {
		case shift == 0:
		case shift > 0:
			value = value << uint(shift)
		case shift < 0:
			value = value >> uint(-shift)
		}
		*out = float64(value)
	case OpBitShift:
		value := uint32(th.deref(operands[argLhs], flag, maskLhs))
		shift := int(th.deref(operands[argRhs], flag, maskRhs))
		switch {
		case shift == 0:
		case shift > 0:
			value = value << uint(shift)
		case shift < 0:
			value = value >> uint(-shift)
		}
		*out = float64(value)
	case OpFloor:
		*out = math.Floor(th.deref(operands[argSrc], flag, maskSrc))
	case OpCeil:
		*out = math.Ceil(th.deref(operands[argSrc], flag, maskSrc))
	case OpRound:
		*out = round(th.deref(operands[argSrc], flag, maskSrc))
	case OpTrunc:
		*out = float64(int64(th.deref(operands[argSrc], flag, maskSrc)))
	case OpEq:
		comp := th.deref(operands[arg1], flag, mask1) == th.deref(operands[arg2], flag, mask2)
		req := th.deref(operands[arg3], flag, mask3) != 0
		if comp != req {
			th.reg[regCounter] += 1
		}
	case OpLE:
		comp := th.deref(operands[arg1], flag, mask1) <= th.deref(operands[arg2], flag, mask2)
		req := th.deref(operands[arg3], flag, mask3) != 0
		if comp != req {
			th.reg[regCounter] += 1
		}
	case OpLT:
		comp := th.deref(operands[arg1], flag, mask1) < th.deref(operands[arg2], flag, mask2)
		req := th.deref(operands[arg3], flag, mask3) != 0
		if comp != req {
			th.reg[regCounter] += 1
		}
	case OpJump:
		th.reg[regCounter] = th.deref(operands[arg1], flag, mask1)
	case OpPush:
		th.push(th.deref(operands[arg1], flag, mask1))
	case OpPop:
		// No peek instruction, as that's covered by Load Dst, -1
		*out, err = th.pop()
	case OpLoad:
		*out = th.deref(operands[argSrc], flag, maskSrc)
	case OpCall:
		panic("Call not implemented")
	case OpReturn:
		panic("Return not implemented")
	case OpRealloc:
		panic("Realloc not implemented")
	case OpFree:
		panic("Free not implemented")
	case OpMemmove:
		panic("Memmove not implemented")
	case OpTrap:
		panic("Trap not implemented")
	case OpMemdup:
		panic("Memdup not implemented")
	case OpMemlen:
		panic("Memlen not implemented")
	case OpPeek:
		panic("Peek not implemented")
	case OpPoke:
		panic("Poke not implemented")
	case OpDefer:
		panic("Defer not implemented")
	case OpForce:
		panic("Force not implemented")
	}

	return err
}
