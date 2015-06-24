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

// Package dfa provides a way to construct deterministic finite automata from non-deterministic finite automata.
package dfa

import (
	"sort"
	"strconv"
	"strings"

	"github.com/opennota/re2dfa/nfa"
	"github.com/opennota/runerange"
)

type Node struct {
	S int  // state
	F bool // final?
	T []T  // transitions

	label string
	cls   []*nfa.Node
}

type T struct {
	R []rune // rune ranges
	N *Node  // node
}

type context struct {
	state        int
	nodesByLabel map[string]*Node
	closureCache map[*nfa.Node][]*nfa.Node
}

func NewFromNFA(nfanode *nfa.Node) *Node {
	ctx := &context{
		nodesByLabel: make(map[string]*Node),
		closureCache: make(map[*nfa.Node][]*nfa.Node),
	}
	node := firstNode(nfanode, ctx)
	constructSubset(node, ctx)
	return node
}

func recursiveClosure(node *nfa.Node, visited map[*nfa.Node]struct{}) []*nfa.Node {
	if visited != nil {
		if _, ok := visited[node]; ok {
			return nil
		}
	}

	cls := []*nfa.Node{node}
	for _, t := range node.T {
		if t.R == nil {
			if visited == nil {
				visited = make(map[*nfa.Node]struct{})
				visited[node] = struct{}{}
			}
			if c := recursiveClosure(t.N, visited); c != nil {
				cls = append(cls, c...)
			}
		}
	}

	if visited != nil {
		delete(visited, node)
	}

	return cls
}

func closure(node *nfa.Node, cache map[*nfa.Node][]*nfa.Node) []*nfa.Node {
	if cache != nil {
		if cls, ok := cache[node]; ok {
			return cls
		}
	}

	cls := recursiveClosure(node, nil)

	if cache != nil {
		cache[node] = cls
	}

	return cls
}

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

func labelFromClosure(cls []*nfa.Node) string {
	m := make(map[int]struct{})
	for _, n := range cls {
		m[n.S] = struct{}{}
	}

	states := make([]int, 0, len(m))
	for n := range m {
		states = append(states, n)
	}

	sort.Ints(states)

	return makeLabel(states)
}

func isFinal(cls []*nfa.Node) bool {
	for _, n := range cls {
		if n.F {
			return true
		}
	}
	return false
}

func union(cls ...[]*nfa.Node) []*nfa.Node {
	if len(cls) == 1 {
		return cls[0]
	}

	size := 0
	for _, c := range cls {
		size += len(c)
	}

	m := make(map[*nfa.Node]struct{}, size)
	for _, c := range cls {
		for _, n := range c {
			m[n] = struct{}{}
		}
	}

	a := make([]*nfa.Node, 0, len(m))
	for n := range m {
		a = append(a, n)
	}

	return a
}

func closuresForRange(n *Node, rr []rune, ctx *context) (closures [][]*nfa.Node) {
	for _, n := range n.cls {
		for _, t := range n.T {
			if runerange.Contains(t.R, rr) {
				cls := closure(t.N, ctx.closureCache)
				closures = append(closures, cls)
			}
		}
	}
	return
}

func constructSubset(root *Node, ctx *context) {
	var ranges [][]rune
	for _, n := range root.cls {
		for _, t := range n.T {
			ranges = append(ranges, t.R)
		}
	}
	pairs := runerange.Split(ranges)

	m := make(map[*Node][]rune)

	for i := 0; i < len(pairs); i += 2 {
		cls := union(closuresForRange(root, pairs[i:i+2], ctx)...)

		label := labelFromClosure(cls)
		var node *Node
		if n, ok := ctx.nodesByLabel[label]; ok {
			node = n
		} else {
			ctx.state++
			node = &Node{
				S:     ctx.state,
				F:     isFinal(cls),
				label: label,
				cls:   cls,
			}
			ctx.nodesByLabel[label] = node
			constructSubset(node, ctx)
		}

		m[node] = runerange.Sum(m[node], pairs[i:i+2])
	}

	for n, rr := range m {
		root.T = append(root.T, T{rr, n})
	}
}

func firstNode(nfanode *nfa.Node, ctx *context) *Node {
	cls := closure(nfanode, ctx.closureCache)
	label := labelFromClosure(cls)

	ctx.state++
	node := &Node{
		S:     ctx.state,
		F:     isFinal(cls),
		label: label,
		cls:   cls,
	}
	ctx.nodesByLabel[label] = node

	return node
}
