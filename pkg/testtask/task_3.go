package testtask

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
)

const SchemaExample = `schema {
    query: Query
}

type Query {
    droid: Droid!
    hero(id: ID!): Character
}

interface Character {
    name: String!
}

type Droid implements Character {
    name: String!
}`

type testSchemaBuilder struct {
	doc      *ast.Document
	typeRefs typeRefs
}

type typeRefs struct {
	characterType     int
	stringNonNullType int
	idNonNullType     int
	droidNonNullType  int
}

func (t *testSchemaBuilder) importSchemaDefinition() {
	t.doc.ImportSchemaDefinition("Query", "", "")
}

func (t *testSchemaBuilder) importQueryDroidFieldDefinition() int {
	// droid (field)
	droidFieldDefRef := t.doc.ImportFieldDefinition(
		"droid", "", t.typeRefs.droidNonNullType, nil, nil)

	return droidFieldDefRef
}

func (t *testSchemaBuilder) importQueryHeroFieldDefinition() int {

	nonNullIDInputValueDefinitionRef := t.doc.ImportInputValueDefinition("id", "", t.typeRefs.idNonNullType, ast.DefaultValue{})

	// hero (field)
	characterFieldDefRef := t.doc.ImportFieldDefinition(
		"hero", "", t.typeRefs.characterType, []int{nonNullIDInputValueDefinitionRef}, nil)

	return characterFieldDefRef
}

func (t *testSchemaBuilder) importQueryDefinition() {
	queryTypeFieldDefRefs := make([]int, 0, 2)

	queryTypeFieldDefRefs = append(queryTypeFieldDefRefs, t.importQueryDroidFieldDefinition())
	queryTypeFieldDefRefs = append(queryTypeFieldDefRefs, t.importQueryHeroFieldDefinition())

	t.doc.ImportObjectTypeDefinition(
		"Query",
		"",
		queryTypeFieldDefRefs,
		nil)
}

func (t *testSchemaBuilder) nameNonNullStringField() int {

	nameFieldDefRef := t.doc.ImportFieldDefinition(
		"name", "", t.typeRefs.stringNonNullType, nil, nil)

	return nameFieldDefRef
}

func (t *testSchemaBuilder) importCharacterDefinition() {

	nameFieldDefRef := t.nameNonNullStringField()

	t.doc.ImportInterfaceTypeDefinition("Character", "", []int{nameFieldDefRef})
}

func (t *testSchemaBuilder) importDroidDefinition() {
	nameFieldDefRef := t.nameNonNullStringField()

	t.doc.ImportObjectTypeDefinition("Droid", "", []int{nameFieldDefRef}, []int{t.typeRefs.characterType})
}

func BuildAst() *ast.Document {
	doc := ast.NewDocument()

	schemaBuilder := testSchemaBuilder{
		doc: doc,
		typeRefs: typeRefs{
			characterType:     doc.AddNamedType([]byte("Character")),
			stringNonNullType: doc.AddNonNullNamedType([]byte("String")),
			idNonNullType:     doc.AddNonNullNamedType([]byte("ID")),
			droidNonNullType:  doc.AddNonNullNamedType([]byte("Droid")),
		},
	}

	schemaBuilder.importSchemaDefinition()

	schemaBuilder.importQueryDefinition()

	schemaBuilder.importCharacterDefinition()

	schemaBuilder.importDroidDefinition()

	return schemaBuilder.doc
}
