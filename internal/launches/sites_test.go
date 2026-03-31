package launches

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupSite_Valid(t *testing.T) {
	site, err := LookupSite("VBG")
	require.NoError(t, err)
	assert.Equal(t, 11, site.LocationID)
	assert.Equal(t, "VBG", site.Abbrev)
	assert.Equal(t, "Vandenberg SFB", site.Name)
}

func TestLookupSite_CaseInsensitive(t *testing.T) {
	site, err := LookupSite("ksc")
	require.NoError(t, err)
	assert.Equal(t, 27, site.LocationID)
}

func TestLookupSite_Invalid(t *testing.T) {
	_, err := LookupSite("INVALID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown site")
}

func TestValidSiteAbbrevs(t *testing.T) {
	abbrevs := ValidSiteAbbrevs()
	assert.Contains(t, abbrevs, "VBG")
	assert.Contains(t, abbrevs, "KSC")
	assert.Contains(t, abbrevs, "CCSFS")
	assert.Equal(t, len(sites), len(abbrevs))

	// Verify sorted.
	for i := 1; i < len(abbrevs); i++ {
		assert.True(t, abbrevs[i-1] < abbrevs[i], "expected sorted order")
	}
}
