package whatsapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "nomor diawali 0 dikonversi ke 62",
			input: "08123456789",
			want:  "628123456789@s.whatsapp.net",
		},
		{
			name:  "nomor sudah diawali 62 tidak berubah",
			input: "628123456789",
			want:  "628123456789@s.whatsapp.net",
		},
		{
			name:  "nomor dengan tanda plus dikonversi ke 62",
			input: "+628123456789",
			want:  "628123456789@s.whatsapp.net",
		},
		{
			name:  "nomor sudah berformat WhatsApp dikembalikan as-is",
			input: "628123456789@s.whatsapp.net",
			want:  "628123456789@s.whatsapp.net",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPhoneNumber(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
