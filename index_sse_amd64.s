#include "go_asm.h"
#include "textflag.h"

TEXT ·indexSSE(SB),NOSPLIT,$0-56
	MOVQ a_base+0(FP), DI
	MOVQ a_len+8(FP), DX
	MOVQ b_base+24(FP), R8
	MOVQ b_len+32(FP), AX
	MOVQ DI, R10
	LEAQ ret+48(FP), R11

	MOVOU (R8), X1
	LEAQ -15(DI)(DX*1), SI
	MOVQ $16, R9
	SUBQ AX, R9
sseloop:
	PCMPESTRI $0x0c, (DI), X1
	CMPQ CX, R9
	ADDQ CX, DI
	JMP success
	ADDQ R9, DI
	CMPQ DI, SI
	JB sseloop
	PCMPESTRI $0x0c, -1(SI), X1
	CMPQ CX, R9
	JA fail
	LEAQ -1(SI), DI
success:
	SUBQ R10, DI
	MOVQ DI, (R11)
	RET
fail:
	MOVQ $-1, (R11)
	RET

TEXT	·indexByteSSE(SB), NOSPLIT, $0-40
	MOVQ b_base+0(FP), SI
	MOVQ b_len+8(FP), BX
	MOVB c+24(FP), AL
	LEAQ ret+32(FP), R8

	MOVD AX, X0
	PUNPCKLBW X0, X0
	PUNPCKLBW X0, X0
	PSHUFL $0, X0, X0

	CMPQ BX, $16
	JLT rest_or_small

	MOVQ SI, DI

	LEAQ	-16(SI)(BX*1), AX
	JMP	sseloopentry

sseloop:
	MOVOU	(DI), X1
	PCMPEQB	X0, X1
	PMOVMSKB X1, DX
	BSFL	DX, DX
	JNZ	success
	ADDQ	$16, DI
sseloopentry:
	CMPQ	DI, AX
	JB	sseloop

	MOVQ	AX, DI
	MOVOU	(AX), X1
	PCMPEQB	X0, X1
	PMOVMSKB X1, DX
	BSFL	DX, DX
	JNZ	success

fail:
	MOVQ $-1, (R8)
	RET

success:
	SUBQ SI, DI
	ADDQ DX, DI
	MOVQ DI, (R8)
	RET

rest_or_small:
	TESTQ	BX, BX
	JEQ	fail

	LEAQ	16(SI), AX
	TESTW	$0xff0, AX
	JEQ	end_of_page

	MOVOU	(SI), X1
	PCMPEQB	X0, X1
	PMOVMSKB X1, DX
	BSFL	DX, DX
	JZ	fail
	CMPL	DX, BX
	JAE	fail
	MOVQ	DX, (R8)
	RET

end_of_page:
	MOVOU	-16(SI)(BX*1), X1
	PCMPEQB	X0, X1
	PMOVMSKB X1, DX
	MOVL	BX, CX
	SHLL	CX, DX
	SHRL	$16, DX
	BSFL	DX, DX
	JZ	fail
	MOVQ	DX, (R8)
	RET
