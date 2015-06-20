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
	"bytes"
	"fmt"
	"go/format"
	"sort"

	"github.com/opennota/re2dfa/nfa"
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

func closuresForRune(n *Node, r rune, ctx *context) (closures [][]*nfa.Node) {
	for _, n := range n.cls {
		for _, t := range n.T {
			if inRange(r, t.R) {
				cls := closure(t.N, ctx.closureCache)
				closures = append(closures, cls)
			}
		}
	}
	return
}

func constructSubset(root *Node, ctx *context) {
	var ranges []rune
	for _, n := range root.cls {
		for _, t := range n.T {
			ranges = foldRanges(ranges, t.R)
		}
	}

	m := make(map[*Node][]rune)

	for i := 0; i < len(ranges); i += 2 {
		for r := ranges[i]; r <= ranges[i+1]; r++ {
			cls := union(closuresForRune(root, r, ctx)...)

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

			m[node] = appendToRange(m[node], r)
		}
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

func allNodes(n *Node, visited map[*Node]struct{}) []*Node {
	if _, ok := visited[n]; ok {
		return nil
	}
	visited[n] = struct{}{}

	nodes := []*Node{n}
	for _, t := range n.T {
		nodes = append(nodes, allNodes(t.N, visited)...)
	}

	return nodes
}

type nodesByState []*Node

func (s nodesByState) Len() int           { return len(s) }
func (s nodesByState) Less(i, j int) bool { return s[i].S < s[j].S }
func (s nodesByState) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func GoGenerate(dfa *Node, packageName, funcName, typ string) string {
	if !(typ == "string" || typ == "[]byte") {
		panic(fmt.Sprintf("invalid type: %s; expected either string or []byte", typ))
	}

	instr := ""
	if typ == "string" {
		instr = "InString"
	}

	nodes := allNodes(dfa, make(map[*Node]struct{}))
	sort.Sort(nodesByState(nodes))

	end := -1
	if nodes[0].F {
		end = 0
	}

	label := ""
outer:
	for _, n := range nodes {
		for _, t := range n.T {
			if t.R[0] < 0 && len(t.N.T) > 0 {
				label = "l0:"
				break outer
			}
		}
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `// Code generated by re2dfa (https://github.com/opennota/re2dfa).

			package %s
			import "unicode/utf8"
			//func isWordChar(r byte) bool {
			//        return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '_'
			//}

			func %s(s %s) int {
				st := 1
				end := %d
				i := 0
				var r rune
				_ = r
				var rlen int

				for {
					r, rlen = utf8.DecodeRune%s(s[i:])
					if rlen == 0 { break }
					i += rlen
			%s
					switch st {
`, packageName, funcName, typ, end, instr, label)

	for _, n := range nodes {
		if len(n.T) == 0 {
			continue
		}
		fmt.Fprintf(&buf, `case %d:
				switch {
		`, n.S)
		for _, t := range n.T {
			i := 0
			for i < len(t.R) && t.R[i] < 0 {
				i++
			}
			if i > 0 {
				fmt.Fprintf(&buf, "case %s:\n", rangesToBoolExpr(t.R[:i], false))
				if t.N.F {
					fmt.Fprintln(&buf, "end = i - rlen")
				}
				if len(t.N.T) > 0 {
					fmt.Fprintf(&buf, "st = %d\n", t.N.S)
					fmt.Fprintln(&buf, "goto l0")
				} else {
					fmt.Fprintln(&buf, "return end")
				}
			}
			if i < len(t.R) {
				fmt.Fprintf(&buf, "case %s:\n", rangesToBoolExpr(t.R[i:], false))
				if t.N.F {
					fmt.Fprintln(&buf, "end = i")
				}
				if len(t.N.T) > 0 {
					fmt.Fprintf(&buf, "st = %d\n", t.N.S)
				} else {
					fmt.Fprintln(&buf, "return end")
				}
			}
		}
		fmt.Fprint(&buf, `default: return end
				}
`)

	}

	fmt.Fprintln(&buf, `}
		}
`)

	var buf2 bytes.Buffer
	for _, n := range nodes {
		hasEndStates := false
		for _, t := range n.T {
			if !t.N.F {
				continue
			}
			if len(t.R) > 0 && t.R[0] < 0 {
				hasEndStates = true
				break
			}
		}
		if !hasEndStates {
			continue
		}

		fmt.Fprintf(&buf2, `case %d:
				switch {
		`, n.S)
		for _, t := range n.T {
			if !t.N.F {
				continue
			}
			var rr []rune
			for i := 0; i < len(t.R) && t.R[i] < 0; i += 2 {
				rr = append(rr, t.R[i:i+2]...)
			}
			if len(rr) > 0 {
				fmt.Fprintf(&buf2, "case %s:\n", rangesToBoolExpr(rr, true))
				fmt.Fprintln(&buf2, "end = i")
			}
		}
		fmt.Fprintln(&buf2, "}")
	}

	if buf2.Len() > 0 {
		fmt.Fprintf(&buf, `switch st {
					%s
				}

`, buf2.String())
	}

	fmt.Fprintln(&buf, `return end
}
`)

	source, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	return string(source)
}
