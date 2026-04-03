package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BradEwing/stellar-status/internal/apod"
	"github.com/BradEwing/stellar-status/internal/astro"
	"github.com/BradEwing/stellar-status/internal/aurora"
	"github.com/BradEwing/stellar-status/internal/deepsky"
	"github.com/BradEwing/stellar-status/internal/launches"
	"github.com/BradEwing/stellar-status/internal/meteors"
	"github.com/BradEwing/stellar-status/internal/moon"
	"github.com/BradEwing/stellar-status/internal/planets"
	"github.com/BradEwing/stellar-status/internal/solar"
	"github.com/BradEwing/stellar-status/internal/twilight"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "stellar-status",
	Short: "Moon phase and rocket launch status for your terminal",
	RunE:  run,
}

func init() {
	rootCmd.Flags().BoolP("no-cache", "n", false, "disable file-based cache for launch API responses")
	rootCmd.Flags().StringP("site", "s", "VBG", "launch site abbreviation ("+strings.Join(launches.ValidSiteAbbrevs(), ", ")+")")
	rootCmd.Flags().BoolP("moon-ascii", "m", false, "show 5x3 ASCII moon art (multi-line output)")
	rootCmd.Flags().Bool("no-moon", false, "disable moon phase display")
	rootCmd.Flags().Bool("no-launch", false, "disable launch tracking")
	rootCmd.Flags().BoolP("solar", "o", false, "show sun altitude")
	rootCmd.Flags().BoolP("twilight", "t", false, "show sunrise/sunset times")
	rootCmd.Flags().BoolP("planets", "p", false, "show visible planets")
	rootCmd.Flags().BoolP("meteors", "e", false, "show meteor shower info")
	rootCmd.Flags().BoolP("deepsky", "d", false, "show best visible deep sky object")
	rootCmd.Flags().BoolP("aurora", "a", false, "show aurora/Kp index")
	rootCmd.Flags().Bool("apod", false, "show NASA Astronomy Picture of the Day title")
	rootCmd.Flags().String("nasa-key", "", "NASA API key (or set NASA_API_KEY env var)")
	rootCmd.Flags().Float64("lat", 34.7420, "observer latitude (degrees, positive north)")
	rootCmd.Flags().Float64("lon", -120.5724, "observer longitude (degrees, positive east)")

	viper.BindPFlags(rootCmd.Flags())
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	useCache := !viper.GetBool("no-cache")
	siteAbbrev := viper.GetString("site")
	moonASCII := viper.GetBool("moon-ascii")
	showMoon := !viper.GetBool("no-moon")
	showLaunch := !viper.GetBool("no-launch")
	showSolar := viper.GetBool("solar")
	showTwilight := viper.GetBool("twilight")
	showPlanets := viper.GetBool("planets")
	showMeteors := viper.GetBool("meteors")
	showDeepSky := viper.GetBool("deepsky")
	showAurora := viper.GetBool("aurora")
	loc := astro.Location{
		Latitude:  viper.GetFloat64("lat"),
		Longitude: viper.GetFloat64("lon"),
	}

	var segments []string

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if showLaunch {
		site, err := launches.LookupSite(siteAbbrev)
		if err != nil {
			return fmt.Errorf("invalid site: %w\nValid sites: %s", err, strings.Join(launches.ValidSiteAbbrevs(), ", "))
		}

		tracker, err := launches.NewTracker(site.Abbrev, site.LocationID, useCache)
		if err != nil {
			segments = append(segments, fmt.Sprintf("🚀[%s] unavailable", site.Abbrev))
		} else {
			result, err := tracker.NextLaunch(ctx)
			if err != nil {
				segments = append(segments, fmt.Sprintf("🚀[%s] unavailable", site.Abbrev))
			} else if result == nil {
				segments = append(segments, fmt.Sprintf("🚀[%s] no upcoming launches", site.Abbrev))
			} else {
				segments = append(segments, result.FormatStatus())
			}
		}
	}

	if showSolar {
		segments = append(segments, solar.Current(loc).FormatStatus())
	}

	if showTwilight {
		segments = append(segments, twilight.Today(loc).FormatStatus(time.Now()))
	}

	if showPlanets {
		segments = append(segments, planets.Current(loc).FormatStatus())
	}

	if showMeteors {
		if s := meteors.Current(loc).FormatStatus(); s != "" {
			segments = append(segments, s)
		}
	}

	if showDeepSky {
		if s := deepsky.Current(loc).FormatStatus(); s != "" {
			segments = append(segments, s)
		}
	}

	if showAurora {
		auroraStatus, err := aurora.Fetch(ctx, useCache)
		if err == nil && auroraStatus != nil {
			segments = append(segments, auroraStatus.FormatStatus())
		}
	}

	if viper.GetBool("apod") {
		nasaKey := viper.GetString("nasa-key")
		if nasaKey == "" {
			nasaKey = os.Getenv("NASA_API_KEY")
		}
		if nasaKey == "" {
			nasaKey = "DEMO_KEY"
		}
		apodResult, err := apod.Fetch(ctx, nasaKey, useCache)
		if err == nil && apodResult != nil {
			if s := apodResult.FormatStatus(); s != "" {
				segments = append(segments, s)
			}
		}
	}

	var p moon.Phase
	if showMoon {
		p = moon.Current()
	}

	printOutput(p, showMoon, moonASCII, segments)
	return nil
}

func printOutput(p moon.Phase, showMoon, ascii bool, segments []string) {
	status := strings.Join(segments, " | ")
	if !showMoon {
		fmt.Fprintln(os.Stdout, status)
		return
	}
	if ascii {
		art := p.ASCII()
		moonInfo := fmt.Sprintf("%s %.0f%%", p.Name, p.Illumination*100)
		fmt.Fprintln(os.Stdout, art[0])
		if status != "" {
			fmt.Fprintf(os.Stdout, "%s %s | %s\n", art[1], moonInfo, status)
		} else {
			fmt.Fprintf(os.Stdout, "%s %s\n", art[1], moonInfo)
		}
		fmt.Fprintln(os.Stdout, art[2])
	} else {
		moonStatus := fmt.Sprintf("%s %s %.0f%%", p.Emoji, p.Name, p.Illumination*100)
		if status != "" {
			fmt.Fprintf(os.Stdout, "%s | %s\n", moonStatus, status)
		} else {
			fmt.Fprintln(os.Stdout, moonStatus)
		}
	}
}
