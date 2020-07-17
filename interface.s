#include "textflag.h"
#include "funcdata.h"

// makeFuncStub is the code half of the function returned by MakeFunc.
// See the comment on the declaration of makeFuncStub in makefunc.go
// for more details.
// No arg size here; runtime pulls arg map out of the func value.
TEXT ·InterfaceCallStub(SB),(NOSPLIT|WRAPPER),$0
	NO_LOCAL_POINTERS
    LEAQ	0(SP), DX
    JMP	reflect·makeFuncStub(SB)
    RET
