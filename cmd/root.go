package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BradEwing/stellar-status/internal/launches"
	"github.com/BradEwing/stellar-status/internal/moon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "stellar-status",
	Short: "Moon phase and rocket launch status for your terminal",
	RunE:  run,
}

func init() {
	rootCmd.Flags().Bool("cache", false, "enable file-based cache for launch API responses")
	rootCmd.Flags().String("site", "VBG", "launch site abbreviation ("+strings.Join(launches.ValidSiteAbbrevs(), ", ")+")")
	rootCmd.Flags().Bool("moon-ascii", false, "show 3x3 ASCII moon art (multi-line output)")

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

	site, err := launches.LookupSite(siteAbbrev)
	if err != nil {
		return fmt.Errorf("invalid site: %w\nValid sites: %s", err, strings.Join(launches.ValidSiteAbbrevs(), ", "))
	}

	p := moon.Current()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tracker, err := launches.NewTracker(site.Abbrev, site.LocationID, useCache)
	if err != nil {
		printOutput(p, moonASCII, fmt.Sprintf("🚀[%s] unavailable", site.Abbrev))
		return nil
	}

	result, err := tracker.NextLaunch(ctx)
	if err != nil {
		printOutput(p, moonASCII, fmt.Sprintf("🚀[%s] unavailable", site.Abbrev))
		return nil
	}

	if result == nil {
		printOutput(p, moonASCII, fmt.Sprintf("🚀[%s] no upcoming launches", site.Abbrev))
		return nil
	}

	printOutput(p, moonASCII, result.FormatStatus())
	return nil
}

func printOutput(p moon.Phase, ascii bool, launchStatus string) {
	if ascii {
		art := p.ASCII()
		moonInfo := fmt.Sprintf("%s %.0f%%", p.Name, p.Illumination*100)
		fmt.Fprintln(os.Stdout, art[0])
		fmt.Fprintf(os.Stdout, "%s %s | %s\n", art[1], moonInfo, launchStatus)
		fmt.Fprintln(os.Stdout, art[2])
	} else {
		moonStatus := fmt.Sprintf("%s %s %.0f%%", p.Emoji, p.Name, p.Illumination*100)
		fmt.Fprintf(os.Stdout, "%s | %s\n", moonStatus, launchStatus)
	}
}
