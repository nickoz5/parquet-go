package parse_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	sch "github.com/parsyl/parquet/generated"
	"github.com/parsyl/parquet/internal/parse"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestParquet(t *testing.T) {
	type testInput struct {
		name     string
		schema   []*sch.SchemaElement
		expected []parse.Field
		errors   []error
	}

	testCases := []testInput{
		{
			name: "single field",
			schema: []*sch.SchemaElement{
				{Name: "root", NumChildren: pint32(1)},
				{Name: "id", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
			},
			expected: []parse.Field{
				{FieldNames: []string{"Id"}, FieldTypes: []string{"int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required}},
			},
		},
		{
			name: "single nested field",
			schema: []*sch.SchemaElement{
				{Name: "root", NumChildren: pint32(1)},
				{Name: "hobby", RepetitionType: prt(sch.FieldRepetitionType_REQUIRED), NumChildren: pint32(1)},
				{Name: "name", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
			},
			expected: []parse.Field{
				{FieldNames: []string{"Hobby", "Name"}, FieldTypes: []string{"Hobby", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
			},
		},
		{
			name: "two nested fields",
			schema: []*sch.SchemaElement{
				{Name: "root", NumChildren: pint32(1)},
				{Name: "hobby", RepetitionType: prt(sch.FieldRepetitionType_REQUIRED), NumChildren: pint32(2)},
				{Name: "name", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "difficulty", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
			},
			expected: []parse.Field{
				{FieldNames: []string{"Hobby", "Name"}, FieldTypes: []string{"Hobby", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
				{FieldNames: []string{"Hobby", "Difficulty"}, FieldTypes: []string{"Hobby", "int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
			},
		},
		{
			name: "nested then not",
			schema: []*sch.SchemaElement{
				{Name: "root", NumChildren: pint32(2)},
				{Name: "hobby", RepetitionType: prt(sch.FieldRepetitionType_REQUIRED), NumChildren: pint32(2)},
				{Name: "name", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "difficulty", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "id", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
			},
			expected: []parse.Field{
				{FieldNames: []string{"Hobby", "Name"}, FieldTypes: []string{"Hobby", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
				{FieldNames: []string{"Hobby", "Difficulty"}, FieldTypes: []string{"Hobby", "int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
				{FieldNames: []string{"Id"}, FieldTypes: []string{"int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required}},
			},
		},
		{
			name: "nested 3 deep",
			schema: []*sch.SchemaElement{
				{Name: "root", NumChildren: pint32(2)},
				{Name: "hobby", RepetitionType: prt(sch.FieldRepetitionType_REQUIRED), NumChildren: pint32(2)},
				{Name: "name", RepetitionType: prt(sch.FieldRepetitionType_REQUIRED), NumChildren: pint32(2)},
				{Name: "first", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "last", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
				{Name: "difficulty", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "id", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
			},
			expected: []parse.Field{
				{FieldNames: []string{"Hobby", "Name", "First"}, FieldTypes: []string{"Hobby", "Name", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Required, parse.Optional}},
				{FieldNames: []string{"Hobby", "Name", "Last"}, FieldTypes: []string{"Hobby", "Name", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Required, parse.Required}},
				{FieldNames: []string{"Hobby", "Difficulty"}, FieldTypes: []string{"Hobby", "int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
				{FieldNames: []string{"Id"}, FieldTypes: []string{"int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required}},
			},
		},
		{
			name: "nested 3 deep v2",
			schema: []*sch.SchemaElement{
				{Name: "root", NumChildren: pint32(2)},
				{Name: "hobby", RepetitionType: prt(sch.FieldRepetitionType_REQUIRED), NumChildren: pint32(2)},
				{Name: "name", RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL), NumChildren: pint32(2)},
				{Name: "first", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "last", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
				{Name: "difficulty", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "id", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
			},
			expected: []parse.Field{
				{FieldNames: []string{"Hobby", "Name", "First"}, FieldTypes: []string{"Hobby", "Name", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional, parse.Optional}},
				{FieldNames: []string{"Hobby", "Name", "Last"}, FieldTypes: []string{"Hobby", "Name", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional, parse.Required}},
				{FieldNames: []string{"Hobby", "Difficulty"}, FieldTypes: []string{"Hobby", "int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
				{FieldNames: []string{"Id"}, FieldTypes: []string{"int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required}},
			},
		},
		{
			name: "nested 3 deep v3",
			schema: []*sch.SchemaElement{
				{Name: "root", NumChildren: pint32(2)},
				{Name: "hobby", RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL), NumChildren: pint32(2)},
				{Name: "name", RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL), NumChildren: pint32(2)},
				{Name: "first", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "last", Type: pt(sch.Type_BYTE_ARRAY), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
				{Name: "difficulty", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_OPTIONAL)},
				{Name: "id", Type: pt(sch.Type_INT32), RepetitionType: prt(sch.FieldRepetitionType_REQUIRED)},
			},
			expected: []parse.Field{
				{FieldNames: []string{"Hobby", "Name", "First"}, FieldTypes: []string{"Hobby", "Name", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Optional, parse.Optional, parse.Optional}},
				{FieldNames: []string{"Hobby", "Name", "Last"}, FieldTypes: []string{"Hobby", "Name", "string"}, RepetitionTypes: []parse.RepetitionType{parse.Optional, parse.Optional, parse.Required}},
				{FieldNames: []string{"Hobby", "Difficulty"}, FieldTypes: []string{"Hobby", "int32"}, RepetitionTypes: []parse.RepetitionType{parse.Optional, parse.Optional}},
				{FieldNames: []string{"Id"}, FieldTypes: []string{"int32"}, RepetitionTypes: []parse.RepetitionType{parse.Required}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%02d %s", i, tc.name), func(t *testing.T) {
			out, err := parse.Parquet(tc.schema)
			if !assert.Nil(t, err, tc.name) {
				return
			}

			if !assert.Equal(t, tc.expected, out.Fields, tc.name) {
				for _, f := range out.Fields {
					fmt.Printf("%+v\n", f)
				}
			}
			if assert.Equal(t, len(tc.errors), len(out.Errors), tc.name) {
				for i, err := range out.Errors {
					assert.EqualError(t, tc.errors[i], err.Error(), tc.name)
				}
			} else {
				for _, err := range out.Errors {
					fmt.Println(err)
				}
			}
		})
	}
}

func TestFields(t *testing.T) {

	type testInput struct {
		name     string
		typ      string
		expected []parse.Field
		errors   []error
	}

	testCases := []testInput{
		{
			name: "flat",
			typ:  "Being",
			expected: []parse.Field{
				{Type: "Being", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Being", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, FieldTypes: []string{"int32"}, ColumnName: "Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
		},
		{
			name: "private fields",
			typ:  "Private",
			expected: []parse.Field{
				{Type: "Private", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Private", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, FieldTypes: []string{"int32"}, ColumnName: "Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
		},
		{
			name: "nested struct",
			typ:  "Nested",
			expected: []parse.Field{
				{Type: "Nested", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Being", "ID"}, FieldTypes: []string{"Being", "int32"}, ColumnName: "Being.ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Required}},
				{Type: "Nested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Being", "Age"}, FieldTypes: []string{"Being", "int32"}, ColumnName: "Being.Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
				{Type: "Nested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, FieldTypes: []string{"uint64"}, ColumnName: "Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
			errors: []error{},
		},
		{
			name: "nested struct with name that doesn't match the struct type",
			typ:  "Nested2",
			expected: []parse.Field{
				{Type: "Nested2", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Info", "ID"}, FieldTypes: []string{"Being", "int32"}, ColumnName: "Info.ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Required}},
				{Type: "Nested2", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Info", "Age"}, FieldTypes: []string{"Being", "int32"}, ColumnName: "Info.Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
				{Type: "Nested2", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, FieldTypes: []string{"uint64"}, ColumnName: "Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
			errors: []error{},
		},
		{
			name: "2 deep nested struct",
			typ:  "DoubleNested",
			expected: []parse.Field{
				{Type: "DoubleNested", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Nested", "Being", "ID"}, FieldTypes: []string{"Nested", "Being", "int32"}, ColumnName: "Nested.Being.ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Required, parse.Required}},
				{Type: "DoubleNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Nested", "Being", "Age"}, FieldTypes: []string{"Nested", "Being", "int32"}, ColumnName: "Nested.Being.Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Required, parse.Optional}},
				{Type: "DoubleNested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Nested", "Anniversary"}, FieldTypes: []string{"Nested", "uint64"}, ColumnName: "Nested.Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
			},
			errors: []error{},
		},
		{
			name: "2 deep optional nested struct",
			typ:  "OptionalDoubleNested",
			expected: []parse.Field{
				{Type: "OptionalDoubleNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"OptionalNested", "Being", "ID"}, FieldTypes: []string{"OptionalNested", "Being", "int32"}, ColumnName: "OptionalNested.Being.ID", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional, parse.Required}},
				{Type: "OptionalDoubleNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"OptionalNested", "Being", "Age"}, FieldTypes: []string{"OptionalNested", "Being", "int32"}, ColumnName: "OptionalNested.Being.Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional, parse.Optional}},
				{Type: "OptionalDoubleNested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"OptionalNested", "Anniversary"}, FieldTypes: []string{"OptionalNested", "uint64"}, ColumnName: "OptionalNested.Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Required, parse.Optional}},
			},
			errors: []error{},
		},
		{
			name: "optional nested struct",
			typ:  "OptionalNested",
			expected: []parse.Field{
				{Type: "OptionalNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Being", "ID"}, FieldTypes: []string{"Being", "int32"}, ColumnName: "Being.ID", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional, parse.Required}},
				{Type: "OptionalNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Being", "Age"}, FieldTypes: []string{"Being", "int32"}, ColumnName: "Being.Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional, parse.Optional}},
				{Type: "OptionalNested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, FieldTypes: []string{"uint64"}, ColumnName: "Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
			errors: []error{},
		},
		{
			name: "optional nested struct v2",
			typ:  "OptionalNested2",
			expected: []parse.Field{
				{Type: "OptionalNested2", FieldType: "StringOptionalField", ParquetType: "StringType", TypeName: "string", FieldNames: []string{"Being", "Name"}, FieldTypes: []string{"Thing", "string"}, ColumnName: "Being.Name", Category: "stringOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional, parse.Required}},
				{Type: "OptionalNested2", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, FieldTypes: []string{"uint64"}, ColumnName: "Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
			errors: []error{},
		},
		{
			name:   "unsupported fields",
			typ:    "Unsupported",
			errors: []error{fmt.Errorf("unsupported type: Time")},
			expected: []parse.Field{
				{Type: "Unsupported", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Unsupported", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, FieldTypes: []string{"int32"}, ColumnName: "Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
		},
		{
			name: "unsupported fields mixed in with supported and embedded",
			typ:  "SupportedAndUnsupported",
			expected: []parse.Field{
				{Type: "SupportedAndUnsupported", FieldType: "Int64Field", ParquetType: "Int64Type", TypeName: "int64", FieldNames: []string{"Happiness"}, FieldTypes: []string{"int64"}, ColumnName: "Happiness", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "SupportedAndUnsupported", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "SupportedAndUnsupported", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, FieldTypes: []string{"int32"}, ColumnName: "Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "SupportedAndUnsupported", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, FieldTypes: []string{"uint64"}, ColumnName: "Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
			errors: []error{
				fmt.Errorf("unsupported type: T1"),
				fmt.Errorf("unsupported type: T2"),
			},
		},
		{
			name: "embedded",
			typ:  "Person",
			expected: []parse.Field{
				{Type: "Person", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Person", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, FieldTypes: []string{"int32"}, ColumnName: "Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "Person", FieldType: "Int64Field", ParquetType: "Int64Type", TypeName: "int64", FieldNames: []string{"Happiness"}, FieldTypes: []string{"int64"}, ColumnName: "Happiness", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Person", FieldType: "Int64OptionalField", ParquetType: "Int64Type", TypeName: "*int64", FieldNames: []string{"Sadness"}, FieldTypes: []string{"int64"}, ColumnName: "Sadness", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "Person", FieldType: "StringField", ParquetType: "StringType", TypeName: "string", FieldNames: []string{"Code"}, FieldTypes: []string{"string"}, ColumnName: "Code", Category: "string", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Person", FieldType: "Float32Field", ParquetType: "Float32Type", TypeName: "float32", FieldNames: []string{"Funkiness"}, FieldTypes: []string{"float32"}, ColumnName: "Funkiness", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Person", FieldType: "Float32OptionalField", ParquetType: "Float32Type", TypeName: "*float32", FieldNames: []string{"Lameness"}, FieldTypes: []string{"float32"}, ColumnName: "Lameness", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "Person", FieldType: "BoolOptionalField", ParquetType: "BoolType", TypeName: "*bool", FieldNames: []string{"Keen"}, FieldTypes: []string{"bool"}, ColumnName: "Keen", Category: "boolOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "Person", FieldType: "Uint32Field", ParquetType: "Uint32Type", TypeName: "uint32", FieldNames: []string{"Birthday"}, FieldTypes: []string{"uint32"}, ColumnName: "Birthday", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Person", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, FieldTypes: []string{"uint64"}, ColumnName: "Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
		},
		{
			name: "embedded preserve order",
			typ:  "NewOrderPerson",
			expected: []parse.Field{
				{Type: "NewOrderPerson", FieldType: "Int64Field", ParquetType: "Int64Type", TypeName: "int64", FieldNames: []string{"Happiness"}, FieldTypes: []string{"int64"}, ColumnName: "Happiness", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "NewOrderPerson", FieldType: "Int64OptionalField", ParquetType: "Int64Type", TypeName: "*int64", FieldNames: []string{"Sadness"}, FieldTypes: []string{"int64"}, ColumnName: "Sadness", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "NewOrderPerson", FieldType: "StringField", ParquetType: "StringType", TypeName: "string", FieldNames: []string{"Code"}, FieldTypes: []string{"string"}, ColumnName: "Code", Category: "string", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "NewOrderPerson", FieldType: "Float32Field", ParquetType: "Float32Type", TypeName: "float32", FieldNames: []string{"Funkiness"}, FieldTypes: []string{"float32"}, ColumnName: "Funkiness", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "NewOrderPerson", FieldType: "Float32OptionalField", ParquetType: "Float32Type", TypeName: "*float32", FieldNames: []string{"Lameness"}, FieldTypes: []string{"float32"}, ColumnName: "Lameness", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "NewOrderPerson", FieldType: "BoolOptionalField", ParquetType: "BoolType", TypeName: "*bool", FieldNames: []string{"Keen"}, FieldTypes: []string{"bool"}, ColumnName: "Keen", Category: "boolOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "NewOrderPerson", FieldType: "Uint32Field", ParquetType: "Uint32Type", TypeName: "uint32", FieldNames: []string{"Birthday"}, FieldTypes: []string{"uint32"}, ColumnName: "Birthday", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "NewOrderPerson", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "ID", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "NewOrderPerson", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, FieldTypes: []string{"int32"}, ColumnName: "Age", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
				{Type: "NewOrderPerson", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, FieldTypes: []string{"uint64"}, ColumnName: "Anniversary", Category: "numericOptional", RepetitionTypes: []parse.RepetitionType{parse.Optional}},
			},
		},
		{
			name: "tags",
			typ:  "Tagged",
			expected: []parse.Field{
				{Type: "Tagged", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "id", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
				{Type: "Tagged", FieldType: "StringField", ParquetType: "StringType", TypeName: "string", FieldNames: []string{"Name"}, FieldTypes: []string{"string"}, ColumnName: "name", Category: "string", RepetitionTypes: []parse.RepetitionType{parse.Required}},
			},
		},
		{
			name: "omit tag",
			typ:  "IgnoreMe",
			expected: []parse.Field{
				{Type: "IgnoreMe", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, FieldTypes: []string{"int32"}, ColumnName: "id", Category: "numeric", RepetitionTypes: []parse.RepetitionType{parse.Required}},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%02d %s", i, tc.name), func(t *testing.T) {
			out, err := parse.Fields(tc.typ, "./parse_test.go")
			assert.Nil(t, err, tc.name)
			assert.Equal(t, tc.expected, out.Fields, tc.name)
			if assert.Equal(t, len(tc.errors), len(out.Errors), tc.name) {
				for i, err := range out.Errors {
					assert.EqualError(t, tc.errors[i], err.Error(), tc.name)
				}
			} else {
				for _, err := range out.Errors {
					fmt.Println(err)
				}
			}
		})
	}
}

func pint32(i int32) *int32 {
	return &i
}

func prt(rt sch.FieldRepetitionType) *sch.FieldRepetitionType {
	return &rt
}

func pt(t sch.Type) *sch.Type {
	return &t
}
