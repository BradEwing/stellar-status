package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BradEwing/stellar-status/internal/astro"
	"github.com/BradEwing/stellar-status/internal/launches"
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
	rootCmd.Flags().BoolP("cache", "c", true, "enable file-based cache for launch API responses")
	rootCmd.Flags().StringP("site", "s", "VBG", "launch site abbreviation ("+strings.Join(launches.ValidSiteAbbrevs(), ", ")+")")
	rootCmd.Flags().BoolP("moon-ascii", "m", false, "show 5x3 ASCII moon art (multi-line output)")
	rootCmd.Flags().BoolP("solar", "o", false, "show sun altitude")
	rootCmd.Flags().BoolP("twilight", "t", false, "show sunrise/sunset times")
	rootCmd.Flags().BoolP("planets", "p", false, "show visible planets")
	rootCmd.Flags().Float64("lat", 34.7420, "observer latitude (degrees, positive north)")
	rootCmd.Flags().Float64("lon", -120.5724, "observer longitude (degrees, positive east)")

	viper.BindPFlags(rootCmd.Flags())
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	useCache := viper.GetBool("cache")
	siteAbbrev := viper.GetString("site")
	moonASCII := viper.GetBool("moon-ascii")
	showSolar := viper.GetBool("solar")
	showTwilight := viper.GetBool("twilight")
	showPlanets := viper.GetBool("planets")
	loc := astro.Location{
		Latitude:  viper.GetFloat64("lat"),
		Longitude: viper.GetFloat64("lon"),
	}

	site, err := launches.LookupSite(siteAbbrev)
	if err != nil {
		return fmt.Errorf("invalid site: %w\nValid sites: %s", err, strings.Join(launches.ValidSiteAbbrevs(), ", "))
	}

	p := moon.Current()

	var segments []string

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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

	if showSolar {
		segments = append(segments, solar.Current(loc).FormatStatus())
	}

	if showTwilight {
		segments = append(segments, twilight.Today(loc).FormatStatus(time.Now()))
	}

	if showPlanets {
		segments = append(segments, planets.Current(loc).FormatStatus())
	}

	printOutput(p, moonASCII, segments)
	return nil
}

func printOutput(p moon.Phase, ascii bool, segments []string) {
	status := strings.Join(segments, " | ")
	if ascii {
		art := p.ASCII()
		moonInfo := fmt.Sprintf("%s %.0f%%", p.Name, p.Illumination*100)
		fmt.Fprintln(os.Stdout, art[0])
		fmt.Fprintf(os.Stdout, "%s %s | %s\n", art[1], moonInfo, status)
		fmt.Fprintln(os.Stdout, art[2])
	} else {
		moonStatus := fmt.Sprintf("%s %s %.0f%%", p.Emoji, p.Name, p.Illumination*100)
		fmt.Fprintf(os.Stdout, "%s | %s\n", moonStatus, status)
	}
}
