package benchmarks

import (
	"regexp"
	"strings"
	"testing"
)

var replacer = strings.NewReplacer(
	"\t", "",
	"\n", "",
	" ", "",
)

var rx1 = regexp.MustCompile(replacer.Replace(`
	^(?:
		<[A-Za-z][A-Za-z0-9\-]*(?:\s+[a-zA-Z_:][a-zA-Z0-9:._-]*(?:\s*=\s*(?:[^"'=<>` + "`" + `\x00-\x20]+|'[^']*'|"[^"]*"))?)*\s*\/?> |

		<\/[A-Za-z][A-Za-z0-9\-]*\s*> |

		<!----> |

		<!--(?:-?[^>-])(?:-?[^-])*--> |

		<[?].*?[?]> |

		<![A-Z]+\s+[^>]*> |

		<!\[CDATA\[[\s\S]*?\]\]>
	)`))

var rx1TestStrings = []string{
	`<a href="http://golang.org" title="The Go Programming Language">golang.org</a>`,
	"</blockquote>",
	"<!---->",
	"<!-- This is a comment. --> <!-- Another comment -->",
	`<?xml-stylesheet alternate="yes" href="alt.css" title="Alternative style"?>`,
	"<!DOCTYPE html>",
	"<![CDATA[ This portion of the document is general character data. ]]>",
}

var rx1MatchLengths = []int{
	64, 13, 7, 27, 75, 15, 69,
}

func TestFSM1(t *testing.T) {
	for i, length := range rx1MatchLengths {
		got := match1(rx1TestStrings[i])
		if got != length {
			t.Errorf("match1(%q) = %d, want %d", rx1TestStrings[i], got, length)
		}
	}
}

func TestRegexp1(t *testing.T) {
	for i, length := range rx1MatchLengths {
		loc := rx1.FindStringIndex(rx1TestStrings[i])
		if loc == nil {
			t.Errorf("rx1.FindStringIndex(%q) = nil, want [0, %d]", rx1TestStrings[i], length)
		} else if loc[0] != 0 || loc[1] != length {
			t.Errorf("rx1.FindStringIndex(%q) = [%d, %d], want [0, %d]", rx1TestStrings[i], loc[0], loc[1], length)
		}
	}
}

func BenchmarkFSM1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range rx1TestStrings {
			match1(s)
		}
	}
}

func BenchmarkRegexp1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range rx1TestStrings {
			rx1.MatchString(s)
		}
	}
}
