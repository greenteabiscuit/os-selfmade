TEXT ·add(SB),$0-0
    MOVB i+0(FP),AX    // 引数iをAXレジスタに
    MOVB i+2(FP),BX
    ADDB $4, AX        // 2を加算
    ADDB BX, AX
    MOVB AX, ret+4(FP) // 計算結果を戻り値として返す
    RET
