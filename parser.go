package main

import (
	"fmt"
	"strconv"
)

// Precedence levels for operators
const (
	_ int = iota
	LOWEST
	PLINE       // |>
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[TokenType]int{
	EQ:       EQUALS,
	NOT_EQ:   EQUALS,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
	MODULO:   PRODUCT,
	PIPELINE: PLINE,
	LPAREN:   CALL,
	LBRACKET: INDEX,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

// Parser holds the lexer, tokens, and parsing functions.
type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

// New creates a new Parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(INT, p.parseIntegerLiteral)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(TRUE, p.parseBooleanLiteral)
	p.registerPrefix(FALSE, p.parseBooleanLiteral)
	p.registerPrefix(NIL, p.parseNilLiteral)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
	p.registerPrefix(LBRACKET, p.parseListLiteral)
	p.registerPrefix(LBRACE, p.parseMapLiteral)
	p.registerPrefix(IF, p.parseIfExpression)
	p.registerPrefix(FOR, p.parseForExpression)
	p.registerPrefix(BANG, p.parsePrefixExpression)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(FAIL, p.parseFailExpression)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(MODULO, p.parseInfixExpression)
	p.registerInfix(EQ, p.parseInfixExpression)
	p.registerInfix(NOT_EQ, p.parseInfixExpression)
	p.registerInfix(LT, p.parseInfixExpression)
	p.registerInfix(GT, p.parseInfixExpression)
	p.registerInfix(PIPELINE, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)
	p.registerInfix(LBRACKET, p.parseIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram is the entry point for parsing.
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case PROC, CONS, SUPP, EPROC, ESUPP:
		return p.parseFunctionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionStatement() *FunctionStatement {
	stmt := &FunctionStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// proc add(a:int, b:int):int -> a + b
	// cons print_message(msg:str) -> print(msg)
	// supp get_random_num:float -> random()
	switch stmt.Token.Type {
	case PROC, EPROC:
		if !p.expectPeek(LPAREN) {
			return nil
		}
		stmt.Parameters = p.parseFunctionParameters()
		if !p.expectPeek(COLON) {
			return nil
		}
		if !p.expectPeek(IDENT) {
			return nil
		}
		stmt.ReturnType = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case CONS:
		if !p.expectPeek(LPAREN) {
			return nil
		}
		stmt.Parameters = p.parseFunctionParameters()
	case SUPP, ESUPP:
		if !p.expectPeek(COLON) {
			return nil
		}
		if !p.expectPeek(IDENT) {
			return nil
		}
		stmt.ReturnType = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(ARROW) {
		return nil
	}
	p.nextToken()
	stmt.Body = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseFunctionParameters() []*Parameter {
	params := []*Parameter{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken() // Consume LPAREN

	// First parameter
	if p.curToken.Type != IDENT {
		p.peekError(IDENT)
		return nil
	}
	param := &Parameter{}
	param.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(COLON) {
		return nil
	}

	if !p.expectPeek(IDENT) {
		return nil
	}
	param.Type = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	params = append(params, param)

	// Subsequent parameters
	for p.peekTokenIs(COMMA) {
		p.nextToken() // Consume COMMA
		p.nextToken() // Move to the next parameter name

		if p.curToken.Type != IDENT {
			p.peekError(IDENT)
			return nil
		}
		param := &Parameter{}
		param.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(COLON) {
			return nil
		}

		if !p.expectPeek(IDENT) {
			return nil
		}
		param.Type = &Identifier{Token: p.curToken, Value: p.curToken.Literal}
		params = append(params, param)
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// --- Prefix Parsing Functions ---

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expression {
	lit := &FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == TRUE}
}

func (p *Parser) parseNilLiteral() Expression {
	return &NilLiteral{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseListLiteral() Expression {
	lit := &ListLiteral{Token: p.curToken}
	lit.Elements = p.parseExpressionList(RBRACKET)
	return lit
}

func (p *Parser) parseMapLiteral() Expression {
	lit := &MapLiteral{Token: p.curToken}
	lit.Pairs = make(map[Expression]Expression)

	for !p.peekTokenIs(RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		lit.Pairs[key] = value

		if !p.peekTokenIs(RBRACE) && !p.expectPeek(COMMA) {
			return nil
		}
	}

	if !p.expectPeek(RBRACE) {
		return nil
	}

	return lit
}

func (p *Parser) parseIfExpression() Expression {
	expression := &IfExpression{Token: p.curToken}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(THEN) {
		return nil
	}

	p.nextToken()
	expression.Consequence = p.parseExpression(LOWEST)

	if p.peekTokenIs(ELSE) {
		p.nextToken()
		p.nextToken()
		expression.Alternative = p.parseExpression(LOWEST)
	}

	return expression
}

func (p *Parser) parseForExpression() Expression {
	expression := &ForExpression{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	expression.Variable = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(IN) {
		return nil
	}

	p.nextToken()
	expression.Collection = p.parseExpression(LOWEST)

	if !p.expectPeek(THEN) {
		return nil
	}

	p.nextToken()
	expression.Body = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseFailExpression() Expression {
	exp := &FailExpression{Token: p.curToken}

	// The 'fail' keyword must be followed by a string literal.
	if !p.expectPeek(STRING) {
		// If not a string, we can't create a valid FailExpression.
		// The error is already recorded by expectPeek.
		return nil
	}

	// The current token is now the string literal.
	// We can get its value directly.
	exp.Message = p.curToken.Literal

	return exp
}

// --- Infix Parsing Functions ---

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	exp := &IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	list := []Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// --- Helper Methods ---

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
