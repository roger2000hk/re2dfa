package test

import "testing"

type testCase struct {
	in   string
	want int
}

func TestMatchAlternatives(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"xxx", -1},
		{"abc", 3},
		{"abd", -1},
		{"acd", -1},
		{"def", 3},
		{"deg", -1},
		{"dfg", -1},
		{"abcdef", 3},
	}
	for _, tc := range testCases {
		got := matchAlternatives(tc.in)
		if got != tc.want {
			t.Errorf("matchAlternatives(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchStar(t *testing.T) {
	testCases := []testCase{
		{"", 0},
		{"x", 0},
		{"a", 1},
		{"aa", 2},
		{"aab", 2},
		{"aaa", 3},
	}
	for _, tc := range testCases {
		got := matchStar(tc.in)
		if got != tc.want {
			t.Errorf("matchStar(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchQuest(t *testing.T) {
	testCases := []testCase{
		{"", 0},
		{"x", 0},
		{"a", 1},
		{"aa", 1},
		{"aaa", 1},
	}
	for _, tc := range testCases {
		got := matchQuest(tc.in)
		if got != tc.want {
			t.Errorf("matchQuest(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchPlus(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"a", 1},
		{"aa", 2},
		{"aab", 2},
		{"aaa", 3},
	}
	for _, tc := range testCases {
		got := matchPlus(tc.in)
		if got != tc.want {
			t.Errorf("matchPlus(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchLiteral(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"abcde", -1},
		{"a#cdef", -1},
		{"ab#def", -1},
		{"abc#ef", -1},
		{"abcd#f", -1},
		{"abcdef", 6},
		{"abcdefg", 6},
		{"abcdeg", -1},
	}
	for _, tc := range testCases {
		got := matchLiteral(tc.in)
		if got != tc.want {
			t.Errorf("matchLiteral(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchCharClass(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"a", 1},
		{"aa", 1},
		{"z", 1},
		{"d", 1},
		{"#", -1},
		{"A", -1},
	}
	for _, tc := range testCases {
		got := matchCharClass(tc.in)
		if got != tc.want {
			t.Errorf("matchCharClass(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchRepeat1(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"a", 1},
		{"ab", 1},
		{"aa", 2},
		{"aab", 2},
		{"aaa", 3},
		{"aaaa", 3},
		{"aaaaa", 3},
	}
	for _, tc := range testCases {
		got := matchRepeat1(tc.in)
		if got != tc.want {
			t.Errorf("matchRepeat1(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchRepeat2(t *testing.T) {
	testCases := []testCase{
		{"", 0},
		{"x", 0},
		{"a", 1},
		{"ab", 1},
		{"aa", 2},
		{"aab", 2},
		{"aaa", 3},
		{"aaaa", 3},
		{"aaaaa", 3},
	}
	for _, tc := range testCases {
		got := matchRepeat2(tc.in)
		if got != tc.want {
			t.Errorf("matchRepeat2(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchConcat(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"ac", -1},
		{"abc", 3},
		{"acd", -1},
		{"abd", -1},
		{"abbc", 4},
		{"abbbc", 5},
		{"abcd", 3},
	}
	for _, tc := range testCases {
		got := matchConcat(tc.in)
		if got != tc.want {
			t.Errorf("matchConcat(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchStartOfText(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"a", 1},
		{"aa", 1},
		{"ab", 1},
		{"\na", -1},
	}
	for _, tc := range testCases {
		got := matchStartOfText(tc.in)
		if got != tc.want {
			t.Errorf("matchStartOfText(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchStartOfTextEmpty(t *testing.T) {
	testCases := []testCase{
		{"", 0},
		{"x", 0},
	}
	for _, tc := range testCases {
		got := matchStartOfTextEmpty(tc.in)
		if got != tc.want {
			t.Errorf("matchStartOfTextEmpty(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchEndOfText(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"a", 1},
		{"aa", -1},
		{"a\n", -1},
	}
	for _, tc := range testCases {
		got := matchEndOfText(tc.in)
		if got != tc.want {
			t.Errorf("matchEndOfText(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchStartOfLine(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"a", 1},
		{"aa", 1},
		{"ab", 1},
		{"\na", -1},
	}
	for _, tc := range testCases {
		got := matchStartOfLine(tc.in)
		if got != tc.want {
			t.Errorf("matchStartOfLine(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchStartOfLineEmpty(t *testing.T) {
	testCases := []testCase{
		{"", 0},
		{"x", 0},
	}
	for _, tc := range testCases {
		got := matchStartOfLineEmpty(tc.in)
		if got != tc.want {
			t.Errorf("matchStartOfLineEmpty(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestMatchEndOfLine(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"a", 1},
		{"aa", -1},
		{"a\n", 1},
	}
	for _, tc := range testCases {
		got := matchEndOfLine(tc.in)
		if got != tc.want {
			t.Errorf("matchEndOfLine(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestWordBoundary(t *testing.T) {
	testCases := []testCase{
		{"", -1},
		{"x", -1},
		{"a", 1},
		{"aa", -1},
		{"aA", -1},
		{"a0", -1},
		{"a_", -1},
		{"a.", 1},
		{"a\n", 1},
	}
	for _, tc := range testCases {
		got := matchWordBoundary(tc.in)
		if got != tc.want {
			t.Errorf("matchWordBoundary(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestLazy1(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	testCases := []testCase{
		{"", 0},
		{"a", 0},
		{"aa", 0},
	}
	for _, tc := range testCases {
		got := matchLazy1(tc.in)
		if got != tc.want {
			t.Errorf("matchLazy1(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestLazy2(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	testCases := []testCase{
		{"", -1},
		{"a", -1},
		{"b", 1},
		{"aa", -1},
		{"ab", 2},
		{"ac", -1},
	}
	for _, tc := range testCases {
		got := matchLazy2(tc.in)
		if got != tc.want {
			t.Errorf("matchLazy2(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestLazy3(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	testCases := []testCase{
		{"", 0},
		{"a", 0},
		{"b", 0},
		{"aa", 0},
		{"aaa", 0},
	}
	for _, tc := range testCases {
		got := matchLazy3(tc.in)
		if got != tc.want {
			t.Errorf("matchLazy3(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestLazy4(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	testCases := []testCase{
		{"", -1},
		{"a", -1},
		{"ab", 2},
		{"b", 1},
		{"aab", 3},
		{"aaab", 4},
	}
	for _, tc := range testCases {
		got := matchLazy4(tc.in)
		if got != tc.want {
			t.Errorf("matchLazy4(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestLazy5(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	testCases := []testCase{
		{"", -1},
		{"a", 1},
		{"b", -1},
		{"aa", 1},
		{"aaa", 1},
	}
	for _, tc := range testCases {
		got := matchLazy5(tc.in)
		if got != tc.want {
			t.Errorf("matchLazy5(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestLazy6(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	testCases := []testCase{
		{"", -1},
		{"a", -1},
		{"ab", 2},
		{"b", -1},
		{"aab", 3},
		{"aaab", 4},
	}
	for _, tc := range testCases {
		got := matchLazy6(tc.in)
		if got != tc.want {
			t.Errorf("matchLazy6(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestLazy7(t *testing.T) {
	type testCase struct {
		in   string
		want int
	}
	testCases := []testCase{
		{"", -1},
		{"a", -1},
		{"b", -1},
		{"c", -1},
		{"ac", 2},
		{"abc", 3},
	}
	for _, tc := range testCases {
		got := matchLazy7(tc.in)
		if got != tc.want {
			t.Errorf("matchLazy7(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
