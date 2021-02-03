package downloader

import "fmt"

func calcLength(L int) string {
	if L < 1024 {
		return fmt.Sprintf("%d Byte", L)
	}
	kb := float32(L) / 1024
	if kb < 1024 {
		return fmt.Sprintf("%g KB", kb)
	}
	mb := kb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%g MB", mb)
	}
	gb := mb / 1024
	if gb < 1024 {
		return fmt.Sprintf("%g GB", gb)
	}
	return fmt.Sprintf("%g PB", gb/1024)
}
