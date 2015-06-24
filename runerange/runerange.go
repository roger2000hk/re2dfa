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

// Package runerange provides operations on rune ranges.
//
// A rune range is a slice of pairs of runes where each pair represents all the runes from the first rune of the pair to the second one (inclusive).
// Thus, the range ['0', '9'] represents the runes '0', '1', '2', '3', '4', '5', '6', '7', '8', and '9'. The range ['0', '9', 'a', 'z'] represents the digits and the lowercase latin letters.
// A range containing a single rune 'a' is represented as ['a', 'a'].
// Rune pairs should be ordered by the unicode value of the first rune of the pair and should not intersect. The slices ['9', '0'], ['a', 'z', '0', '9'], ['a', 'o', 'b', 'z'] and ['a', 'b', 'c', 'z'] are not valid ranges.
package runerange

import (
	"sort"
	"unicode"
)

// In returns true if a rune is in the range.
func In(ranges []rune, r rune) bool {
	for i := 0; i < len(ranges); i += 2 {
		if r >= ranges[i] && r <= ranges[i+1] {
			return true
		}
	}
	return false
}

// Contains returns true if the range a contains the range b.
func Contains(a, b []rune) bool {
outer:
	for i := 0; i < len(b); i += 2 {
		for j := 0; j < len(a); j += 2 {
			if b[i] >= a[j] && b[i+1] <= a[j+1] {
				continue outer
			}
		}
		return false
	}
	return true
}

// Add adds a rune to the range (maybe modifying the original range) and returns the new range.
func Add(ranges []rune, r rune) []rune {
	if len(ranges) == 0 {
		return []rune{r, r}
	}

	i := 0
	for i < len(ranges) && r >= ranges[i] {
		if r == ranges[i+1]+1 {
			if i+2 < len(ranges) && r+1 == ranges[i+2] {
				ranges[i+1] = ranges[i+3]
				return append(ranges[:i+2], ranges[i+4:]...)
			}
			ranges[i+1]++
			return ranges
		}
		if r <= ranges[i+1] {
			return ranges
		}
		i += 2
	}
	if i >= len(ranges) {
		ranges = append(ranges, r, r)
	} else if r == ranges[i]-1 {
		ranges[i]--
	} else {
		ranges = append(ranges[:i], append([]rune{r, r}, ranges[i:]...)...)
	}

	return ranges
}

type pairs []rune

func (p pairs) Len() int           { return len(p) / 2 }
func (p pairs) Less(i, j int) bool { return p[i*2] < p[j*2] }
func (p pairs) Swap(i, j int) {
	i *= 2
	j *= 2
	p[i], p[i+1], p[j], p[j+1] = p[j], p[j+1], p[i], p[i+1]
}

// Sum returns a range containing all the runes from the ranges a and b. The a and b ranges are not modified.
func Sum(a, b []rune) []rune {
	if len(a) == 0 {
		return append([]rune(nil), b...)
	}
	if len(b) == 0 {
		return append([]rune(nil), a...)
	}

	c := make([]rune, 0, len(a)+len(b))
	c = append(c, a...)
	c = append(c, b...)
	sort.Sort(pairs(c))

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

// Fold returns a range containing all the runes from the original range and all the runes that can be obtained from them by using unicode case folding. The original range is not modified.
func Fold(ranges []rune) []rune {
	if len(ranges) == 0 {
		return nil
	}

	rr := make([]rune, len(ranges))
	copy(rr, ranges)
	for i := 0; i < len(ranges); i += 2 {
		for r := ranges[i]; r <= ranges[i+1]; r++ {
			r0 := r
			for r := unicode.SimpleFold(r); r != r0; r = unicode.SimpleFold(r) {
				rr = Add(rr, r)
			}
		}
	}
	return rr
}

// Split splits a set of ranges into a set of non-intersecting pairs so that each range in the set is a sum of some of the pairs.
func Split(rs [][]rune) []rune {
	if len(rs) == 0 {
		return nil
	}

	size := 0
	for _, rr := range rs {
		size += len(rr)
	}
	if size == 0 {
		return nil
	}

	queue := make([]rune, 0, size)
	for _, rr := range rs {
		queue = append(queue, rr...)
	}

	result := make([]rune, 0, size)
	var r []rune
outer:
	for len(queue) > 0 {
		r, queue = queue[len(queue)-2:], queue[:len(queue)-2]
		r0 := r[0]
		r1 := r[1]

		for i := 0; i < len(result); i += 2 {
			if r0 == result[i] && r1 == result[i+1] {
				continue outer
			}

			if r0 <= result[i] {
				if r1 >= result[i] {
					if r1 <= result[i+1] {
						if r0 <= result[i]-1 {
							queue = append(queue, r0, result[i]-1)
						}
						if result[i] <= r1 {
							queue = append(queue, result[i], r1)
						}
						if r1+1 <= result[i+1] {
							queue = append(queue, r1+1, result[i+1])
						}
						result = append(result[:i], result[i+2:]...)
					} else {
						if r0 <= result[i]-1 {
							queue = append(queue, r0, result[i]-1)
						}
						if result[i+1]+1 <= r1 {
							queue = append(queue, result[i+1]+1, r1)
						}
					}

					continue outer
				}
			} else if r1 >= result[i+1] {
				if r0 <= result[i+1] {
					if result[i] <= r0-1 {
						queue = append(queue, result[i], r0-1)
					}
					if r0 <= result[i+1] {
						queue = append(queue, r0, result[i+1])
					}
					if result[i+1]+1 <= r1 {
						queue = append(queue, result[i+1]+1, r1)
					}
					result = append(result[:i], result[i+2:]...)

					continue outer
				}
			} else {
				if result[0] <= r0-1 {
					queue = append(queue, result[0], r0-1)
				}
				queue = append(queue, r0, r1)
				if r1+1 <= result[i+1] {
					queue = append(queue, r1+1, result[i+1])
				}
				result = append(result[:i], result[i+2:]...)

				continue outer
			}
		}
		result = append(result, r...)
	}
	sort.Sort(pairs(result))
	return result
}
