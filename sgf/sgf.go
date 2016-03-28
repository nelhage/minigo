// Package sgf contains a parser for the SGF "Smart Game Format"
// format. See http://www.red-bean.com/sgf/ for the format
// specification

//go:generate -command yacc go tool yacc
//go:generate yacc -o parser.go parser.y

package sgf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const eof = 0

type lexer struct {
	c        *Collection
	unread   rune
	r        io.RuneReader
	ioErr    error
	parseErr []string
}

func (l *lexer) next() rune {
	if l.unread != 0 {
		r := l.unread
		l.unread = 0
		return r
	}
	r, _, err := l.r.ReadRune()
	if err != nil {
		l.ioErr = err
		return eof
	}
	return r
}

func (l *lexer) putback(r rune) {
	if l.unread != 0 {
		panic("double putback")
	}
	l.unread = r
}

func (l *lexer) Lex(lval *yySymType) int {
	for {
		r := l.next()
		if r == eof {
			return eof
		}
		if r >= 'A' && r <= 'Z' {
			return l.lexPropName(lval, r)
		}
		switch r {
		case '(', ')', ';':
			return int(r)
		case '[':
			return l.lexPropVal(lval)
		case '\v', ' ', '\t', '\r', '\n':
		default:
			l.Error(fmt.Sprintf("Unexpected character `%c'", r))
		}
	}
}

func (l *lexer) lexPropName(lval *yySymType, r rune) int {
	rs := []rune{r}
	for {
		r := l.next()
		if r == eof {
			break
		}
		if r < 'A' || r > 'Z' {
			l.putback(r)
			break
		}
		rs = append(rs, r)
	}
	lval.name = string(rs)
	return TokPropName
}

func (l *lexer) lexPropVal(lval *yySymType) int {
	var b bytes.Buffer
L:
	for {
		r := l.next()
		switch r {
		case eof:
			break L
		case ']':
			break L
		case '\\':
			rr := l.next()
			switch rr {
			case eof:
				break L
			case '\n':
				continue L
			}
			r = rr
		}
		b.WriteRune(r)
	}
	lval.v = PropValue(b.String())
	return TokPropValue
}

func (l *lexer) Error(s string) {
	l.parseErr = append(l.parseErr, s)
}

// ParseError represents one or more parse errors parsing an SGF file
type ParseError []string

func (p ParseError) Error() string {
	return strings.Join([]string(p), "\n")
}

// ParseSGF parses an SGF file from the provided reader and returns it
// in tree form
func ParseSGF(in io.Reader) (*Collection, error) {
	l := &lexer{
		r: bufio.NewReader(in),
	}
	p := yyNewParser()
	p.Parse(l)
	if l.ioErr != nil && l.ioErr != io.EOF {
		return nil, l.ioErr
	}

	if len(l.parseErr) > 0 {
		return nil, ParseError(l.parseErr)
	}
	return l.c, nil
}
