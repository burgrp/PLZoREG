package flash

import (
	"runtime"
	"runtime/volatile"
	"unsafe"

	// Use the import path where your generated py32f030xx.go lives.
	// The *package name* in that file is: package py32
	"device/py32"
)

const (
	MainFlashBase = 0x08000000
	MainFlashSize = 64 * 1024

	PageSize     = 128
	WordsPerPage = PageSize / 4

	key1 = 0x45670123
	key2 = 0xCDEF89AB
)

func AlignDownToPage(addr uint32) uint32 { return addr &^ uint32(PageSize-1) }

func inMainFlash(addr uint32) bool {
	return addr >= MainFlashBase && addr < (MainFlashBase+MainFlashSize)
}

func waitNotBusy() {
	for py32.FLASH.SR.Get()&py32.Flash_SR_BSY_Msk != 0 {
	}
}

func clearFlags() {
	// SR flags are RC_W1: write 1 to clear (do NOT use the SetSR_* RMW helpers).
	py32.FLASH.SR.Set(py32.Flash_SR_EOP | py32.Flash_SR_WRPERR)
}

func unlock() {
	py32.FLASH.KEYR.Set(key1)
	py32.FLASH.KEYR.Set(key2)
}

func lock() { py32.FLASH.CR.Set(py32.Flash_CR_LOCK) }

//go:section .ramfunc
func ErasePage(addr uint32) bool {

	if !inMainFlash(addr) {
		return false
	}
	page := AlignDownToPage(addr)

	waitNotBusy()
	unlock()
	clearFlags()

	// PER=1, EOPIE=1
	py32.FLASH.CR.Set(py32.Flash_CR_PER | py32.Flash_CR_EOPIE)

	// Trigger page erase by writing any 32-bit word into the page
	(*volatile.Register32)(unsafe.Pointer(uintptr(page))).Set(0xFFFFFFFF)

	waitNotBusy()

	// Check WRPERR
	ok := (py32.FLASH.SR.Get() & py32.Flash_SR_WRPERR) == 0

	clearFlags()
	py32.FLASH.CR.Set(0) // clear PER/EOPIE/others we set
	lock()

	return ok
	// return true
}

//go:section .ramfunc
func ProgramPage(addr uint32, data *[WordsPerPage]uint32) bool {
	if !inMainFlash(addr) {
		return false
	}
	page := AlignDownToPage(addr)

	waitNotBusy()
	unlock()
	clearFlags()

	// PG=1, EOPIE=1
	py32.FLASH.CR.Set(py32.Flash_CR_PG | py32.Flash_CR_EOPIE)

	base := uintptr(page)

	// Write first 31 words
	for i := 0; i < WordsPerPage-1; i++ {
		*(*uint32)(unsafe.Pointer(base + uintptr(i*4))) = data[i]
	}

	// PGSTRT then write last word to start programming
	py32.FLASH.CR.Set(py32.Flash_CR_PG | py32.Flash_CR_EOPIE | py32.Flash_CR_PGSTRT)
	*(*uint32)(unsafe.Pointer(base + uintptr((WordsPerPage-1)*4))) = data[WordsPerPage-1]

	runtime.Gosched()
	waitNotBusy()

	ok := (py32.FLASH.SR.Get() & py32.Flash_SR_WRPERR) == 0

	clearFlags()
	py32.FLASH.CR.Set(0) // clear PG/EOPIE/PGSTRT
	lock()
	return ok
}

func ReadWords(addr uint32, n int) []uint32 {
	return unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(addr))), n)
}
