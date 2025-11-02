// go:build darwin && arm64
// +build darwin,arm64

#include "textflag.h"

GLOBL	·libc_proc_listallpids_trampoline_addr(SB), RODATA, $8
DATA	·libc_proc_listallpids_trampoline_addr(SB)/8, $libc_proc_listallpids_trampoline<>(SB)

TEXT libc_proc_listallpids_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_proc_listallpids(SB)

GLOBL	·libc_proc_name_trampoline_addr(SB), RODATA, $8
DATA	·libc_proc_name_trampoline_addr(SB)/8, $libc_proc_name_trampoline<>(SB)

TEXT libc_proc_name_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_proc_name(SB)

GLOBL	·libc_proc_pidinfo_trampoline_addr(SB), RODATA, $8
DATA	·libc_proc_pidinfo_trampoline_addr(SB)/8, $libc_proc_pidinfo_trampoline<>(SB)

TEXT libc_proc_pidinfo_trampoline<>(SB),NOSPLIT,$0-0
	JMP	libc_proc_pidinfo(SB)
