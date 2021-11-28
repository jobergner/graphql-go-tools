package testtask

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astvisitor"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

type DocumentStats struct {
	uniqueFieldNames []string
	objectTypesNames []string
	stringFieldCount int
	enumValues       []string
}

func GatherDocumentStats(doc *ast.Document, report *operationreport.Report) *DocumentStats {
	walker := astvisitor.NewWalker(48)
	visitor := &DocumentStatsVisitor{
		Walker: &walker,
	}

	walker.RegisterEnterDocumentVisitor(visitor)

	// register additional walk methods here
	walker.RegisterEnterEnumValueDefinitionVisitor(visitor)
	walker.RegisterEnterFieldDefinitionVisitor(visitor)
	walker.RegisterEnterObjectTypeDefinitionVisitor(visitor)
	walker.RegisterEnterInputValueDefinitionVisitor(visitor)

	// run walker
	walker.Walk(doc, nil, report)

	// obtain results

	uniqueFieldNames := make([]string, 0, len(visitor.uniqueFieldNames))
	for fieldName := range visitor.uniqueFieldNames {
		uniqueFieldNames = append(uniqueFieldNames, fieldName)
	}

	return &DocumentStats{
		uniqueFieldNames: uniqueFieldNames,
		objectTypesNames: visitor.objectTypesNames,
		stringFieldCount: visitor.stringFieldCount,
		enumValues:       visitor.enumValues,
	}
}

type DocumentStatsVisitor struct {
	*astvisitor.Walker
	definition       *ast.Document
	enumValues       []string
	uniqueFieldNames map[string]struct{}
	stringFieldCount int
	objectTypesNames []string
}

func (v *DocumentStatsVisitor) EnterEnumValueDefinition(ref int) {
	v.enumValues = append(v.enumValues, v.definition.EnumValueDefinitionNameString(ref))
}

func (v *DocumentStatsVisitor) EnterFieldDefinition(ref int) {

	fieldName := v.definition.FieldDefinitionNameString(ref)
	v.uniqueFieldNames[fieldName] = struct{}{}

	fieldTypeRef := v.definition.FieldDefinitionType(ref)
	v.countStringType(fieldTypeRef)
}

func (v *DocumentStatsVisitor) EnterObjectTypeDefinition(ref int) {
	definitionName := v.definition.ObjectTypeDefinitionNameString(ref)
	v.objectTypesNames = append(v.objectTypesNames, definitionName)
}

func (v *DocumentStatsVisitor) EnterDocument(operation, _ *ast.Document) {
	v.definition = operation
	v.enumValues = make([]string, 0, 3)
	v.uniqueFieldNames = make(map[string]struct{})
}

func (v *DocumentStatsVisitor) EnterInputValueDefinition(ref int) {
	valueTypeRef := v.definition.InputValueDefinitionType(ref)
	v.countStringType(valueTypeRef)
}

func (v *DocumentStatsVisitor) countStringType(typeRef int) {
	fieldType := v.definition.Types[typeRef]

	switch fieldType.TypeKind {
	case ast.TypeKindNamed:
		if v.definition.TypeNameString(typeRef) == "String" {
			v.stringFieldCount++
		}
	case ast.TypeKindNonNull:
		if v.definition.TypeNameString(fieldType.OfType) == "String" {
			v.stringFieldCount++
		}
	}
}
