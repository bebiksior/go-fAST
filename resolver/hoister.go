package resolver

import (
	"github.com/t14raptor/go-fast/ast"
	"github.com/t14raptor/go-fast/token"
)

// Loosely inspired from https://rustdoc.swc.rs/swc_ecma_transforms_base/fn.resolver.html

type Hoister struct {
	ast.NoopVisitor

	resolver *Resolver
	kind     DeclKind
	inBlock  bool

	inCatchBody bool

	excludedFromCatch map[string]struct{}
	catchParamDecls   map[string]struct{}
}

func NewHoister(resolver *Resolver) *Hoister {
	return &Hoister{
		resolver:          resolver,
		kind:              DeclKindVar,
		excludedFromCatch: make(map[string]struct{}),
		catchParamDecls:   make(map[string]struct{}),
	}
}

func (h *Hoister) addIdent(id *ast.Identifier) {
	if h.inCatchBody {
		if _, ok := h.catchParamDecls[id.Name]; ok {
			if r, _ := h.resolver.lookupContext(id.Name); r != UnresolvedMark {
				return
			}
		}

		h.excludedFromCatch[id.Name] = struct{}{}
	} else if _, ok := h.catchParamDecls[id.Name]; ok {
		if _, excluded := h.excludedFromCatch[id.Name]; !excluded {
			return
		}
	}

	h.resolver.modify(id, h.kind)
}

func (h *Hoister) VisitBlockStatement(n *ast.BlockStatement) {
	old := h.inBlock
	h.inBlock = true
	n.VisitChildrenWith(h)
	h.inBlock = old
}

func (h *Hoister) VisitStatements(n *ast.Statements) {
	others := make(ast.Statements, 0, len(*n))
	for i := range *n {
		switch it := (*n)[i].Stmt.(type) {
		case *ast.VariableDeclaration:
			it.VisitWith(h)
		case *ast.FunctionDeclaration:
			it.VisitWith(h)
		default:
			others = append(others, (*n)[i])
		}
	}

	for i := range others {
		others[i].VisitWith(h)
	}
}

func (h *Hoister) VisitVariableDeclaration(n *ast.VariableDeclaration) {
	if h.inBlock {
		if n.Token != token.Var {
			return
		}
	}

	oldKind := h.kind
	h.kind = DeclKindVar
	n.VisitChildrenWith(h)
	h.kind = oldKind
}

func (h *Hoister) VisitBindingTarget(n *ast.BindingTarget) {
	if ident, ok := n.Target.(*ast.Identifier); ok {
		h.addIdent(ident)
		return
	}
	n.VisitChildrenWith(h)
}

func (h *Hoister) VisitFunctionDeclaration(n *ast.FunctionDeclaration) {
	if _, ok := h.catchParamDecls[n.Function.Name.Name]; ok {
		return
	}

	if h.inBlock {
		if symbol := h.resolver.current.findVarInfo(n.Function.Name.Name); symbol != nil {
			if symbol.DeclKind != DeclKindVar && symbol.DeclKind != DeclKindFunction {
				return
			}
		}
	}

	h.resolver.modify(n.Function.Name, DeclKindFunction)
}

func (h *Hoister) VisitSwitchStatement(n *ast.SwitchStatement) {
	n.Discriminant.VisitWith(h)

	old := h.inBlock
	h.inBlock = true
	n.Body.VisitWith(h)
	h.inBlock = old
}

func (h *Hoister) VisitArrowFunctionLiteral(n *ast.ArrowFunctionLiteral) {}
func (h *Hoister) VisitExpression(n *ast.Expression)                     {}
func (h *Hoister) VisitFunctionLiteral(n *ast.FunctionLiteral)           {}
