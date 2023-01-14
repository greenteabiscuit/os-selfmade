#include "textflag.h"

#define NUM_IDT_ENTRIES 256
#define IDT_ENTRY_SIZE 8 // size 8 for i386

// The 32-bit SIDT consists of 6 bytes and has the following layout:
//   BYTE
// [00 - 01] size of IDT minus 1
// [02 - 05] address of the IDT, for 32-bit mode
GLOBL ·idtDescriptor<>(SB), NOPTR, $10

TEXT ·add(SB),$24-16
    MOVB i+0(FP),AX    // first arg to AX reg
    MOVB i+2(FP),BX
    ADDB $4, AX        // add number
    ADDB BX, AX
    MOVB AX, ret1+4(FP) // 計算結果を戻り値として返す

    SUBL $12, SP    // neg関数の引数と戻り値サイズ+BPレジスタの退避先を確保
    MOVB BP, 8(SP) // 現在のBPレジスタをpush
    LEAL 8(SP), BP // BPレジスタを新しいスタックに更新
    MOVB AX, (SP)   // 最初の引数iを渡す
    CALL ·sub(SB)   // main.negを呼ぶ
    MOVB 4(SP), AX  // main.negの戻り値をAXレジスタに取り出す
    MOVB 8(SP), BP // 退避していたBPレジスタをpop
    ADDL $12, SP    // スタックサイズを戻す
    MOVB AX, ret2+6(FP) // 2番目の戻り値として返す
    RET

TEXT ·load_idtr(SB),$0-0
    MOVL i+0(FP),AX    // first arg to AX reg
    MOVL i+4(FP),BX
	LEAL ·idtDescriptor<>(SB), AX
	MOVW $(NUM_IDT_ENTRIES*IDT_ENTRY_SIZE)-1, 0(AX)
	MOVL BX, 2(AX)
	MOVL (AX), IDTR 	// LIDT[RAX]
	// MOVW 0(AX), AX
	// MOVW AX, ret+8(FP) // return address for debugging: if returning 0(AX) (as WORD), this should return 2048 - 1 = 2047
    MOVL BX, ret+8(FP) // return address for debugging: if returning 0(AX) (as WORD), this should return 2048 - 1 = 2047
	RET


TEXT ·PortWriteByte(SB),NOSPLIT,$0
	MOVW port+0(FP), DX
	MOVB val+2(FP), AX
	BYTE $0xee // out al, dx
	RET

TEXT ·PortReadByte(SB),NOSPLIT,$0
	MOVW port+0(FP), DX
	BYTE $0xec  // in al, dx
	MOVB AX, ret+0(FP)
	RET
