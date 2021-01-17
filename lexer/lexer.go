package lexer

type Lexer struct {
	input         string
	postition     int
	readPosistion int
	ch            byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	return l
}
