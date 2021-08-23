package jsonparser

import (
	"fmt"
	"strconv"
)

// Lexer ...
type Lexer struct {
	src     []byte
	line    int
	current int
	start   int
	tokens  []Token
}

var keywords = map[string]TokenType{
	"true":  True,
	"false": False,
	"null":  Null,
}

// ScanTokens ...
func (l *Lexer) ScanTokens() ([]Token, error) {
	l.line = 1
	for !l.isAtEnd() {
		l.start = l.current
		if err := l.ScanToken(); err != nil {
			return nil, err
		}
	}

	return l.tokens, nil
}

// ScanToken ...
func (l *Lexer) ScanToken() error {
	cur := l.advance()
	switch cur {
	case '{':
		t := Token{
			Lexeme: "{",
			Line:   l.line,
			Type:   LeftBrace,
		}
		l.add(t)
	case '}':
		t := Token{
			Lexeme: "}",
			Line:   l.line,
			Type:   RightBrace,
		}
		l.add(t)
	case '[':
		t := Token{
			Lexeme: "[",
			Line:   l.line,
			Type:   LeftParen,
		}
		l.add(t)
	case ']':
		t := Token{
			Lexeme: "]",
			Line:   l.line,
			Type:   RightParen,
		}
		l.add(t)
	case ',':
		t := Token{
			Lexeme: ",",
			Line:   l.line,
			Type:   Comma,
		}
		l.add(t)
	case ':':
		t := Token{
			Lexeme: ":",
			Line:   l.line,
			Type:   Colon,
		}
		l.add(t)
	case '\n':
		l.line++
	case ' ', '\t', '\r':
		//ignore whitespace
	case '"':
		if err := l.getString(); err != nil {
			return err
		}
	default:
		if l.isAlpha(cur) {
			if err := l.identifier(); err != nil {
				return err
			}
		} else if l.isDigit(cur) {
			if err := l.number(); err != nil {
				return err
			}
		} else {
			return LexerError{
				msg:  fmt.Sprintf("Unexpected character %c.", cur),
				line: l.line,
			}
		}
	}

	return nil
}

func (l *Lexer) advance() byte {
	val := l.src[l.current]
	l.current++
	return val
}

func (l *Lexer) add(t Token) {
	l.tokens = append(l.tokens, t)
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.src)
}

func (l *Lexer) identifier() error {
	for l.isAlpha(l.peek()) {
		l.advance()
	}

	lexeme := l.src[l.start:l.current]
	if tokenType, ok := keywords[string(lexeme)]; ok {
		var literal interface{}
		switch tokenType {
		case True:
			literal = true
		case False:
			literal = false
		case Null:
			literal = nil
		}
		l.add(Token{
			Lexeme:  string(lexeme),
			Literal: literal,
			Line:    l.line,
			Type:    tokenType,
		})
		return nil
	}

	return LexerError{
		msg:  fmt.Sprintf("Unexpected identifier '%s'.", lexeme),
		line: l.line,
	}
}

func (l *Lexer) number() error {
	for l.isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && l.isDigit(l.peekNext()) {
		l.advance()
		for l.isDigit(l.peek()) {
			l.advance()
		}
	}

	val := l.src[l.start:l.current]
	number, err := strconv.ParseFloat(string(val), 64)
	if err != nil {
		return err
	}

	l.add(Token{
		Lexeme:  string(val),
		Literal: number,
		Line:    l.line,
		Type:    Number,
	})
	return nil
}

func (l *Lexer) getString() error {
	for l.peek() != '"' && !l.isAtEnd() && l.peek() != '\n' {
		l.advance()
	}

	if l.isAtEnd() {
		return LexerError{
			msg:  "Unterminated string",
			line: l.line,
		}
	}

	if l.peek() == '\n' {
		return LexerError{
			msg:  "JSON doesn't allow newlines in strings",
			line: l.line,
		}
	}

	// Consume the '"'
	l.advance()

	l.add(Token{
		Lexeme:  string(l.src[l.start:l.current]),
		Literal: string(l.src[l.start+1 : l.current-1]),
		Line:    l.line,
		Type:    String,
	})
	return nil
}

func (l *Lexer) isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_'
}

func (l *Lexer) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (l *Lexer) peek() byte {
	return l.src[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.src) {
		return '\000'
	}

	return l.src[l.current+1]
}

// LexerError error encountered when lexing
type LexerError struct {
	msg  string
	line int
}

func (err LexerError) Error() string {
	return fmt.Sprintf("[line %d] Error %s", err.line, err.msg)
}
