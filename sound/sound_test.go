package sound

import (
	"encoding/binary"
	"testing"
	"time"
)

// TestBellWAV guards the embedded asset: a malformed or wrong-length file
// would only surface when a timer ends.
func TestBellWAV(t *testing.T) {
	if string(bell[0:4]) != "RIFF" || string(bell[8:12]) != "WAVE" {
		t.Fatalf("bell.wav: not a RIFF/WAVE file")
	}

	var rate, channels, bits, dataLen int

	for i := 12; i+8 <= len(bell); {
		id := string(bell[i : i+4])
		size := int(binary.LittleEndian.Uint32(bell[i+4 : i+8]))
		body := i + 8

		switch id {
		case "fmt ":
			channels = int(binary.LittleEndian.Uint16(bell[body+2 : body+4]))
			rate = int(binary.LittleEndian.Uint32(bell[body+4 : body+8]))
			bits = int(binary.LittleEndian.Uint16(bell[body+14 : body+16]))
		case "data":
			dataLen = size
		}

		i = body + size + (size & 1) // chunks are word-aligned
	}

	if rate == 0 || channels == 0 || bits == 0 || dataLen == 0 {
		t.Fatalf("bell.wav: malformed (rate=%d channels=%d bits=%d data=%d)", rate, channels, bits, dataLen)
	}

	got := time.Duration(float64(dataLen) / float64(rate*channels*bits/8) * float64(time.Second))
	if got != BellDuration {
		t.Fatalf("bell.wav duration = %v, want %v", got, BellDuration)
	}
}
