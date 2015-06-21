// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package dfa

import (
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/opennota/re2dfa/nfa"
)

func writeToFile(fn, s string) (err error) {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer func() {
		err2 := f.Close()
		if err2 != nil {
			err = err2
		}
	}()
	_, err = f.WriteString(s)
	return err
}

func uppercaseInitial(s string) string {
	for i, r := range s {
		return string(unicode.ToUpper(r)) + s[i+1:]
	}
	return ""
}

func TestGenerateTests(t *testing.T) {
	type test struct {
		pattern string
		name    string
	}
	tests := []test{
		{"abcdef", "literal"},
		{"[a-z]", "CharClass"},
		{"a*", "star"},
		{"a?", "quest"},
		{"a+", "plus"},
		{"(abc|def)", "alternatives"},
		{"a{1,3}", "repeat1"},
		{"a{0,3}", "repeat2"},
		{"ab+c", "concat"},
		{"^a", "StartOfText"},
		{"^", "StartOfTextEmpty"},
		{"a$", "EndOfText"},
		{"(?m)^a", "StartOfLine"},
		{"(?m)^", "StartOfLineEmpty"},
		{"(?m)a$", "EndOfLine"},
		{`a\b`, "WordBoundary"},
		{`a??`, "lazy1"},
		{`a??b`, "lazy2"},
		{`a*?`, "lazy3"},
		{`a*?b`, "lazy4"},
		{`a+?`, "lazy5"},
		{`a+?b`, "lazy6"},
		{`ab??c`, "lazy7"},
	}
	for _, tst := range tests {
		nfanode, err := nfa.New(tst.pattern)
		if err != nil {
			t.Error(err)
		} else {
			node := NewFromNFA(nfanode)
			source := GoGenerate(node, "test", "match"+uppercaseInitial(tst.name), "string")
			err := writeToFile("test/"+strings.ToLower(tst.name)+".go", source)
			if err != nil {
				t.Error(err)
			}
		}
	}
}
