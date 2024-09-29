package ast

// Idx is a compact encoding of a source position within JS code.
type Idx int

type Node interface {
	// Idx0 returns the index of the first character belonging to the node.
	Idx0() Idx
	// Idx1 returns the index of the first character immediately after the node.
	Idx1() Idx
}

type Program struct {
	Body Statements
}

func (n *Optional) Idx0() Idx              { return (*n.Expr).Expr.Idx0() }
func (n *OptionalChain) Idx0() Idx         { return (*n.Base).Expr.Idx0() }
func (n *ObjectPattern) Idx0() Idx         { return n.LeftBrace }
func (n *ParameterList) Idx0() Idx         { return n.Opening }
func (n *ArrayLiteral) Idx0() Idx          { return n.LeftBracket }
func (n *ArrayPattern) Idx0() Idx          { return n.LeftBracket }
func (n *YieldExpression) Idx0() Idx       { return n.Yield }
func (n *AwaitExpression) Idx0() Idx       { return n.Await }
func (n *AssignExpression) Idx0() Idx      { return (*n.Left).Expr.Idx0() }
func (n *BinaryExpression) Idx0() Idx      { return (*n.Left).Expr.Idx0() }
func (n *BooleanLiteral) Idx0() Idx        { return n.Idx }
func (n *CallExpression) Idx0() Idx        { return (*n.Callee).Expr.Idx0() }
func (n *ConditionalExpression) Idx0() Idx { return (*n.Test).Expr.Idx0() }
func (n *PrivateDotExpression) Idx0() Idx  { return (*n.Left).Expr.Idx0() }
func (n *FunctionLiteral) Idx0() Idx       { return n.Function }
func (n *ClassLiteral) Idx0() Idx          { return n.Class }
func (n *ArrowFunctionLiteral) Idx0() Idx  { return n.Start }
func (n *Identifier) Idx0() Idx            { return n.Idx }
func (n *InvalidExpression) Idx0() Idx     { return n.From }
func (n *NewExpression) Idx0() Idx         { return n.New }
func (n *NullLiteral) Idx0() Idx           { return n.Idx }
func (n *NumberLiteral) Idx0() Idx         { return n.Idx }
func (n *ObjectLiteral) Idx0() Idx         { return n.LeftBrace }
func (n *RegExpLiteral) Idx0() Idx         { return n.Idx }
func (n *SequenceExpression) Idx0() Idx    { return n.Sequence[0].Expr.Idx0() }
func (n *StringLiteral) Idx0() Idx         { return n.Idx }
func (n *TemplateElement) Idx0() Idx       { return n.Idx }
func (n *TemplateLiteral) Idx0() Idx       { return n.OpenQuote }
func (n *ThisExpression) Idx0() Idx        { return n.Idx }
func (n *SuperExpression) Idx0() Idx       { return n.Idx }
func (n *UnaryExpression) Idx0() Idx       { return n.Idx }
func (n *UpdateExpression) Idx0() Idx      { return n.Idx }
func (n *MetaProperty) Idx0() Idx          { return n.Idx }
func (n *MemberExpression) Idx0() Idx      { return 0 }
func (n *MemberExpression) Idx1() Idx      { return 0 }
func (n *SpreadElement) Idx0() Idx {
	return n.Expression.Expr.Idx0()
}
func (n *SpreadElement) Idx1() Idx {
	return n.Expression.Expr.Idx1()
}

func (n *BadStatement) Idx0() Idx        { return n.From }
func (n *BlockStatement) Idx0() Idx      { return n.LeftBrace }
func (n *BreakStatement) Idx0() Idx      { return n.Idx }
func (n *ContinueStatement) Idx0() Idx   { return n.Idx }
func (n *CaseStatement) Idx0() Idx       { return n.Case }
func (n *CatchStatement) Idx0() Idx      { return n.Catch }
func (n *DebuggerStatement) Idx0() Idx   { return n.Debugger }
func (n *DoWhileStatement) Idx0() Idx    { return n.Do }
func (n *EmptyStatement) Idx0() Idx      { return n.Semicolon }
func (n *ExpressionStatement) Idx0() Idx { return (*n.Expression).Expr.Idx0() }
func (n *ForInStatement) Idx0() Idx      { return n.For }
func (n *ForOfStatement) Idx0() Idx      { return n.For }
func (n *ForStatement) Idx0() Idx        { return n.For }
func (n *IfStatement) Idx0() Idx         { return n.If }
func (n *LabelledStatement) Idx0() Idx   { return n.Label.Idx0() }
func (n *Program) Idx0() Idx             { return n.Body[0].Idx0() }
func (n *ReturnStatement) Idx0() Idx     { return n.Return }
func (n *SwitchStatement) Idx0() Idx     { return n.Switch }
func (n *ThrowStatement) Idx0() Idx      { return n.Throw }
func (n *TryStatement) Idx0() Idx        { return n.Try }
func (n *WhileStatement) Idx0() Idx      { return n.While }
func (n *WithStatement) Idx0() Idx       { return n.With }
func (n *VariableDeclaration) Idx0() Idx { return n.Idx }
func (n *FunctionDeclaration) Idx0() Idx { return n.Function.Idx0() }
func (n *ClassDeclaration) Idx0() Idx    { return n.Class.Idx0() }
func (n *VariableDeclarator) Idx0() Idx  { return n.Target.Idx0() }

func (n *PropertyShort) Idx0() Idx { return n.Name.Idx }
func (n *PropertyKeyed) Idx0() Idx { return (*n.Key).Expr.Idx0() }

func (n *FieldDefinition) Idx0() Idx  { return n.Idx }
func (n *MethodDefinition) Idx0() Idx { return n.Idx }
func (n *ClassStaticBlock) Idx0() Idx { return n.Static }

func (n *ForDeclaration) Idx0() Idx    { return n.Idx }
func (n *ForIntoVar) Idx0() Idx        { return n.Binding.Idx0() }
func (n *ForIntoExpression) Idx0() Idx { return (*n.Expression).Expr.Idx0() }

func (n *ForLoopInitializer) Idx0() Idx { return 0 }

func (n *Optional) Idx1() Idx              { return (*n.Expr).Expr.Idx1() }
func (n *OptionalChain) Idx1() Idx         { return (*n.Base).Expr.Idx1() }
func (n *ArrayLiteral) Idx1() Idx          { return n.RightBracket + 1 }
func (n *ArrayPattern) Idx1() Idx          { return n.RightBracket + 1 }
func (n *AssignExpression) Idx1() Idx      { return (*n.Right).Expr.Idx1() }
func (n *AwaitExpression) Idx1() Idx       { return (*n.Argument).Expr.Idx1() }
func (n *InvalidExpression) Idx1() Idx     { return n.To }
func (n *BinaryExpression) Idx1() Idx      { return (*n.Right).Expr.Idx1() }
func (n *BooleanLiteral) Idx1() Idx        { return Idx(int(n.Idx) + 4) }
func (n *CallExpression) Idx1() Idx        { return n.RightParenthesis + 1 }
func (n *ConditionalExpression) Idx1() Idx { return (*n.Test).Expr.Idx1() }
func (n *PrivateDotExpression) Idx1() Idx  { return n.Identifier.Idx1() }
func (n *FunctionLiteral) Idx1() Idx       { return n.Body.Idx1() }
func (n *ClassLiteral) Idx1() Idx          { return n.RightBrace + 1 }
func (n *ArrowFunctionLiteral) Idx1() Idx  { return n.Body.Body.Idx1() }
func (n *Identifier) Idx1() Idx            { return Idx(int(n.Idx) + len(n.Name)) }
func (n *NewExpression) Idx1() Idx {
	if n.ArgumentList != nil {
		return n.RightParenthesis + 1
	} else {
		return (*n.Callee).Expr.Idx1()
	}
}
func (n *NullLiteral) Idx1() Idx        { return Idx(int(n.Idx) + 4) } // "null"
func (n *NumberLiteral) Idx1() Idx      { return Idx(int(n.Idx) + len(n.Literal)) }
func (n *ObjectLiteral) Idx1() Idx      { return n.RightBrace + 1 }
func (n *ObjectPattern) Idx1() Idx      { return n.RightBrace + 1 }
func (n *ParameterList) Idx1() Idx      { return n.Closing + 1 }
func (n *RegExpLiteral) Idx1() Idx      { return Idx(int(n.Idx) + len(n.Literal)) }
func (n *SequenceExpression) Idx1() Idx { return n.Sequence[len(n.Sequence)-1].Expr.Idx1() }
func (n *StringLiteral) Idx1() Idx      { return Idx(int(n.Idx) + len(n.Literal)) }
func (n *TemplateElement) Idx1() Idx    { return Idx(int(n.Idx) + len(n.Literal)) }
func (n *TemplateLiteral) Idx1() Idx    { return n.CloseQuote + 1 }
func (n *ThisExpression) Idx1() Idx     { return n.Idx + 4 }
func (n *SuperExpression) Idx1() Idx    { return n.Idx + 5 }
func (n *UnaryExpression) Idx1() Idx {
	return (*n.Operand).Expr.Idx1()
}
func (n *UpdateExpression) Idx1() Idx {
	if n.Postfix {
		return (*n.Operand).Expr.Idx1() + 2 // ++ --
	}
	return (*n.Operand).Expr.Idx1()
}
func (n *MetaProperty) Idx1() Idx {
	return n.Property.Idx1()
}

func (n *BadStatement) Idx1() Idx        { return n.To }
func (n *BlockStatement) Idx1() Idx      { return n.RightBrace + 1 }
func (n *BreakStatement) Idx1() Idx      { return n.Idx }
func (n *ContinueStatement) Idx1() Idx   { return n.Idx }
func (n *CaseStatement) Idx1() Idx       { return n.Consequent[len(n.Consequent)-1].Idx1() }
func (n *CatchStatement) Idx1() Idx      { return n.Body.Idx1() }
func (n *DebuggerStatement) Idx1() Idx   { return n.Debugger + 8 }
func (n *DoWhileStatement) Idx1() Idx    { return (*n.Test).Expr.Idx1() }
func (n *EmptyStatement) Idx1() Idx      { return n.Semicolon + 1 }
func (n *ExpressionStatement) Idx1() Idx { return (*n.Expression).Expr.Idx1() }
func (n *ForInStatement) Idx1() Idx      { return (*n.Body).Idx1() }
func (n *ForOfStatement) Idx1() Idx      { return (*n.Body).Idx1() }
func (n *ForStatement) Idx1() Idx        { return (*n.Body).Idx1() }
func (n *IfStatement) Idx1() Idx {
	if n.Alternate != nil {
		return (*n.Alternate).Idx1()
	}
	return (*n.Consequent).Idx1()
}
func (n *LabelledStatement) Idx1() Idx { return n.Colon + 1 }
func (n *Program) Idx1() Idx           { return n.Body[len(n.Body)-1].Idx1() }
func (n *ReturnStatement) Idx1() Idx   { return n.Return + 6 }
func (n *SwitchStatement) Idx1() Idx   { return n.Body[len(n.Body)-1].Idx1() }
func (n *ThrowStatement) Idx1() Idx    { return (*n.Argument).Expr.Idx1() }
func (n *TryStatement) Idx1() Idx {
	if n.Finally != nil {
		return n.Finally.Idx1()
	}
	if n.Catch != nil {
		return n.Catch.Idx1()
	}
	return n.Body.Idx1()
}
func (n *WhileStatement) Idx1() Idx      { return (*n.Body).Idx1() }
func (n *WithStatement) Idx1() Idx       { return (*n.Body).Idx1() }
func (n *VariableDeclaration) Idx1() Idx { return n.List[len(n.List)-1].Idx1() }
func (n *FunctionDeclaration) Idx1() Idx { return n.Function.Idx1() }
func (n *ClassDeclaration) Idx1() Idx    { return n.Class.Idx1() }
func (n *VariableDeclarator) Idx1() Idx {
	if n.Initializer != nil {
		return (*n.Initializer).Expr.Idx1()
	}
	return n.Target.Idx1()
}

func (n *PropertyShort) Idx1() Idx {
	if n.Initializer != nil {
		return (*n.Initializer).Expr.Idx1()
	}
	return n.Name.Idx1()
}

func (n *PropertyKeyed) Idx1() Idx { return (*n.Value).Expr.Idx1() }

func (n *FieldDefinition) Idx1() Idx {
	if n.Initializer != nil {
		return (*n.Initializer).Expr.Idx1()
	}
	return (*n.Key).Expr.Idx1()
}

func (n *MethodDefinition) Idx1() Idx {
	return n.Body.Idx1()
}

func (n *ClassStaticBlock) Idx1() Idx {
	return n.Block.Idx1()
}

func (n *YieldExpression) Idx1() Idx {
	if n.Argument != nil {
		return (*n.Argument).Expr.Idx1()
	}
	return n.Yield + 5
}

func (n *ForDeclaration) Idx1() Idx     { return n.Target.Idx1() }
func (n *ForIntoVar) Idx1() Idx         { return n.Binding.Idx1() }
func (n *ForIntoExpression) Idx1() Idx  { return (*n.Expression).Expr.Idx1() }
func (n *ForLoopInitializer) Idx1() Idx { return 0 }
func (n *ConciseBody) Idx0() Idx        { return n.Body.Idx0() }
func (n *ConciseBody) Idx1() Idx        { return n.Body.Idx1() }
