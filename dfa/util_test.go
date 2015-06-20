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
	"reflect"
	"testing"
)

func TestAppendToRange(t *testing.T) {
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
		{[]rune{'c', 'z'}, 'b', []rune{'b', 'z'}},
		{[]rune{'c', 'z'}, 'q', []rune{'c', 'z'}},
	}
	for _, tc := range testCases {
		got := appendToRange(tc.a, tc.r)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("appendToRange(%q, %c) = %q, want %q", string(tc.a), tc.r, string(got), string(tc.want))
		}
	}
}

func TestFoldRanges(t *testing.T) {
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
		got := foldRanges(tc.a, tc.b)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("foldRanges(%q, %q) = %q, want %q", string(tc.a), string(tc.b), string(got), string(tc.want))
		}
	}
}
