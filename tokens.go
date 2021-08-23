package jsonparser

import "fmt"

/*
	json				↣ array | object ;
	array 			↣ "[" ( value ( "," value )* )? "]" ;
	object  		↣ "{" ( objectEntry ( "," objectEntry ) ) "}" ;
	objectEntry	↣ STRING ":" value ;
	value       ↣ ( primary | object | array ) ;
	primary  		↣ ( null | false | true | STRING | NUMBER )
*/

/*
	Token Types
	String
	Number
	False
	True
	Null
	LeftBrace
	RightBrace
	LeftParen
	RightParen
	Comma
	Colon
*/

//go:generate stringer -type=TokenType

// TokenType ...
type TokenType int

const (
	String TokenType = iota
	Number
	False
	True
	Null
	LeftBrace
	RightBrace
	LeftParen
	RightParen
	Comma
	Colon
)

// Token ...
type Token struct {
	Lexeme  string
	Literal interface{}
	Line    int
	Type    TokenType
}

// String ...
func (t Token) String() string {
	return fmt.Sprintf("Token<%s, %s>\n", t.Type, t.Lexeme)
}
