package jsonparser

import "fmt"

// Parser ...
type Parser struct{}

// Parse ...
func Parse(src []byte) error {
	lexer := &Lexer{
		src:  src,
		line: 1,
	}

	tokens, err := lexer.ScanTokens()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}
