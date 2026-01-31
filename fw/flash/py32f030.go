//go:build py32f030

package flash

const (
	MainFlashBase = 0x08000000
	MainFlashSize = 64 * 1024

	PageSize     = 128
	WordsPerPage = PageSize / 4
)
