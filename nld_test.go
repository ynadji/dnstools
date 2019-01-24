package nld

import (
	//	"fmt"
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
		testCase{"google.com", "www.google.com", "google.com"},
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
	}

	for _, tc := range testCases {
		nld, err := extractNLD(tc.domain, tc.n, true)
		if nld != tc.nld {
			t.Fatalf("Found %v, expected %v). Error: %+v", nld, tc.nld, err)
		}
	}
}
