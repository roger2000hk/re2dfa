re2dfa
======

re2dfa transforms regular expressions into deterministic finite state machines and outputs Go source code containing the matching function.

# Install

    go get github.com/opennota/re2dfa

# Usage

    re2dfa ^a+$ main.matchAPlus string

# TODO

* Support non-greedy matches.

# License

GNU GPL v3+
