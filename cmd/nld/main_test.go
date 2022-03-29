package main

import (
	"testing"
)

func sliceEqual(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, _ := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestReverse(t *testing.T) {
	empty := []string{}
	one := []string{"one"}
	two := []string{"one", "two"}
	three := []string{"one", "two", "three"}

	if len(empty) != 0 || len(reverse(empty)) != 0 {
		t.Fatal("empty adding elements")
	}
	if !sliceEqual(empty, reverse(reverse(empty))) {
		t.Fatal("empty not good")
	}
	if !sliceEqual(one, reverse(reverse(one))) {
		t.Fatal("one not good")
	}
	if !sliceEqual(two, reverse(reverse(two))) {
		t.Fatal("two not good")
	}
	if !sliceEqual(three, reverse(reverse(three))) {
		t.Fatal("three not good")
	}
	if !sliceEqual(reverse(empty), []string{}) {
		t.Fatal("empty not good")
	}
	if !sliceEqual(reverse(one), []string{"one"}) {
		t.Fatal("one not good")
	}
	if !sliceEqual(reverse(two), []string{"two", "one"}) {
		t.Fatal("two not good")
	}
	if !sliceEqual(reverse(three), []string{"three", "two", "one"}) {
		t.Fatal("three not good")
	}
}

func TestRemoveSuffix(t *testing.T) {
	type testCase struct {
		domain string
		e2ld   string
		rest   string
	}

	testCases := []testCase{
		testCase{"www.google.com", "google.com", "www"},
		testCase{"foo.bar.google.co.uk", "google.co.uk", "foo.bar"},
		testCase{"google.com", "www.google.com", ""},
		testCase{"google.com", "google.com", ""},
	}

	for _, tc := range testCases {
		noSuffix := removeSuffix(tc.domain, tc.e2ld)
		if noSuffix != tc.rest {
			t.Fatalf("Failure. %v - %v != %v (was %v)", tc.domain, tc.e2ld, tc.rest, noSuffix)
		}
	}
}

func TestExtractNLD(t *testing.T) {
	type testCase struct {
		domain string
		n      int
		nld    string
	}

	testCases := []testCase{
		testCase{"foo.www.google.com", 1, "com"},
		testCase{"foo.www.google.com", 2, "google.com"},
		testCase{"foo.www.google.com", 3, "www.google.com"},
		testCase{"foo.www.google.com", 4, "foo.www.google.com"},
		testCase{"foo.www.google.com", 5, "foo.www.google.com"},
		testCase{"foo.www.google.com", 100, "foo.www.google.com"},
		testCase{"google.com", 2, "google.com"},
		testCase{"google.com", 3, "google.com"},
		// ignore trailing dot
		testCase{"foo.www.google.com.", 1, "com"},
		testCase{"foo.www.google.com.", 2, "google.com"},
		testCase{"foo.www.google.com.", 3, "www.google.com"},
	}

	for i, tc := range testCases {
		nld, err := ExtractNLD(tc.domain, tc.n, true, true)
		if nld != tc.nld {
			t.Fatalf("[%d] Found '%v', expected '%v'. Error: %+v", i, nld, tc.nld, err)
		}
	}

	// Ensure ICANN check works
	testCases = []testCase{
		testCase{"www.google.invalid", 1, ""},
		testCase{"www.google.test", 1, ""},
		testCase{"www.google.local", 1, ""},
		testCase{"www.google.localhost", 1, ""},
		testCase{"www.google.example", 1, ""},
		testCase{"djhkusahvuih", 1, ""},
		testCase{"www.google.onion", 1, "onion"},
	}

	for i, tc := range testCases {
		nld, err := ExtractNLD(tc.domain, tc.n, true, true)
		if tc.nld == "" && err == nil {
			t.Fatalf("[%d] Found '%v', expected '%v'. Error: %+v", i, nld, tc.nld, err)
		}
		if tc.nld != "" && tc.nld != nld {
			t.Fatalf("[%d] Found '%v', expected '%v'. Error: %+v", i, nld, tc.nld, err)
		}
	}
}

func TestHasListedSuffix(t *testing.T) {
	var hasListedSuffixTestCases = []struct {
		domain string
		want   bool
	}{
		{"foo.com", true},
		{"test", false},
		{"com", true},
		{"foo.test", false}, // Reserved TLDs see https://tools.ietf.org/html/rfc2606#page-2
		{"foo.example", false},
		{"foo.invalid", false},
		{"foo.localhost", false},
		{"example", false},
		{"invalid", false},
		{"localhost", false},
		{"foo.baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaar", false}, // too long, can never be valid TLD
		{"万岁.中国", true},                        // Unicode
		{"xn--chqu66a.xn--fiqs8s", true},       // Above in punycode
		{"ésta.bien.es", true},                 // Unicode
		{"xn--sta-9la.bien.es", true},          // Above in punycode
		{"ياسين.الجزائر", true},                // Works with reverse directional unicode
		{"xn--mgby9cnc.xn--lgbbat1ad8j", true}, // Above in punycode (parts reversed)
		{"!@#$%^&*.com", true},                 // Does not check for invalid characters
		{"dyndns-at-work.com", true},
	}

	for _, tc := range hasListedSuffixTestCases {
		got := HasListedSuffix(tc.domain)
		if got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.domain, got, tc.want)
		}
	}
}
