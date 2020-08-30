// +build !bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomain(t *testing.T) {
	tdata := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	funcs := []struct {
		name string
		fun  func(io.Reader, string) (DomainStat, error)
	}{
		{"inline", GetDomainStat},
		{"brute", GetDomainStatBruteForce},
		{"string", GetDomainStatString},
		{"concurrent", GetDomainStatConcurrent},
		{"original", GetDomainStatOriginal},
	}

	for _, v := range funcs {
		f := v
		t.Run("find 'com'", func(t *testing.T) {
			result, err := f.fun(bytes.NewBufferString(tdata), "com")
			require.NoError(t, err)
			require.Equal(t, DomainStat{
				"browsecat.com": 2,
				"linktype.com":  1,
			}, result)
		})

		t.Run("find 'gov'", func(t *testing.T) {
			result, err := f.fun(bytes.NewBufferString(tdata), "gov")
			require.NoError(t, err)
			require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
		})

		t.Run("find 'unknown'", func(t *testing.T) {
			result, err := f.fun(bytes.NewBufferString(tdata), "unknown")
			require.NoError(t, err)
			require.Equal(t, DomainStat{}, result)
		})
	}
}

func TestGetEmail(t *testing.T) {
	funcs := []struct {
		name string
		fun  func(data []byte) string
	}{
		{"json", getEmailJSON},
		{"regexp", getEmailRegExp},
		{"brute force", getEmailBruteForce},
	}

	emails := []struct {
		data  string
		email string
	}{
		{
			data:  `{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}`,
			email: "mLynch@broWsecat.com",
		},
		{
			data:  `{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}`,
			email: "5Moore@Teklist.net",
		},
		{
			data:  "no-email no-json ",
			email: "",
		},
		{
			data:  `{"Email":"5Moore@Teklist.net"}`,
			email: "5Moore@Teklist.net",
		},
	}

	for _, v := range funcs {
		f := v
		t.Run(f.name, func(t *testing.T) {
			for _, e := range emails {
				require.Equal(t, e.email, f.fun([]byte(e.data)))
			}
		})
	}
	for _, e := range emails {
		require.Equal(t, e.email, getEmailBruteForceString(e.data))
	}
}
