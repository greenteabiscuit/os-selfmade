TEXT ·add(SB),$0-0
    MOVB i+0(FP),AX    // first arg to AX reg
    MOVB i+2(FP),BX
    ADDB $4, AX        // add number
    ADDB BX, AX
    MOVB AX, ret+4(FP) // 計算結果を戻り値として返す

    SUBL $24, SP    // neg関数の引数と戻り値サイズ+BPレジスタの退避先を確保
    // MOVB BP, 16(SP) // 現在のBPレジスタをpush
    // LEAQ 16(SP), BP // BPレジスタを新しいスタックに更新
    // MOVB AX, (SP)   // 最初の引数iを渡す
    // CALL ·sub(SB)   // main.negを呼ぶ
    // MOVB 8(SP), AX  // main.negの戻り値をAXレジスタに取り出す
    // MOVB 16(SP), BP // 退避していたBPレジスタをpop
    ADDL $24, SP    // スタックサイズを戻す
    // MOVB AX, ret2+12(FP) // 2番目の戻り値として返す
    RET
