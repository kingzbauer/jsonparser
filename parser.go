package jsonparser

import "fmt"

// Parser ...
type Parser struct {
	tokens  []Token
	current int
}

// Parse ...
func Parse(src []byte) error {
	lexer := &Lexer{
		src:  src,
		line: 1,
	}

	tokens, err := lexer.ScanTokens()
	if err != nil {
		return err
	}

	p := &Parser{tokens: tokens}
	return p.parse()
}

func (p *Parser) parse() error {
	if p.match(LeftBrace) {
		return p.object()
	} else if p.match(LeftParen) {
		return p.array()
	}

	return ParserError{
		token: p.peek(),
		msg:   "Unexpected token.",
	}
}

func (p *Parser) advance() Token {
	t := p.tokens[p.current]
	p.current++
	return t
}

func (p *Parser) match(typ ...TokenType) bool {
	for _, t := range typ {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == typ
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) object() error {
	if !p.check(RightBrace) {
		if err := p.objectEntries(); err != nil {
			return err
		}
	}

	if _, err := p.consume(RightBrace, "Expected a '}' or separator between values."); err != nil {
		return err
	}
	return nil
}

func (p *Parser) objectEntries() error {
	if p.check(String) {
		if err := p.objectEntry(); err != nil {
			return err
		}

		for p.match(Comma) {
			if err := p.objectEntry(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Parser) objectEntry() error {
	if _, err := p.consume(String, "Expected a key of type 'String'."); err != nil {
		return err
	}

	if _, err := p.consume(Colon, "Expected ':' after key"); err != nil {
		return err
	}

	if err := p.value(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) value() error {
	next := p.advance()

	switch next.Type {
	case String:
		return nil
	case Number:
		return nil
	case Null:
		return nil
	case True:
		return nil
	case False:
		return nil
	case LeftBrace:
		return p.object()
	case LeftParen:
		return p.array()
	case EOF:
		return ParserError{
			msg:   "Missing key-value pair",
			token: next,
		}
	default:
		return ParserError{
			msg:   "Unexpected token.",
			token: next,
		}
	}
}

func (p *Parser) array() error {
	if !p.check(RightParen) {
		if err := p.value(); err != nil {
			return err
		}

		for p.match(Comma) {
			if err := p.value(); err != nil {
				return err
			}
		}
	}

	if _, err := p.consume(RightParen, "Expected a ']' or separator between values."); err != nil {
		return err
	}

	return nil
}

func (p *Parser) consume(typ TokenType, err string) (*Token, error) {
	if p.peek().Type != typ {
		return nil, ParserError{
			token: p.peek(),
			msg:   err,
		}
	}

	t := p.advance()
	return &t, nil
}

// ParserError ...
type ParserError struct {
	token Token
	msg   string
}

func (p ParserError) Error() string {
	return fmt.Sprintf("Error [line %d]. %s Got token '%v'.", p.token.Line, p.msg, p.token.Lexeme)
}
