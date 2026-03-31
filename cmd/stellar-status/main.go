package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bradewing/stellar-status/internal/launches"
	"github.com/bradewing/stellar-status/internal/moon"
)

func main() {
	p := moon.Current()
	moonStatus := fmt.Sprintf("%s %s %.0f%%", p.Emoji, p.Name, p.Illumination*100)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tracker, err := launches.NewTracker()
	if err != nil {
		fmt.Printf("%s | 🚀[VBG] unavailable\n", moonStatus)
		os.Exit(0)
	}

	result, err := tracker.NextLaunch(ctx)
	if err != nil {
		fmt.Printf("%s | 🚀[VBG] unavailable\n", moonStatus)
		os.Exit(0)
	}

	if result == nil {
		fmt.Printf("%s | 🚀[VBG] no upcoming launches\n", moonStatus)
		os.Exit(0)
	}

	fmt.Printf("%s | %s\n", moonStatus, result.FormatStatus())
}
