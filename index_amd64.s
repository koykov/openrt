#include "go_asm.h"
#include "textflag.h"

TEXT ·indexSSE(SB),NOSPLIT,$0-56
	MOVQ a_base+0(FP), DI
	MOVQ a_len+8(FP), DX
	MOVQ b_base+24(FP), R8
	MOVQ b_len+32(FP), AX
	MOVQ DI, R10
	LEAQ ret+48(FP), R11
	JMP  indexbodySSE<>(SB)

TEXT indexbodySSE<>(SB),NOSPLIT,$0
	MOVOU (R8), X1
	LEAQ -15(DI)(DX*1), SI
	MOVQ $16, R9
	SUBQ AX, R9 // We advance by 16-len(sep) each iteration, so precalculate it into R9
loop_sse:
	// 0x0c means: unsigned byte compare (bits 0,1 are 00)
	// for equality (bits 2,3 are 11)
	// result is not masked or inverted (bits 4,5 are 00)
	// and corresponds to first matching byte (bit 6 is 0)
	PCMPESTRI $0x0c, (DI), X1
	// CX == 16 means no match,
	// CX > R9 means partial match at the end of the string,
	// otherwise sep is at offset CX from X1 start
	CMPQ CX, R9
	JBE success_sse
	ADDQ R9, DI
	CMPQ DI, SI
	JB loop_sse
	PCMPESTRI $0x0c, -1(SI), X1
	CMPQ CX, R9
	JA fail
	LEAQ -1(SI), DI
success_sse:
	ADDQ CX, DI
	JMP success
success:
	SUBQ R10, DI
	MOVQ DI, (R11)
	RET
fail:
	MOVQ $-1, (R11)
	RET
