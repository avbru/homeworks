package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"runtime"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var (
	jsonIter = jsoniter.ConfigCompatibleWithStandardLibrary                            //nolint:unused,varcheck,deadcode
	re       = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}`) //nolint:unused
)

type DomainStat map[string]int

// GetDomainStat inline version is winner by time and memory usage.
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	str := bufio.NewScanner(r)
	for str.Scan() {
		i := bytes.Index(str.Bytes(), []byte(`Email":"`))
		for j := i + 8; j < len(str.Bytes()); j++ {
			if str.Bytes()[j] == '"' {
				if bytes.HasSuffix(str.Bytes()[i+8:j], []byte("."+domain)) {
					subdomain := bytes.ToLower(bytes.SplitN(str.Bytes()[i+8:j], []byte("@"), 2)[1])
					result[string(subdomain)]++
				}
				break
			}
		}
	}
	return result, nil
}

func GetDomainStatBruteForce(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	s := bufio.NewScanner(r)
	for s.Scan() {
		email := getEmailBruteForce(s.Bytes())
		if strings.HasSuffix(email, "."+domain) {
			subdomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
			result[subdomain]++
		}
	}
	return result, nil
}

func GetDomainStatString(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	s := bufio.NewScanner(r)
	for s.Scan() {
		email := getEmailBruteForceString(s.Text())
		if strings.HasSuffix(email, "."+domain) {
			subdomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
			result[subdomain]++
		}
	}
	return result, nil
}

func GetDomainStatConcurrent(r io.Reader, domain string) (DomainStat, error) {
	mu := sync.Mutex{}
	result := make(DomainStat)

	nProc := runtime.NumCPU() + 1/2
	linesChan := make(chan string, nProc)

	var wg sync.WaitGroup
	wg.Add(nProc)
	for i := 0; i < nProc; i++ {
		go func() {
			defer wg.Done()
			for {
				line, ok := <-linesChan
				if ok {
					email := getEmailBruteForceString(line)
					if strings.HasSuffix(email, "."+domain) {
						subdomain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
						mu.Lock()
						result[subdomain]++
						mu.Unlock()
					}
				} else {
					return
				}
			}
		}()
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		linesChan <- s.Text()
	}
	close(linesChan)
	wg.Wait()
	return result, nil
}

func getEmailJSON(data []byte) string { //nolint:unused,deadcode
	var user struct{ Email string }
	if err := json.Unmarshal(data, &user); err != nil {
		return ""
	}
	return user.Email
}

func getEmailJSONiter(data []byte) string { //nolint:unused,deadcode
	var user struct{ Email string }
	if err := jsoniter.Unmarshal(data, &user); err != nil {
		return ""
	}
	return user.Email
}

func getEmailRegExp(data []byte) string { //nolint:unused,deadcode
	return re.FindString(string(data))
}

func getEmailBruteForce(data []byte) string {
	i := bytes.Index(data, []byte(`Email":`))
	for j := i + 8; j < len(data); j++ {
		if data[j] == '"' {
			return string(data[i+8 : j])
		}
	}
	return ""
}

// getEmailBruteForceString is used by concurrent version of Stats.
// chan []byte is a mess...
func getEmailBruteForceString(data string) string {
	i := strings.Index(data, `Email":`)
	if i == -1 {
		return ""
	}
	for j := i + 8; true; j++ {
		if data[j] == '"' {
			return data[i+8 : j]
		}
	}
	return ""
}

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

func GetDomainStatOriginal(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
