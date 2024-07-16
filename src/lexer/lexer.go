package lexer

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/NikoMalik/Tod-go-compiler/src/print2"
	"github.com/NikoMalik/Tod-go-compiler/src/token"
)

// Lexer : Lexer struct for lexing :GentlemenSphere:
type Lexer struct {
	Code                  []rune
	File                  string
	Line                  int
	Column                int
	Index                 int
	Tokens                []token.Token
	TreatHashtagAsComment bool
}

// Lex takes a filename and converts it into its respective lexical tokens
func Lex(code []rune, filename string) []token.Token {
	return LexInternal(code, filename, true)
}

// LexInternal converts code into tokens based on the filename and treatment of hashtags as comments
func LexInternal(code []rune, filename string, treatHashtagsAsComments bool) []token.Token {
	scanner := &Lexer{
		Code:                  code,
		File:                  filename,
		Line:                  1,
		Column:                1,
		Index:                 0,
		Tokens:                make([]token.Token, 0),
		TreatHashtagAsComment: treatHashtagsAsComments,
	}

	RememberSourceFile(code, filename)

	for scanner.Index < len(scanner.Code) {
		c := scanner.Code[scanner.Index]

		peek := func(offset int) rune {
			if scanner.Index+offset < len(scanner.Code) {
				return scanner.Code[scanner.Index+offset]
			}
			return '\000'
		}

		if unicode.IsLetter(c) {
			scanner.getId()
		} else if unicode.IsNumber(c) {
			scanner.getNumber()
		} else if c == '"' || c == '\'' {
			scanner.getString()
		} else if c == '/' && peek(1) == '/' ||
			(scanner.TreatHashtagAsComment && c == '#') {
			scanner.getComment()
		} else if c != ' ' && c != '\n' && c != '\t' && c != '\v' {
			scanner.getOperator()
		} else {
			scanner.Increment()
		}
	}

	scanner.Tokens = append(scanner.Tokens, token.CreateToken("\000", token.EOF))
	return scanner.Tokens
}

// getNumber keeps getting bytes until it finds a non-number
// then it generates an integer (or a float) token and slaps it back to the lexer.
func (lxr *Lexer) getNumber() {
	buffer := string(lxr.Code[lxr.Index])
	lxr.Increment()

	if lxr.Code[lxr.Index] == 'x' {
		lxr.Increment()
		lxr.getNumberHex()
		return
	}

	if lxr.Code[lxr.Index] == 'b' {
		lxr.Increment()
		lxr.getNumberBinary()
		return
	}

	isDigitOrDotOrUnderScore := func(c rune) bool {
		return unicode.IsDigit(c) || c == '.' || c == '_'
	}

	for lxr.Index < len(lxr.Code) && isDigitOrDotOrUnderScore(lxr.Code[lxr.Index]) {
		if lxr.Code[lxr.Index] != '_' {
			buffer += string(lxr.Code[lxr.Index])
		}
		lxr.Increment()
	}

	if strings.Contains(buffer, ".") {
		realValueBuffer, err := strconv.ParseFloat(strings.ReplaceAll(buffer, "_", ""), 32)
		if err != nil {
			print2.Error(
				"LEXER",
				print2.RealValueConversionError,
				lxr.GetCurrentTextSpan(len(buffer)),
				"value \"%s\" could not be converted to real value [float] (NumberToken)!",
				buffer,
			)
		}

		lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, float32(realValueBuffer), token.FLOAT32, lxr.GetCurrentTextSpan(len(buffer))))

	} else {
		realValueBuffer, err := strconv.Atoi(strings.ReplaceAll(buffer, "_", ""))
		if err != nil {
			realerValueBuffer, err := strconv.ParseInt(strings.ReplaceAll(buffer, "_", ""), 10, 64)
			if err != nil {
				print2.Error(
					"LEXER",
					print2.RealValueConversionError,
					lxr.GetCurrentTextSpan(len(buffer)),
					"value \"%s\" could not be converted to real value [int] (NumberToken)!",
					buffer,
				)
			}
			lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, realerValueBuffer, token.INT, lxr.GetCurrentTextSpan(len(buffer))))
		}
		lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, realValueBuffer, token.INT, lxr.GetCurrentTextSpan(len(buffer))))
	}
}

// getNumberHex keeps getting bytes until it finds a character that isn't 0-9 or A-F
// then it generates an integer token and slaps it back to the lexer.
func (lxr *Lexer) getNumberHex() {
	allowedChars := "abcdefABCDEF"

	buffer := string(lxr.Code[lxr.Index])
	lxr.Increment()

	isDigitOrAllowedLetter := func(c rune) bool {
		return unicode.IsDigit(c) || strings.ContainsRune(allowedChars, c)
	}

	for lxr.Index < len(lxr.Code) && isDigitOrAllowedLetter(lxr.Code[lxr.Index]) {
		if lxr.Code[lxr.Index] != '_' {
			buffer += string(lxr.Code[lxr.Index])
		}
		lxr.Increment()
	}

	realValueBuffer, err := strconv.ParseInt(strings.ReplaceAll(buffer, "_", ""), 16, 32)
	if err != nil {
		print2.Error(
			"LEXER",
			print2.RealValueConversionError,
			lxr.GetCurrentTextSpan(len(buffer)),
			"hex value \"%s\" could not be converted to real value [int] (NumberToken)!",
			buffer,
		)
	}
	lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, int(realValueBuffer), token.INT, lxr.GetCurrentTextSpan(len(buffer))))
}

// getNumberBinary keeps getting bytes until it finds a character that isn't 0 or 1
// then it generates an integer token and slaps it back to the lexer.
func (lxr *Lexer) getNumberBinary() {
	buffer := string(lxr.Code[lxr.Index])
	lxr.Increment()

	for lxr.Index < len(lxr.Code) && (lxr.Code[lxr.Index] == '0' ||
		lxr.Code[lxr.Index] == '1') {
		if lxr.Code[lxr.Index] != '_' {
			buffer += string(lxr.Code[lxr.Index])
		}
		lxr.Increment()
	}

	realValueBuffer, err := strconv.ParseInt(strings.ReplaceAll(buffer, "_", ""), 2, 32)
	if err != nil {
		print2.Error(
			"LEXER",
			print2.RealValueConversionError,
			lxr.GetCurrentTextSpan(len(buffer)),
			"binary value \"%s\" could not be converted to real value [int] (NumberToken)!",
			buffer,
		)
	}
	lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, int(realValueBuffer), token.INT, lxr.GetCurrentTextSpan(len(buffer))))
}

// getString keeps getting bytes until it finds the end of the string
// then it generates a string token and slaps it back to the lexer.
func (lxr *Lexer) getString() {
	buffer := string(lxr.Code[lxr.Index])
	startingCharacter := buffer
	lxr.Increment()

	for lxr.Index < len(lxr.Code) && string(lxr.Code[lxr.Index]) != startingCharacter {
		buffer += string(lxr.Code[lxr.Index])
		lxr.Increment()
	}
	buffer += string(lxr.Code[lxr.Index])
	lxr.Increment()

	realValueBuffer := buffer[1 : len(buffer)-1]

	lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, realValueBuffer, token.STRING, lxr.GetCurrentTextSpan(len(buffer))))
}

// getId checks if an identifier is a keyword or a regular identifier
// then it generates a token and slaps it back to the lexer.
func (lxr *Lexer) getId() {
	buffer := string(lxr.Code[lxr.Index])
	lxr.Increment()

	for lxr.Index < len(lxr.Code) &&
		(unicode.IsLetter(lxr.Code[lxr.Index]) || unicode.IsNumber(lxr.Code[lxr.Index]) ||
			lxr.Code[lxr.Index] == '_' || lxr.Code[lxr.Index] == '$') {
		buffer += string(lxr.Code[lxr.Index])
		lxr.Increment()
	}

	switch buffer {
	case "fn":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.FN))
	case "return":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.RETURN))
	case "var":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.VAR))
	case "while":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.WHILE))
	case "for":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.FOR))
	case "if":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.IF))
	case "else":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.ELSE))
	case "true":
		lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, true, token.TRUE, lxr.GetCurrentTextSpan(len(buffer))))
	case "false":
		lxr.Tokens = append(lxr.Tokens, token.CreateTokenReal(buffer, false, token.FALSE, lxr.GetCurrentTextSpan(len(buffer))))
	default:
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.IDENT))
	}
}

// getComment skips the comments
func (lxr *Lexer) getComment() {
	lxr.Increment()
	for lxr.Index < len(lxr.Code) && lxr.Code[lxr.Index] != '\n' {
		lxr.Increment()
	}
	lxr.Increment()
}

// getOperator gets all the operators (symbol-like things)
// it also makes sure that they are correctly lexed in case of combination (==, etc)
func (lxr *Lexer) getOperator() {
	buffer := string(lxr.Code[lxr.Index])

	for lxr.Index < len(lxr.Code) && (len(buffer) == 1 || buffer == ":") {
		lxr.Increment()
		if lxr.Index < len(lxr.Code) {
			buffer += string(lxr.Code[lxr.Index])
		}
	}

	switch buffer {
	case "=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.ASSIGN))
	case "+":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.ADD))
	case "-":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SUB))
	case "*":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.MUL))
	case "/":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.QUO))
	case "%":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.REM))
	case "&":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.AND))
	case "|":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.OR))
	case "^":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.XOR))
	case "<<":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SHL))
	case ">>":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SHR))
	case "&^":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.AND_NOT))
	case "&&":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.LAND))
	case "||":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.LOR))
	case "==":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.EQ))
	case "!=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.NOT_EQ))
	case "<":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.LT))
	case "<=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.LEQ))
	case ">":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.GT))
	case ">=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.GEQ))
	case "<=>":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SPACESHIP))
	case "!":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.BANG))
	case "+=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.ADD_ASSIGN))
	case "-=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SUB_ASSIGN))
	case "*=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.MUL_ASSIGN))
	case "/=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.QUO_ASSIGN))
	case "%=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.REM_ASSIGN))
	case "&=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.AND_ASSIGN))
	case "|=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.OR_ASSIGN))
	case "^=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.XOR_ASSIGN))
	case "<<=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SHL_ASSIGN))
	case ">>=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SHR_ASSIGN))
	case "&^=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.AND_NOT_ASSIGN))
	case ",":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.COMMA))
	case "(":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.LPAREN))
	case ")":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.RPAREN))
	case "{":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.LBRACE))
	case "}":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.RBRACE))
	case ".":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.PERIOD))
	case "[":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.LBRACK))
	case "]":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.RBRACK))
	case ";":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.SEMICOLON))
	case ":=":
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.DEFINE))
	default:
		lxr.Tokens = append(lxr.Tokens, token.CreateToken(buffer, token.ILLEGAL))
	}
}

// Increment increments the lexer index, column, and line
func (lxr *Lexer) Increment() {
	if lxr.Code[lxr.Index] == '\n' {
		lxr.Line++
		lxr.Column = 1
	} else {
		lxr.Column++
	}
	lxr.Index++
}

// GetCurrentTextSpan gets the current span of text being processed
func (lxr *Lexer) GetCurrentTextSpan(buffer int) print2.TextSpan {
	return print2.TextSpan{
		File: lxr.File,

		StartIndex: lxr.Index - buffer,
		EndIndex:   lxr.Index,

		StartLine: lxr.Line,
		EndLine:   lxr.Line,

		StartColumn: lxr.Column - buffer,
		EndColumn:   lxr.Column,
	}
}

// RememberSourceFile remembers the source code for the given file
func RememberSourceFile(contents []rune, filename string) {
	// Offload a copy of contents for error handling
	// Also split at new lines because that makes referencing easier
	print2.CodeReference = make([]string, 0)
	print2.CodeReference = strings.Split(string(contents), "\n")
	print2.SourceFiles[filename] = string(contents)
}
