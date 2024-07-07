package tabp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Parser define a Tabp parser.
type Parser struct {
	reader  *bufio.Reader
	cursor  Position
	symbols map[Symbol]any
	unread  bool
}

type ParseError struct {
	Cause    error
	Position Position
}

// NewParser returns a new initialized parser that will read from the given
// reader.
func NewParser(r io.Reader) Parser {
	return Parser{
		reader:  bufio.NewReader(r),
		cursor:  NewPosition(),
		symbols: map[Symbol]any{},
	}
}

// Error implements error.
func (pe ParseError) Error() string {
	return fmt.Sprintf("failed to parse tabp expression at %v: %v", pe.Position, pe.Cause)
}

// Unwrap returns underlying cause of this error.
func (pe ParseError) Unwrap() error {
	return pe.Cause
}

func (p *Parser) readRune() (rune, ParseError) {
	r, size, err := p.reader.ReadRune()
	if err != nil {
		return 0, ParseError{
			Cause:    err,
			Position: p.cursor,
		}
	}

	if !p.unread {
		p.cursor.byte += size
		p.cursor.col++

		if r == '\n' {
			p.cursor.line++
			p.cursor.col = 0
		}
	}
	p.unread = false

	return r, ParseError{}
}

func (p *Parser) mustUnreadRune() {
	err := p.reader.UnreadRune()
	if err != nil {
		panic(err)
	}
	p.unread = true
}

func (p *Parser) peekRune() (rune, ParseError) {
	rune, err := p.readRune()
	if err.Cause != nil {
		return rune, err
	}

	p.mustUnreadRune()

	return rune, ParseError{}
}

func (p *Parser) skipWhile(predicate func(r rune) bool) (r rune, err ParseError) {
	for {
		r, err = p.readRune()
		if err.Cause != nil {
			return
		}

		if !predicate(r) {
			break
		}
	}

	return
}

func (p *Parser) collectBytesWhile(predicate func(r rune) bool, buf []byte) ([]byte, rune, ParseError) {
	var lastRune rune

	for {
		r, err := p.readRune()
		if err.Cause != nil {
			// EOF error after reading some rune.
			if err.Cause == io.EOF && lastRune != 0 {
				// Return collected runes as EOF will be returned on next call.
				break
			}
			return nil, lastRune, err
		}

		lastRune = r
		if !predicate(r) {
			p.mustUnreadRune()
			break
		}

		buf = utf8.AppendRune(buf, r)
	}

	return buf, lastRune, ParseError{}
}

func (p *Parser) mustSkip(n int) {
	for i := 0; i < n; i++ {
		_, err := p.readRune()
		if err.Cause != nil {
			panic(err)
		}
	}
}

// Parse parses a single S-Expression and returns it as one of the following
// types:
//
//   - Table
//
//   - int / float64
//
//   - String
//
//   - Symbol
func (p *Parser) Parse() (Value, ParseError) {
	r, err := p.skipWhile(unicode.IsSpace)
	if err.Cause != nil {
		return nil, err
	}

	// Comment.
	if r == ';' || r == '/' { // Regular lisp comments.
		err = p.parseComment(r)
		if err.Cause != nil {
			return nil, err
		}

		// Skip whitespaces again.
		r, err = p.skipWhile(unicode.IsSpace)
		if err.Cause != nil {
			return nil, err
		}
	}

	// Table.
	if r == '(' {
		table, err := p.parseTable()
		if err.Cause != nil {
			return nil, err
		}
		if table == nil {
			return nil, ParseError{}
		}

		return table, ParseError{}
	}

	// Number.
	if r == '+' || r == '-' || r == '.' || unicode.IsDigit(r) {
		number, err := p.parseNumber(r)
		if err.Cause != nil {
			return nil, err
		}

		return number, ParseError{}
	}

	// String.
	if r == '"' || r == '\'' || r == '`' {
		str, err := p.parseString(r)
		if err.Cause != nil {
			return nil, err
		}

		return str, ParseError{}
	}

	// Symbol.
	symbol, err := p.parseSymbol(r)
	if err.Cause != nil {
		return nil, err
	}

	return symbol, ParseError{}
}

func (p *Parser) parseTable() (*Table, ParseError) {
	var table Table

	for {
		// Skip whitespaces.
		r, err := p.skipWhile(unicode.IsSpace)
		if err.Cause != nil {
			if err.Cause == io.EOF {
				return nil, ParseError{
					Cause:    fmt.Errorf("unexpected EOF, table closing parenthesis missing: %w", err),
					Position: p.cursor,
				}
			}

			return nil, err
		}

		// End of table.
		if r == ')' {
			return &table, ParseError{}
		}

		p.mustUnreadRune()

		// Parse values.
		value, err := p.Parse()
		if err.Cause != nil {
			return nil, err
		}

		// Skip whitespaces.
		r, err = p.skipWhile(unicode.IsSpace)
		if err.Cause != nil {
			return nil, err
		}

		// Value is a key.
		if r == ':' {
			key := value
			value, err = p.Parse()
			if err.Cause != nil {
				return nil, err
			}

			table.Set(key, value)
		} else {
			table.Append(value)
			p.mustUnreadRune()
		}
	}
}

func (p *Parser) parseNumber(r rune) (any, ParseError) {
	var (
		buf []byte
		err ParseError
	)

	buf = utf8.AppendRune(buf, r)

	buf, r, err = p.collectBytesWhile(unicode.IsDigit, buf)
	if err.Cause != nil {
		return nil, err
	}

	// Float.
	if r == '.' {
		buf = utf8.AppendRune(buf, r)
		p.mustSkip(1)

		buf, _, err = p.collectBytesWhile(unicode.IsDigit, buf)
		if err.Cause != nil {
			return nil, err
		}

		// Parse float.
		f, parseFloatErr := strconv.ParseFloat(UnsafeString(buf), 64)
		if parseFloatErr != nil {
			return 0.0, ParseError{
				Cause:    parseFloatErr,
				Position: p.cursor,
			}
		}

		return f, ParseError{}
	}

	// Integer.
	i, parseIntErr := strconv.Atoi(UnsafeString(buf))
	if parseIntErr != nil {
		return 0, ParseError{
			Cause:    parseIntErr,
			Position: p.cursor,
		}
	}

	return i, ParseError{}
}

func (p *Parser) parseString(r rune) (string, ParseError) {
	var (
		buf      []byte
		parseErr ParseError
	)

	buf = utf8.AppendRune(buf, r)

	buf, r, parseErr = p.collectBytesWhile(func(r rune) bool {
		return r != '"'
	}, buf)
	if parseErr.Cause != nil {
		return "", parseErr
	}

	buf = utf8.AppendRune(buf, r)
	// Move cursor after closing quote
	p.mustSkip(1)

	str, err := strconv.Unquote(UnsafeString(buf))
	if err != nil {
		return "", ParseError{
			Cause:    err,
			Position: p.cursor,
		}
	}

	return str, ParseError{}
}

func (p *Parser) parseSymbol(r rune) (Symbol, ParseError) {
	var (
		buf []byte
		err ParseError
	)

	buf = utf8.AppendRune(buf, r)

	if r == '|' {
		buf, _, err = p.collectBytesWhile(func(r rune) bool {
			return r != '|' && r != '(' && r != ')' && r != ':'
		}, buf)
		if err.Cause != nil {
			return Symbol(UnsafeString(buf)), err
		}
		buf = utf8.AppendRune(buf, '|')
	} else {
		buf, _, err = p.collectBytesWhile(func(r rune) bool {
			return unicode.IsPrint(r) && r != ' ' && r != '(' && r != ')' && r != ':'
		}, buf)
		if err.Cause != nil {
			return Symbol(UnsafeString(buf)), err
		}
	}

	return Symbol(UnsafeString(bytes.ToUpper(buf))), ParseError{}
}

func (p *Parser) parseComment(r rune) (err ParseError) {
	if r == ';' { // Lisp comment.
		r, err = p.skipWhile(func(r rune) bool {
			return r != '\n'
		})
	} else if r == '/' { // C style comments.
		r, err = p.peekRune()
		if err.Cause != nil {
			return err
		}

		switch r {
		// Single line comment.
		case '/':
			r, err = p.skipWhile(func(r rune) bool {
				return r != '\n'
			})
		// Multiline.
		case '*':
			previousR := '/'
			r, err = p.skipWhile(func(r rune) bool {
				if previousR == '*' && r == '/' {
					return false
				}

				previousR = r
				return true
			})
		}
	}

	return ParseError{}
}
