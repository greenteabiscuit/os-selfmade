#include "textflag.h"

#define NUM_IDT_ENTRIES 256
#define IDT_ENTRY_SIZE 8 // size 8 for i386

#define ENTRY_TYPE_INTERRUPT_GATE 0x8e

// The 64-bit SIDT consists of 10 bytes and has the following layout:
//   BYTE
// [00 - 01] size of IDT minus 1
// [02 - 09] address of the IDT
GLOBL ·idtDescriptor<>(SB), NOPTR, $10

// The 32-bit IDT consists of NUM_IDT_ENTRIES slots containing 16-byte entries
// with the following layout: TODO
//   BYTE
GLOBL ·idt<>(SB), NOPTR, $NUM_IDT_ENTRIES*IDT_ENTRY_SIZE

// A list of 256 function pointers for installed gate handlers. These pointers
// serve as the jump targets for the trap/int/task dispatchers.
GLOBL ·gateHandlers<>(SB), NOPTR, $NUM_IDT_ENTRIES*8

// installIDT populates idtDescriptor with the address of IDT and loads it to
// the CPU. All gate entries are initially marked as non-present and must be
// explicitly enabled by invoking HandleInterrupt.
TEXT ·installIDT(SB),NOSPLIT,$0
	LEAL ·idtDescriptor<>(SB), AX
	MOVW $(NUM_IDT_ENTRIES*IDT_ENTRY_SIZE)-1, 0(AX)
	LEAL ·idt<>(SB), BX
	MOVL BX, 2(AX)
	MOVL 0(AX), IDTR 	// LIDT[RAX]
	RET

TEXT ·asmIntHandler21(SB),$0-0
	// Save GP regs. The push order MUST match the field layout in the 
	// Registers struct.
    SUBL $12, SP    // neg関数の引数と戻り値サイズ+BPレジスタの退避先を確保
    MOVB BP, 8(SP) // 現在のBPレジスタをpush
    LEAL 8(SP), BP // BPレジスタを新しいスタックに更新
    MOVB AX, (SP)   // 最初の引数iを渡す
    CALL ·IntHandler21(SB)
    MOVB 4(SP), AX  // main.negの戻り値をAXレジスタに取り出す
    MOVB 8(SP), BP // 退避していたBPレジスタをpop
    ADDL $12, SP    // スタックサイズを戻す

	IRETL

// HandleInterrupt ensures that the provided handler will be invoked when a
// particular interrupt number occurs. The value of the istOffset argument
// specifies the offset in the interrupt stack table (if 0 then IST is not
// used).
TEXT ·HandleInterrupt(SB),NOSPLIT,$0-10
	// Dereference pointer to trap handler and copy it into gateHandlers
	MOVQ handler+0(FP), BX
	MOVQ 0(BX), BX
	LEAQ ·gateHandlers<>+0(SB), DI
	MOVQ BX, (DI)(CX*8)

	// Calculate IDT entry address
	LEAQ ·idt<>+0(SB), DI
	MOVQ CX, BX
	SHLQ $IDT_ENTRY_SIZE_SHIFT, BX
	ADDQ BX, DI

	// The trap gate entries have variable lengths depending on whether
	// the CPU pushes an exception code or not. Each generated entry ends
	// with a sequence of 4 NOPs (0x90). The code below uses this information
	// to locate the correct entry point address.
	LEAQ ·interruptGateEntries(SB), SI // SI points to entry for trap 0
update_idt_entry:
	// IDT entry layout (bytes)
	// ------------------------
	// [00-01] bits 0-15 of 32-bit handler address
	// [02-03] CS selector
	// [04-04] RESERVED, DW_COUNT
	// [05-05] gate type/attributes, ACCESS_RIGHT
	// [06-07] bits 16-31 of 32-bit handler address
	//-------------------------