	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 12, 0
	.globl	_main                           ; -- Begin function main
	.p2align	2
_main:                                  ; @main
	.cfi_startproc
; %bb.0:
	sub	sp, sp, #32                     ; =32
	stp	x29, x30, [sp, #16]             ; 16-byte Folded Spill
	.cfi_def_cfa_offset 32
	.cfi_offset w30, -8
	.cfi_offset w29, -16
	mov	x9, #25928
	movk	x9, #27756, lsl #16
	mov	w8, #27762
	movk	x9, #8303, lsl #32
	movk	w8, #100, lsl #16
	movk	x9, #28535, lsl #48
	add	x0, sp, #4                      ; =4
	str	w8, [sp, #12]
	stur	x9, [sp, #4]
	bl	_puts
	ldp	x29, x30, [sp, #16]             ; 16-byte Folded Reload
	mov	w0, wzr
	add	sp, sp, #32                     ; =32
	ret
	.cfi_endproc
                                        ; -- End function
.subsections_via_symbols
