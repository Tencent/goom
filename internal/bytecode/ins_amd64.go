package bytecode

import (
	"git.woa.com/goom/mocker/internal/arch/x86asm"
	"git.woa.com/goom/mocker/internal/logger"
)

// ParseIns parse instruction
func ParseIns(pos int, copyOrigin []byte) (*x86asm.Inst, []byte, error) {
	if pos >= len(copyOrigin) {
		return nil, nil, nil
	}
	// read 16 bytes at most each time
	endPos := pos + 16
	if endPos > len(copyOrigin) {
		endPos = len(copyOrigin)
	}

	code := copyOrigin[pos:endPos]
	ins, err := x86asm.Decode(code, 64)
	if err != nil {
		logger.Error("decode assembly code err:", err)
	}
	return &ins, code, err
}
