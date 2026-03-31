package launches

import (
	"fmt"
	"sort"
	"strings"
)

// Site maps a human-friendly abbreviation to a Launch Library 2 location ID.
type Site struct {
	Abbrev     string
	LocationID int
	Name       string
}

var sites = map[string]Site{
	"VBG":        {Abbrev: "VBG", LocationID: 11, Name: "Vandenberg SFB"},
	"CCSFS":      {Abbrev: "CCSFS", LocationID: 12, Name: "Cape Canaveral SFS"},
	"KSC":        {Abbrev: "KSC", LocationID: 27, Name: "Kennedy Space Center"},
	"CSG":        {Abbrev: "CSG", LocationID: 13, Name: "Guiana Space Centre"},
	"WFF":        {Abbrev: "WFF", LocationID: 21, Name: "Wallops Flight Facility"},
	"STARBASE":   {Abbrev: "STARBASE", LocationID: 143, Name: "SpaceX Starbase"},
	"BAIKONUR":   {Abbrev: "BAIKONUR", LocationID: 15, Name: "Baikonur Cosmodrome"},
	"PLESETSK":   {Abbrev: "PLESETSK", LocationID: 6, Name: "Plesetsk Cosmodrome"},
	"JSLC":       {Abbrev: "JSLC", LocationID: 17, Name: "Jiuquan Satellite Launch Center"},
	"XSLC":       {Abbrev: "XSLC", LocationID: 16, Name: "Xichang Satellite Launch Center"},
	"SDSC":       {Abbrev: "SDSC", LocationID: 14, Name: "Satish Dhawan Space Centre"},
	"TANEGASHIMA": {Abbrev: "TANEGASHIMA", LocationID: 26, Name: "Tanegashima Space Center"},
	"LC1":        {Abbrev: "LC1", LocationID: 10, Name: "Rocket Lab LC-1"},
	"CORNRANCH":  {Abbrev: "CORNRANCH", LocationID: 29, Name: "Blue Origin Corn Ranch"},
	"PSCA":       {Abbrev: "PSCA", LocationID: 25, Name: "Pacific Spaceport Complex Alaska"},
}

// LookupSite returns the site for the given abbreviation (case-insensitive).
func LookupSite(abbrev string) (Site, error) {
	s, ok := sites[strings.ToUpper(abbrev)]
	if !ok {
		return Site{}, fmt.Errorf("unknown site %q", abbrev)
	}
	return s, nil
}

// ValidSiteAbbrevs returns a sorted list of valid site abbreviations.
func ValidSiteAbbrevs() []string {
	abbrevs := make([]string, 0, len(sites))
	for k := range sites {
		abbrevs = append(abbrevs, k)
	}
	sort.Strings(abbrevs)
	return abbrevs
}
