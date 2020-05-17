package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

type word struct {
	word  string
	count int
}

type top []word

func Top10(s string) []string {
	if s == "" {
		return nil
	}

	fields := strings.Fields(s)
	r := regexp.MustCompile(`(^[\D]+$)`) //Exclude non-words: digits, symbols, etc.

	counter := make(map[string]int)
	for _, word := range fields {
		if !r.MatchString(word) {
			continue
		}

		if v, ok := counter[word]; ok {
			counter[word] += v
		} else {
			counter[word] = 1
		}
	}

	var words top = make([]word, len(counter))

	i := 0
	for k, v := range counter {
		words[i] = word{k, v}
		i++
	}

	sort.Sort(words)

	topTen := make([]string, 0, 10)
	for k, v := range words {
		if k > 9 { //nolint:gomnd
			break
		}
		topTen = append(topTen, v.word)
	}

	return topTen
}

func (a top) Len() int           { return len(a) }
func (a top) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a top) Less(i, j int) bool { return a[i].count > a[j].count }
