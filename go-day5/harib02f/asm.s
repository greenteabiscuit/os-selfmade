TEXT ·add(SB),$0-0
    MOVB i+0(FP),AX    // 引数iをAXレジスタに
    ADDB $4, AX        // 2を加算
    MOVB AX, ret+4(FP) // 計算結果を戻り値として返す
    RET
