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

package nfa

import "regexp/syntax"

const (
	// Last unicode rune
	RuneLast = 0x10ffff

	// Pseudo-runes
	RuneBeginText = -100 * iota
	RuneEndText
	RuneBeginLine
	RuneEndLine
	RuneWordBoundary
	RuneNoWordBoundary
)

type Node struct {
	S int  // state
	F bool // final?
	T []T  // transitions
}

type T struct {
	R []rune // rune ranges
	N *Node  // node
}

type context struct {
	state int
}

func (c *context) node() *Node {
	c.state++
	return &Node{S: c.state}
}

func (n *Node) copy() *Node {
	nn := Node{
		S: n.S,
		F: n.F,
		T: make([]T, len(n.T)),
	}
	copy(nn.T, n.T)
	return &nn
}

func New(pattern string) (*Node, error) {
	r, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return nil, err
	}

	return NewFromRegexp(r.Simplify()), nil
}

func NewFromRegexp(r *syntax.Regexp) *Node {
	begin, end := recursiveNewFromRegexp(r, &context{})
	end.F = true
	return begin
}

func opString(op syntax.Op) string {
	switch op {
	case syntax.OpNoMatch:
		return "OpNoMatch"
	case syntax.OpEmptyMatch:
		return "OpEmptyMatch"
	case syntax.OpLiteral:
		return "OpLiteral"
	case syntax.OpCharClass:
		return "OpCharClass"
	case syntax.OpAnyCharNotNL:
		return "OpAnyCharNotNL"
	case syntax.OpAnyChar:
		return "OpAnyChar"
	case syntax.OpBeginLine:
		return "OpBeginLine"
	case syntax.OpEndLine:
		return "OpEndLine"
	case syntax.OpBeginText:
		return "OpBeginText"
	case syntax.OpEndText:
		return "OpEndText"
	case syntax.OpWordBoundary:
		return "OpWordBoundary"
	case syntax.OpNoWordBoundary:
		return "OpNoWordBoundary"
	case syntax.OpCapture:
		return "OpCapture"
	case syntax.OpStar:
		return "OpStar"
	case syntax.OpPlus:
		return "OpPlus"
	case syntax.OpQuest:
		return "OpQuest"
	case syntax.OpRepeat:
		return "OpRepeat"
	case syntax.OpConcat:
		return "OpConcat"
	case syntax.OpAlternate:
		return "OpAlternate"
	}
	return "OpUnknown"
}

func recursiveNewFromRegexp(r *syntax.Regexp, ctx *context) (begin *Node, end *Node) {
	switch r.Op {
	case syntax.OpEmptyMatch:
		begin = ctx.node()
		end = begin

	case syntax.OpLiteral:
		begin = ctx.node()
		cur := begin
		for _, r := range r.Rune {
			end = ctx.node()
			cur.T = append(cur.T, T{R: []rune{r, r}, N: end})
			cur = end
		}

	case syntax.OpCharClass:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: r.Rune, N: end})

	case syntax.OpAnyCharNotNL:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{0, 9, 11, RuneLast}, N: end})

	case syntax.OpAnyChar:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{0, RuneLast}, N: end})

	case syntax.OpBeginLine:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{RuneBeginLine, RuneBeginLine}, N: end})

	case syntax.OpEndLine:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{RuneEndLine, RuneEndLine}, N: end})

	case syntax.OpBeginText:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{RuneBeginText, RuneBeginText}, N: end})

	case syntax.OpEndText:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{RuneEndText, RuneEndText}, N: end})

	case syntax.OpWordBoundary:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{RuneWordBoundary, RuneWordBoundary}, N: end})

	case syntax.OpNoWordBoundary:
		begin = ctx.node()
		end = ctx.node()
		begin.T = append(begin.T, T{R: []rune{RuneNoWordBoundary, RuneNoWordBoundary}, N: end})

	case syntax.OpCapture:
		return recursiveNewFromRegexp(r.Sub[0], ctx)

	case syntax.OpStar:
		begin = ctx.node()
		end = ctx.node()
		b, e := recursiveNewFromRegexp(r.Sub[0], ctx)
		begin.T = append(begin.T, T{N: b})
		begin.T = append(begin.T, T{N: end})
		e.T = append(e.T, T{N: b})
		e.T = append(e.T, T{N: end})

	case syntax.OpPlus:
		begin = ctx.node()
		end = ctx.node()
		b, e := recursiveNewFromRegexp(r.Sub[0], ctx)
		begin.T = append(begin.T, T{N: b})
		e.T = append(e.T, T{N: b})
		e.T = append(e.T, T{N: end})

	case syntax.OpQuest:
		begin = ctx.node()
		end = ctx.node()
		b, e := recursiveNewFromRegexp(r.Sub[0], ctx)
		begin.T = append(begin.T, T{N: b})
		begin.T = append(begin.T, T{N: end})
		e.T = append(e.T, T{N: end})

	case syntax.OpRepeat:
		toRepeat, e := recursiveNewFromRegexp(r.Sub[0], ctx)

		var prev *Node
		for i := 0; i < r.Min; i++ {
			node := toRepeat.copy()
			if begin == nil {
				begin = node
			}
			if prev != nil {
				prev.T = append(prev.T, T{N: node})
			}
			prev = node
		}

		b := ctx.node()
		if begin == nil {
			begin = b
		} else {
			prev.T = append(prev.T, T{N: b})
		}
		end = ctx.node()
		b.T = append(b.T, T{N: toRepeat})
		b.T = append(b.T, T{N: end})
		e.T = append(e.T, T{N: toRepeat})
		e.T = append(e.T, T{N: end})

	case syntax.OpConcat:
		var cur *Node
		for _, r := range r.Sub {
			var b *Node
			b, end = recursiveNewFromRegexp(r, ctx)
			if begin == nil {
				begin = b
			}
			if cur != nil {
				cur.T = append(cur.T, T{N: b})
			}
			cur = end
		}

	case syntax.OpAlternate:
		begin = ctx.node()
		end = ctx.node()
		for _, r := range r.Sub {
			b, e := recursiveNewFromRegexp(r, ctx)
			begin.T = append(begin.T, T{N: b})
			e.T = append(e.T, T{N: end})
		}

	default:
		panic("unsupported op: " + opString(r.Op))
	}

	return
}
