package asm

import (
	"encoding/binary"
	"math"
)

// Mnemonic        Condition tested      Description
// ----------------------------------------------------------------------------
// jo              OF = 1                overflow
// jno             OF = 0                not overflow
// jc, jb, jnae    CF = 1                carry / below / not above nor equal
// jnc, jae, jnb   CF = 0                not carry / above or equal / not below
// je, jz          ZF = 1                equal / zero
// jne, jnz        ZF = 0                not equal / not zero
// jbe, jna        CF or ZF = 1          below or equal / not above
// ja, jnbe        CF or ZF = 0          above / not below or equal
// js              SF = 1                sign
// jns             SF = 0                not sign
// jp, jpe         PF = 1                parity / parity even
// jnp, jpo        PF = 0                not parity / parity odd
// jl, jnge        SF xor OF = 1         less / not greater nor equal
// jge, jnl        SF xor OF = 0         greater or equal / not less
// jle, jng        (SF xor OF) or ZF = 1 less or equal / not greater
// jg, jnle        (SF xor OF) or ZF = 0 greater / not less nor equal

// Jump continues program flow at the new address.
// The address is relative to the next instruction.
func (a *Assembler) Jump(label string) {
	a.jump(0xeb, []byte{0xe9}, label)
}

// JumpIfLess jumps if the result was less.
func (a *Assembler) JumpIfLess(label string) {
	a.jump(0x7c, []byte{0x0f, 0x8c}, label)
}

// JumpIfLessOrEqual jumps if the result was less or equal.
func (a *Assembler) JumpIfLessOrEqual(label string) {
	a.jump(0x7e, []byte{0x0f, 0x8e}, label)
}

// JumpIfGreater jumps if the result was greater.
func (a *Assembler) JumpIfGreater(label string) {
	a.jump(0x7f, []byte{0x0f, 0x8f}, label)
}

// JumpIfGreaterOrEqual jumps if the result was greater or equal.
func (a *Assembler) JumpIfGreaterOrEqual(label string) {
	a.jump(0x7d, []byte{0x0f, 0x8d}, label)
}

// JumpIfEqual jumps if the result was equal.
func (a *Assembler) JumpIfEqual(label string) {
	a.jump(0x74, []byte{0x0f, 0x84}, label)
}

// JumpIfNotEqual jumps if the result was not equal.
func (a *Assembler) JumpIfNotEqual(label string) {
	a.jump(0x75, []byte{0x0f, 0x85}, label)
}

// jump implements program flow jumps.
func (a *Assembler) jump(shortCode byte, nearCode []byte, label string) {
	instructionPosition := a.Len()
	pointerPosition := instructionPosition + 1
	pointerSize := uint8(1)
	absoluteAddress, exists := a.Labels[label]

	if !exists {
		// TODO: Support 32-bit jumps for unknown labels
		pointer := jumpPointer{pointerPosition, pointerSize}
		a.undefinedJumpLabels[label] = append(a.undefinedJumpLabels[label], pointer)
		a.WriteBytes(shortCode)
		a.WriteBytes(0)
		return
	}

	offset := int32(absoluteAddress - (pointerPosition + uint32(pointerSize)))

	// Near jump (32-bit)
	if offset < math.MinInt8 || offset > math.MaxInt8 {
		a.WriteBytes(nearCode...)
		_ = binary.Write(a, binary.LittleEndian, offset)
		return
	}

	// Short jump (8-bit)
	a.WriteBytes(shortCode)
	a.WriteBytes(byte(offset))
}
