#include "textflag.h"
#include "funcdata.h"

// InterfaceCallStub 恢复DX
TEXT ·InterfaceCallStub(SB),(NOSPLIT),$0
	NO_LOCAL_POINTERS
	MOVQ	8(SP), DX
    JMP	reflect·makeFuncStub(SB)
    RET
