package parser

import (
	"os"

	"github.com/NikoMalik/Tod-go-compiler/src/ast"
	"github.com/NikoMalik/Tod-go-compiler/src/print2"
	"github.com/NikoMalik/Tod-go-compiler/src/token"
)

type Parser struct {
	Tokens []token.Token
	Index  int
}

func (p *Parser) current() token.Token {
	return p.peek(0)

}

func (p *Parser) peek(offset int) token.Token {
	if p.Index+offset < 0 || p.Index+offset >= len(p.Tokens) {
		return token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	}
	return p.Tokens[p.Index+offset]
}

func (p *Parser) consume(except token.TokenType) token.Token {
	if p.current().Type == except {
		additionalInfo := ""

		if p.current().Type == token.IDENT {
			additionalInfo = "ident may be " + p.current().Literal
		}

		print2.Error(
			"PARSER",
			print2.UnexpectedTokenError,
			p.current().Span,
			"unexpected Token \"%s\"! Expected \"%s\"!"+additionalInfo,
			p.current().Type,
			except,
		)
	}
	p.Index++
	return p.peek(-1)
}

func (p *Parser) rewind(to token.Token) {
	for p.current().String(false) != to.String(false) {
		p.Index--
	}
}

func Parse(tokens []token.Token) []ast.MemberNode {
	parser := Parser{
		Tokens: tokens,
		Index:  0,
	}
	return parser.parseMembers()

}

func (p *Parser) parseMembers() []ast.MemberNode {
	members := make([]ast.MemberNode, 0)

	for p.current().Type != token.EOF {
		startToken := p.current()

		// parse all

		member := p.parseMember(true, true)

		members = append(members, member)

		// if we got stuck

		if startToken == p.current() {
			p.Index++
		}
	}
	return members
}

func (p *Parser) parseMember(allow bool, allowPackages bool) ast.MemberNode {

	if p.current().Type == token.FN {
		return p.parseFunctionDeclaration()
	}

	if p.current().Type == token.PACKAGE && allowPackages {
		return p.parsePackageUse()

	}

	if p.current().Type == token.STRUCT {
		return p.parseStructDeclaration()
	}

	// any stetements outside of a function
	return p.parseGlobalStatement()
}

func (p *Parser) parseGlobalStatement() ast.GlobalStatementMember {
	statement := p.parseStatement()
	return ast.CreateGlobalStatementMember(statement)
}

func (p *Parser) parseFunctionDeclaration() ast.FunctionDeclarationMember {

	isPublic := false
	if p.current().Type == token.SET {
		p.consume(token.SET)
		isPublic = true
	}

	kw := p.consume(token.FN) // fn yo(opa string) string {?????}

	identifier := p.consume(token.IDENT)

	p.consume(token.LPAREN)

	params := p.parseParameterList() // we need only arguments

	p.consume(token.RPAREN)

	typeClause := p.parseOptionalTypeClause()

	body := p.parseBlockStatement()

	return ast.CreateFunctionDeclarationMember(kw, identifier, params, typeClause, body, isPublic)
}

func (p *Parser) parseBlockStatement() ast.BlockStatementNode {
	statements := make([]ast.Statement, 0)

	openBrace := p.consume(token.LBRACE)
	for p.current().Type != token.EOF && p.current().Type != token.RBRACE {
		startToken := p.current()

		statement := p.parseStatement()
		statements = append(statements, statement)

		if startToken == p.current() {
			p.Index++
		}
	}

	closeBrace := p.consume(token.RBRACE)

	return ast.CreateBlockStatementNode(openBrace, statements, closeBrace)
}

func (p *Parser) parseExternalFunctionDeclaration() ast.ExternalFunctionDeclarationMember {
	kw := p.consume(token.EXTERNAL) // external fn yo(opa string) string {?????}

	identifier := p.consume(token.IDENT)

	p.consume(token.LPAREN)
	params := p.parseParameterList() // we need only arguments
	p.consume(token.RPAREN)

	closing := p.consume(token.RPAREN)
	typeClause := p.parseOptionalTypeClause()

	if p.current().Type == token.SEMICOLON {
		p.consume(token.SEMICOLON)
	}

	return ast.CreateExternalFunctionDeclarationMember(kw, identifier, params, typeClause, closing)
}

func (p *Parser) parsePackageReference() ast.PackageReferenceMember {
	kw := p.consume(token.PACKAGE)
	id := p.consume(token.IDENT)

	if p.current().Type == token.SEMICOLON {
		p.consume(token.SEMICOLON)
	}
	return ast.CreatePackageReferenceMember(kw, id)
}

func (p *Parser) parsePackageUse() ast.PackageUseMember {
	kw := p.consume(token.USING)
	id := p.consume(token.IDENT)

	if p.current().Type == token.SEMICOLON {
		p.consume(token.SEMICOLON)
	}
	return ast.CreatePackageUseMember(kw, id)
}

func (p *Parser) parseParameterList() []ast.ParameterNode {
	params := make([]ast.ParameterNode, 0)

	for p.current().Type != token.RPAREN && p.current().Type != token.EOF {
		param := p.parseParameter()

		params = append(params, param)

		if p.current().Type == token.COMMA {
			p.consume(token.COMMA)
		} else {
			break
		}
	}
	return params

}

func (p *Parser) parseStructDeclaration() ast.StructDeclarationMember {
	kw := p.consume(token.STRUCT)
	id := p.consume(token.IDENT)

	// begin struct

	p.consume(token.LBRACE)

	// list of the fields

	fields := make([]ast.ParameterNode, 0)
	for p.current().Type != token.EOF && p.current().Type != token.RBRACE {
		field := p.parseParameter() // name + type

		fields = append(fields, field)

		if p.current().Type != token.EOF && p.current().Type != token.RBRACE {
			p.consume(token.COMMA)
		}

	}

	closing := p.consume(token.RBRACE)

	return ast.CreateStructDeclarationMember(kw, id, fields, closing)
}

func (p *Parser) parseTypeClause() ast.TypeClauseNode {
	var pack *token.Token = nil
	if p.peek(1).Type == token.PACKAGE {
		pck := p.consume(token.IDENT)
		pack = &pck
		p.consume(token.PACKAGE)
	}

	id := p.consume(token.IDENT)
	subTypes := make([]ast.TypeClauseNode, 0)
	var closing token.Token

	if p.current().Type == token.LBRACK {
		p.consume(token.LBRACK)

		for true {
			subTypes = append(subTypes, p.parseTypeClause())

			if p.current().Type != token.COMMA {
				break
			}
			p.consume(token.COMMA)
		}

		closing = p.consume(token.RBRACK)
	}

	return ast.CreateTypeClauseNode(pack, id, subTypes, closing)
}

func (p *Parser) parseParameter() ast.ParameterNode {
	identifier := p.consume(token.IDENT)
	typeClause := p.parseTypeClause()

	return ast.CreateParameterNode(identifier, typeClause)

}

func (p *Parser) parseOptionalTypeClause() ast.TypeClauseNode {
	if p.current().Type != token.IDENT {
		return ast.TypeClauseNode{}
	}

	return p.parseTypeClause()
}

//---------------------------------------------------------------

func (p *Parser) parseStatement() ast.Statement {
	var statement ast.Statement = nil

	cur := p.current().Type
	// { ... }

	if cur == token.VAR || cur == token.SET {
		statement = p.parseVariableDeclaration()

	} else if cur == token.LBRACE {
		statement = p.parseBlockStatement()

	} else if cur == token.IF {
		statement = p.parseIfStatement()
	} else if cur == token.RETURN {
		statement = p.parseReturnStatement()
	} else if cur == token.FOR {
		statement = p.parseForStatement()
	} else if cur == token.BREAK {
		statement = p.parseBreakStatement()
	} else if cur == token.CONTINUE {
		statement = p.parseContinueStatement()

	} else if cur == token.WHILE {
		statement = p.parseWhileStatement()
	} else if cur == token.FN {
		statement = p.parseForStatement()
	} else if cur == token.ELSE {
		statement = p.parseElseClause()
	} else {
		statement = p.parseExpressionStatement()
	}

	if p.current().Type == token.SEMICOLON {
		p.consume(token.SEMICOLON)
	}
	return statement
}

func (p *Parser) parseVariableDeclaration() ast.VariableDeclarationStatementNode {
	keyword := p.consume(p.current().Type)

	typeClause := ast.TypeClauseNode{}

	if p.current().Type == token.IDENT &&
		(p.peek(1).Type == token.IDENT || p.peek(1).Type == token.LBRACK) {
		typeClause = p.parseTypeClause()
	}

	identifier := p.consume(token.IDENT)

	if p.current().Type == token.SEMICOLON {
		p.consume(token.SEMICOLON)

		initializer := p.parseExpression()

		return ast.CreateVariableDeclarationStatementNode(keyword, typeClause, identifier, initializer)

	} else {
		return ast.CreateVariableDeclarationStatementNode(keyword, typeClause, identifier, nil)
	}
}

func (p *Parser) parseExpression() ast.Expression {
	if p.current().Type == token.TokenType(token.IDENT) && p.peek(1).Type == token.TokenType(token.ADD_ASSIGN) && !p.peek(1).SpaceAfter && token.GetBinaryOperatorPrecedence(p.peek(1)) != 0 {
		return p.parseVariableEditorExpression()
	}

	if p.current().Type == token.IDENT && ((p.peek(1).Type == token.ADD) || p.peek(2).Type == token.ADD) ||
		(p.peek(1).Type == token.SUB || p.peek(2).Type == token.SUB) {
		identifier := p.consume(token.IDENT)
		operator := p.consume(p.current().Type)
		p.consume(p.current().Type)
		return ast.CreateVariableEditorExpressionNode(identifier, operator, nil, true)
	}

	if p.current().Type == token.IDENT && p.peek(1).Type == token.ADD_ASSIGN {
		return p.parseAssignmentExpression()
	}

	// Add more conditions as needed
	return p.parseBinaryExpression(0)
}

func (p *Parser) parseAssignmentExpression() ast.Expression {
	identifier := p.consume(token.IDENT)

	p.consume(token.ASSIGN)

	value := p.parseExpression()

	return ast.CreateAssignmentExpressionNode(identifier, value)

}

func (p *Parser) parseBinaryExpression(parentPrecedence int) ast.Expression {
	var left ast.Expression

	// check if this is a unary expression

	unaryPrecedence := token.GetUnaryOperatorPrecedence(p.current())

	if unaryPrecedence != 0 && unaryPrecedence > parentPrecedence {
		operator := p.consume(p.current().Type)
		operand := p.parseBinaryExpression(unaryPrecedence)
		return ast.CreateUnaryExpressionNode(operator, operand)

		// if not, start by parsing our left expression

	} else {
		left = p.parsePrimaryExpression()

		if p.current().Type == token.LBRACK {
			left = p.parseArrayAccessExpressionFromValue(left)
		}

		for {
			precedence := token.GetBinaryOperatorPrecedence(p.current())

			if precedence == 0 || precedence <= parentPrecedence {
				break
			}

			operator := p.consume(p.current().Type)

			right := p.parseBinaryExpression(precedence)

			left = ast.CreateBinaryExpressionNode(operator, left, right)
		}
	}

	return left

}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	cur := p.current().Type

	if cur == token.STRING {
		return p.parseStringLiteral()
	} else if cur == token.INT || cur == token.UINT || cur == token.FLOAT32 || cur == token.FLOAT64 {
		return p.parseNumberLiteral()
	} else if cur == token.TRUE || cur == token.FALSE {
		return p.parseBoolLiteral()
	} else if cur == token.LPAREN {
		return p.parseParanthesisedExpression()

	} else if p.current().Type == token.IDENT &&
		p.peek(1).Type == token.ASSIGN {
		return p.parseTypeCallExpression()
	} else if cur == token.IDENT {
		return p.parseNameOrCallExpression()
	} else if p.current().Type == token.POINTER {
		return p.parseReferenceExpression()

	} else if cur == token.ADDRESS {
		return p.parseDereferenceExpression()
	} else if cur == token.MAIN {
		return p.parseMainExpression()
	}

	additionalInfo := ""

	print2.Error(
		"PARSER",
		print2.UnexpectedTokenError,
		p.current().Span,
		"unexpected token: "+p.current().Type.String()+"\n"+additionalInfo,
		p.current().Type,
	)
	os.Exit(1)

	return nil

}

func (p *Parser) parseReferenceExpression() ast.ReferenceExpressionNode {
	keyword := p.consume(token.POINTER)
	expression := p.parseNameExpression()

	return ast.CreateReferenceExpressionNode(keyword, expression)
}

func (p *Parser) parseDereferenceExpression() ast.DereferenceExpressionNode {
	keyword := p.consume(token.ADDRESS)
	expression := p.parsePrimaryExpression()

	return ast.CreateDereferenceExpressionNode(keyword, expression)
}

func (p *Parser) parseIfStatement() ast.IfStatementNode {
	// if ( ... ) { ... }

	keyword := p.consume(token.IF)

	p.consume(token.LPAREN)

	condition := p.parseExpression()

	p.consume(token.RPAREN)

	statement := p.parseStatement()

	elseClause := p.parseElseClause()

	return ast.CreateIfStatementNode(keyword, condition, statement, elseClause)

}

func (p *Parser) parseNameOrCallExpression() ast.Expression {
	if p.peek(1).Type == token.LPAREN {
		return p.parseCallExpression()
	}
	if p.peek(1).Type == token.LBRACK {
		return p.parseComplexCastExpression()

	}

	if p.peek(1).Type == token.USING && p.peek(2).Type == token.PACKAGE {
		return p.parsePackageCallExpression()
	}

	return p.parseNameExpression()
}

func (p *Parser) parseElseClause() ast.ElseClauseNode {
	if p.current().Type != token.ELSE {
		return ast.ElseClauseNode{}
	}

	keyword := p.consume(token.ELSE)

	statement := p.parseStatement()

	return ast.CreateElseClauseNode(keyword, statement)

}

func (p *Parser) parseReturnStatement() ast.ReturnStatementNode {

	keyword := p.consume(token.RETURN)

	var expression ast.Expression = nil

	if p.current().Type != token.SEMICOLON {
		expression = p.parseExpression()
	}

	return ast.CreateReturnStatementNode(keyword, expression)

}

func (p *Parser) parseForStatement() ast.ForStatementNode {
	keyword := p.consume(token.FOR)

	p.consume(token.LPAREN)

	initializer := p.parseVariableDeclaration()

	p.consume(token.SEMICOLON)

	condition := p.parseExpression()

	p.consume(token.SEMICOLON)

	updation := p.parseStatement()

	p.consume(token.RPAREN)

	statement := p.parseStatement()

	return ast.CreateForStatementNode(keyword, initializer, condition, updation, statement)

}

func (p *Parser) parseWhileStatement() ast.WhileStatementNode {

	keyword := p.consume(token.WHILE)

	p.consume(token.LPAREN) // (

	condition := p.parseExpression()

	p.consume(token.RPAREN) // )

	statement := p.parseStatement()

	return ast.CreateWhileStatementNode(keyword, condition, statement)

}

func (p *Parser) parseBreakStatement() ast.BreakStatementNode {

	keyword := p.consume(token.BREAK)

	return ast.CreateBreakStatementNode(keyword)

}

func (p *Parser) parseContinueStatement() ast.ContinueStatementNode {
	keyword := p.consume(token.CONTINUE)

	return ast.CreateContinueStatementNode(keyword)
}

func (p *Parser) parseExpressionStatement() ast.ExpressionStatementNode {
	expression := p.parseExpression()

	return ast.CreateExpressionStatementNode(expression)
}

// -----------------------------------------------------

func (p *Parser) parseAssighmentExpression() ast.AssignmentExpressionNode {
	identifier := p.consume(token.IDENT)
	p.consume(token.ASSIGN)
	value := p.parseExpression() // new value of variable
	return ast.CreateAssignmentExpressionNode(identifier, value)

}

func (p *Parser) parseVariableEditorExpression() ast.VariableEditorExpressionNode {
	identifier := p.consume(token.IDENT)
	p.consume(token.ASSIGN)
	operator := p.consume(p.current().Type)
	expression := p.parseExpression() // new value of variable
	return ast.CreateVariableEditorExpressionNode(identifier, operator, expression, false)
}

func (p *Parser) parseCallExpression() ast.CallExpressionNode {
	identifier := p.consume(token.IDENT)

	p.consume(token.LPAREN)    // (
	args := p.parseArguments() //  we get arguments

	closing := p.consume(token.RPAREN) // )

	return ast.CreateCallExpressionNode(identifier, args, ast.TypeClauseNode{}, closing)
}

// for example complexType[string, int](someAny)
func (p *Parser) parseComplexCastExpression() ast.Expression {
	// We store the identifier, so we can rewind in case we need to
	identifier := p.current()
	typeClause, ok := p.parseUncertainTypeClause()

	if !ok { // rewind this is actually an array access
		p.rewind(identifier)
		return p.parseArrayAccessExpression()
	}
	// if there's no open parenthesis for the cast also rewind
	if p.current().Type != token.LPAREN {
		p.rewind(identifier)
		return p.parseArrayAccessExpression()

	}

	p.consume(token.LPAREN)           // (
	expression := p.parseExpression() // We get the expression we want to cast

	closing := p.consume(token.RPAREN) // )
	return ast.CreateCallExpressionNode(identifier, []ast.Expression{expression}, typeClause, closing)
}

func (p *Parser) parseUncertainTypeClause() (ast.TypeClauseNode, bool) {
	var pack *token.Token = nil

	if p.peek(1).Type == token.PACKAGE {
		pck := p.consume(token.IDENT)
		pack = &pck
		p.consume(token.PACKAGE)
	}

	identifier := p.consume(token.IDENT)

	subTypes := make([]ast.TypeClauseNode, 0)

	var closing token.Token

	if p.current().Type == token.LBRACK {
		p.consume(token.LBRACK)

		for {

			if p.current().Type != token.IDENT {
				return ast.TypeClauseNode{}, false

			}

			subClause, ok := p.parseUncertainTypeClause()

			if !ok {
				return ast.TypeClauseNode{}, false
			}

			subTypes = append(subTypes, subClause)

			if p.current().Type != token.COMMA {
				break
			}
			p.consume(token.COMMA)
		}

		if p.current().Type != token.RBRACK {
			return ast.TypeClauseNode{}, false
		}

		closing = p.consume(token.RBRACK)

	}

	return ast.CreateTypeClauseNode(pack, identifier, subTypes, closing), true
}

func (p *Parser) parsePackageCallExpression() ast.PackageCallExpressionNode {
	// we need the identifier to know which package to select
	pack := p.consume(token.IDENT)
	p.consume(token.USING)

	identifier := p.consume(token.IDENT)

	p.consume(token.LPAREN)            // (
	args := p.parseArguments()         //  we get arguments
	closing := p.consume(token.RPAREN) // )

	return ast.CreatePackageCallExpressionNode(pack, identifier, args, closing)

}

func (p *Parser) parseArrayAccessExpression() ast.Expression {
	base := p.parseNameExpression()

	return p.parseArrayAccessExpressionFromValue(base)

}

func (p *Parser) parseArrayAccessExpressionFromValue(base ast.Expression) ast.Expression {
	// we need the identifier to know which package to select
	p.consume(token.LBRACK)      // [
	index := p.parseExpression() //  we get arguments
	p.consume(token.RBRACK)      // ]

	if p.current().Type == token.ASSIGN {
		p.consume(token.ASSIGN)
		value := p.parseExpression()
		return ast.CreateArrayAssignmentExpressionNode(base, index, value)
	}

	return ast.CreateArrayAccessExpressionNode(base, index)

}

func (p *Parser) parseMakeExpression() ast.Expression {
	makeKeyword := p.consume(token.MAKE) // make

	var pack *token.Token = nil
	if p.peek(1).Type == token.PACKAGE {
		pck := p.consume(token.IDENT)
		pack = &pck
		p.consume(token.PACKAGE)
	}

	baseType := p.consume(token.IDENT)

	if p.current().Type != token.LPAREN && p.current().Type != token.LBRACE {
		p.rewind(makeKeyword)
		return p.parseMakeArrayExpression()
	}

	if p.current().Type == token.LBRACE {
		return p.parseMakeStructExpression(makeKeyword, baseType)
	}

	p.consume(token.LPAREN)    // (
	args := p.parseArguments() //  we get arguments

	closing := p.consume(token.RPAREN) // )

	return ast.CreateMakeExpressionNode(pack, baseType, args, makeKeyword, closing)

}

func (p *Parser) parseMakeStructExpression(makeKeyword token.Token, baseType token.Token) ast.Expression {
	literals := make([]ast.Expression, 0)

	p.consume(token.LBRACE) // {

	for p.current().Type != token.RBRACE &&
		p.current().Type != token.EOF {
		expression := p.parseExpression()
		literals = append(literals, expression)

		if p.current().Type == token.COMMA {
			p.consume(token.COMMA)
		} else {
			break
		}
	}

	closing := p.consume(token.RBRACE) // }

	return ast.CreateMakeStructExpressionNode(baseType, literals, makeKeyword, closing)
}

func (p *Parser) parseMakeArrayExpression() ast.MakeArrayExpressionNode {
	keyword := p.consume(token.MAKE) // make

	baseType := p.parseTypeClause()

	p.consume(token.IDENT)

	if p.current().Type == token.LBRACE {

		literals := make([]ast.Expression, 0)

		p.consume(token.LBRACE) // {

		for p.current().Type != token.RBRACE &&
			p.current().Type != token.EOF {
			expression := p.parseExpression()
			literals = append(literals, expression)

			if p.current().Type == token.COMMA {
				p.consume(token.COMMA)
			} else {
				break
			}
		}

		closing := p.consume(token.RBRACE) // }
		return ast.CreateMakeArrayExpressionNodeLiteral(baseType, literals, keyword, closing)
	}

	p.consume(token.LPAREN)       // (
	length := p.parseExpression() //  we get arguments

	closing := p.consume(token.RPAREN) // )

	return ast.CreateMakeArrayExpressionNode(baseType, length, keyword, closing)

}

func (p *Parser) parseMainExpression() ast.Expression {
	p.consume(token.MAIN) // main

	callIdentifier := p.consume(token.IDENT)

	if p.current().Type != token.LPAREN {
		return p.parseMainAccessExpression(callIdentifier)
	}

	p.consume(token.LPAREN) // (

	args := p.parseArguments() //  we get arguments

	closing := p.consume(token.RPAREN) // )

	return ast.CreateMainCallExpressionNode(callIdentifier, args, ast.TypeClauseNode{}, closing)
}

func (p *Parser) parseMainAccessExpression(id token.Token) ast.Expression {
	if p.current().Type == token.ASSIGN {
		return p.parseMainAssignmentExpression(id)
	}

	return ast.CreateMainNameExpressionNode(id)
}

func (p *Parser) parseMainAssignmentExpression(id token.Token) ast.Expression {
	p.consume(token.ASSIGN)      // =
	value := p.parseExpression() // the value

	return ast.CreateMainAssignmentExpressionNode(id, value)
}

func (p *Parser) parseTypeCallExpression() ast.Expression {
	base := p.parseNameExpression()

	return p.parseTypeCallExpressionFromValue(base)

}

func (p *Parser) parseNameExpression() ast.NameExpressNode {

	identifier := p.consume(token.IDENT)

	return ast.CreateNameExpressionNode(identifier)
}

func (p *Parser) parseTypeCallExpressionFromValue(expr ast.Expression) ast.Expression {

	p.consume(token.ASSIGN)

	callIdentifier := p.consume(token.IDENT)

	p.consume(token.LPAREN)            // (
	args := p.parseArguments()         //  we get arguments
	closing := p.consume(token.RPAREN) // )

	return ast.CreateTypeCallExpressionNode(expr, callIdentifier, args, closing)
}

func (p *Parser) parseParanthesisedExpression() ast.ParanthesisedExpressionNode {

	opening := p.consume(token.LPAREN) // (

	expression := p.parseExpression()

	closing := p.consume(token.RPAREN) // )

	return ast.CreateParanthesisedExpressionNode(expression, opening, closing)
}

func (p *Parser) parseArguments() []ast.Expression {
	args := make([]ast.Expression, 0)

	for p.current().Type != token.RPAREN &&
		p.current().Type != token.EOF {
		expression := p.parseExpression()
		args = append(args, expression)
		if p.current().Type == token.COMMA {
			p.consume(token.COMMA)
		} else {
			break
		}
	}
	return args

}

func (p *Parser) parseStringLiteral() ast.LiteralExpressionNode {

	str := p.consume(token.STRING)
	return ast.CreateLiteralExpressionNode(str)
}

func (p *Parser) parseNumberLiteral() ast.LiteralExpressionNode {
	if p.current().Type == token.INT {
		integer := p.consume(token.INT)
		return ast.CreateLiteralExpressionNode(integer)
	} else if p.current().Type == token.FLOAT32 {
		float := p.consume(token.FLOAT64)
		return ast.CreateLiteralExpressionNode(float)
	} else if p.current().Type == token.FLOAT64 {
		float := p.consume(token.FLOAT64)
		return ast.CreateLiteralExpressionNode(float)
	} else if p.current().Type == token.UINT {
		char := p.consume(token.UINT)
		return ast.CreateLiteralExpressionNode(char)
	}

	return ast.CreateLiteralExpressionNode(token.Token{})
}

func (p *Parser) parseBoolLiteral() ast.LiteralExpressionNode {
	_bool := p.consume(p.current().Type)

	return ast.CreateLiteralExpressionNode(_bool)
}
