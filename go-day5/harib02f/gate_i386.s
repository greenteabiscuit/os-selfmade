#include "textflag.h"

#define NUM_IDT_ENTRIES 256
#define IDT_ENTRY_SIZE 8 // size 8 for i386

#define ENTRY_TYPE_INTERRUPT_GATE 0x8e

#define IDT_ENTRY_SIZE_SHIFT 4

// The 32-bit SIDT consists of 6 bytes and has the following layout:
//   BYTE
// [00 - 01] size of IDT minus 1
// [02 - 05] address of the IDT, for 32-bit mode
GLOBL ·idtDescriptor<>(SB), NOPTR, $10

// The 32-bit IDT consists of NUM_IDT_ENTRIES slots containing 16-byte entries
// with the following layout:
	// [00-01] bits 0-15 of 32-bit handler address
	// [02-03] CS selector
	// [04-04] RESERVED, DW_COUNT
	// [05-05] gate type/attributes, ACCESS_RIGHT
	// [06-07] bits 16-31 of 32-bit handler address
	//-------------------------
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
	MOVL (AX), IDTR 	// LIDT[RAX]
	MOVL 2(AX), AX
	MOVL AX, ret+0(FP) // return address for debugging: if returning 0(AX) (as WORD), this should return 2048 - 1 = 2047
	// if returning address, use below:
	// MOVL 2(AX), AX
	// MOVL AX, ret+0(FP)
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

TEXT ·GetIDTAddr(SB),$0-0
    LEAL ·idt<>+0(SB), DI
    MOVL DI, AX
    MOVL AX, ret+0(FP)
    RET

// HandleInterrupt ensures that the provided handler will be invoked when a
// particular interrupt number occurs. The value of the istOffset argument
// specifies the offset in the interrupt stack table (if 0 then IST is not
// used).
TEXT ·HandleInterrupt(SB),NOSPLIT,$0-10
    MOVL $21, CX

	// Dereference pointer to trap handler and copy it into gateHandlers
	MOVL handler+0(FP), BX
	MOVL 0(BX), BX
	LEAL ·gateHandlers<>+0(SB), DI
	MOVL BX, (DI)(CX*8)

	// Calculate IDT entry address
	LEAL ·idt<>+0(SB), DI
	MOVL CX, BX
	SHLL $IDT_ENTRY_SIZE_SHIFT, BX
	ADDL BX, DI

	// The trap gate entries have variable lengths depending on whether
	// the CPU pushes an exception code or not. Each generated entry ends
	// with a sequence of 4 NOPs (0x90). The code below uses this information
	// to locate the correct entry point address.
	LEAL ·interruptGateEntries(SB), SI // SI points to entry for trap 0
update_idt_entry:
	// IDT entry layout (bytes)
	// ------------------------
	// [00-01] bits 0-15 of 32-bit handler address
	// [02-03] CS selector
	// [04-04] RESERVED, DW_COUNT
	// [05-05] gate type/attributes, ACCESS_RIGHT
	// [06-07] bits 16-31 of 32-bit handler address
	//-------------------------

	// Mark entry as non-present while updating the handler address
	MOVB $0, 5(DI)

	MOVW $16, 2(DI) // selector

    // Copy the entrypoint address from SI
    MOVW SI, 0(DI)
    SHRL $16, SI
    MOVW SI, 6(DI)

    // Mark entry as a present, 32-bit interrupt gate
    MOVB $ENTRY_TYPE_INTERRUPT_GATE, 5(DI) // gd->access_right = ar & 0xff

    RET

// Emit interrupt dispatching code for traps where the CPU pushes an exception
// code to the stack. The code below just pushes the handler's address to the
// stack and jumps to dispatchInterrupt.
//
// This code uses some tricks to bypass Go assembler limitations:
// - replace PUSH with: SUBQ $8, RSP; MOVQ X, 0(RSP). This prevents the Go
//   assembler from complaining about unbalanced PUSH/POP statements.
// - use a PUSH/RET (0xc3 byte instead of RET mnemonic) trick instead of a
//   "JMP dispatchInterrupt" to prevent the optimizer from optimizing away all
//   but the first entry in interruptGateEntries.
//
// Finally, each entry block ends with a series of 4 NOP instructions. This
// delimiter is used by the HandleInterrupt implementation to locate the correct
// entrypoint address for a particular interrupt.
#define INT_ENTRY_WITH_CODE(num)

// Emit interrupt dispatching code for traps where the CPU does not push an
// exception code to the stack. The implementation is identical with the
// INT_ENTRY_WITH_CODE above with the exception that the interrupt number is
// manually pushed to the stack before the handler address so both entry
// variants can use the same dispatching code.
#define INT_ENTRY_WITHOUT_CODE(num) \
	SUBL $24, SP;                          \
	MOVL BX, 0(SP);                       \
	MOVL ·gateHandlers<>+4*num(SB), BX;   \
    MOVL BX, 4(SP);                       \
    MOVL $num, 8(SP);                     \
    LEAL ·dispatchInterrupt(SB), BX;      \
    XCHGL BX, 0(SP);                      \
    JMP ·dispatchInterrupt(SB);                            \
    BYTE $0x90; BYTE $0x90;

// dispatchInterrupt is invoked by the interrupt gate entrypoints to route 
// an incoming interrupt to the selected handler.
//
// Callers MUST ensure that the stack has the following layout before calling 
// dispatchInterrupt:
//
// |-----------------| <=== SP after jumping to dispatchInterrupt
// | handler address | <- pushed by the interrupt entry code
// |-----------------|
// | exception code  | <- pushed by CPU or a dummy code pushed by the gate entry
// |-----------------|
// | RIP             | <- pushed by CPU (exception frame)
// | CS              |
// | RFLAGS          |
// | RSP             |
// | SS              |
// |-----------------|
//
// Once the handler returns, the GP regs are restored and the stack is unwinded
// so that the CPU can resume excecution of the code that triggered the
// interrupt.
//
// Interrupts are automatically disabled by the CPU upon entry and re-enabled
// when this function returns.
//--------------------------- -----------------------------------------
TEXT ·dispatchInterrupt(SB),$0-0
	IRETL

// interruptGateEntries contains a list of generated entries for each possible
// interrupt number. Depending on the
TEXT ·interruptGateEntries(SB),NOSPLIT,$0
	// For a list of gate numbers that push an error code see: http://wiki.osdev.org/Exceptions
	INT_ENTRY_WITHOUT_CODE(0)
    RET
