// Package sound plays the timer start beep and the end-of-timer bell.
package sound

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gen2brain/beeep"
)

//go:embed bell.wav
var bell []byte

// BellDuration is how long the end-of-timer bell rings.
const BellDuration = 5 * time.Second

// Start plays a short beep when a session starts.
// It blocks for the beep duration, so call it from a goroutine.
func Start() {
	if err := beeep.Beep(880, 200); err != nil {
		log.Printf("start beep failed: %v", err)
	}
}

// Bell plays the end-of-timer bell. It blocks until the bell finishes or
// ctx is cancelled, so call it from a goroutine.
func Bell(ctx context.Context) {
	if err := play(ctx, bell); err != nil {
		log.Printf("bell failed: %v", err)
	}
}

func play(ctx context.Context, wav []byte) error {
	path := filepath.Join(os.TempDir(), "pomo-bell.wav")
	if err := os.WriteFile(path, wav, 0o644); err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "powershell",
			"-NoProfile", "-Command",
			"(New-Object Media.SoundPlayer '"+path+"').PlaySync()",
		)
	case "darwin":
		cmd = exec.CommandContext(ctx, "afplay", path)
	default:
		player, err := exec.LookPath("paplay")
		if err != nil {
			if player, err = exec.LookPath("aplay"); err != nil {
				return err
			}
		}

		cmd = exec.CommandContext(ctx, player, path)
	}

	return cmd.Run()
}
