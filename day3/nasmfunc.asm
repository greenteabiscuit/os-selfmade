; naskfunc
; TAB=4

;[FORMAT "WOOFF"]							; オブジェクトファイルを作るモード
;[BITS 32]									; 32ビットモード用の機械語を作らせる

; オブジェクトファイルのための情報

;[FILE "naskfunc.nas"]						; ソースファイル名情報

;		GLOBAL _io_hlt						; このプログラムに含まれる関数名

; 以下は実際の関数

;[SECTION .text] 							; オブジェクトファイルではこれを書いてからプログラムを書く
section .data
section .text
GLOBAL main.io_hlt

main.io_hlt:	; void io_hlt(void);
		RET


