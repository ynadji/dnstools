package main

import (
	"fmt"
	//"github.com/influxdata/influxdb/kit/cli"
	"bufio"
	"golang.org/x/net/publicsuffix"
	"os"
	"strings"
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

func ExtractNLD(domain string, n int, public bool) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("n must be greated than 0")
	}
	suffix, _ := publicsuffix.PublicSuffix(domain)
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

//var flags struct {
//	n                  int
//	ignoreErrors       bool
//	usePrivateSuffixes bool
//}

func main() {
	//	fmt.Println("大家好！")
	//	cmd := cli.NewCommand(&cli.Program{
	//		Run:  run,
	//		Name: "myprogram",
	//		Opts: []cli.Opt{
	//			{
	//				DestP:   &flags.n,
	//				Flag:    "n",
	//				Default: 2,
	//				Desc:    "Domain level to extract",
	//			},
	//			{
	//				DestP:   &flags.ignoreErrors,
	//				Flag:    "ignoreErrors",
	//				Default: false,
	//				Desc:    "Silently ignore errors",
	//			},
	//			{
	//				DestP:   &flags.usePrivateSuffixes,
	//				Flag:    "usePrivateSuffixes",
	//				Default: false,
	//				Desc:    "Treat known private suffixes as TLDs",
	//			},
	//		},
	//	})
	//
	//	if err := cmd.Execute(); err != nil {
	//		fmt.Fprintln(os.Stderr, err)
	//		os.Exit(1)
	//	}
	//
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := scanner.Text()
		// Should be domain, flags.n, !flags.usePrivateSuffixes
		eld, err := ExtractNLD(domain, 1, true)
		// err != nil && !flags.ignoreErrors
		if err != nil {
			fmt.Printf("%v,%v\n", eld, err)
		} else {
			fmt.Println(eld)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
		os.Exit(1)
	}
	os.Exit(0)
}
