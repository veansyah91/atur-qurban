package utils

import (
	"strings"
)

// NormalizePhoneNumber mengkonversi nomor HP dari format 08xxx menjadi 628xxxx
// Contoh: 081234567890 → 6281234567890
func NormalizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)
	
	// Hapus semua karakter non-digit
	cleaned := ""
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			cleaned += string(c)
		}
	}
	phone = cleaned

	// Jika diawali 08, ubah ke 62
	if len(phone) >= 2 && phone[0] == '0' && phone[1] == '8' {
		phone = "62" + phone[2:]
	}

	return phone
}
