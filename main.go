package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	//"github.com/urfave/cli"
	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

var flags struct {
	n                  int
	ignoreErrors       bool
	usePrivateSuffixes bool
	ignoreNonIcann     bool
	exactOnly          bool
}

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
	if n <= 0 {
		return "", fmt.Errorf("n must be greated than 0")
	}
	// Trim trailing dot
	domain = strings.TrimRight(domain, ".")
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

func run(c *cli.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		eld, err := ExtractNLD(domain, c.Int("numLevels"), !c.Bool("usePrivateSuffixes"), c.Bool("ignoreNonIcann"))
		if c.Bool("exactOnly") && domain != eld {
			err = fmt.Errorf("no exact match between domain/eld '%v' != '%v'. More: %v", domain, eld, err)
			eld = ""
		}
		if err != nil && !c.Bool("ignoreErrors") {
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
	app := &cli.App{
		Name:   "nld",
		Usage:  "nld < domains.txt > nlds.txt",
		Action: run,
	}
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "numLevels",
			Usage:   "Domain level to extract",
			Value:   2,
			Aliases: []string{"n"},
		},
		&cli.BoolFlag{
			Name:  "ignoreErrors",
			Usage: "Silently ignore errors",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "usePrivateSuffixes",
			Usage: "Treat known private suffixes as TLDs",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "ignoreNonIcann",
			Usage: "Ignore domains with TLDs not known by the PSL",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "exactOnly",
			Usage: "Only keep domains if exactly matches label level",
			Value: false,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
