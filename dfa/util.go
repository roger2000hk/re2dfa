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

func inRange(r rune, ranges []rune) bool {
	for i := 0; i < len(ranges); i += 2 {
		if r >= ranges[i] && r <= ranges[i+1] {
			return true
		}
	}
	return false
}

func addToRange(ranges []rune, r rune) []rune {
	if len(ranges) == 0 {
		return []rune{r, r}
	}

	i := 0
	for i < len(ranges) && r >= ranges[i] {
		if r == ranges[i+1]+1 {
			ranges[i+1]++
			return ranges
		}
		if r <= ranges[i+1] {
			return ranges
		}
		i += 2
	}
	if i >= len(ranges) {
		ranges = append(ranges, []rune{r, r}...)
	} else if r == ranges[i]-1 {
		ranges[i]--
	} else {
		ranges = append([]rune{r, r}, ranges...)
	}

	return ranges
}

func copyOf(r []rune) []rune {
	rr := make([]rune, len(r))
	copy(rr, r)
	return rr
}

func foldRanges(a, b []rune) []rune {
	if len(a) == 0 {
		return copyOf(b)
	}
	if len(b) == 0 {
		return copyOf(a)
	}

	c := make([]rune, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		if i < len(a) {
			if j < len(b) {
				if a[i] < b[j] {
					c = append(c, a[i:i+2]...)
					i += 2
				} else {
					c = append(c, b[j:j+2]...)
					j += 2
				}
				continue
			}
			c = append(c, a[i:i+2]...)
			i += 2
			continue
		}
		c = append(c, b[j:j+2]...)
		j += 2
	}

	d := make([]rune, 0, len(c))
	d = append(d, c[:2]...)
	for i := 2; i < len(c); i += 2 {
		if c[i] <= d[len(d)-1] {
			if c[i+1] > d[len(d)-1] {
				d[len(d)-1] = c[i+1]
			}
		} else {
			d = append(d, c[i:i+2]...)
		}
	}

	return d
}

func rangesToBoolExpr(rr []rune, atEnd bool) string {
	s := make([]string, 0, len(rr))
	for i := 0; i < len(rr); i += 2 {
		if rr[i] < 0 {
			switch rr[i] {
			case -100: // BeginText
				s = append(s, "i == rlen")
			case -200: // EndText
				s = append(s, "i == len(s)")
			case -300: // BeginLine
				s = append(s, `i == rlen || r == '\n'`)
			case -400: // EndLine
				s = append(s, `i == len(s) || s[i] == '\n'`)
			case -500: // WordBoundary
				s = append(s, "(i >= rlen && isWordChar(s[i-rlen])) != (rlen > 0 && i < len(s) && isWordChar(s[i]))")
			case -600: // NoWordBoundary
				s = append(s, "(i >= rlen && isWordChar(s[i-rlen])) == (rlen > 0 && i < len(s) && isWordChar(s[i]))")
			}
		} else if rr[i] == rr[i+1] {
			s = append(s, fmt.Sprintf("r == %d", rr[i]))
		} else if rr[i] == 0 {
			s = append(s, fmt.Sprintf("r <= %d", rr[i+1]))
		} else if rr[i+1] == '\U0010ffff' {
			s = append(s, fmt.Sprintf("r >= %d", rr[i]))
		} else {
			s = append(s, fmt.Sprintf("r >= %d && r <= %d", rr[i], rr[i+1]))
		}
	}
	return strings.Join(s, "||")
}
