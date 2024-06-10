package vm

import "encoding/binary"

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGt
	OpMinus
	OpBang
	OpJump
	OpJumpNotTruthy
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpGetLocal
	OpSetLocal
	OpGetBuiltin
	OpClosure
	OpGetFree
	OpCurrentClosure
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]Definition{
	OpConstant:       {"OpConstant", []int{2}},
	OpAdd:            {"OpAdd", []int{}},
	OpPop:            {"OpPop", []int{}},
	OpSub:            {"OpSub", []int{}},
	OpMul:            {"OpMul", []int{}},
	OpDiv:            {"OpDiv", []int{}},
	OpTrue:           {"OpTrue", []int{}},
	OpFalse:          {"OpFalse", []int{}},
	OpEqual:          {"OpEqual", []int{}},
	OpNotEqual:       {"OpNotEqual", []int{}},
	OpGt:             {"OpGt", []int{}},
	OpMinus:          {"OpMinus", []int{}},
	OpBang:           {"OpBang", []int{}},
	OpJump:           {"OpJump", []int{2}},
	OpJumpNotTruthy:  {"OpJumpNotTruthy", []int{2}},
	OpNull:           {"OpNull", []int{}},
	OpGetGlobal:      {"OpGetGlobal", []int{2}},
	OpSetGlobal:      {"OpSetGlobal", []int{2}},
	OpArray:          {"OpArray", []int{2}},
	OpHash:           {"OpHash", []int{}},
	OpIndex:          {"OpIndex", []int{}},
	OpCall:           {"OpCall", []int{1}},
	OpReturnValue:    {"OpReturnValue", []int{}},
	OpReturn:         {"OpReturn", []int{}},
	OpGetLocal:       {"OpGetLocal", []int{1}},
	OpSetLocal:       {"OpSetLocal", []int{1}},
	OpGetBuiltin:     {"OpGetBuiltin", []int{1}},
	OpClosure:        {"OpClosure", []int{}},
	OpGetFree:        {"OpGetFree", []int{}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

func MakeInstruction(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instructions := make(Instructions, instructionLen)
	instructions[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 1:
			instructions[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instructions[offset:], uint16(o))
		case 4:
			binary.BigEndian.PutUint32(instructions[offset:], uint32(o))
		case 8:
			binary.BigEndian.PutUint64(instructions[offset:], uint64(o))
		}
		offset += width
	}
	return instructions
}
