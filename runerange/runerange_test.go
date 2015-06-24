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

package runerange

import (
	"reflect"
	"testing"
)

func TestIn(t *testing.T) {
	type testCase struct {
		a    []rune
		r    rune
		want bool
	}
	testCases := []testCase{
		{nil, 'a', false},
		{[]rune{'a', 'a'}, 'a', true},
		{[]rune{'a', 'z'}, 'o', true},
		{[]rune{'a', 'z'}, '0', false},
		{[]rune{'0', '9', 'a', 'z'}, '1', true},
		{[]rune{'0', '9', 'a', 'z'}, 'b', true},
		{[]rune{'0', '9', 'a', 'z'}, '@', false},
	}
	for _, tc := range testCases {
		got := In(tc.a, tc.r)
		if got != tc.want {
			t.Errorf("In(%q, '%c') = %v, want %v", string(tc.a), tc.r, got, tc.want)
		}
	}
}

func TestAdd(t *testing.T) {
	type testCase struct {
		a    []rune
		r    rune
		want []rune
	}
	testCases := []testCase{
		{[]rune{}, 'a', []rune{'a', 'a'}},
		{[]rune{'a', 'a'}, 'a', []rune{'a', 'a'}},
		{[]rune{'a', 'a'}, 'b', []rune{'a', 'b'}},
		{[]rune{'a', 'a'}, 'c', []rune{'a', 'a', 'c', 'c'}},
		{[]rune{'c', 'z'}, 'a', []rune{'a', 'a', 'c', 'z'}},
		{[]rune{'b', 'z'}, 'a', []rune{'a', 'z'}},
		{[]rune{'a', 'y'}, 'z', []rune{'a', 'z'}},
		{[]rune{'c', 'z'}, 'q', []rune{'c', 'z'}},
		{[]rune{'a', 'z'}, 'A', []rune{'A', 'A', 'a', 'z'}},
		{[]rune{'a', 'n', 'p', 'z'}, 'o', []rune{'a', 'z'}},
		{[]rune{'a', 'n', 'q', 'z'}, 'o', []rune{'a', 'o', 'q', 'z'}},
		{[]rune{'a', 'n', 'q', 'z'}, 'p', []rune{'a', 'n', 'p', 'z'}},
		{[]rune{'A', 'J', 'a', 'j', 'l', 'r'}, 'L', []rune{'A', 'J', 'L', 'L', 'a', 'j', 'l', 'r'}},
	}
	for _, tc := range testCases {
		got := Add(tc.a, tc.r)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Add(%q, '%c') = %q, want %q", string(tc.a), tc.r, string(got), string(tc.want))
		}
	}
}

func TestSum(t *testing.T) {
	type testCase struct {
		a, b []rune
		want []rune
	}
	testCases := []testCase{
		{[]rune{'a', 'z'}, []rune{}, []rune{'a', 'z'}},
		{[]rune{}, []rune{'0', '9'}, []rune{'0', '9'}},
		{[]rune{'0', '9'}, []rune{'a', 'z'}, []rune{'0', '9', 'a', 'z'}},
		{[]rune{'a', 'z'}, []rune{'0', '9'}, []rune{'0', '9', 'a', 'z'}},
		{[]rune{'a', 'z'}, []rune{'b', 'y'}, []rune{'a', 'z'}},
		{[]rune{'b', 'y'}, []rune{'a', 'z'}, []rune{'a', 'z'}},
		{[]rune{'a', 't'}, []rune{'o', 'z'}, []rune{'a', 'z'}},
		{[]rune{'o', 'z'}, []rune{'a', 't'}, []rune{'a', 'z'}},
		{[]rune{'a', 't'}, []rune{'t', 'z'}, []rune{'a', 'z'}},
		{[]rune{'t', 'z'}, []rune{'a', 't'}, []rune{'a', 'z'}},
		{[]rune{'a', 't'}, []rune{'x', 'z'}, []rune{'a', 't', 'x', 'z'}},
	}
	for _, tc := range testCases {
		got := Sum(tc.a, tc.b)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Sum(%q, %q) = %q, want %q", string(tc.a), string(tc.b), string(got), string(tc.want))
		}
	}
}

func TestFold(t *testing.T) {
	type testCase struct {
		in   []rune
		want []rune
	}
	testCases := []testCase{
		{nil, nil},
		{[]rune{'0', '9'}, []rune{'0', '9'}},
		{[]rune{'a', 'j'}, []rune{'A', 'J', 'a', 'j'}},
		{[]rune{'a', 'j', 'l', 'r'}, []rune{'A', 'J', 'L', 'R', 'a', 'j', 'l', 'r'}},
		{[]rune{'a', 'j', 'l', 'r', 't', 'z'}, []rune{'A', 'J', 'L', 'R', 'T', 'Z', 'a', 'j', 'l', 'r', 't', 'z'}},
		{[]rune{'0', '9', 'a', 'z'}, []rune{'0', '9', 'A', 'Z', 'a', 'z', 'ſ', 'ſ', 'K', 'K'}},
	}
	for _, tc := range testCases {
		got := Fold(tc.in)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Fold(%q) = %q, want %q", string(tc.in), string(got), string(tc.want))
		}
	}
}

func TestSplit(t *testing.T) {
	type testCase struct {
		in   [][]rune
		want []rune
	}
	testCases := []testCase{
		{nil, nil},
		{[][]rune{{}, {}}, nil},
		{[][]rune{{'0', '9'}, {'a', 'z'}}, []rune{'0', '9', 'a', 'z'}},
		{[][]rune{{'a', 'z'}, {'0', '9'}}, []rune{'0', '9', 'a', 'z'}},
		{[][]rune{{'0', '9'}, {'0', '9'}}, []rune{'0', '9'}},
		{[][]rune{{'a', 'z'}, {'b', 'y'}}, []rune{'a', 'a', 'b', 'y', 'z', 'z'}},
		{[][]rune{{'b', 'y'}, {'a', 'z'}}, []rune{'a', 'a', 'b', 'y', 'z', 'z'}},
		{[][]rune{{'a', 'y'}, {'b', 'z'}}, []rune{'a', 'a', 'b', 'y', 'z', 'z'}},
		{[][]rune{{'b', 'z'}, {'a', 'y'}}, []rune{'a', 'a', 'b', 'y', 'z', 'z'}},
		{[][]rune{{'a', 'o'}, {'o', 'z'}}, []rune{'a', 'n', 'o', 'o', 'p', 'z'}},
		{[][]rune{{'o', 'z'}, {'a', 'o'}}, []rune{'a', 'n', 'o', 'o', 'p', 'z'}},
		{[][]rune{{'a', 'z'}, {'n', 'p'}}, []rune{'a', 'm', 'n', 'p', 'q', 'z'}},
		{[][]rune{{'n', 'p'}, {'a', 'z'}}, []rune{'a', 'm', 'n', 'p', 'q', 'z'}},
		{[][]rune{{'a', 'p'}, {'n', 'z'}}, []rune{'a', 'm', 'n', 'p', 'q', 'z'}},
		{[][]rune{{'a', 'c'}, {'d', 'f'}, {'g', 'i'}}, []rune{'a', 'c', 'd', 'f', 'g', 'i'}},
		{[][]rune{{'a', 'd'}, {'d', 'f'}, {'f', 'i'}}, []rune{'a', 'c', 'd', 'd', 'e', 'e', 'f', 'f', 'g', 'i'}},
	}
	for _, tc := range testCases {
		got := Split(tc.in)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Split(%#v) = %q, want %#v", tc.in, string(got), string(tc.want))
		}
	}
}
