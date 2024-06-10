package vm

import (
	"encoding/binary"

	"github.com/alwaifu/monkey/pkg/object"
)

type Frame struct {
	fn          *object.CompiledFunction
	pc          int //program counter
	basePointer int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{
		fn:          fn,
		pc:          0,
		basePointer: basePointer,
	}
}
func (f *Frame) readInsOprandUint8() uint8 {
	oprand := f.fn.Instructions[f.pc]
	f.pc++ //skip oprand
	return oprand
}
func (f *Frame) readInsOprandUint16() uint16 {
	oprand := binary.BigEndian.Uint16(f.fn.Instructions[f.pc:])
	f.pc += 2 //skip oprand
	return oprand
}
