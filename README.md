re2dfa [![License](http://img.shields.io/:license-gpl3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0.html) [![Build Status](https://travis-ci.org/opennota/re2dfa.png?branch=master)](https://travis-ci.org/opennota/re2dfa)
======

re2dfa transforms regular expressions into deterministic finite state machines and outputs Go source code containing the matching function.

# Installation

    go get github.com/opennota/re2dfa

# Usage

    re2dfa ^a+$ main.matchAPlus string

# TODO

* Optimize generation from regular expressions containing multiple broad unicode ranges, such as `.*` or `[^a-z]`.
