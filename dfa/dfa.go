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

			m[node] = addToRange(m[node], r)
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

	labelFirstState := false
	enableLazy := false
	var lazyStates map[int]struct{}
	for i, n := range nodes {
		for _, t := range n.T {
			if t.N == nodes[0] {
				labelFirstState = true
			}
			hasLazy := false
			for i := 0; i < len(t.R); i += 2 {
				if t.R[i] == nfa.RuneLazy {
					hasLazy = true
					break
				}
			}
			if hasLazy {
				enableLazy = true
				if i == 0 {
					labelFirstState = true
				}
			}
		}
	}
	returnOrBacktrack := "return"
	if enableLazy {
		lazyStates = make(map[int]struct{})
		returnOrBacktrack = "goto lazy"
	}

	needUtf8 := false
	atLeastOneSwitch := false
	usesIsWordChar := false

	var buf bytes.Buffer

	for ni, n := range nodes {
		if len(n.T) == 0 {
			continue
		}

		if n.S != 1 || labelFirstState {
			fmt.Fprintf(&buf, "s%d:\n", n.S)
		}

		hasEmpty := false
		hasNonEmpty := false
		hasLazy := false
		for _, t := range n.T {
			for i := 0; i < len(t.R); i += 2 {
				if t.R[i] < 0 {
					if t.R[i] == nfa.RuneLazy {
						hasLazy = true
					} else {
						hasEmpty = true
					}
				} else {
					hasNonEmpty = true
				}
			}
		}

		if hasLazy {
			lazyStates[n.S] = struct{}{}
			for _, t := range n.T {
				for i := 0; i < len(t.R) && t.R[i] < 0; i += 2 {
					if t.R[i] != nfa.RuneLazy {
						continue
					}
					fmt.Fprintf(&buf, `if lazyOn {
								lazyOn = false
								goto s%d
							}
							lazyStack = append(lazyStack, jmp{s:%d, i:i})
							`, t.N.S, n.S)
					break
				}
			}
		}

		if hasEmpty {
			atLeastOneSwitch = true
			fmt.Fprintln(&buf, "switch {")
			for _, t := range n.T {
				for i := 0; i < len(t.R) && t.R[i] < 0; i += 2 {
					if t.R[i] == nfa.RuneLazy {
						continue
					}
					if t.R[i] == nfa.RuneWordBoundary || t.R[i] == nfa.RuneNoWordBoundary {
						usesIsWordChar = true
					}
					fmt.Fprintf(&buf, "case %s:\n", rangesToBoolExpr(t.R[i:i+2], false))
					if t.N.F {
						fmt.Fprintln(&buf, "end = i")
					}
					if len(t.N.T) > 0 {
						fmt.Fprintf(&buf, "goto s%d\n", t.N.S)
					} else if hasNonEmpty {
						fmt.Fprintln(&buf, returnOrBacktrack)
					}
				}
			}
			fmt.Fprintln(&buf, "}")
		}

		if hasNonEmpty {
			atLeastOneSwitch = true
			needUtf8 = true
			fmt.Fprintf(&buf, `r, rlen = utf8.DecodeRune%s(s[i:])
						if rlen == 0 { %s }
						i += rlen
						switch {
						`, instr, returnOrBacktrack)
			for _, t := range n.T {
				i := 0
				for i < len(t.R) && t.R[0] < 0 {
					i++
				}
				if i >= len(t.R) {
					continue
				}

				fmt.Fprintf(&buf, "case %s:\n", rangesToBoolExpr(t.R[i:], false))
				if t.N.F {
					fmt.Fprintln(&buf, "end = i")
				}
				if len(t.N.T) > 0 {
					fmt.Fprintf(&buf, "goto s%d\n", t.N.S)
				}
			}
			fmt.Fprintln(&buf, "}")
		}
		if !enableLazy || ni != len(nodes)-1 {
			fmt.Fprintln(&buf, returnOrBacktrack)
		}
	}

	if enableLazy {
		fmt.Fprintln(&buf, `lazy:
					if end >= 0 || len(lazyStack) == 0 { return }
					var to jmp
					to, lazyStack = lazyStack[len(lazyStack)-1], lazyStack[:len(lazyStack)-1]
					lazyOn = true
					i = to.i
					switch to.s {`)
		states := make([]int, 0, len(lazyStates))
		for s := range lazyStates {
			states = append(states, s)
		}
		sort.Ints(states)
		for _, s := range states {
			fmt.Fprintf(&buf, "case %d: goto s%[1]d\n", s)
		}
		fmt.Fprintln(&buf, "}")
		fmt.Fprintln(&buf, "return")
	}

	imports := ""
	if needUtf8 {
		imports = `import "unicode/utf8"`
	}

	helperFuncs := ""
	if usesIsWordChar {
		helperFuncs = `
			//func isWordChar(r byte) bool {
			//        return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '_'
			//}`
	}

	end := -1
	if nodes[0].F {
		end = 0
	}

	decls := `var r rune
		var rlen int
		i := 0`
	if enableLazy {
		decls += `
			lazyOn := false
			type jmp struct { s, i int }
			var lazyStack []jmp`
	}

	var buf2 bytes.Buffer
	fmt.Fprintf(&buf2, `// Code generated by re2dfa (https://github.com/opennota/re2dfa).

			package %s
			%s
			%s

			func %s(s %s) (end int) {
				end = %d
				%s
				_, _, _ = r, rlen, i
`, packageName, imports, helperFuncs, funcName, typ, end, decls)
	buf2.Write(buf.Bytes())
	if !atLeastOneSwitch {
		fmt.Fprintln(&buf2, "return")
	}
	fmt.Fprintln(&buf2, "}")

	source, err := format.Source(buf2.Bytes())
	if err != nil {
		panic(err)
	}

	return string(source)
}
