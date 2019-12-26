package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/influxdata/influxdb/kit/cli"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Does not check for matches. Only works on length. If the e2ld suffix is longer
// than the original domain or they are identical, returns "".
func removeSuffix(domain, e2ld string) string {
	if domain == e2ld {
		return ""
	}
	startOfSuffix := len(domain) - len(e2ld) - 1
	if startOfSuffix < 0 {
		return ""
	}
	restAndDot := domain[0:startOfSuffix]
	return strings.Trim(restAndDot, ". ")
}

// HasListedSuffix returns true if the domain has a TLD that appears on the
// public suffix list and false otherwise. Converts to ASCII to ensure suffix
// check succeeds but does no other normalization.
func HasListedSuffix(domain string) bool {
	domain, err := idna.ToASCII(domain)
	if err != nil {
		return false
	}
	ps, icann := publicsuffix.PublicSuffix(domain)
	// Only ICANN-managed domains can have a single label and
	// privately-managed domains must have multiple labels. If there is no
	// known suffix, `PublicSuffix` just returns the last label to `ps`
	// (e.g., single label). If it isn't managed by ICANN and does not
	// contain a '.', it must not be present on the list.
	return icann || (strings.IndexByte(ps, '.') >= 0)
}

func ExtractNLD(domain string, n int, public bool, onlyIcann bool) (string, error) {
	// Trim trailing dot
	domain = strings.TrimRight(domain, ".")
	if n <= 0 {
		return "", fmt.Errorf("n must be greated than 0")
	}
	suffix, _ := publicsuffix.PublicSuffix(domain)
	if onlyIcann && !HasListedSuffix(domain) {
		return "", fmt.Errorf("'%v' has unknown TLD '%v'", domain, suffix)
	}
	if n == 1 {
		return suffix, nil
	} else {
		e2ld, err := publicsuffix.EffectiveTLDPlusOne(domain)
		if err != nil {
			return "", err
		}
		rest := removeSuffix(domain, e2ld)
		var restLevels []string
		if rest != "" {
			restLevels = reverse(strings.Split(rest, "."))
		}
		extra := []string{}
		for i, ld := range restLevels {
			if i == n-2 {
				break
			}
			extra = append(extra, ld)
		}
		extra = reverse(extra)
		if len(extra) > 0 {
			e2ld = fmt.Sprintf("%v.%v", strings.Join(extra, "."), e2ld)
		}
		return e2ld, nil
	}

	return "", fmt.Errorf("Unknown error :(")
}

var flags struct {
	n                  int
	ignoreErrors       bool
	usePrivateSuffixes bool
	ignoreNonIcann     bool
}

func run() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		eld, err := ExtractNLD(domain, flags.n, !flags.usePrivateSuffixes, flags.ignoreNonIcann)
		if err != nil && !flags.ignoreErrors {
			fmt.Printf("%v,%v\n", eld, err)
		} else {
			if eld != "" {
				fmt.Println(eld)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading standard input: %+v\n", err)
	}
	return nil
}

func main() {
	cmd := cli.NewCommand(&cli.Program{
		Run:  run,
		Name: "nld",
		Opts: []cli.Opt{
			{
				DestP:   &flags.n,
				Flag:    "n",
				Default: 2,
				Desc:    "Domain level to extract",
			},
			{
				DestP:   &flags.ignoreErrors,
				Flag:    "ignoreErrors",
				Default: false,
				Desc:    "Silently ignore errors",
			},
			{
				DestP:   &flags.usePrivateSuffixes,
				Flag:    "usePrivateSuffixes",
				Default: false,
				Desc:    "Treat known private suffixes as TLDs",
			},
			{
				DestP:   &flags.ignoreNonIcann,
				Flag:    "ignoreNonIcann",
				Default: false,
				Desc:    "Ignore domains with TLDs not known by the PSL",
			},
		},
	})

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
