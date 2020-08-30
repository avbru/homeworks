package hw10_program_optimization

import (
	"archive/zip"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

var data = []string{
	`{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`,
	`{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}`,
	`{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}`,
	`{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}`,
	`{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`,
}

func BenchmarkGetDomainStat(b *testing.B) {
	funcs := []struct {
		name string
		fun  func(io.Reader, string) (DomainStat, error)
	}{
		{"inline", GetDomainStat},
		{"brute force", GetDomainStatBruteForce},
		{"string", GetDomainStatString},
		{"concurrent", GetDomainStatConcurrent},
		{"original", GetDomainStatOriginal},
	}

	for _, f := range funcs {
		b.Run(fmt.Sprintf("%s", f.name), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r, err := zip.OpenReader("testdata/users.dat.zip")
				require.NoError(b, err)
				require.Equal(b, 1, len(r.File))
				datafile, err := r.File[0].Open()
				require.NoError(b, err)
				b.StartTimer()
				_, _ = f.fun(datafile, "com")
				b.StopTimer()
				r.Close()
			}
		})
	}
}

func BenchmarkGetEMail(b *testing.B) {
	funcs := []struct {
		name string
		fun  func(data []byte) (email string)
	}{
		{"json", getEmailJSON},
		{"jsoniter", getEmailJSONiter},
		{"regexp", getEmailRegExp},
		{"brute force", getEmailBruteForce},
	}
	for _, f := range funcs {
		b.Run(fmt.Sprintf("%s", f.name), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				for _, line := range data {
					_ = f.fun([]byte(line))
				}
			}
		})
	}
}

func BenchmarkGetEMailString(b *testing.B) {
	funcs := []struct {
		name string
		fun  func(data string) (email string)
	}{
		{"brute force strings", getEmailBruteForceString},
	}
	for _, f := range funcs {
		b.Run(fmt.Sprintf("%s", f.name), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				for _, line := range data {
					_ = f.fun(line)
				}
			}
		})
	}
}
