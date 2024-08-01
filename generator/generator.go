package generator

import (
	"fmt"
	"github.com/t14raptor/go-trump/ast"
	"github.com/t14raptor/go-trump/token"
	"strings"
	"unicode"
)

func Generate(node ast.Node) string {
	s := &state{
		out:    &strings.Builder{},
		node:   node,
		parent: &state{},
	}
	gen(s)
	return s.out.String()
}

func gen(s *state) {
	switch n := s.node.(type) {
	case nil:
	case *ast.ArrowFunctionLiteral:
		s.out.WriteString(n.Source)
	case *ast.ArrayLiteral:
		s.out.WriteString("[")
		for i, ex := range n.Value {
			if ex.Expr != nil {
				gen(s.wrap(ex))
				if i < len(n.Value)-1 {
					s.out.WriteString(", ")
				}
			}
		}
		s.out.WriteString("]")
	case *ast.AssignExpression:
		if _, ok := s.parent.node.(*ast.BinaryExpression); ok {
			s.out.WriteString("(")
			defer s.out.WriteString(")")
		}
		gen(s.wrap(*n.Left))

		s.out.WriteString(" ")
		s.out.WriteString(n.Operator.String())
		if n.Operator != token.Assign {
			s.out.WriteString("=")
		}
		s.out.WriteString(" ")

		gen(s.wrap(*n.Right))
	case *ast.InvalidExpression:
	case *ast.BadStatement:
	case *ast.BinaryExpression:
		if pn, ok := s.parent.node.(*ast.BinaryExpression); ok {
			operatorPrecedence := n.Operator.Precedence(true)
			parentOperatorPrecedence := pn.Operator.Precedence(true)
			if operatorPrecedence < parentOperatorPrecedence || operatorPrecedence == parentOperatorPrecedence && pn.Right.Expr == n {
				s.out.WriteString("(")
				defer s.out.WriteString(")")
			}
		}
		gen(s.wrap(*n.Left))
		s.out.WriteString(" " + n.Operator.String() + " ")
		gen(s.wrap(*n.Right))
	case *ast.BlockStatement:
		s.out.WriteString("{")

		s.indent++
		for _, st := range n.List {
			s.lineAndPad()
			gen(s.wrap(st))
		}
		s.indent--

		s.lineAndPad()
		s.out.WriteString("}")
	case *ast.BooleanLiteral:
		s.out.WriteString(n.Literal)
	case *ast.BranchStatement:
		s.out.WriteString(n.Token.String())
		if n.Label != nil {
			s.out.WriteString(" ")
			gen(s.wrap(n.Label))
		}
		s.out.WriteString(";")
	case *ast.CallExpression:
		if _, ok := n.Callee.Expr.(ast.FunctionLiteral); ok {
			s.out.WriteString("(")
			gen(s.wrap(*n.Callee))
			s.out.WriteString(")")
		} else {
			gen(s.wrap(*n.Callee))
		}
		s.out.WriteString("(")
		for i, a := range n.ArgumentList {
			gen(s.wrap(a))
			if i < len(n.ArgumentList)-1 {
				s.out.WriteString(", ")
			}
		}
		s.out.WriteString(")")
	case *ast.CaseStatement:
		if n.Test != nil {
			s.out.WriteString("case ")
			gen(s.wrap(*n.Test))
			s.out.WriteString(": ")
		} else {
			s.out.WriteString("default: ")
		}
		gen(s.wrap(&ast.BlockStatement{List: n.Consequent}))
	case *ast.CatchStatement:
		gen(s.wrap(*n.Parameter))
		gen(s.wrap(n.Body))
	case *ast.FunctionDeclaration:
		s.lineAndPad()
		gen(s.wrap(n.Function))
	case *ast.ConditionalExpression:
		if _, ok := s.parent.node.(*ast.BinaryExpression); ok {
			s.out.WriteString("(")
			defer s.out.WriteString(")")
		}
		gen(s.wrap(*n.Test))
		s.out.WriteString(" ? ")
		gen(s.wrap(*n.Consequent))
		s.out.WriteString(" : ")
		gen(s.wrap(*n.Alternate))
	case *ast.DebuggerStatement:
		s.out.WriteString("debugger;")
	case *ast.DoWhileStatement:
		gen(s.wrap(*n.Test))
		gen(s.wrap(*n.Body))
	case *ast.MemberExpression:
		gen(s.wrap(*n.Object))
		if st, ok := n.Property.Expr.(*ast.StringLiteral); ok && valid(st.Value.String()) {
			s.out.WriteString(".")
			s.out.WriteString(st.Value.String())
		} else {
			s.out.WriteString("[")
			gen(s.wrap(*n.Property))
			s.out.WriteString("]")
		}
	case *ast.DotExpression:
		gen(s.wrap(*n.Left))
		s.out.WriteString(".")
		s.out.WriteString(n.Identifier.Name.String())
	case *ast.EmptyStatement:
		s.out.WriteString(";")
	case *ast.ExpressionStatement:
		gen(s.wrap(*n.Expression))
		s.out.WriteString(";")
		if len(n.Comment) > 0 {
			s.out.WriteString(" // " + n.Comment)
		}
	case *ast.ExpressionBody:
		gen(s.wrap(n.Expression))
	case *ast.ForInStatement:
		s.out.WriteString("for (")
		gen(s.wrap(*n.Into))
		s.out.WriteString(" in ")
		gen(s.wrap(*n.Source))
		s.out.WriteString(") ")
		gen(s.wrap(*n.Body))
	case *ast.ForIntoExpression:
		gen(s.wrap(*n.Expression))
	case *ast.ForStatement:
		s.out.WriteString("for (")
		gen(s.wrap(*n.Initializer))
		s.out.WriteString("; ")
		gen(s.wrap(*n.Test))
		s.out.WriteString("; ")
		gen(s.wrap(*n.Update))
		s.out.WriteString(") ")

		switch n.Body.Stmt.(type) {
		case ast.EmptyStatement, ast.BlockStatement:
		default:
			n.Body = &ast.Statement{ast.BlockStatement{List: ast.Statements{*n.Body}}}
		}
		gen(s.wrap(*n.Body))
	case *ast.ForLoopInitializerExpression:
		gen(s.wrap(*n.Expression))
	case *ast.FunctionLiteral:
		s.out.WriteString("function ")
		gen(s.wrap(n.Name))
		s.out.WriteString("(")
		for i, p := range n.ParameterList.List {
			gen(s.wrap(p))
			if i < len(n.ParameterList.List)-1 {
				s.out.WriteString(", ")
			}
		}
		s.out.WriteString(") ")
		gen(s.wrap(n.Body))
	case *ast.Identifier:
		if n != nil {
			s.out.WriteString(n.Name.String())
		}
	case *ast.IfStatement:
		s.out.WriteString("if (")
		gen(s.wrap(*n.Test))
		s.out.WriteString(") ")

		switch n.Consequent.Stmt.(type) {
		case ast.EmptyStatement, ast.BlockStatement:
		default:
			n.Consequent = &ast.Statement{Stmt: ast.BlockStatement{List: ast.Statements{*n.Consequent}}}
		}
		gen(s.wrap(*n.Consequent))

		if n.Alternate != nil {
			s.out.WriteString(" else ")

			switch n.Alternate.Stmt.(type) {
			case ast.EmptyStatement, ast.BlockStatement, ast.IfStatement:
			default:
				n.Alternate = &ast.Statement{Stmt: ast.BlockStatement{List: ast.Statements{*n.Alternate}}}
			}
			gen(s.wrap(*n.Alternate))
		}
	case *ast.LabelledStatement:
		gen(s.wrap(n.Label))
		gen(s.wrap(*n.Statement))
	case *ast.NewExpression:
		s.out.WriteString("new ")
		gen(s.wrap(*n.Callee))
		s.out.WriteString("(")
		for i, a := range n.ArgumentList {
			gen(s.wrap(a))
			if i < len(n.ArgumentList)-1 {
				s.out.WriteString(", ")
			}
		}
		s.out.WriteString(")")
	case *ast.NullLiteral:
		s.out.WriteString(n.Literal)
	case *ast.NumberLiteral:
		s.out.WriteString(n.Literal)
	case *ast.ObjectLiteral:
		s.out.WriteString("{")

		s.indent++
		for i, p := range n.Value {
			s.lineAndPad()
			gen(s.wrap(p))
			if i < len(n.Value)-1 {
				s.out.WriteString(", ")
			}
		}
		s.indent--

		if len(n.Value) > 0 {
			s.lineAndPad()
		}
		s.out.WriteString("}")
	case *ast.PropertyKeyed:
		gen(s.wrap(*n.Key))
		s.out.WriteString(": ")
		gen(s.wrap(*n.Value))
	case *ast.Program:
		if n != nil {
			for _, b := range n.Body {
				gen(s.wrap(b))
				s.line()
			}
		}
	case *ast.RegExpLiteral:
		s.out.WriteString(n.Literal)
	case *ast.ReturnStatement:
		if n != nil {
			s.out.WriteString("return")
			if n.Argument != nil {
				s.out.WriteString(" ")
				gen(s.wrap(*n.Argument))
			}
			s.out.WriteString(";")
		}
	case *ast.SequenceExpression:
		switch s.parent.node.(type) {
		case *ast.BinaryExpression, *ast.ConditionalExpression, *ast.AssignExpression, *ast.CallExpression:
			s.out.WriteString("(")
			defer s.out.WriteString(")")
		}
		for i, e := range n.Sequence {
			gen(s.wrap(e))
			if i < len(n.Sequence)-1 {
				s.out.WriteString(", ")
			}
		}
	case *ast.StringLiteral:
		s.out.WriteString(n.Literal)
	case *ast.SwitchStatement:
		s.out.WriteString("switch (")
		gen(s.wrap(*n.Discriminant))
		s.out.WriteString(") {")

		s.indent++
		for _, c := range n.Body {
			s.lineAndPad()
			gen(s.wrap(c))
		}
		s.indent--

		if len(n.Body) > 0 {
			s.lineAndPad()
		}
		s.out.WriteString("}")
	case *ast.ThisExpression:
		s.out.WriteString("this")
	case *ast.ThrowStatement:
		s.out.WriteString("throw ")
		gen(s.wrap(*n.Argument))
		s.out.WriteString(";")
	case *ast.TryStatement:
		s.out.WriteString("try ")

		gen(s.wrap(n.Body))

		if n.Catch != nil {
			s.out.WriteString(" catch (")
			gen(s.wrap(*n.Catch.Parameter))
			s.out.WriteString(") ")
			gen(s.wrap(n.Catch.Body))
		}
		if n.Finally != nil {
			gen(s.wrap(n.Finally))
		}
	case *ast.UnaryExpression:
		if !n.Postfix {
			s.out.WriteString(n.Operator.String())
			if len(n.Operator.String()) > 2 {
				s.out.WriteString(" ")
			}
		}

		wrap := false
		switch n.Operand.Expr.(type) {
		case *ast.BinaryExpression, *ast.ConditionalExpression, *ast.AssignExpression, *ast.UnaryExpression:
			wrap = true
		}

		if wrap {
			s.out.WriteString("(")
		}
		gen(s.wrap(*n.Operand))
		if wrap {
			s.out.WriteString(")")
		}

		if n.Postfix {
			s.out.WriteString(n.Operator.String())
		}
	case *ast.VariableStatement:
		s.out.WriteString("var ")
		for _, v := range n.List {
			gen(s.wrap(v.Target))
			if v.Initializer != nil {
				s.out.WriteString(" = ")
				gen(s.wrap(*v.Initializer))
			}
		}
		s.out.WriteString(";")
		s.line()
	case *ast.WhileStatement:
		s.out.WriteString("while (")
		gen(s.wrap(*n.Test))
		s.out.WriteString(") ")
		gen(s.wrap(*n.Body))
	case *ast.WithStatement:
		gen(s.wrap(*n.Object))
		gen(s.wrap(*n.Body))
	case *ast.VariableDeclarator:
		gen(s.wrap(n.Target))
		if n.Initializer != nil {
			s.out.WriteString(" = ")
			gen(s.wrap(*n.Initializer))
		}
	case *ast.ForLoopInitializerLexicalDecl:
		s.out.WriteString(n.LexicalDeclaration.Token.String())
		s.out.WriteString(" ")
		for _, b := range n.LexicalDeclaration.List {
			gen(s.wrap(b))
		}
	case *ast.LexicalDeclaration:
		s.out.WriteString(n.Token.String())
		s.out.WriteString(" ")
		for i, b := range n.List {
			gen(s.wrap(b))
			if i < len(n.List)-1 {
				s.out.WriteString(", ")
			}
		}
		s.out.WriteString(";")
		if len(n.Comment) > 0 {
			s.out.WriteString(" // " + n.Comment)
		}
	default:
		panic(fmt.Sprintf("gen: unexpected node type %T", n))
	}
}

func valid(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
