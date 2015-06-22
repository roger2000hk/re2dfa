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

package rr

import "unicode"

func In(ranges []rune, r rune) bool {
	for i := 0; i < len(ranges); i += 2 {
		if r >= ranges[i] && r <= ranges[i+1] {
			return true
		}
	}
	return false
}

func Add(ranges []rune, r rune) []rune {
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

func Sum(a, b []rune) []rune {
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

func FoldCase(ranges []rune) []rune {
	rr := make([]rune, len(ranges))
	copy(rr, ranges)
	for i := 0; i < len(ranges); i += 2 {
		for r := ranges[i]; r <= ranges[i+1]; r++ {
			lc := unicode.ToLower(r)
			if lc != r {
				rr = Add(rr, lc)
			}
			uc := unicode.ToUpper(r)
			if uc != r {
				rr = Add(rr, uc)
			}
		}
	}
	return rr
}
