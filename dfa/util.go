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
	"fmt"
	"strconv"
	"strings"

	"github.com/opennota/re2dfa/nfa"
)

func intsToStrings(a []int) []string {
	s := make([]string, 0, len(a))
	for _, i := range a {
		s = append(s, strconv.Itoa(i))
	}
	return s
}

func makeLabel(a []int) string {
	return strings.Join(intsToStrings(a), ",")
}

func rangesToBoolExpr(rr []rune, atEnd bool) string {
	s := make([]string, 0, len(rr))
	for i := 0; i < len(rr); i += 2 {
		if rr[i] < 0 {
			switch rr[i] {
			case nfa.RuneBeginText:
				s = append(s, "i == 0")
			case nfa.RuneEndText:
				s = append(s, "i == len(s)")
			case nfa.RuneBeginLine:
				s = append(s, `i == 0 || r == '\n'`)
			case nfa.RuneEndLine:
				s = append(s, `i == len(s) || s[i] == '\n'`)
			case nfa.RuneWordBoundary:
				s = append(s, "(i >= rlen && isWordChar(s[i-rlen])) != (rlen > 0 && i < len(s) && isWordChar(s[i]))")
			case nfa.RuneNoWordBoundary:
				s = append(s, "(i >= rlen && isWordChar(s[i-rlen])) == (rlen > 0 && i < len(s) && isWordChar(s[i]))")
			}
		} else if rr[i] == rr[i+1] {
			s = append(s, fmt.Sprintf("r == %d", rr[i]))
		} else if rr[i] == 0 {
			s = append(s, fmt.Sprintf("r <= %d", rr[i+1]))
		} else if rr[i+1] == nfa.RuneLast {
			s = append(s, fmt.Sprintf("r >= %d", rr[i]))
		} else {
			s = append(s, fmt.Sprintf("r >= %d && r <= %d", rr[i], rr[i+1]))
		}
	}
	return strings.Join(s, "||")
}
