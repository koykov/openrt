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

TEXT	·indexByteSSE(SB), NOSPLIT, $0-40
	MOVQ b_base+0(FP), SI
	MOVQ b_len+8(FP), BX
	MOVB c+24(FP), AL
	LEAQ ret+32(FP), R8
	JMP  indexbytebodySSE<>(SB)

TEXT	indexbytebodySSE<>(SB), NOSPLIT, $0
	// Shuffle X0 around so that each byte contains
	// the character we're looking for.
	MOVD AX, X0
	PUNPCKLBW X0, X0
	PUNPCKLBW X0, X0
	PSHUFL $0, X0, X0

	CMPQ BX, $16
	JLT small

	MOVQ SI, DI

	LEAQ	-16(SI)(BX*1), AX	// AX = address of last 16 bytes
	JMP	sseloopentry

sseloop:
	// Move the next 16-byte chunk of the data into X1.
	MOVOU	(DI), X1
	// Compare bytes in X0 to X1.
	PCMPEQB	X0, X1
	// Take the top bit of each byte in X1 and put the result in DX.
	PMOVMSKB X1, DX
	// Find first set bit, if any.
	BSFL	DX, DX
	JNZ	ssesuccess
	// Advance to next block.
	ADDQ	$16, DI
sseloopentry:
	CMPQ	DI, AX
	JB	sseloop

	// Search the last 16-byte chunk. This chunk may overlap with the
	// chunks we've already searched, but that's ok.
	MOVQ	AX, DI
	MOVOU	(AX), X1
	PCMPEQB	X0, X1
	PMOVMSKB X1, DX
	BSFL	DX, DX
	JNZ	ssesuccess

failure:
	MOVQ $-1, (R8)
	RET

ssesuccess:
	SUBQ SI, DI	// Compute offset of chunk within data.
	ADDQ DX, DI	// Add offset of byte within chunk.
	MOVQ DI, (R8)
	RET

small:
	TESTQ	BX, BX
	JEQ	failure

	// Check if we'll load across a page boundary.
	LEAQ	16(SI), AX
	TESTW	$0xff0, AX
	JEQ	endofpage

	MOVOU	(SI), X1 // Load data
	PCMPEQB	X0, X1	// Compare target byte with each byte in data.
	PMOVMSKB X1, DX	// Move result bits to integer register.
	BSFL	DX, DX	// Find first set bit.
	JZ	failure	// No set bit, failure.
	CMPL	DX, BX
	JAE	failure	// Match is past end of data.
	MOVQ	DX, (R8)
	RET

endofpage:
	MOVOU	-16(SI)(BX*1), X1	// Load data into the high end of X1.
	PCMPEQB	X0, X1	// Compare target byte with each byte in data.
	PMOVMSKB X1, DX	// Move result bits to integer register.
	MOVL	BX, CX
	SHLL	CX, DX
	SHRL	$16, DX	// Shift desired bits down to bottom of register.
	BSFL	DX, DX	// Find first set bit.
	JZ	failure	// No set bit, failure.
	MOVQ	DX, (R8)
	RET
