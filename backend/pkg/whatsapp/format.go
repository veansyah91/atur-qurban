package whatsapp

import (
	"regexp"
	"strings"
)

// FormatPhoneNumber memformat nomor telepon ke format GOWA: 628xxx@s.whatsapp.net
func FormatPhoneNumber(phone string) string {
	// Hapus semua karakter non-digit
	re := regexp.MustCompile("[^0-9]")
	cleanPhone := re.ReplaceAllString(phone, "")

	// Jika sudah dalam format @s.whatsapp.net, kembalikan as-is
	if strings.HasSuffix(phone, "@s.whatsapp.net") {
		return phone
	}

	// Jika dimulai dengan 0, ganti dengan 62 (kode Indonesia)
	if strings.HasPrefix(cleanPhone, "0") {
		cleanPhone = "62" + cleanPhone[1:]
	}

	// Jika tidak dimulai dengan 62, tambahkan
	if !strings.HasPrefix(cleanPhone, "62") {
		cleanPhone = "62" + cleanPhone
	}

	return cleanPhone + "@s.whatsapp.net"
}
