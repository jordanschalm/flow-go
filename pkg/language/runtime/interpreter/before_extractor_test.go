package interpreter

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dapperlabs/flow-go/pkg/language/runtime/ast"
	"github.com/dapperlabs/flow-go/pkg/language/runtime/parser"
)

func TestBeforeExtractor(t *testing.T) {
	RegisterTestingT(t)

	expression, inputIsComplete, err := parser.ParseExpression(`
        before(x + before(y)) + z
	`)

	Expect(inputIsComplete).To(BeTrue())

	Expect(err).
		To(Not(HaveOccurred()))

	extractor := NewBeforeExtractor()

	identifier1 := ast.Identifier{
		Identifier: extractor.ExpressionExtractor.FormatIdentifier(0),
	}
	identifier2 := ast.Identifier{
		Identifier: extractor.ExpressionExtractor.FormatIdentifier(1),
	}

	result := extractor.ExtractBefore(expression)

	Expect(result).
		To(Equal(ast.ExpressionExtraction{
			RewrittenExpression: &ast.BinaryExpression{
				Operation: ast.OperationPlus,
				Left: &ast.IdentifierExpression{
					Identifier: identifier2,
				},
				Right: &ast.IdentifierExpression{
					Identifier: ast.Identifier{
						Identifier: "z",
						Pos:        ast.Position{Offset: 33, Line: 2, Column: 32},
					},
				},
			},
			ExtractedExpressions: []ast.ExtractedExpression{
				{
					Identifier: identifier1,
					Expression: &ast.IdentifierExpression{
						Identifier: ast.Identifier{
							Identifier: "y",
							Pos:        ast.Position{Offset: 27, Line: 2, Column: 26},
						},
					},
				},
				{
					Identifier: identifier2,
					Expression: &ast.BinaryExpression{
						Operation: ast.OperationPlus,
						Left: &ast.IdentifierExpression{
							Identifier: ast.Identifier{
								Identifier: "x",
								Pos:        ast.Position{Offset: 16, Line: 2, Column: 15},
							},
						},
						Right: &ast.IdentifierExpression{
							Identifier: identifier1,
						},
					},
				},
			},
		}))
}
