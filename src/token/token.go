package token

import (
	"fmt"
	"strconv"

	"github.com/NikoMalik/Tod-go-compiler/src/print2"
)

type TokenType int

type Token struct {
	Type       TokenType
	Literal    string
	RealValue  interface{}
	SpaceAfter bool
	Span       print2.TextSpan
}

const (
	ILLEGAL TokenType = iota
	EOF
	IDENT
	COMMENT

	literal_beg
	INT
	UINT
	FLOAT32
	FLOAT64
	STRING
	BOOLEAN
	literal_end

	operator_beg

	ASSIGN         // =
	ADD            // +
	SUB            // -
	MUL            // *
	QUO            // /
	REM            // %
	AND            // &
	OR             // |
	XOR            // ^
	SHL            // <<
	SHR            // >>
	AND_NOT        // &^
	LAND           // &&
	LOR            // ||
	EQ             // ==
	NOT_EQ         // !=
	LT             // <
	LEQ            // <=
	GT             // >
	GEQ            // >=
	SPACESHIP      // <=>
	BANG           // !
	ADD_ASSIGN     // +=
	SUB_ASSIGN     // -=
	MUL_ASSIGN     // *=
	QUO_ASSIGN     // /=
	REM_ASSIGN     // %=
	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=
	COMMA          // ,
	LPAREN         // (
	RPAREN         // )
	LBRACE         // {
	RBRACE         // }
	PERIOD         // .
	LBRACK         // [
	RBRACK         // ]
	SEMICOLON      // ;
	DEFINE         // :=
	POINTER        // *
	ADDRESS        // &
	operator_end

	keyword_beg
	IMPORT
	MAIN
	EXTERNAL
	CONTINUE
	SET
	MAKE
	FN
	VAR
	IF
	BREAK
	FOR
	TRUE
	FALSE
	ELSE
	WHILE
	RETURN
	MAP
	PACKAGE
	RANGE
	STRUCT
	TYPE
	USING
	keyword_end
)

var tokens = map[TokenType]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	IDENT:      "IDENT",
	COMMENT:    "COMMENT",
	INT:        "INT",
	FLOAT32:    "FLOAT32",
	FLOAT64:    "FLOAT64",
	STRING:     "STRING",
	BOOLEAN:    "BOOLEAN",
	ASSIGN:     "=",
	ADD:        "+",
	SUB:        "-",
	MUL:        "*",
	QUO:        "/",
	LPAREN:     "(",
	RPAREN:     ")",
	LBRACE:     "{",
	RBRACE:     "}",
	COMMA:      ",",
	SEMICOLON:  ";",
	PERIOD:     ".",
	LBRACK:     "[",
	RBRACK:     "]",
	EQ:         "==",
	NOT_EQ:     "!=",
	LT:         "<",
	LEQ:        "<=",
	GT:         ">",
	GEQ:        ">=",
	DEFINE:     ":=",
	POINTER:    "*",
	ADDRESS:    "&",
	ADD_ASSIGN: "+=",
	SUB_ASSIGN: "-=",
	MUL_ASSIGN: "*=",
	QUO_ASSIGN: "/=",
	REM_ASSIGN: "%=",
	IMPORT:     "import",
	LOR:        "||",
	LAND:       "&&",
	AND:        "&",
	OR:         "|",
	XOR:        "^",
	SHL:        "<<",
	SHR:        ">>",
	AND_NOT:    "&^",
	FN:         "fn",
	VAR:        "var",
	IF:         "if",
	ELSE:       "else",
	WHILE:      "while",
	RETURN:     "return",
	MAP:        "map",
	PACKAGE:    "package",
	RANGE:      "range",
	STRUCT:     "struct",
	TYPE:       "type",
	TRUE:       "true",
	FALSE:      "false",
	BANG:       "!",
	SET:        "set",
	CONTINUE:   "continue",
	USING:      "using",
	EXTERNAL:   "external",
	MAKE:       "make",
	MAIN:       "main",
	BREAK:      "break",
	REM:        "%",
}

func GetUnaryOperatorPrecedence(tok Token) int {
	switch tokens[tok.Type] {
	case "+", "-", "!":
		return 6 // always one higher than the highest binary operator
	default:
		return 0
	}
}

func GetBinaryOperatorPrecedence(tok Token) int {
	switch tokens[tok.Type] {
	case "*", "/", "%":
		return 5
	case "+", "-":
		return 4
	case "==", "!=", "<", ">", "<=", ">=", "<<", ">>":
		return 3
	case "&":
		return 2
	case "|", "^":
		return 1
	default:
		return 0
	}
}

const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

func (t Token) String(pretty bool) string {
	if !pretty {
		return fmt.Sprintf("Token { value: %s, kind: %s, position: (%d, %d), real: %v}", t.Literal, t.Type, t.Span.StartLine, t.Span.StartColumn, t.RealValue)
	} else {
		return fmt.Sprintf("Token { \n\tvalue: %s, \n\tkind: %s, \n\tposition: (L%d, SC%d, EC%d, Len %d)\n}", t.Literal, t.Type, t.Span.StartLine, t.Span.StartColumn, t.Span.EndColumn, t.Span.EndIndex-t.Span.StartIndex)
	}
}

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.

func (tok TokenType) String() string {
	s := ""
	if 0 <= tok && tok < TokenType(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for tok := keyword_beg + 1; tok < keyword_end; tok++ {
		keywords[tokens[tok]] = tok
	}
}

func LookupIdent(ident string) TokenType {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func CreateToken(literal string, Type TokenType) Token {
	return Token{
		Type:      Type,
		Literal:   literal,
		RealValue: nil,
		Span:      print2.TextSpan{},
	}
}

func CreateTokenSpaced(literal string, Type TokenType, spaced bool, span print2.TextSpan) Token {
	return Token{
		Type:       Type,
		Literal:    literal,
		RealValue:  nil,
		SpaceAfter: spaced,
		Span:       span,
	}
}

func CreateTokenReal(buffer string, real interface{}, Type TokenType, span print2.TextSpan) Token {
	return Token{
		Type:      Type,
		Literal:   buffer,
		RealValue: real,
		Span:      span,
	}
}

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (tok TokenType) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (tok TokenType) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok TokenType) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }

// IsExported reports whether name is exported...
