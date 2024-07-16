package ast

import (
	"fmt"

	"github.com/NikoMalik/Tod-go-compiler/src/print2"
	"github.com/NikoMalik/Tod-go-compiler/src/token"
)

type NodeType string

const (

	// Members
	// -------
	GlobalStatement             NodeType = "Global Statement"
	FunctionDeclaration         NodeType = "Function Declaration"
	ExternalFunctionDeclaration NodeType = "External Function Declaration"
	ClassDeclaration            NodeType = "Class Declaration"
	StructDeclaration           NodeType = "Struct Declaration"

	PackageReference NodeType = "Package Reference"

	PackageUsing NodeType = "Package Using"

	// General
	// -------
	Parameter  NodeType = "Parameter"
	TypeClause NodeType = "Type Clause"

	// Statements
	// ----------
	BlockStatement      NodeType = "Block Statement"
	VariableDeclaration NodeType = "Variable Declaration"
	IfStatement         NodeType = "If Statement"
	ElseClause          NodeType = "Else Clause"
	ReturnStatement     NodeType = "Return Statement"
	ForStatement        NodeType = "For Statement"
	WhileStatement      NodeType = "While Statement"
	BreakStatement      NodeType = "Break Statement"
	ContinueStatement   NodeType = "Continue Statement"
	FromToStatement     NodeType = "FromTo Statement"
	ExpressionStatement NodeType = "Expression Statement"

	// Expressions
	// -----------
	LiteralExpression              NodeType = "Literal Expression"
	ParenthesisedExpression        NodeType = "Parenthesised Expression"
	NameExpression                 NodeType = "Name Expression"
	AssignmentExpression           NodeType = "Assignment Expression"
	CallExpression                 NodeType = "Call Expression"
	PackageCallExpression          NodeType = "PackageCall Expression"
	UnaryExpression                NodeType = "Unary Expression"
	BinaryExpression               NodeType = "Binary Expression"
	VariableEditorExpression       NodeType = "VariableEditor Expression"
	TypeCallExpression             NodeType = "TypeCall Expression"
	ClassFieldAccessExpression     NodeType = "ClassFieldAccess Expression"
	ClassFieldAssignmentExpression NodeType = "ClassFieldAssignment Expression"
	ArrayAccessExpression          NodeType = "ArrayAccess Expression"
	ArrayAssignmentExpression      NodeType = "ArrayAssignment Expression"
	MakeExpression                 NodeType = "Make Expression"
	MakeArrayExpression            NodeType = "MakeArray Expression"
	ReferenceExpression            NodeType = "Reference Expression"
	DereferenceExpression          NodeType = "Dereference Expression "

	MakeStructExpression NodeType = "MakeStruct Expression"

	ThisExpression NodeType = "This Expression"
)

type Node interface {
	NodeType() NodeType
	Span() print2.TextSpan // exact text position of this node
	Print(indent string)
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type MemberNode interface {
	Node
}

// type call

type TypeCallExpressionNode struct {
	Expression
	Base           Expression
	CallIdentifier token.Token
	Arguments      []Expression
	ClosingToken   token.Token
}

func (TypeCallExpressionNode) NodeType() NodeType { return TypeCallExpression }

func (node TypeCallExpressionNode) Span() print2.TextSpan {
	closingSpan := node.ClosingToken.Span
	baseSpan := node.Base.Span()
	return baseSpan.SpanBetween(closingSpan)
}

func (node TypeCallExpressionNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ TypeCallExpressionNode")
	fmt.Println(indent + "  └ Base: ")
	node.Base.Print(indent + "    ")
	fmt.Println(indent + "  └ CallIdentifier: ")

	fmt.Println(indent + "  └ Arguments: ")
	for _, val := range node.Arguments {
		val.Print(indent + "    ")
	}
}

func CreateTypeCallExpressionNode(base Expression, callIdentifier token.Token, args []Expression, closing token.Token) TypeCallExpressionNode {
	return TypeCallExpressionNode{
		Base:           base,
		CallIdentifier: callIdentifier,
		Arguments:      args,
		ClosingToken:   closing,
	}
}

// literal exression node

type LiteralExpressionNode struct {
	Expression
	LiteralToken token.Token
	LiteralValue interface{}
	IsNative     bool
}

func (LiteralExpressionNode) NodeType() NodeType { return LiteralExpression }

func (node LiteralExpressionNode) Span() print2.TextSpan {
	return node.LiteralToken.Span
}

func (node LiteralExpressionNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ LiteralExpressionNode")
	fmt.Printf("%s  └ LiteralToken: %s\n", indent, node.LiteralToken.Type)
}

func CreateLiteralExpressionNode(tok token.Token) LiteralExpressionNode {
	return LiteralExpressionNode{
		LiteralToken: tok,
		LiteralValue: tok.RealValue,
		IsNative:     false,
	}
}

func CreateNativeLiteralExpressionNode(tok token.Token) LiteralExpressionNode {
	return LiteralExpressionNode{
		LiteralToken: tok,
		LiteralValue: tok.RealValue,
		IsNative:     true,
	}
}

// parathesised expression node

type ParanthesisedExpressionNode struct {
	Expression
	OpenParenthesis   token.Token
	ExpressionNode    Expression
	ClosedParenthesis token.Token
}

func (ParanthesisedExpressionNode) NodeType() NodeType { return ParenthesisedExpression }

func (node ParanthesisedExpressionNode) Span() print2.TextSpan {
	return node.OpenParenthesis.Span.SpanBetween(node.ClosedParenthesis.Span)
}

func (node ParanthesisedExpressionNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ ParanthesisedExpressionNode")
	fmt.Printf("%s  └ OpenParenthesis: %s\n", indent, node.OpenParenthesis.Type)
	node.ExpressionNode.Print(indent + "  ")
	fmt.Printf("%s  └ ClosedParenthesis: %s\n", indent, node.ClosedParenthesis.Type)
}

func CreateParanthesisedExpressionNode(expression Expression, open token.Token, close token.Token) ParanthesisedExpressionNode {
	return ParanthesisedExpressionNode{
		ExpressionNode:    expression,
		OpenParenthesis:   open,
		ClosedParenthesis: close,
	}
}

// name expression node

type NameExpressNode struct {
	Expression
	InMain     bool
	Identifier token.Token
}

func (NameExpressNode) NodeType() NodeType { return NameExpression }

func (node NameExpressNode) Span() print2.TextSpan {
	return node.Identifier.Span
}

func (node NameExpressNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ NameExpressNode")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
}

func CreateNameExpressionNode(id token.Token) NameExpressNode {
	return NameExpressNode{
		Identifier: id,
	}
}

func CreateMainNameExpressionNode(id token.Token) NameExpressNode {
	return NameExpressNode{
		Identifier: id,
		InMain:     true,
	}
}

// expression statement

type ExpressionStatementNode struct {
	Statement
	Expression Expression
}

func (ExpressionStatementNode) NodeType() NodeType { return ExpressionStatement }

func (node ExpressionStatementNode) Span() print2.TextSpan {
	return node.Expression.Span()
}

func (node ExpressionStatementNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ ExpressionStatementNode")
	fmt.Println(indent + "  └ Expression: ")
	node.Expression.Print(indent + "    ")
}

func CreateExpressionStatementNode(expression Expression) ExpressionStatementNode {
	return ExpressionStatementNode{
		Expression: expression,
	}
}

//binary expression

type BinaryExpressionNode struct {
	Expression

	Left     Expression
	Operator token.Token
	Right    Expression
}

func (BinaryExpressionNode) NodeType() NodeType { return BinaryExpression }

func (node BinaryExpressionNode) Span() print2.TextSpan {
	return node.Left.Span().SpanBetween(node.Right.Span())
}

func (node BinaryExpressionNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ BinaryExpressionNode")
	fmt.Printf("%s  └ Operator: %s\n", indent, node.Operator.Type)
	fmt.Println(indent + "  └ Left: ")
	node.Left.Print(indent + "    ")
	fmt.Println(indent + "  └ Right: ")
	node.Right.Print(indent + " ")
}

func CreateBinaryExpressionNode(op token.Token, left Expression, right Expression) BinaryExpressionNode {
	return BinaryExpressionNode{
		Left:     left,
		Operator: op,
		Right:    right,
	}
}

// global

type GlobalStatementMember struct {
	MemberNode
	Statement Statement
}

func (GlobalStatementMember) NodeType() NodeType { return GlobalStatement }

func (node GlobalStatementMember) Span() print2.TextSpan {
	return node.Statement.Span()
}

func CreateGlobalStatementMember(stmt Statement) GlobalStatementMember {
	return GlobalStatementMember{
		Statement: stmt,
	}
}

// Type

// basic global member

type TypeClauseNode struct {
	Node
	ClauseIsSet    bool
	Package        *token.Token
	TypeIdentifier token.Token
	SubClauses     []TypeClauseNode
	ClosingBracket token.Token
}

func (TypeClauseNode) NodeType() NodeType { return TypeClause }

func (node TypeClauseNode) Span() print2.TextSpan {
	return node.TypeIdentifier.Span

}

func (node TypeClauseNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ TypeClauseNode")
	fmt.Printf("%s  └ TypeIdentifier: %s\n", indent, node.TypeIdentifier.Type)
	fmt.Printf("%s  └ ClosingBracket: %s\n", indent, node.ClosingBracket.Type)
	if node.ClauseIsSet {
		fmt.Printf("%s  └ SubClauses: \n", indent)
		for _, subClause := range node.SubClauses {
			subClause.Print(indent + "    ")
		}
	}
}

func CreateTypeClauseNode(pack *token.Token, id token.Token, subtypes []TypeClauseNode, bracket token.Token) TypeClauseNode {
	return TypeClauseNode{
		ClauseIsSet:    true,
		Package:        pack,
		TypeIdentifier: id,
		SubClauses:     subtypes,
		ClosingBracket: bracket,
	}
}

// block statement Node

type BlockStatementNode struct {
	Statement
	OpenBrace  token.Token
	Statements []Statement
	CloseBrace token.Token
}

func (BlockStatementNode) NodeType() NodeType { return BlockStatement }

func (node BlockStatementNode) Span() print2.TextSpan {
	return node.OpenBrace.Span.SpanBetween(node.CloseBrace.Span)
}

func (node BlockStatementNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ BlockStatementNode")
	fmt.Printf("%s  └ OpenBrace: %s\n", indent, node.OpenBrace.Type)
	fmt.Printf("%s  └ CloseBrace: %s\n", indent, node.CloseBrace.Type)
	fmt.Println(indent + "  └ Statements: ")

	for _, stmt := range node.Statements {
		stmt.Print(indent + "    ")
	}

}

func CreateBlockStatementNode(openBrace token.Token, statements []Statement, closeBrace token.Token) BlockStatementNode {
	return BlockStatementNode{
		OpenBrace:  openBrace,
		Statements: statements,
		CloseBrace: closeBrace,
	}
}

// function declaration

type FunctionDeclarationMember struct {
	MemberNode
	FunctionKeyword token.Token
	Identifier      token.Token
	Parameters      []ParameterNode
	TypeClause      TypeClauseNode
	Body            BlockStatementNode
	IsPublic        bool
}

func (FunctionDeclarationMember) NodeType() NodeType { return FunctionDeclaration }

func (node FunctionDeclarationMember) Span() print2.TextSpan {
	span := node.FunctionKeyword.Span

	if node.FunctionKeyword.Type == token.FN {
		span = span.SpanBetween(node.Body.Span())
	} else {
		span = span.SpanBetween(node.Identifier.Span)
	}

	return span
}

func (node FunctionDeclarationMember) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- FunctionDeclarationMember")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
	fmt.Printf("%s  └ IsPublic: %t\n", indent, node.IsPublic)
	fmt.Println(indent + "  └ Parameters: ")

	for _, param := range node.Parameters {
		param.Print(indent + "    ")
	}

	if !node.TypeClause.ClauseIsSet {
		fmt.Printf("%s  └ TypeClause: \n", indent)
	} else {
		fmt.Printf("%s  └ TypeClause: %s\n", indent, node.TypeClause.TypeIdentifier.Type)
		node.TypeClause.Print(indent + "    ")
	}

	fmt.Println(indent + "  └ Body: ")

	node.Body.Print(indent + "    ")

}

func CreateFunctionDeclarationMember(kw token.Token, id token.Token, params []ParameterNode, typeClause TypeClauseNode, body BlockStatementNode, public bool) FunctionDeclarationMember {
	return FunctionDeclarationMember{
		FunctionKeyword: kw,
		Identifier:      id,
		Parameters:      params,
		TypeClause:      typeClause,
		Body:            body,
		IsPublic:        public,
	}
}

// external function declaration

type ExternalFunctionDeclarationMember struct {
	MemberNode
	FunctionKeyword token.Token
	Identifier      token.Token
	Parameters      []ParameterNode
	TypeClause      TypeClauseNode
	IsPublic        bool

	ClosingToken token.Token
}

func (ExternalFunctionDeclarationMember) NodeType() NodeType { return ExternalFunctionDeclaration }

func (node ExternalFunctionDeclarationMember) Span() print2.TextSpan {
	span := node.FunctionKeyword.Span.SpanBetween(node.ClosingToken.Span)
	if node.TypeClause.ClauseIsSet {
		span = span.SpanBetween(node.TypeClause.Span())
	}
	return span

}

func (node ExternalFunctionDeclarationMember) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- ExternalFunctionDeclarationMember")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
	fmt.Println(indent + "  └ Parameters: ")

	for _, param := range node.Parameters {
		param.Print(indent + "    ")
	}

	if !node.TypeClause.ClauseIsSet {
		fmt.Printf("%s  └ TypeClause: \n", indent)
	} else {

		node.TypeClause.Print(indent + "    ")
	}

}

func CreateExternalFunctionDeclarationMember(kw token.Token, id token.Token, params []ParameterNode, typeClause TypeClauseNode, closing token.Token) ExternalFunctionDeclarationMember {
	return ExternalFunctionDeclarationMember{
		FunctionKeyword: kw,
		Identifier:      id,
		Parameters:      params,
		TypeClause:      typeClause,
		IsPublic:        true,

		ClosingToken: closing,
	}
}

// struct

type StructDeclarationMember struct {
	MemberNode
	StructKeyword token.Token
	Identifier    token.Token
	Fields        []ParameterNode
	ClosingToken  token.Token
}

func (StructDeclarationMember) NodeType() NodeType { return StructDeclaration }

func (node StructDeclarationMember) Span() print2.TextSpan {
	return node.StructKeyword.Span.SpanBetween(node.ClosingToken.Span)
}

func (node StructDeclarationMember) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- StructDeclarationMember")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
	fmt.Println(indent + "  └ Fields: ")

	for _, param := range node.Fields {
		param.Print(indent + "    ")
	}
}

func CreateStructDeclarationMember(kw token.Token, id token.Token, fields []ParameterNode, closing token.Token) StructDeclarationMember {
	return StructDeclarationMember{
		StructKeyword: kw,
		Identifier:    id,
		Fields:        fields,
		ClosingToken:  closing,
	}
}

// parameters

type ParameterNode struct {
	Node
	Identifier token.Token

	TypeClause TypeClauseNode
}

func (ParameterNode) NodeType() NodeType { return Parameter }

func (node ParameterNode) Span() print2.TextSpan {
	return node.Identifier.Span.SpanBetween(node.TypeClause.Span())
}

func (node ParameterNode) Print(indent string) {
	print2.PrintC(print2.Green, indent+"- ParameterNode")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Literal)

	if !node.TypeClause.ClauseIsSet {
		fmt.Printf("%s  └ TypeClause: none\n", indent)
	} else {
		fmt.Printf("%s  └ TypeClause: ", indent)
		node.TypeClause.Print(indent + "    ")
	}
}

func CreateParameterNode(id token.Token, typeClause TypeClauseNode) ParameterNode {
	return ParameterNode{
		Identifier: id,
		TypeClause: typeClause,
	}

}

// package reference

type PackageReferenceMember struct {
	MemberNode
	PackageKeyword token.Token
	Package        token.Token
}

func (PackageReferenceMember) NodeType() NodeType { return PackageReference }

func (node PackageReferenceMember) Span() print2.TextSpan {
	return node.PackageKeyword.Span.SpanBetween(node.Package.Span)

}

func (node PackageReferenceMember) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- PackageReferenceMember")
	fmt.Printf("%s  └ Package: %s\n", indent, node.Package.Literal)
}

func CreatePackageReferenceMember(kw token.Token, pkg token.Token) PackageReferenceMember {
	return PackageReferenceMember{
		PackageKeyword: kw,
		Package:        pkg,
	}
}

// package use

type PackageUseMember struct {
	MemberNode
	PackageKeyword token.Token
	Package        token.Token
}

func (PackageUseMember) NodeType() NodeType { return PackageUsing }

func (node PackageUseMember) Span() print2.TextSpan {
	return node.PackageKeyword.Span.SpanBetween(node.Package.Span)
}

func (node PackageUseMember) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- PackageUseMember")
	fmt.Printf("%s  └ Package: %s\n", indent, node.Package.Literal)
}

func CreatePackageUseMember(kw token.Token, pkg token.Token) PackageUseMember {
	return PackageUseMember{
		PackageKeyword: kw,
		Package:        pkg,
	}
}

// variable declaration

type VariableDeclarationStatementNode struct {
	Statement

	Keyword    token.Token
	TypeClause TypeClauseNode

	Identifier token.Token

	Initializer Expression
}

func (VariableDeclarationStatementNode) NodeType() NodeType { return VariableDeclaration }

func (node VariableDeclarationStatementNode) Span() print2.TextSpan {
	span := node.Keyword.Span.SpanBetween(node.Identifier.Span)

	if node.Initializer != nil {
		span = span.SpanBetween(node.Initializer.Span())
	}

	return span

}

func (node VariableDeclarationStatementNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- VariableDeclarationStatementNode")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Literal)

	if !node.TypeClause.ClauseIsSet {
		fmt.Printf("%s  └ TypeClause: none\n", indent)
	} else {
		fmt.Printf("%s  └ TypeClause: ", indent)
		node.TypeClause.Print(indent + "    ")
	}

	if node.Initializer != nil {
		fmt.Println(indent + "  └ Initializer: ")
		node.Initializer.Print(indent + "    ")
	}
}

func CreateVariableDeclarationStatementNode(kw token.Token, typeClause TypeClauseNode, id token.Token, initializer Expression) VariableDeclarationStatementNode {
	return VariableDeclarationStatementNode{
		Keyword:     kw,
		TypeClause:  typeClause,
		Identifier:  id,
		Initializer: initializer,
	}
}

// if statement

type IfStatementNode struct {
	Statement

	IfKeyword token.Token

	Condition     Expression
	ThenStatement Statement

	ElseClause ElseClauseNode
}

func (IfStatementNode) NodeType() NodeType { return IfStatement }

func (node IfStatementNode) Span() print2.TextSpan {
	return node.IfKeyword.Span.SpanBetween(node.ThenStatement.Span().SpanBetween(node.ElseClause.Span()))

}

func (node IfStatementNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- IfStatementNode")
	fmt.Printf("%s  └ Condition: \n", indent)
	node.Condition.Print(indent + "    ")
	fmt.Printf("%s  └ ThenStatement: \n", indent)
	node.ThenStatement.Print(indent + "    ")
	if node.ElseClause.ClauseIsSet {
		fmt.Printf("%s  └ ElseClause: \n", indent)
		node.ElseClause.Print(indent + "    ")
	}

}

func CreateIfStatementNode(kw token.Token, condition Expression, then Statement, elseClause ElseClauseNode) IfStatementNode {
	return IfStatementNode{
		IfKeyword:     kw,
		Condition:     condition,
		ThenStatement: then,
		ElseClause:    elseClause,
	}
}

// else ----------------------------------------

type ElseClauseNode struct {
	Node
	ClauseIsSet   bool
	ElseKeyword   token.Token
	ElseStatement Statement
}

func (ElseClauseNode) NodeType() NodeType { return ElseClause }

func (node ElseClauseNode) Span() print2.TextSpan {
	if node.ClauseIsSet {
		return node.ElseKeyword.Span.SpanBetween(node.ElseStatement.Span())
	} else {
		return print2.TextSpan{}
	}
}

func (node ElseClauseNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- ElseClauseNode")

	fmt.Printf("%s  └ ElseKeyword: %s\n", indent, node.ElseKeyword.Literal)

	fmt.Println(indent + "  └ ElseStatement: ")

	node.ElseStatement.Print(indent + "    ")

}

func CreateElseClauseNode(kw token.Token, statement Statement) ElseClauseNode {
	return ElseClauseNode{
		ClauseIsSet:   true,
		ElseKeyword:   kw,
		ElseStatement: statement,
	}
}

// return

type ReturnStatementNode struct {
	Statement
	Keyword    token.Token
	Expression Expression
}

func (ReturnStatementNode) NodeType() NodeType { return ReturnStatement }

func (node ReturnStatementNode) Span() print2.TextSpan {
	return node.Keyword.Span.SpanBetween(node.Expression.Span())

}

func (node ReturnStatementNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- ReturnStatementNode")

	fmt.Printf("%s  └ Keyword: %s\n", indent, node.Keyword.Type)
	fmt.Printf("%s  └ Expression: \n", indent)
	node.Expression.Print(indent + "    ")

	if node.Expression == nil {
		fmt.Printf("%s  └ Expression: none\n", indent)

	} else {
		fmt.Printf("%s  └ Expression: \n", indent)
		node.Expression.Print(indent + "    ")
	}

}

func CreateReturnStatementNode(keyword token.Token, expression Expression) ReturnStatementNode {
	return ReturnStatementNode{
		Keyword:    keyword,
		Expression: expression,
	}
}

//for statement

type ForStatementNode struct {
	StatementNode Statement
	Keyword       token.Token
	Initializer   VariableDeclarationStatementNode
	Condition     Expression
	Updation      Statement
	Statement
}

func (ForStatementNode) NodeType() NodeType { return ForStatement }

func (node ForStatementNode) Span() print2.TextSpan {
	return node.Keyword.Span.SpanBetween(node.Condition.Span())
}

func (node ForStatementNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- ForStatementNode")

	fmt.Printf("%s  └ Keyword: %s\n", indent, node.Keyword.Type)

	fmt.Println(indent + "  └ Initializer: ")

	node.Initializer.Print(indent + "    ")
	fmt.Println(indent + "  └ Condition: ")

	node.Condition.Print(indent + "    ")
	fmt.Println(indent + "  └ Updation: ")

	node.Updation.Print(indent + "    ")
	fmt.Println(indent + "  └ Statement: ")

	node.Statement.Print(indent + "    ")
}

func CreateForStatementNode(keyword token.Token, initializer VariableDeclarationStatementNode, condition Expression, updation Statement, statement Statement) ForStatementNode {
	return ForStatementNode{
		Keyword:     keyword,
		Initializer: initializer,
		Condition:   condition,
		Updation:    updation,
		Statement:   statement,
	}
}

// while statement

type WhileStatementNode struct {
	Statement
	Keyword       token.Token
	Condition     Expression
	StatementNode Statement
}

func (WhileStatementNode) NodeType() NodeType { return WhileStatement }

func (node WhileStatementNode) Span() print2.TextSpan {
	return node.Keyword.Span.SpanBetween(node.Statement.Span())
}

func (node WhileStatementNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- WhileStatementNode")
	fmt.Printf("%s  └ Keyword: %s\n", indent, node.Keyword.Type)

	fmt.Println(indent + "  └ Condition: ")
	node.Condition.Print(indent + "    ")

	fmt.Println(indent + "  └ Statement: ")

	node.StatementNode.Print(indent + "    ")
}

func CreateWhileStatementNode(keyword token.Token, condition Expression, statement Statement) WhileStatementNode {
	return WhileStatementNode{
		Keyword:       keyword,
		Condition:     condition,
		StatementNode: statement,
	}
}

// break

type BreakStatementNode struct {
	Statement
	Keyword token.Token
}

func (BreakStatementNode) NodeType() NodeType { return BreakStatement }

func (node BreakStatementNode) Span() print2.TextSpan {
	return node.Keyword.Span
}

func (node BreakStatementNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- BreakStatemenNode")
	fmt.Printf("%s  └ Keyword: %s\n", indent, node.Keyword.Type)
}

func CreateBreakStatementNode(keyword token.Token) BreakStatementNode {
	return BreakStatementNode{
		Keyword: keyword,
	}
}

// continue

type ContinueStatementNode struct {
	Statement
	Keyword token.Token
}

func (ContinueStatementNode) NodeType() NodeType { return ContinueStatement }

func (node ContinueStatementNode) Span() print2.TextSpan {
	return node.Keyword.Span
}

func (node ContinueStatementNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- ContinueStatemenNode")
	fmt.Printf("%s  └ Keyword: %s\n", indent, node.Keyword.Type)
}

func CreateContinueStatementNode(keyword token.Token) ContinueStatementNode {
	return ContinueStatementNode{
		Keyword: keyword,
	}
}

// assignment expression

type AssignmentExpressionNode struct {
	Expression
	InMain         bool
	Identifier     token.Token
	ExpressionNode Expression
}

func (AssignmentExpressionNode) NodeType() NodeType { return AssignmentExpression }

func (node AssignmentExpressionNode) Span() print2.TextSpan {
	return node.Identifier.Span.SpanBetween(node.ExpressionNode.Span())

}

func (node AssignmentExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- AssignmentExpressionNode")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
	fmt.Println(indent + "  └ Expression: ")
	node.ExpressionNode.Print(indent + "    ")
}

func CreateAssignmentExpressionNode(identifier token.Token, expressionNode Expression) AssignmentExpressionNode {
	return AssignmentExpressionNode{
		Identifier:     identifier,
		ExpressionNode: expressionNode,
		InMain:         true,
	}
}

func CreateMainAssignmentExpressionNode(id token.Token, expr Expression) AssignmentExpressionNode {
	return AssignmentExpressionNode{
		Identifier:     id,
		ExpressionNode: expr,
		InMain:         false,
	}
}

// variable epidor expression

type VariableEditorExpressionNode struct {
	Expression
	Identifier     token.Token
	Operator       token.Token
	IsSingleStep   bool
	ExpressionNode Expression
}

func (VariableEditorExpressionNode) NodeType() NodeType { return VariableEditorExpression }

func (node VariableEditorExpressionNode) Span() print2.TextSpan {
	span := node.Identifier.Span.SpanBetween(node.Operator.Span)
	if !node.IsSingleStep {
		span.SpanBetween(node.ExpressionNode.Span())
	}
	return span
}

func (node VariableEditorExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- VariableEditorExpressionNode")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
	fmt.Printf("%s  └ Operator: %s\n", indent, node.Operator.Type)
	fmt.Println(indent + "  └ Expression: ")
	node.ExpressionNode.Print(indent + "    ")
}

func CreateVariableEditorExpressionNode(identifier token.Token, operator token.Token, expressionNode Expression, isSingleStep bool) VariableEditorExpressionNode {
	return VariableEditorExpressionNode{
		Identifier:     identifier,
		Operator:       operator,
		IsSingleStep:   isSingleStep,
		ExpressionNode: expressionNode,
	}
}

// call expression node

type CallExpressionNode struct {
	Expression
	InMain bool

	Identifier         token.Token
	ClosingParenthesis token.Token

	Arguments   []Expression
	CastingType TypeClauseNode // if this call is actually a complex cast

}

func (CallExpressionNode) NodeType() NodeType { return CallExpression }

func (node CallExpressionNode) Span() print2.TextSpan {
	if node.CastingType.ClauseIsSet {
		return node.CastingType.Span().SpanBetween(node.ClosingParenthesis.Span)
	}
	return node.Identifier.Span.SpanBetween(node.ClosingParenthesis.Span)

}

func (node CallExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- CallExpressionNode")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
	fmt.Println(indent + "  └ Arguments: ")
	for _, arg := range node.Arguments {
		arg.Print(indent + "    ")
	}

}

func CreateCallExpressionNode(identifier token.Token, arguments []Expression, clause TypeClauseNode, closingParenthesis token.Token) CallExpressionNode {
	return CallExpressionNode{
		Identifier:         identifier,
		Arguments:          arguments,
		CastingType:        clause,
		ClosingParenthesis: closingParenthesis,
	}
}

func CreateMainCallExpressionNode(id token.Token, args []Expression, castClause TypeClauseNode, parenthesis token.Token) CallExpressionNode {
	return CallExpressionNode{
		Identifier:         id,
		Arguments:          args,
		CastingType:        castClause,
		ClosingParenthesis: parenthesis,
		InMain:             true,
	}
}

// package expression

type PackageCallExpressionNode struct {
	Expression
	Package      token.Token
	Identifier   token.Token
	Arguments    []Expression
	ClosingToken token.Token
}

func (PackageCallExpressionNode) NodeType() NodeType { return PackageCallExpression }

func (node PackageCallExpressionNode) Span() print2.TextSpan {
	return node.Identifier.Span.SpanBetween(node.ClosingToken.Span)
}

func (node PackageCallExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- PackageCallExpressionNode")
	fmt.Printf("%s  └ Identifier: %s\n", indent, node.Identifier.Type)
	fmt.Println(indent + "  └ Arguments: ")
	for _, arg := range node.Arguments {
		arg.Print(indent + "    ")
	}
}

func CreatePackageCallExpressionNode(pck token.Token, id token.Token, arguments []Expression, closing token.Token) PackageCallExpressionNode {
	return PackageCallExpressionNode{
		Package:      pck,
		Identifier:   id,
		Arguments:    arguments,
		ClosingToken: closing,
	}
}

// make expression node

type MakeArrayExpressionNode struct {
	Expression
	IsLiteral     bool
	MakeKeyword   token.Token
	ClosingToken  token.Token
	Type          TypeClauseNode
	Length        Expression
	LiteralValues []Expression
}

func (MakeArrayExpressionNode) NodeType() NodeType { return MakeArrayExpression }

func (node MakeArrayExpressionNode) Span() print2.TextSpan {
	return node.MakeKeyword.Span.SpanBetween(node.ClosingToken.Span)
}

func (node MakeArrayExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- MakeArrayExpressionNode")
	fmt.Println(indent + "  └ Type: ")
	node.Type.Print(indent + "    ")
	fmt.Println(indent + "  └ Lenght: ")
	node.Length.Print(indent + "    ")
	fmt.Println(indent + "  └ LiteralValues: ")
	for _, val := range node.LiteralValues {
		val.Print(indent + "    ")
	}
}

func CreateMakeArrayExpressionNode(typ TypeClauseNode, length Expression, makeKw token.Token, closing token.Token) MakeArrayExpressionNode {
	return MakeArrayExpressionNode{
		Type:   typ,
		Length: length,

		IsLiteral:    false,
		MakeKeyword:  makeKw,
		ClosingToken: closing,
	}

}

func CreateMakeArrayExpressionNodeLiteral(typ TypeClauseNode, literals []Expression, makeKw token.Token, closing token.Token) MakeArrayExpressionNode {
	return MakeArrayExpressionNode{
		Type:          typ,
		LiteralValues: literals,
		IsLiteral:     true,
		MakeKeyword:   makeKw,
		ClosingToken:  closing,
	}
}

// make struct

type MakeStructExpressionNode struct {
	Expression
	MakeKeyword  token.Token
	ClosingToken token.Token

	Type token.Token

	LiteralValues []Expression
}

func (node MakeStructExpressionNode) NodeType() NodeType { return MakeStructExpression }

func (node MakeStructExpressionNode) Span() print2.TextSpan {
	return node.MakeKeyword.Span.SpanBetween(node.ClosingToken.Span)
}

func (node MakeStructExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- MakeStructExpressionNode")
	fmt.Println(indent + "  └ LiteralValues: ")
	for _, val := range node.LiteralValues {
		val.Print(indent + "    ")
	}

}

func CreateMakeStructExpressionNode(typ token.Token, literals []Expression, makeKw token.Token, closing token.Token) MakeStructExpressionNode {
	return MakeStructExpressionNode{
		Type:          typ,
		LiteralValues: literals,
		MakeKeyword:   makeKw,
		ClosingToken:  closing,
	}
}

// make expression

type MakeExpressionNode struct {
	Expression
	MakeKeyword  token.Token
	ClosingToken token.Token

	Package  *token.Token
	BaseType token.Token

	Arguments []Expression
}

func (node MakeExpressionNode) NodeType() NodeType { return MakeExpression }

func (node MakeExpressionNode) Span() print2.TextSpan {
	return node.MakeKeyword.Span.SpanBetween(node.ClosingToken.Span)
}

func (node MakeExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- MakeExpressionNode")
	fmt.Println(indent + "  └ Arguments: ")
	for _, arg := range node.Arguments {
		arg.Print(indent + "    ")
	}

}

func CreateMakeExpressionNode(pack *token.Token, typ token.Token, args []Expression, makeKw token.Token, closing token.Token) MakeExpressionNode {
	return MakeExpressionNode{
		Package:      pack,
		BaseType:     typ,
		Arguments:    args,
		MakeKeyword:  makeKw,
		ClosingToken: closing,
	}
}

// array assignment expression

type ArrayAssignmentExpressionNode struct {
	Expression
	Base  Expression
	Index Expression
	Value Expression
}

func (ArrayAssignmentExpressionNode) NodeType() NodeType { return ArrayAssignmentExpression }

func (node ArrayAssignmentExpressionNode) Span() print2.TextSpan {
	return node.Base.Span().SpanBetween(node.Value.Span())

}

func (node ArrayAssignmentExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- ArrayAssignmentExpressionNode")
	fmt.Println(indent + "  └ Index: ")
	node.Index.Print(indent + "    ")
	fmt.Println(indent + "  └ Value: ")
	node.Value.Print(indent + "    ")

}

func CreateArrayAssignmentExpressionNode(base Expression, index Expression, value Expression) ArrayAssignmentExpressionNode {
	return ArrayAssignmentExpressionNode{
		Base:  base,
		Index: index,
		Value: value,
	}
}

// array access expression

type ArrayAccessExpressionNode struct {
	Expression

	Base           Expression
	Index          Expression
	ClosingBracket token.Token
}

// implement node type from interface
func (ArrayAccessExpressionNode) NodeType() NodeType { return ArrayAccessExpression }

// Position returns the starting line and column, and the total length of the statement
// The starting line and column aren't always the absolute beginning of the statement just what's most
// convenient.
func (node ArrayAccessExpressionNode) Span() print2.TextSpan {
	return node.Base.Span().SpanBetween(node.ClosingBracket.Span)
}

// node print function
func (node ArrayAccessExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"└ ArrayAccessExpressionNode")
	fmt.Println(indent + "  └ Base: ")
	node.Base.Print(indent + "    ")
	fmt.Println(indent + "  └ Index: ")
	node.Index.Print(indent + "    ")
}

// "constructor" / ooga booga OOP cave man brain
func CreateArrayAccessExpressionNode(base Expression, index Expression) ArrayAccessExpressionNode {
	return ArrayAccessExpressionNode{
		Base:  base,
		Index: index,
	}
}

// reference expression

type ReferenceExpressionNode struct {
	Expression
	ExpressionNode NameExpressNode
	Reference      token.Token
}

func (ReferenceExpressionNode) NodeType() NodeType { return ReferenceExpression }

func (node ReferenceExpressionNode) Span() print2.TextSpan {
	return node.Reference.Span.SpanBetween(node.Expression.Span())
}

func (node ReferenceExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- ReferenceExpressionNode")
	fmt.Println(indent + "  └ Expression: ")
	node.Expression.Print(indent + "    ")
}

func CreateReferenceExpressionNode(kw token.Token, expr NameExpressNode) ReferenceExpressionNode {
	return ReferenceExpressionNode{
		Reference:  kw,
		Expression: expr,
	}
}

// deference expression

type DereferenceExpressionNode struct {
	ExpressionNode Expression

	DerefKeyword token.Token
	Expression
}

// implement node type from interface
func (DereferenceExpressionNode) NodeType() NodeType { return DereferenceExpression }

func (node DereferenceExpressionNode) Span() print2.TextSpan {
	return node.DerefKeyword.Span.SpanBetween(node.Expression.Span())
}

// node print function
func (node DereferenceExpressionNode) Print(indent string) {
	print2.PrintC(print2.Yellow, indent+"└ DereferenceExpressionNode")
	fmt.Println(indent + "  └ Expression: ")
	node.Expression.Print(indent + "    ")
}

// "constructor" / ooga booga OOP cave man brain
func CreateDereferenceExpressionNode(kw token.Token, expr Expression) DereferenceExpressionNode {
	return DereferenceExpressionNode{
		DerefKeyword: kw,
		Expression:   expr,
	}
}

// Unary Expression

type UnaryExpressionNode struct {
	Expression

	Operator token.Token
	Operand  Expression
}

func (UnaryExpressionNode) NodeType() NodeType { return UnaryExpression }

func (node UnaryExpressionNode) Span() print2.TextSpan {
	return node.Operator.Span.SpanBetween(node.Operand.Span())
}

func (node UnaryExpressionNode) Print(indent string) {
	print2.PrintC(print2.Cyan, indent+"- UnaryExpressionNode")
	fmt.Println(indent + "  └ Operator: ")
	node.Operand.Print(indent + "    ")
}

func CreateUnaryExpressionNode(kw token.Token, expr Expression) UnaryExpressionNode {
	return UnaryExpressionNode{
		Operator: kw,
		Operand:  expr,
	}
}
