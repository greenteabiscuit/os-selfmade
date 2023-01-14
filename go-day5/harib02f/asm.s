#include "textflag.h"

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
