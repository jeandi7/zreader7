package zinterpreter

// zanzibar restricted BNF grammar

/**

<Zschema> ::= <Zdef>*
<Zdef> ::= "definition" <Zname> "{" <Zbody> "}"
<Zname> ::= <identifier>
<Zbody> ::= (<Zrelation> | <Zpermission>)* // * means zero or more <Zrelation> or <Zpermission>
<Zrelation> ::= "relation" <Rname> ":" <Sname> ("|" <Sname>)*
<Zpermission> ::= "permission" <Rname> "=" <RPexpr>
<RPexpr> ::= <RPterm> <RPexpr1>
<RPexpr1> ::= <Zop> <RPterm> <RPexpr1> | ''   // '' means RPexpr1 can be empty
<RPterm> ::= <Rname> | "(" <RPexpr> ")" | <Rname> "." "any" "(" <Rname> ")" | <Rname> "." "all" "(" <Rname> ")"
<Rname> ::= <identifier>
<Sname> ::= <Zname> | <Zname> "#" <Rname> | <Zname> ":" "*"
<Zop> ::= "+" | "&" | "-" | "->"     // I handle operator levels as the same priority.
<identifier> ::= [a-zA-Z_][a-zA-Z0-9_]*

*/

import (
	"fmt"
	"strings"
	"unicode"
)

// Token represents the different tokens
type Token int

const (
	DefinitionToken       Token = iota // "definition"
	RelationToken                      // "relation"
	PermissionToken                    // "permission"
	AllToken                           // "all"
	AnyToken                           // "any"
	ArrowToken                         // "->"
	ColonToken                         // ":"
	OrToken                            // "|"
	LeftBraceToken                     // "{"
	RightBraceToken                    // "}"
	LeftParenthesisToken               // "("
	RightParenthesisToken              // ")"
	HashToken                          // "#"
	DotToken                           // "."
	EqualToken                         // "="
	UnionToken                         // "+"
	IntersectionToken                  // "&"
	ExclusionToken                     // "-"
	IdentifierToken                    // [a-zA-Z_][a-zA-Z0-9_]*
	WildCardToken                      // *
	EOFToken                           // ''
	InvalidToken                       //
)

// Item représente un token avec sa valeur
type Item struct {
	Token Token
	Value string
}

// Lexer parses input text and generates tokens
type Lexer struct {
	input       string
	pos         int
	length      int
	currentItem *Item
}

// for Lexer message
func TokenToString(t Token) string {
	switch t {
	case DefinitionToken:
		return "definition"
	case RelationToken:
		return "relation"
	case PermissionToken:
		return "permission"
	case AnyToken:
		return "any"
	case AllToken:
		return "all"
	case ArrowToken:
		return "->"
	case ColonToken:
		return ":"
	case OrToken:
		return "|"
	case LeftBraceToken:
		return "{"
	case RightBraceToken:
		return "}"
	case LeftParenthesisToken:
		return "("
	case RightParenthesisToken:
		return ")"
	case UnionToken:
		return "+"
	case IntersectionToken:
		return "&"
	case ExclusionToken:
		return "-"
	case HashToken:
		return "#"
	case IdentifierToken:
		return "Identifier"
	case WildCardToken:
		return "*"
	case DotToken:
		return "."
	case EqualToken:
		return "="
	case EOFToken:
		return ""
	case InvalidToken:
		return "invalid"
	default:
		return "unknown"
	}
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		length: len(input),
		currentItem: &Item{
			Token: InvalidToken,
			Value: "",
		},
	}
}

// We eat up the white spaces and // lines

func (l *Lexer) eatSpace() {
	for l.pos < l.length {
		// Check the spaces
		if unicode.IsSpace(rune(l.input[l.pos])) {
			l.pos++
			continue
		}
		// Check the comments starting with '//'
		if l.pos+1 < l.length && l.input[l.pos:l.pos+2] == "//" {
			l.pos += 2 // Move forward after '//'
			// Move to the end of the line or EOF
			for l.pos < l.length && l.input[l.pos] != '\n' {
				l.pos++
			}
			// Move forward after '\n' if '\n' is present
			if l.pos < l.length && l.input[l.pos] == '\n' {
				l.pos++
			}
			continue
		}
		// If it is neither a space nor a comment, exit.
		break
	}
}

// Lexer returns the next token to read
func (l *Lexer) NextToken() *Item {
	l.eatSpace()

	if l.pos >= l.length {
		l.currentItem.Token = EOFToken
		l.currentItem.Value = ""
		return l.currentItem
	}

	switch {
	case strings.HasPrefix(l.input[l.pos:], "definition"):
		l.currentItem.Token = DefinitionToken
		l.currentItem.Value = "definition"
		l.pos += len("definition")
	case strings.HasPrefix(l.input[l.pos:], "relation"):
		l.currentItem.Token = RelationToken
		l.currentItem.Value = "relation"
		l.pos += len("relation")
	case strings.HasPrefix(l.input[l.pos:], "permission"):
		l.currentItem.Token = PermissionToken
		l.currentItem.Value = "permission"
		l.pos += len("permission")
	case strings.HasPrefix(l.input[l.pos:], "->"):
		l.currentItem.Token = ArrowToken
		l.currentItem.Value = "->"
		l.pos += len("->")
	case strings.HasPrefix(l.input[l.pos:], "any"):
		l.currentItem.Token = AnyToken
		l.currentItem.Value = "any"
		l.pos += len("any")
	case strings.HasPrefix(l.input[l.pos:], "all"):
		l.currentItem.Token = AllToken
		l.currentItem.Value = "all"
		l.pos += len("all")
	case l.input[l.pos] == ':':
		l.currentItem.Token = ColonToken
		l.currentItem.Value = ":"
		l.pos++
	case l.input[l.pos] == '|':
		l.currentItem.Token = OrToken
		l.currentItem.Value = "|"
		l.pos++
	case l.input[l.pos] == '{':
		l.currentItem.Token = LeftBraceToken
		l.currentItem.Value = "{"
		l.pos++
	case l.input[l.pos] == '}':
		l.currentItem.Token = RightBraceToken
		l.currentItem.Value = "}"
		l.pos++
	case l.input[l.pos] == '#':
		l.currentItem.Token = HashToken
		l.currentItem.Value = "#"
		l.pos++
	case l.input[l.pos] == '*':
		l.currentItem.Token = WildCardToken
		l.currentItem.Value = "*"
		l.pos++
	case l.input[l.pos] == '=':
		l.currentItem.Token = EqualToken
		l.currentItem.Value = "="
		l.pos++
	case l.input[l.pos] == '+':
		l.currentItem.Token = UnionToken
		l.currentItem.Value = "+"
		l.pos++
	case l.input[l.pos] == '-':
		l.currentItem.Token = ExclusionToken
		l.currentItem.Value = "-"
		l.pos++
	case l.input[l.pos] == '&':
		l.currentItem.Token = IntersectionToken
		l.currentItem.Value = "&"
		l.pos++
	case l.input[l.pos] == '(':
		l.currentItem.Token = LeftParenthesisToken
		l.currentItem.Value = "("
		l.pos++
	case l.input[l.pos] == ')':
		l.currentItem.Token = RightParenthesisToken
		l.currentItem.Value = ")"
		l.pos++
	case l.input[l.pos] == '.':
		l.currentItem.Token = DotToken
		l.currentItem.Value = "."
		l.pos++
	default:
		if unicode.IsLetter(rune(l.input[l.pos])) {
			start := l.pos
			for l.pos < l.length && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '_') {
				l.pos++
			}
			l.currentItem.Token = IdentifierToken
			l.currentItem.Value = l.input[start:l.pos]
		} else {
			l.currentItem.Token = InvalidToken
			l.currentItem.Value = string(l.input[l.pos])
			l.pos++
		}
	}
	return l.currentItem
}

func (l *Lexer) readAndMatchToken(expected Token) error {
	if l.currentItem.Token == expected {
		return nil
	}
	return fmt.Errorf("expected token '%v', but got '%v'", TokenToString(expected), l.currentItem.Value)
}

func (l *Lexer) readAndMatchTokenList(expected ...Token) error {
	for _, exp := range expected {
		if l.currentItem.Token == exp {
			return nil
		}
	}

	var expectedStrings []string
	for _, exp := range expected {
		expectedStrings = append(expectedStrings, TokenToString(exp))
	}
	return fmt.Errorf("expected one of tokens %v, but got '%v'", expectedStrings, l.currentItem.Value)
}

// Syntaxic Analyser
type ZDef struct {
	Name        string
	Relations   []*ZRelation
	Permissions []*ZPermission
	ID          string
}

// <Zrelation> ::= "relation" <Rname> ":" <Sname> ("|" <Sname>)*
// <Sname> ::= <Zname> | <Zname> "#" <Rname> | <Zname> ":" "*"

type ZRelation struct {
	Name             string
	Zobjects         []*Zobject
	ZobjectSets      []*ZobjectSet
	ZobjectWildCards []*ZobjectWildCard
	ID               string
	myZDef           *ZDef
}

// object
type Zobject struct {
	Name   string
	ID     string
	myZDef *ZDef
	Unique bool
}

// object#relation
type ZobjectSet struct {
	Name         string
	Relation     string // can be a permission
	ID           string
	IDRelation   string
	Unique       bool
	IsPermission bool
}

// object:*
type ZobjectWildCard struct {
	Name   string
	ID     string
	Unique bool
}

// ZPermission
type ZPermission struct {
	Name   string
	ID     string
	RPexpr *RPexpr
	myZDef *ZDef
}

// <RPexpr> ::= <RPterm> (<Zop> <RPterm>)*
// <RPexpr> ::= <RPterm> <RPexpr1>

type RPexpr struct {
	Term    *RPterm
	RPexpr1 *RPexpr1
}

// <RPexpr1> ::= <Zop> <RPterm> <RPexpr1> | ”
type RPexpr1 struct {
	Zop     Token // +, &, -,->
	RPTerm  *RPterm
	RPexpr1 *RPexpr1
	ID      string // for drawing it
}

type RPterm struct {
	Name       string // relation/permission name if single
	Relation   *ZRelation
	Permission *ZPermission
}

// object#relation
// <Zschema> ::= <Zdef>*
func (l *Lexer) ReadZSchema() ([]*ZDef, error) {
	var zdefs []*ZDef

	for l.currentItem.Token != EOFToken {
		_zdef, _err := l.readZDef()
		if _err != nil {
			return zdefs, _err
		}
		zdefs = append(zdefs, &_zdef)
		l.NextToken()

	}
	return zdefs, nil
}

// <Zdef> ::= "definition" <Zname> "{" <Zbody> "}"
func (l *Lexer) readZDef() (ZDef, error) {
	var zdef ZDef

	// read "definition"
	err := l.readAndMatchToken(DefinitionToken)
	if err != nil {
		return zdef, err
	}
	l.NextToken()

	// read <Zname>
	err = l.readAndMatchToken(IdentifierToken)
	if err != nil {
		return zdef, err
	}
	zdef.Name = l.currentItem.Value
	l.NextToken()

	// read '{'
	err = l.readAndMatchToken(LeftBraceToken)
	if err != nil {
		return zdef, err
	}
	l.NextToken()

	// read ZBody
	// ZBody is not a token
	// no need to call NextToken after

	zdef, err = l.readZBody(zdef)
	if err != nil {
		return zdef, err
	}

	// read '}'
	err = l.readAndMatchToken(RightBraceToken)
	if err != nil {
		return zdef, err
	}

	return zdef, nil
}

// <Zbody> ::= (<Zrelation>|<Zpermission>)*
// * means zero or more <Zrelation> or <Zpermission>

func (l *Lexer) readZBody(zdef ZDef) (ZDef, error) {
	for l.currentItem.Token == RelationToken || l.currentItem.Token == PermissionToken {
		if l.currentItem.Token == RelationToken {
			relation, err := l.readZRelation()
			if err != nil {
				return zdef, err
			}
			zdef.Relations = append(zdef.Relations, &relation)
		} else if l.currentItem.Token == PermissionToken {
			permission, err := l.readZPermission()
			if err != nil {
				return zdef, err
			}
			zdef.Permissions = append(zdef.Permissions, &permission)
			permission.myZDef = &zdef
		}
	}
	return zdef, nil
}

// <Zpermission> ::= "permission" <Rname> "=" <RPexpr>
// <RPexpr> ::= <RPterm> <RPexpr1>
// <RPexpr1> ::= <Zop> <RPterm> <RPexpr1> | ''
// <RPterm> ::= "(" <RPexpr> ")" | <Rname> "." "any" "(" <Rname> ")" | <Rname> "." "all" "(" <Rname> ")"
// <Rname> ::= <identifier>
// <Zop> ::= "+" | "&" | "-" | "->"

func (l *Lexer) readZPermission() (ZPermission, error) {
	var zpermission ZPermission
	if l.currentItem.Value != "permission" {
		return zpermission, fmt.Errorf("expected 'permission', but got '%s'", l.currentItem.Value)
	}
	l.NextToken()

	err := l.readAndMatchToken(IdentifierToken)

	if err != nil {
		return zpermission, err
	}
	zpermission.Name = l.currentItem.Value
	l.NextToken()

	err = l.readAndMatchToken(EqualToken)
	if err != nil {
		return zpermission, err
	}
	l.NextToken()

	RPexpr, err := l.readRPexpr()
	if err != nil {
		return zpermission, err
	}
	zpermission.RPexpr = &RPexpr

	return zpermission, nil

}

// <RPexpr> ::= <RPterm> <RPexpr1>
// <RPexpr1> ::= <Zop> <RPterm> <RPexpr1> | ''
// <RPterm> ::= "(" <RPexpr> ")" | <Rname> "." "any" "(" <Rname> ")" | <Rname> "." "all" "(" <Rname> ")"

func (l *Lexer) readRPexpr() (RPexpr, error) {
	rpexpr := RPexpr{}

	rpterm, err := l.readRPterm()
	if err != nil {
		return rpexpr, err
	}

	rpexpr1, err := l.readRPexpr1()
	if err != nil {
		return rpexpr, err
	}

	rpexpr.Term = &rpterm
	rpexpr.RPexpr1 = &rpexpr1

	return rpexpr, nil
}

// <RPexpr1> ::= <Zop> <RPterm> <RPexpr1> | ”
// <Zop> ::= "+" | "&" | "-" | "->"

func (l *Lexer) readRPexpr1() (RPexpr1, error) {
	rpexpr1 := RPexpr1{}
	rpexpr1.Zop = EOFToken
	err := l.readAndMatchTokenList(UnionToken, IntersectionToken, ExclusionToken, ArrowToken)
	if err != nil {
		return rpexpr1, nil // perhaps it's the void expression
	}
	rpexpr1.Zop = l.currentItem.Token
	l.NextToken()
	rpterm, err2 := l.readRPterm()
	if err2 != nil {
		return rpexpr1, err2
	}
	rpexpr1.RPTerm = &rpterm

	rpexpr2, err3 := l.readRPexpr1()
	if err3 != nil {
		return rpexpr1, err3
	}
	rpexpr1.RPexpr1 = &rpexpr2

	return rpexpr1, nil
}

// <RPterm> ::= <Rname> | "(" <RPexpr> ")" | <Rname> "." "any" "(" <Rname> ")" | <Rname> "." "all" "(" <Rname> ")"

func (l *Lexer) readRPterm() (RPterm, error) {
	rpterm := RPterm{}
	switch l.currentItem.Token {
	case LeftParenthesisToken:
		l.NextToken()
		_, err := l.readRPexpr()
		if err == nil {
			err = l.readAndMatchToken(RightParenthesisToken)
			if err == nil {
				l.NextToken()
				return rpterm, nil
			} else {
				return rpterm, err
			}
		} else {
			return rpterm, err
		}

	case IdentifierToken:
		name, err := l.readRName()
		if err == nil {
			rpterm.Name = name
			err = l.readAndMatchToken(DotToken)
			if err != nil {
				return rpterm, nil // perhaps end of rname
			} else {
				l.NextToken()
				err = l.readAndMatchTokenList(AnyToken, AllToken)
				if err == nil {
					l.NextToken()
					err = l.readAndMatchToken(LeftParenthesisToken)
					if err == nil {
						l.NextToken()
						rname, err := l.readRName()
						rpterm.Name = rname
						if err == nil {
							err = l.readAndMatchToken(RightParenthesisToken)
							if err == nil {
								l.NextToken()
								return rpterm, nil
							} else {
								return rpterm, fmt.Errorf("error, wait for ( token  but got %v at pos: %v", l.currentItem.Value, l.pos)
							}
						} else {
							return rpterm, err
						}

					} else {
						return rpterm, fmt.Errorf("error, wait for ( token  but got %v at pos: %v", l.currentItem.Value, l.pos)
					}
				} else {
					return rpterm, fmt.Errorf("error, wait for any or all token but got %v at pos: %v", l.currentItem.Value, l.pos)
				}
			}
		} else {
			return rpterm, err
		}
	default:
		return rpterm, fmt.Errorf("error, got %v  at pos: %v", l.currentItem.Value, l.pos)
	}
}

func (l *Lexer) readRName() (string, error) {
	var Rname string
	switch l.currentItem.Token {
	case IdentifierToken:
		Rname = l.currentItem.Value
		l.NextToken()
		return Rname, nil
	default:
		return "", fmt.Errorf("expected identifier token, got %v", l.currentItem.Value)
	}
}

// <Zrelation> ::= "relation" <Rname> ":" <Sname> ("|" <Sname)*
// <Sname> ::= <Zname> | <Zname> "#" <Rname> | <Zname> ":" "*"

func (l *Lexer) readZRelation() (ZRelation, error) {
	var zrelation ZRelation

	if l.currentItem.Value != "relation" {
		return zrelation, fmt.Errorf("expected 'relation', but got '%s'", l.currentItem.Value)
	}
	l.NextToken()

	err := l.readAndMatchToken(IdentifierToken)
	if err != nil {
		return zrelation, err
	}
	zrelation.Name = l.currentItem.Value
	l.NextToken()

	err = l.readAndMatchToken(ColonToken)
	if err != nil {
		return zrelation, err
	}
	l.NextToken()

	for l.currentItem.Token == IdentifierToken {
		_name := l.currentItem.Value
		l.NextToken()

		//  <Sname> ::= <Zname> "#" <Rname>
		if l.currentItem.Token == HashToken {
			l.NextToken()
			err := l.readAndMatchToken(IdentifierToken) // <Rname>
			if err != nil {
				return zrelation, err
			} else {
				zrelation.ZobjectSets = append(zrelation.ZobjectSets, &ZobjectSet{Name: _name, Relation: l.currentItem.Value})
				l.NextToken()
			}

		} else {
			// <Sname> ::= <Zname> ":" "*"
			if l.currentItem.Token == ColonToken {
				l.NextToken()
				err := l.readAndMatchToken(WildCardToken) // "*"
				if err != nil {
					return zrelation, err
				} else {
					// zrelation.ZobjectAll =
					zrelation.ZobjectWildCards = append(zrelation.ZobjectWildCards, &ZobjectWildCard{Name: _name})
					l.NextToken()
				}
			} else { // <Sname> ::= <Zname>
				zrelation.Zobjects = append(zrelation.Zobjects, &Zobject{Name: _name})
			}
		}

		if l.currentItem.Token == OrToken {
			l.NextToken()
		} else {
			break
			// it's the last
		}
	}

	return zrelation, nil
}

// Generation Code

type PlantUMLArchimateSchema struct {
	Zdefs   []*ZDef
	ZdefMap map[string]*ZDef

	SchemaDpi   int
	SchemaScale float64
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) GenerateForArchimate(modelname string) string {
	plantUMLArchimateSchema.createIDforZdef()

	type Node struct {
		ID         string
		Name       string
		Typ        string
		IsRelation bool // <-- pour distinguer la couleur du business objet relation
		X          int  // Ajout de la coordonnée X
		Y          int  // Ajout de la coordonnée Y
	}
	type Edge struct {
		ID         string
		Source     string
		Target     string
		Typ        string
		AccessType string
		Name       string
	}

	var nodes []Node
	var edges []Edge
	edgeCounter := 1

	addEdge := func(src, tgt, typ, accessType, name string) {
		edges = append(edges, Edge{
			ID:         fmt.Sprintf("rel_%d", edgeCounter),
			Source:     src,
			Target:     tgt,
			Typ:        typ,
			AccessType: accessType,
			Name:       name,
		})
		edgeCounter++
	}

	// 1. Zdefs (Business_Object)
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		nodes = append(nodes, Node{ID: zdef.ID, Name: zdef.Name, Typ: "BusinessObject", IsRelation: false})
	}

	// 2. Relations (Business_Object <<relation>>)
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			if zrel.ID != "NOTDRAW" {
				nodes = append(nodes, Node{ID: zrel.ID, Name: zrel.Name, Typ: "BusinessObject", IsRelation: true})
				addEdge(zdef.ID, zrel.ID, "Association", "", "")

				for _, zobject := range zrel.Zobjects {
					if zobject.ID != "NOTDRAW" && zobject.Unique {
						addEdge(zrel.ID, zobject.ID, "Access", "Write", "")
					}
				}
			}
		}
	}

	// 3. Permissions (Application_Service <<permission>>) et expressions
	var extractRPexpr func(sourceID string, expr *RPexpr)
	extractRPexpr = func(sourceID string, expr *RPexpr) {
		if expr == nil || expr.Term == nil {
			return
		}
		term := expr.Term
		next := expr.RPexpr1

		if next == nil || next.Zop == EOFToken {
			if term.Relation != nil && term.Relation.ID != "NOTDRAW" {
				addEdge(sourceID, term.Relation.ID, "Access", "Write", "")
			} else if term.Permission != nil && term.Permission.ID != "NOTDRAW" {
				addEdge(sourceID, term.Permission.ID, "Access", "Write", "")
			}
			return
		}

		op := next.Zop
		opID := next.ID
		opSymbol := TokenToString(op)

		// Operator node
		nodes = append(nodes, Node{ID: opID, Name: opSymbol, Typ: "BusinessObject"})
		addEdge(sourceID, opID, "Access", "Write", "")

		// Left branch
		if term.Relation != nil && term.Relation.ID != "NOTDRAW" {
			addEdge(opID, term.Relation.ID, "Access", "Write", "")
		} else if term.Permission != nil && term.Permission.ID != "NOTDRAW" {
			addEdge(opID, term.Permission.ID, "Access", "Write", "")
		}

		// Right branch
		rightTerm := next.RPTerm
		if op == ArrowToken {
			plantUMLArchimateSchema.validateAndComplete(rightTerm, term)
			if rightTerm != nil {
				if rightTerm.Relation != nil && rightTerm.Relation.ID != "NOTDRAW" {
					addEdge(opID, rightTerm.Relation.ID, "Access", "Write", "")
				} else if rightTerm.Permission != nil && rightTerm.Permission.ID != "NOTDRAW" {
					addEdge(opID, rightTerm.Permission.ID, "Access", "Write", "")
				}
			}
			return
		}

		if rightTerm != nil {
			if rightTerm.Relation != nil && rightTerm.Relation.ID != "NOTDRAW" {
				addEdge(opID, rightTerm.Relation.ID, "Access", "Write", "")
			} else if rightTerm.Permission != nil && rightTerm.Permission.ID != "NOTDRAW" {
				addEdge(opID, rightTerm.Permission.ID, "Access", "Write", "")
			}
		}

		// Recursion
		if next.RPexpr1 != nil && next.RPexpr1.Zop != EOFToken {
			subExpr := &RPexpr{
				Term:    rightTerm,
				RPexpr1: next.RPexpr1,
			}
			extractRPexpr(opID, subExpr)
		}
	}

	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zperm := range zdef.Permissions {
			if zperm.ID != "NOTDRAW" {
				nodes = append(nodes, Node{ID: zperm.ID, Name: zperm.Name, Typ: "ApplicationService"})
				addEdge(zdef.ID, zperm.ID, "Association", "", "")

				if zperm.RPexpr != nil {
					extractRPexpr(zperm.ID, zperm.RPexpr)
				}
			}
		}
	}

	// 4. ObjectSets (#)
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			if zrel.ID == "NOTDRAW" {
				continue
			}
			for _, zobjectSet := range zrel.ZobjectSets {
				if zobjectSet.ID != "NOTDRAW" && zobjectSet.IDRelation != "NOTDRAW" && zobjectSet.Unique {
					addEdge(zrel.ID, zobjectSet.IDRelation, "Access", "Write", fmt.Sprintf("%s#%s", zobjectSet.Name, zobjectSet.Relation))
				}
			}
		}
	}

	// 5. Wildcards (*)
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			if zrel.ID == "NOTDRAW" {
				continue
			}
			for _, zobjectWildCard := range zrel.ZobjectWildCards {
				if zobjectWildCard.ID != "NOTDRAW" && zobjectWildCard.Unique {
					addEdge(zrel.ID, zobjectWildCard.ID, "Access", "Write", "ALL")
				}
			}
		}
	}

	// --- Début de l'algorithme hiérarchique (Light Sugiyama) ---
	adj := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialisation des degrés entrants pour trouver les racines
	for i := range nodes {
		inDegree[nodes[i].ID] = 0
	}
	for _, e := range edges {
		adj[e.Source] = append(adj[e.Source], e.Target)
		inDegree[e.Target]++
	}

	// Étape 1 : Assigner des couches (layers) via BFS (gère les cycles SpiceDB)
	layerMap := make(map[string]int)
	visited := make(map[string]bool)
	var queue []string

	// Trouver les racines (nœuds non ciblés par une flèche)
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
			layerMap[id] = 0
			visited[id] = true
		}
	}

	// Sécurité anti-boucle : s'il n'y a pas de racine (cycle pur), on force le premier nœud
	if len(queue) == 0 && len(nodes) > 0 {
		rootID := nodes[0].ID
		queue = append(queue, rootID)
		layerMap[rootID] = 0
		visited[rootID] = true
	}

	maxLayer := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, neighbor := range adj[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				layerMap[neighbor] = layerMap[curr] + 1
				if layerMap[neighbor] > maxLayer {
					maxLayer = layerMap[neighbor]
				}
				queue = append(queue, neighbor)
			}
		}
	}

	// Rattrapage pour les nœuds isolés ou non atteints
	for i := range nodes {
		if !visited[nodes[i].ID] {
			layerMap[nodes[i].ID] = maxLayer + 1
		}
	}

	// Étape 2 : Regrouper les nœuds physiquement en mémoire
	layers := make([][]*Node, maxLayer+2)
	for i := range nodes {
		l := layerMap[nodes[i].ID]
		layers[l] = append(layers[l], &nodes[i])
	}

	// Étape 3 : Calculer et assigner X et Y par couche
	startY := 80
	nodeW, nodeH := 120, 55
	horizontalGap, verticalGap := 60, 80 // Espace entre les éléments

	for l, layerNodes := range layers {
		if len(layerNodes) == 0 {
			continue
		}

		yPos := startY + l*(nodeH+verticalGap)
		startX := 80

		for _, n := range layerNodes {
			n.X = startX
			n.Y = yPos
			startX += nodeW + horizontalGap
		}
	}
	// --- Fin de l'algorithme hiérarchique ---

	// --- XML Construction ---
	var out []string
	out = append(out, `<?xml version="1.0" encoding="UTF-8"?>`)
	out = append(out, `<model xmlns="http://www.opengroup.org/xsd/archimate/3.0/" `+
		`xmlns:dc="http://purl.org/dc/elements/1.1/" `+
		`xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" `+
		`xsi:schemaLocation="http://www.opengroup.org/xsd/archimate/3.0/ `+
		`http://www.opengroup.org/xsd/archimate/3.1/archimate3_Diagram.xsd `+
		`http://purl.org/dc/elements/1.1/ `+
		`http://www.opengroup.org/xsd/archimate/3.1/dc.xsd" `+
		`identifier="zanzibar-archimate">`)

	out = append(out, fmt.Sprintf(`  <name xml:lang="fr">%s</name>`, escapeXML(modelname)))
	out = append(out, `  <metadata>`)
	out = append(out, `    <schema>Dublin Core</schema>`)
	out = append(out, `    <schemaversion>1.1</schemaversion>`)
	out = append(out, `    <dc:creator>zreader</dc:creator>`)
	out = append(out, `  </metadata>`)

	// Elements
	out = append(out, `  <elements>`)
	for _, n := range nodes {
		out = append(out, fmt.Sprintf(`    <element identifier="%s" xsi:type="%s">`, n.ID, n.Typ))
		out = append(out, fmt.Sprintf(`      <name xml:lang="fr">%s</name>`, escapeXML(n.Name)))
		out = append(out, `    </element>`)
	}
	out = append(out, `  </elements>`)

	// Relationships
	out = append(out, `  <relationships>`)
	for _, e := range edges {
		accessStr := ""
		if e.AccessType != "" {
			accessStr = fmt.Sprintf(` accessType="%s"`, e.AccessType)
		}
		out = append(out, fmt.Sprintf(`    <relationship identifier="%s" source="%s" target="%s" xsi:type="%s"%s>`, e.ID, e.Source, e.Target, e.Typ, accessStr))
		if e.Name != "" {
			out = append(out, fmt.Sprintf(`      <name xml:lang="fr">%s</name>`, escapeXML(e.Name)))
		}
		out = append(out, `    </relationship>`)
	}
	out = append(out, `  </relationships>`)

	// Diagram Views
	out = append(out, `  <views>`)
	out = append(out, `    <diagrams>`)
	out = append(out, `      <view identifier="zanzibar-view" xsi:type="Diagram">`)
	out = append(out, fmt.Sprintf(`        <name xml:lang="fr">%s_diagram</name>`, escapeXML(modelname)))

	// Layout / Coordinates (Basé sur l'algorithme par couches)
	for _, n := range nodes {
		out = append(out, fmt.Sprintf(`        <node identifier="node_%s" elementRef="%s" xsi:type="Element" x="%d" y="%d" w="120" h="55">`, n.ID, n.ID, n.X, n.Y))
		out = append(out, `          <style>`)

		// Colors mapping : BusinessObject (jaune), ApplicationService (bleu cyan)
		r, g, b := 255, 255, 181

		if n.Typ == "ApplicationService" {
			r, g, b = 181, 255, 255 // ApplicationService (bleu cyan)
		} else if n.IsRelation {
			r, g, b = 255, 181, 181 // Relation (rouge pastel)
		}

		out = append(out, fmt.Sprintf(`            <fillColor r="%d" g="%d" b="%d" a="100" />`, r, g, b))
		out = append(out, `            <lineColor r="92" g="92" b="92" a="100" />`)
		out = append(out, `            <font name="Segoe UI" size="9">`)
		out = append(out, `              <color r="0" g="0" b="0" />`)
		out = append(out, `            </font>`)
		out = append(out, `          </style>`)
		out = append(out, `        </node>`)
	}

	// Diagram Connections
	for _, e := range edges {
		out = append(out, fmt.Sprintf(`        <connection identifier="conn_%s" relationshipRef="%s" xsi:type="Relationship" source="node_%s" target="node_%s">`, e.ID, e.ID, e.Source, e.Target))
		out = append(out, `          <style>`)
		out = append(out, `            <lineColor r="0" g="0" b="0" />`)
		out = append(out, `            <font name="Segoe UI" size="9">`)
		out = append(out, `              <color r="0" g="0" b="0" />`)
		out = append(out, `            </font>`)
		out = append(out, `          </style>`)
		out = append(out, `        </connection>`)
	}

	out = append(out, `      </view>`)
	out = append(out, `    </diagrams>`)
	out = append(out, `  </views>`)
	out = append(out, `</model>`)

	return strings.Join(out, "\n")
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) Generate(pngfilename string) string {
	var out []string
	plantUMLArchimateSchema.createIDforZdef()

	out = append(out, "@startuml "+pngfilename)
	out = append(out, "!include <archimate/Archimate>")

	out = append(out, "scale 1.0")
	out = append(out, "skinparam dpi 96")

	// Generate a row for each businessObject

	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		line := fmt.Sprintf("Business_Object(%s,\"%s\")", zdef.ID, zdef.Name)
		out = append(out, line)
	}

	// Generate a relationship line as a business object for each zdef
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			switch zrel.ID {
			case "NOTDRAW":
				line := fmt.Sprintf("rectangle \"relation %s is duplicated in definition %s \" #red", zrel.Name, zdef.Name)
				out = append(out, line)
			default:
				line := fmt.Sprintf("Business_Object(%s,\"%s\") <<relation>>", zrel.ID, zrel.Name)
				line2 := fmt.Sprintf("Rel_Association(%s,%s)", zdef.ID, zrel.ID)
				out = append(out, line)
				out = append(out, line2)
				for _, zobject := range zrel.Zobjects {
					switch zobject.ID {
					case "NOTDRAW":
						line3 := fmt.Sprintf("rectangle \"definition %s does not exist \" #red", zobject.Name)
						out = append(out, line3)
					default:
						switch zobject.Unique {
						case true:
							line4 := fmt.Sprintf("Rel_Access_w(%s,%s)", zrel.ID, zobject.ID)
							out = append(out, line4)
						default:
							line4 := fmt.Sprintf("rectangle \" %s is declared more that one in relation %s of definition %s\" #red ", zdef.Name, zrel.Name, zdef.Name)
							out = append(out, line4)
						}
					}
				}
			}
		}
	}

	// Generate a permissionship line as a business object for each zdef
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zperm := range zdef.Permissions {
			switch zperm.ID {
			case "NOTDRAW":
				line := fmt.Sprintf("rectangle \"permission %s is duplicated in definition %s \" #red", zperm.Name, zdef.Name)
				out = append(out, line)
			default:
				line := fmt.Sprintf("Application_Service(%s,\"%s\") <<permission>>", zperm.ID, zperm.Name)
				line2 := fmt.Sprintf("Rel_Association(%s,%s)", zdef.ID, zperm.ID)
				out = append(out, line)
				out = append(out, line2)

				if zperm.RPexpr != nil {
					plantUMLArchimateSchema.generateRPexpr(zperm.ID, zperm.RPexpr, &out)
				}

			}
		}
	}

	// Generate a relationshipSet row on a relation
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			for _, zobjectSet := range zrel.ZobjectSets {
				switch zobjectSet.ID {
				case "NOTDRAW":
					line := fmt.Sprintf("rectangle \"definition %s does not exist in \" #red", zobjectSet.Name)
					out = append(out, line)
				default:
					switch zobjectSet.IDRelation {
					case "NOTDRAW":
						line := fmt.Sprintf("rectangle \"  %s#%s in definition %s  : relation %s does not exist in %s \"  #red", zobjectSet.Name, zobjectSet.Relation, zdef.Name, zobjectSet.Relation, zobjectSet.Name)
						out = append(out, line)
					default:
						switch zobjectSet.Unique {
						case true:
							// line2 := fmt.Sprintf("Rel_Access_w(%s,%s,\"%s#%s\")", zobjectSet.IDRelation, zrel.ID, zobjectSet.Name, zobjectSet.Relation)
							// inversion
							line2 := fmt.Sprintf("Rel_Access_w(%s,%s,\"%s#%s\")", zrel.ID, zobjectSet.IDRelation, zobjectSet.Name, zobjectSet.Relation)
							out = append(out, line2)
						case false:
							line2 := fmt.Sprintf("rectangle \"  %s#%s declared more that one in relation %s of definition %s \"  #red", zobjectSet.Name, zobjectSet.Relation, zrel.Name, zdef.Name)
							out = append(out, line2)
						}
					}
				}
			}
		}
	}

	// Generate a relationWildCard row on a relation
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			for _, zobjectWildCard := range zrel.ZobjectWildCards {
				switch zobjectWildCard.ID {
				case "NOTDRAW":
					line := fmt.Sprintf("rectangle \"definition %s does not exist in \" #red", zobjectWildCard.Name)
					out = append(out, line)
				default:
					switch zobjectWildCard.Unique {

					case true:
						line2 := fmt.Sprintf("Rel_Access_w(%s,%s,\"%s\")", zrel.ID, zobjectWildCard.ID, "ALL")
						out = append(out, line2)

					case false:
						line3 := fmt.Sprintf("rectangle \"wildcard  %s:* is declared more than one in relation %s of definition %s\" #red", zobjectWildCard.Name, zrel.Name, zdef.Name)
						out = append(out, line3)
					}

				}
			}
		}
	}

	out = append(out, "@enduml")
	return strings.Join(out, "\n")
}

// validateAndComplete resolves the right-hand side of a tuple-to-child arrow (->)
//
// Example:  members->viewer
//   - contextTerm = left side ==> relation "members" whose type is "user | group | team"
//   - termRight   = right side ==> "viewer" to resolve
//
// In Zanzibar, "viewer" is valid if it exists as a relation or permission
// in **any** of the types listed in the relation definition.

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) validateAndComplete(termRight, contextTerm *RPterm) {
	if termRight == nil || termRight.Name == "" || contextTerm == nil {
		return
	}
	if contextTerm.Relation == nil || len(contextTerm.Relation.Zobjects) == 0 {
		return
	}

	toFind := termRight.Name

	// Try EVERY possible target type defined in the relation (because of | in the grammar)
	for _, zobj := range contextTerm.Relation.Zobjects {
		if zobj.myZDef == nil {
			continue // should not happen if schema validation passed
		}

		targetDef := zobj.myZDef

		// First: look for a relation named "viewer" in this type
		for _, rel := range targetDef.Relations {
			if rel.Name == toFind {
				termRight.Relation = rel
				return // success → stop here
			}
		}

		// Second: look for a permission named "viewer" in this type
		for _, perm := range targetDef.Permissions {
			if perm.Name == toFind {
				termRight.Permission = perm
				return // success → stop here
			}
		}
	}

	// Not found in ANY of the allowed types → leave empty
	// NOTDRAW system will show a red error → perfect!
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) generateRPexpr(sourceID string, expr *RPexpr, out *[]string) {
	if expr == nil || expr.Term == nil {
		return
	}

	term := expr.Term
	next := expr.RPexpr1

	// no operator
	if next == nil || next.Zop == EOFToken {
		if term.Relation != nil {
			*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", sourceID, term.Relation.ID))
		} else if term.Permission != nil {
			*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", sourceID, term.Permission.ID))
		}
		return
	}

	op := next.Zop
	opID := next.ID
	opSymbol := TokenToString(op)

	// operator node
	*out = append(*out, fmt.Sprintf("Business_Object(%s,\"%s\") <<op>>", opID, opSymbol))
	*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", sourceID, opID))

	// left branch
	if term.Relation != nil {
		*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", opID, term.Relation.ID))
	} else if term.Permission != nil {
		*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", opID, term.Permission.ID))
	}

	// right branch
	rightTerm := next.RPTerm

	if op == ArrowToken {
		// resolve ->
		plantUMLArchimateSchema.validateAndComplete(rightTerm, term)

		if rightTerm.Relation != nil {
			*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", opID, rightTerm.Relation.ID))
		} else if rightTerm.Permission != nil {
			*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", opID, rightTerm.Permission.ID))
		}
		// -> est terminal → pas de récursion
		return
	}

	// Normal Case: +, &, -
	if rightTerm != nil {
		if rightTerm.Relation != nil {
			*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", opID, rightTerm.Relation.ID))
		} else if rightTerm.Permission != nil {
			*out = append(*out, fmt.Sprintf("Rel_Access_w(%s,%s)", opID, rightTerm.Permission.ID))
		}
	}

	//  recursion except for ->
	if next.RPexpr1 != nil && next.RPexpr1.Zop != EOFToken {
		subExpr := &RPexpr{
			Term:    rightTerm,
			RPexpr1: next.RPexpr1,
		}
		plantUMLArchimateSchema.generateRPexpr(opID, subExpr, out)
	}
}

// utility
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) createIDforZdef() {
	zdefMapNameToVarName := make(map[string]string)

	for index, zdef := range plantUMLArchimateSchema.Zdefs {
		varname := fmt.Sprintf("b%d", index+1)
		if _, exists := zdefMapNameToVarName[zdef.Name]; exists {
			fmt.Printf("definition %s is declared more that one  \n", zdef.Name)
			continue
		}
		zdefMapNameToVarName[zdef.Name] = varname
		zdef.ID = varname
	}

	plantUMLArchimateSchema.createIDforZdefRelations()
	plantUMLArchimateSchema.initZdefMap()
	plantUMLArchimateSchema.createIDforZdefPermissions()

	plantUMLArchimateSchema.verifyAndAssignIDforZobjectInRelations()
	plantUMLArchimateSchema.verifyUniqueObjectForEachRelation()
	plantUMLArchimateSchema.verifyAndAssignIDInRelationsforZobjectSet()
	plantUMLArchimateSchema.verifyUniqueSetObjectForEachRelation()
	plantUMLArchimateSchema.verifyAndAssignIDInRelationsforZobjectWildCard()
	plantUMLArchimateSchema.verifyUniqueObjectWildCardForEachRelation()

	plantUMLArchimateSchema.verifyAndAssignTypeforRPTermInPermissions()
	plantUMLArchimateSchema.createIDforRPexpr1InPermissions()

	// plantUMLArchimateSchema.printRPTerms()

}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) createIDforZdefPermissions() {
	var permCount int = 0
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		PermNameSlice := []string{}
		for _, zperm := range zdef.Permissions {
			permCount++
			varname := fmt.Sprintf("p%d", permCount)
			if contains(PermNameSlice, zperm.Name) {
				zperm.ID = "NOTDRAW"
				fmt.Printf("permission %s is declared more that one in definition %s \n", zperm.Name, zdef.Name)

			} else {
				PermNameSlice = append(PermNameSlice, zperm.Name)
				zperm.ID = varname
				zperm.myZDef = zdef
			}
		}
	}
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) verifyAndAssignTypeforRPTermInPermissions() {
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zperm := range zdef.Permissions {
			if zperm.RPexpr != nil {
				// Résoudre le premier terme
				plantUMLArchimateSchema.processRPterm(zdef.Name, zperm.RPexpr.Term, nil)

				// Résoudre expr1 avec contexte
				var context *ZDef = nil
				if zperm.RPexpr.RPexpr1 != nil && zperm.RPexpr.RPexpr1.Zop == ArrowToken {
					if zperm.RPexpr.Term.Relation != nil {
						context = zperm.RPexpr.Term.Relation.myZDef
					} else if zperm.RPexpr.Term.Permission != nil {
						context = zperm.RPexpr.Term.Permission.myZDef
					}
				}
				plantUMLArchimateSchema.processRPexpr1(zdef.Name, zperm.RPexpr.RPexpr1, context)
			}
		}
	}
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) processRPterm(zdefName string, term *RPterm, arrowContext *ZDef) {
	if term == nil || term.Name == "" {
		return
	}

	// --- CAS 1 : contexte -> (ex: member dans usergroup) ---
	if arrowContext != nil {
		// 1. Chercher d'abord comme RELATION
		if rel, err := plantUMLArchimateSchema.findZRelation(arrowContext.Name, term.Name); err == nil {
			term.Relation = rel
			return
		}
		// 2. Puis comme PERMISSION
		if perm, err := plantUMLArchimateSchema.findZPermission(arrowContext.Name, term.Name); err == nil {
			term.Permission = perm
			return
		}
		return
	}

	// --- CAS 2 : normal (dans zdef courant) ---
	if rel, err := plantUMLArchimateSchema.findZRelation(zdefName, term.Name); err == nil {
		term.Relation = rel
		return
	}
	if perm, err := plantUMLArchimateSchema.findZPermission(zdefName, term.Name); err == nil {
		term.Permission = perm
	}
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) processRPexpr1(zdefName string, expr1 *RPexpr1, arrowContext *ZDef) {
	if expr1 == nil {
		return
	}

	var context *ZDef = nil
	if expr1.Zop == ArrowToken {
		// Contexte = objet du terme gauche (ex: group)
		context = arrowContext
	}

	plantUMLArchimateSchema.processRPterm(zdefName, expr1.RPTerm, context)
	plantUMLArchimateSchema.processRPexpr1(zdefName, expr1.RPexpr1, nil) // pas de -> imbriqué
}

// RPexpr1 contains operator Zop : +, &, -, ->

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) createIDforRPexpr1InPermissions() {
	var RPExpr1Count int = 1
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zperm := range zdef.Permissions {
			// varname := fmt.Sprintf("o%d", RPExpr1Count) // o for operator
			if zperm.RPexpr != nil {
				// Parcourir récursivement RPexpr1
				plantUMLArchimateSchema.createIDforRPexpr1InRPexpr1(zperm, zperm.RPexpr.RPexpr1, &RPExpr1Count)
			}
		}
	}

	// fmt.Printf("PRINT count %d\n", RPExpr1Count)

}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) createIDforRPexpr1InRPexpr1(zperm *ZPermission, rpexpr1 *RPexpr1, rpexpr1count *int) {
	if rpexpr1 == nil {
		return
	}

	if rpexpr1.Zop != EOFToken {
		varname := fmt.Sprintf("o%d", *rpexpr1count) // o pour opérateur
		rpexpr1.ID = varname
		*rpexpr1count++ // Incrémenter le compteur

		// Appel récursif correct sur la même méthode
		plantUMLArchimateSchema.createIDforRPexpr1InRPexpr1(zperm, rpexpr1.RPexpr1, rpexpr1count)

	}
}

// temp

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) printRPTerms() {
	fmt.Printf("Parcours...\n")
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zperm := range zdef.Permissions {
			// Print the permission itself
			// fmt.Printf("PRINT perm %s of %s is %s\n", zperm.Name, zperm.myZDef.Name, zperm.ID)
			fmt.Printf("PRINT Type: Permission, Name: %s, ID: %s in %s\n", zperm.Name, zperm.ID, zperm.myZDef.Name)

			// Then process the RPexpr components
			if zperm.RPexpr != nil {
				// Parcourir récursivement tous les RPterm dans RPexpr
				plantUMLArchimateSchema.printRPterm(zperm.RPexpr.Term)
				plantUMLArchimateSchema.printRPexpr1(zdef.Name, zperm.RPexpr.RPexpr1)
			}
		}
	}
	fmt.Printf("fin Parcours...\n")
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) printRPterm(term *RPterm) {
	if term == nil {
		return
	}

	if term.Relation != nil {
		// fmt.Printf("PRINT Type: Relation, Name: %s, ID : %s in %s\n", term.Relation.Name, term.Relation.ID, term.Relation.myZDef.Name)
	} else if term.Permission != nil {
		fmt.Printf("PRINT Type: Permission, Name: %s, ID: %s in %s\n", term.Permission.Name, term.Permission.ID, term.Permission.myZDef.Name)
	} else {
		fmt.Printf("PRINT Type: Unknown\n")
	}

}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) printRPexpr1(zdefName string, expr1 *RPexpr1) {
	if expr1 == nil {
		return
	}
	plantUMLArchimateSchema.printRPterm(expr1.RPTerm)
	plantUMLArchimateSchema.printRPexpr1(zdefName, expr1.RPexpr1)
}

// fin temp

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) findZPermission(objectName string, permissionName string) (*ZPermission, error) {
	myZf, err := plantUMLArchimateSchema.findZDef(objectName)
	if err != nil {
		return nil, err
	} else {
		zpermMap := make(map[string]*ZPermission)
		for _, zperm := range myZf.Permissions {
			zpermMap[zperm.Name] = zperm
		}
		if myZperm, exists := zpermMap[permissionName]; exists {
			return myZperm, nil
		} else {
			err = fmt.Errorf("permission %s does not exist in definition %s. ", permissionName, objectName)
			return nil, err
		}
	}

}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) createIDforZdefRelations() {
	var relCount int = 0
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		RelNameSlice := []string{}
		for _, zrel := range zdef.Relations {
			relCount++
			varname := fmt.Sprintf("r%d", relCount)
			if contains(RelNameSlice, zrel.Name) {
				zrel.ID = "NOTDRAW"
				fmt.Printf("relation %s is declared more that one in definition %s \n", zrel.Name, zdef.Name)

			} else {
				RelNameSlice = append(RelNameSlice, zrel.Name)
				zrel.ID = varname
				zrel.myZDef = zdef
			}
		}
	}
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) initZdefMap() {
	plantUMLArchimateSchema.ZdefMap = make(map[string]*ZDef)
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		plantUMLArchimateSchema.ZdefMap[zdef.Name] = zdef
	}

}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) verifyAndAssignIDforZobjectInRelations() {
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			for _, zobject := range zrel.Zobjects {
				if myZDef, exists := plantUMLArchimateSchema.ZdefMap[zobject.Name]; exists {
					zobject.ID = myZDef.ID
					zobject.myZDef = myZDef
				} else {
					fmt.Printf("%s declared in relation %s of definition %s does not exist.  \n", zobject.Name, zrel.Name, zdef.Name)
					zobject.ID = "NOTDRAW"
				}

			}
		}
	}
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) findZDef(objectName string) (*ZDef, error) {
	if myZDef, exists := plantUMLArchimateSchema.ZdefMap[objectName]; exists {
		return myZDef, nil
	} else {
		err := fmt.Errorf("definition %s does not exist. ", objectName)
		return nil, err
	}

}

// with tuple example like object:id#relation1@objectSet#relation2 (with Zanzibar notation)
// objectName is objectSet, relationName is relation2

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) findZRelation(objectName string, relationName string) (*ZRelation, error) {
	myZf, err := plantUMLArchimateSchema.findZDef(objectName)
	if err != nil {
		return nil, err
	} else {
		zrelMap := make(map[string]*ZRelation)
		for _, zrel := range myZf.Relations {
			zrelMap[zrel.Name] = zrel
		}
		if myZrel, exists := zrelMap[relationName]; exists {
			return myZrel, nil
		} else {
			err = fmt.Errorf("relation %s does not exist in definition %s. ", relationName, objectName)
			return nil, err
		}
	}

}

// tuple like resource:id#relation@group#relation (with Zanzibar notation)
// group#relation must exist or group#permission must exist

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) verifyAndAssignIDInRelationsforZobjectSet() {
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			for _, zobjectSet := range zrel.ZobjectSets {
				if myZDef, exists := plantUMLArchimateSchema.ZdefMap[zobjectSet.Name]; exists {
					zobjectSet.ID = myZDef.ID
					myZel, error := plantUMLArchimateSchema.findZRelation(myZDef.Name, zobjectSet.Relation)
					if error != nil {
						zobjectSet.IDRelation = "NOTDRAW"
						// check if zobjectSet.Relation is not the name of a permission declared inside the zobjetSet.Name
						myZel2, error2 := plantUMLArchimateSchema.findZPermission(zobjectSet.Name, zobjectSet.Relation)
						if error2 != nil {
							fmt.Printf("relation/permission %s declared in %s does not exist in  %s.  \n", zobjectSet.Relation, zdef.Name, myZDef.Name)

						} else {
							zobjectSet.IDRelation = myZel2.ID
							zobjectSet.IsPermission = true
						}

					} else {
						//double ?
						zobjectSet.IDRelation = myZel.ID
						zobjectSet.IsPermission = false
					}

				} else {
					fmt.Printf("%s declared in %s does not exist.  \n", zobjectSet.Name, zdef.Name)
					zobjectSet.ID = "NOTDRAW"
				}

			}
		}
	}
}

// tuple like resource:id#relation@user:* (with Zanzibar notation)

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) verifyAndAssignIDInRelationsforZobjectWildCard() {
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			for _, zobjectWildCard := range zrel.ZobjectWildCards {
				if myZDef, exists := plantUMLArchimateSchema.ZdefMap[zobjectWildCard.Name]; exists {
					zobjectWildCard.ID = myZDef.ID

				} else {
					fmt.Printf("%s declared in relation %s of definition %s does not exist.  \n", zobjectWildCard.Name, zrel.Name, zdef.Name)
					zobjectWildCard.ID = "NOTDRAW"
				}
			}
		}
	}
}

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) verifyUniqueObjectForEachRelation() {
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			keyObjectSlice := []string{}
			for _, zobject := range zrel.Zobjects {
				varname := zobject.Name
				if contains(keyObjectSlice, varname) {
					zobject.Unique = false
					fmt.Printf("%s is declared more that one in relation %s of definition %s\n", varname, zrel.Name, zdef.Name)
				} else {
					zobject.Unique = true
					keyObjectSlice = append(keyObjectSlice, varname)
				}
			}
		}
	}
}

// 'objectName#relationName' must be unique in zrel
func (plantUMLArchimateSchema *PlantUMLArchimateSchema) verifyUniqueSetObjectForEachRelation() {
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			keySetObjectSlice := []string{}
			for _, zobjectSet := range zrel.ZobjectSets {
				varname := fmt.Sprintf("%s#%s", zobjectSet.Name, zobjectSet.Relation)
				if contains(keySetObjectSlice, varname) {
					zobjectSet.Unique = false
					fmt.Printf("%s is declared more that one in relation %s of definition %s\n", varname, zrel.Name, zdef.Name)
				} else {
					zobjectSet.Unique = true
					keySetObjectSlice = append(keySetObjectSlice, varname)
				}
			}
		}
	}
}

// 'objectname:*' must be unique in zrel

func (plantUMLArchimateSchema *PlantUMLArchimateSchema) verifyUniqueObjectWildCardForEachRelation() {
	for _, zdef := range plantUMLArchimateSchema.Zdefs {
		for _, zrel := range zdef.Relations {
			keySetObjectSlice := []string{}
			for _, zobjectWildCard := range zrel.ZobjectWildCards {
				varname := zobjectWildCard.Name
				if contains(keySetObjectSlice, varname) {
					zobjectWildCard.Unique = false
					fmt.Printf("wildcard %s:* is declared more that one in relation %s of definition %s\n", varname, zrel.Name, zdef.Name)
				} else {
					zobjectWildCard.Unique = true
					keySetObjectSlice = append(keySetObjectSlice, varname)
				}
			}
		}
	}
}
