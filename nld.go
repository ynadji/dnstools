package nld

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"strings"
)

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Does not check for matches. Only works on length.x
func removeSuffix(domain, e2ld string) string {
	startOfSuffix := len(domain) - len(e2ld) - 1
	if startOfSuffix < 0 {
		return domain
	}
	restAndDot := domain[0:startOfSuffix]
	return strings.Trim(restAndDot, ".")
}

func extractNLD(domain string, n int, public bool) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("n must be greated than 0")
	}
	suffix, _ := publicsuffix.PublicSuffix(domain)
	fmt.Printf("domain: %v, n: %v, public: %v\n", domain, n, public)
	if n == 1 {
		return suffix, nil
	} else {
		e2ld, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			return "", err
		}
		rest := removeSuffix(domain, e2ld)
		restLevels := reverse(strings.Split(rest, "."))
		extra := []string{}
		for i := 0; i < len(restLevels) && i < n-2; i++ {
			extra = append(extra, restLevels[i])
		}
		extra = reverse(extra)
		if len(extra) > 0 {
			e2ld = fmt.Sprintf("%v.%v", strings.Join(extra, "."), e2ld)
		}
		return e2ld, nil
	}

	return "", fmt.Errorf("Unknown error :(")
}
