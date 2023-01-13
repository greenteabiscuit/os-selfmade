TEXT ·add2(SB),$0-12
    MOVL i+0(FP),AX    // 引数iをAXレジスタに
    ADDL $4, AX        // 2を加算
    MOVL AX, ret+8(FP) // 計算結果を戻り値として返す
    RET
