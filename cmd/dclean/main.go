package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/ynadji/dnstrie/dns"
)

func filter(c *cli.Context, p func(domain string) bool) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := strings.TrimRight(scanner.Text(), ".")
		// xor
		if p(domain) != c.Bool("complement") {
			fmt.Println(domain)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading standard input: %+v\n", err)
	}
	return nil
}

func foreach(c *cli.Context, t func(domain string) (string, error)) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		domain := strings.TrimRight(scanner.Text(), ".")
		if processed, err := t(domain); processed != "" && err == nil {
			fmt.Println(processed)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading standard input: %+v\n", err)
	}
	return nil
}

func valid(c *cli.Context) error {
	return filter(c, dns.Valid)
}

func registerable(c *cli.Context) error {
	return filter(c, dns.IsRegisterableDomain)
}

func ascii(c *cli.Context) error {
	return foreach(c, dns.Normalize)
}

func main() {
	app := &cli.App{
		Name:  "dclean",
		Usage: "dclean [subcommand] < domains.txt > processed-domains.txt",
		Commands: []*cli.Command{
			{
				Name:    "valid",
				Aliases: []string{"v"},
				Usage:   "Filter to domains that are syntactically valid.",
				Action:  valid,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "complement", Aliases: []string{"c"}},
				},
			},
			{
				Name:    "registerable",
				Aliases: []string{"r"},
				Usage:   "Filter to domains that could be registered and used. Excludes TLDs.",
				Action:  registerable,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "complement", Aliases: []string{"c"}},
				},
			},
			{
				Name:    "ascii",
				Aliases: []string{"a"},
				Usage:   "Map input domains to ASCII (IDN->punycode).",
				Action:  ascii,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
