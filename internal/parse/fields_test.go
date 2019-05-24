package parse_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/parsyl/parquet/internal/parse"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
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
				{Type: "Being", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "ID", Category: "numeric", Optionals: []bool{false}},
				{Type: "Being", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, ColumnName: "Age", Category: "numericOptional", Optionals: []bool{true}},
			},
		},
		{
			name: "private fields",
			typ:  "Private",
			expected: []parse.Field{
				{Type: "Private", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "ID", Category: "numeric", Optionals: []bool{false}},
				{Type: "Private", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, ColumnName: "Age", Category: "numericOptional", Optionals: []bool{true}},
			},
		},
		{
			name: "nested struct",
			typ:  "Nested",
			expected: []parse.Field{
				{Type: "Nested", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Being", "ID"}, ColumnName: "Being.ID", Category: "numeric", Optionals: []bool{false, false}},
				{Type: "Nested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Being", "Age"}, ColumnName: "Being.Age", Category: "numericOptional", Optionals: []bool{false, true}},
				{Type: "Nested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, ColumnName: "Anniversary", Category: "numericOptional", Optionals: []bool{true}},
			},
			errors: []error{},
		},
		{
			name: "nested struct with name that doesn't match the struct type",
			typ:  "Nested2",
			expected: []parse.Field{
				{Type: "Nested", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Being", "ID"}, ColumnName: "Being.ID", Category: "numeric", Optionals: []bool{false, false}},
				{Type: "Nested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Being", "Age"}, ColumnName: "Being.Age", Category: "numericOptional", Optionals: []bool{false, true}},
				{Type: "Nested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, ColumnName: "Anniversary", Category: "numericOptional", Optionals: []bool{true}},
			},
			errors: []error{},
		},
		{
			name: "2 deep nested struct",
			typ:  "DoubleNested",
			expected: []parse.Field{
				{Type: "DoubleNested", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Nested", "Being", "ID"}, ColumnName: "Nested.Being.ID", Category: "numeric", Optionals: []bool{false, false, false}},
				{Type: "DoubleNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Nested", "Being", "Age"}, ColumnName: "Nested.Being.Age", Category: "numericOptional", Optionals: []bool{false, false, true}},
				{Type: "DoubleNested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Nested", "Anniversary"}, ColumnName: "Nested.Anniversary", Category: "numericOptional", Optionals: []bool{false, true}},
			},
			errors: []error{},
		},
		{
			name: "2 deep optional nested struct",
			typ:  "OptionalDoubleNested",
			expected: []parse.Field{
				{Type: "OptionalDoubleNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"OptionalNested", "Being", "ID"}, ColumnName: "OptionalNested.Being.ID", Category: "numericOptional", Optionals: []bool{false, true, false}},
				{Type: "OptionalDoubleNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"OptionalNested", "Being", "Age"}, ColumnName: "OptionalNested.Being.Age", Category: "numericOptional", Optionals: []bool{false, true, true}},
				{Type: "OptionalDoubleNested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"OptionalNested", "Anniversary"}, ColumnName: "OptionalNested.Anniversary", Category: "numericOptional", Optionals: []bool{false, true}},
			},
			errors: []error{},
		},
		{
			name: "optional nested struct",
			typ:  "OptionalNested",
			expected: []parse.Field{
				{Type: "OptionalNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"Being", "ID"}, ColumnName: "Being.ID", Category: "numericOptional", Optionals: []bool{true, false}},
				{Type: "OptionalNested", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Being", "Age"}, ColumnName: "Being.Age", Category: "numericOptional", Optionals: []bool{true, true}},
				{Type: "OptionalNested", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, ColumnName: "Anniversary", Category: "numericOptional", Optionals: []bool{true}},
			},
			errors: []error{},
		},
		{
			name:   "unsupported fields",
			typ:    "Unsupported",
			errors: []error{fmt.Errorf("unsupported type: Time")},
			expected: []parse.Field{
				{Type: "Unsupported", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "ID", Category: "numeric", Optionals: []bool{false}},
				{Type: "Unsupported", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, ColumnName: "Age", Category: "numericOptional", Optionals: []bool{true}},
			},
		},
		{
			name: "unsupported fields mixed in with supported and embedded",
			typ:  "SupportedAndUnsupported",
			expected: []parse.Field{
				{Type: "SupportedAndUnsupported", FieldType: "Int64Field", ParquetType: "Int64Type", TypeName: "int64", FieldNames: []string{"Happiness"}, ColumnName: "Happiness", Category: "numeric", Optionals: []bool{false}},
				{Type: "SupportedAndUnsupported", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "ID", Category: "numeric", Optionals: []bool{false}},
				{Type: "SupportedAndUnsupported", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, ColumnName: "Age", Category: "numericOptional", Optionals: []bool{true}},
				{Type: "SupportedAndUnsupported", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, ColumnName: "Anniversary", Category: "numericOptional", Optionals: []bool{true}},
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
				{Type: "Person", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "ID", Category: "numeric", Optionals: []bool{false}},
				{Type: "Person", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, ColumnName: "Age", Category: "numericOptional", Optionals: []bool{true}},
				{Type: "Person", FieldType: "Int64Field", ParquetType: "Int64Type", TypeName: "int64", FieldNames: []string{"Happiness"}, ColumnName: "Happiness", Category: "numeric", Optionals: []bool{false}},
				{Type: "Person", FieldType: "Int64OptionalField", ParquetType: "Int64Type", TypeName: "*int64", FieldNames: []string{"Sadness"}, ColumnName: "Sadness", Category: "numericOptional", Optionals: []bool{true}},
				{Type: "Person", FieldType: "StringField", ParquetType: "StringType", TypeName: "string", FieldNames: []string{"Code"}, ColumnName: "Code", Category: "string", Optionals: []bool{false}},
				{Type: "Person", FieldType: "Float32Field", ParquetType: "Float32Type", TypeName: "float32", FieldNames: []string{"Funkiness"}, ColumnName: "Funkiness", Category: "numeric", Optionals: []bool{false}},
				{Type: "Person", FieldType: "Float32OptionalField", ParquetType: "Float32Type", TypeName: "*float32", FieldNames: []string{"Lameness"}, ColumnName: "Lameness", Category: "numericOptional", Optionals: []bool{true}},
				{Type: "Person", FieldType: "BoolOptionalField", ParquetType: "BoolType", TypeName: "*bool", FieldNames: []string{"Keen"}, ColumnName: "Keen", Category: "boolOptional", Optionals: []bool{true}},
				{Type: "Person", FieldType: "Uint32Field", ParquetType: "Uint32Type", TypeName: "uint32", FieldNames: []string{"Birthday"}, ColumnName: "Birthday", Category: "numeric", Optionals: []bool{false}},
				{Type: "Person", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, ColumnName: "Anniversary", Category: "numericOptional", Optionals: []bool{true}},
			},
		},
		{
			name: "embedded preserve order",
			typ:  "NewOrderPerson",
			expected: []parse.Field{
				{Type: "NewOrderPerson", FieldType: "Int64Field", ParquetType: "Int64Type", TypeName: "int64", FieldNames: []string{"Happiness"}, ColumnName: "Happiness", Category: "numeric", Optionals: []bool{false}},
				{Type: "NewOrderPerson", FieldType: "Int64OptionalField", ParquetType: "Int64Type", TypeName: "*int64", FieldNames: []string{"Sadness"}, ColumnName: "Sadness", Category: "numericOptional", Optionals: []bool{true}},
				{Type: "NewOrderPerson", FieldType: "StringField", ParquetType: "StringType", TypeName: "string", FieldNames: []string{"Code"}, ColumnName: "Code", Category: "string", Optionals: []bool{false}},
				{Type: "NewOrderPerson", FieldType: "Float32Field", ParquetType: "Float32Type", TypeName: "float32", FieldNames: []string{"Funkiness"}, ColumnName: "Funkiness", Category: "numeric", Optionals: []bool{false}},
				{Type: "NewOrderPerson", FieldType: "Float32OptionalField", ParquetType: "Float32Type", TypeName: "*float32", FieldNames: []string{"Lameness"}, ColumnName: "Lameness", Category: "numericOptional", Optionals: []bool{true}},
				{Type: "NewOrderPerson", FieldType: "BoolOptionalField", ParquetType: "BoolType", TypeName: "*bool", FieldNames: []string{"Keen"}, ColumnName: "Keen", Category: "boolOptional", Optionals: []bool{true}},
				{Type: "NewOrderPerson", FieldType: "Uint32Field", ParquetType: "Uint32Type", TypeName: "uint32", FieldNames: []string{"Birthday"}, ColumnName: "Birthday", Category: "numeric", Optionals: []bool{false}},
				{Type: "NewOrderPerson", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "ID", Category: "numeric", Optionals: []bool{false}},
				{Type: "NewOrderPerson", FieldType: "Int32OptionalField", ParquetType: "Int32Type", TypeName: "*int32", FieldNames: []string{"Age"}, ColumnName: "Age", Category: "numericOptional", Optionals: []bool{true}},
				{Type: "NewOrderPerson", FieldType: "Uint64OptionalField", ParquetType: "Uint64Type", TypeName: "*uint64", FieldNames: []string{"Anniversary"}, ColumnName: "Anniversary", Category: "numericOptional", Optionals: []bool{true}},
			},
		},
		{
			name: "tags",
			typ:  "Tagged",
			expected: []parse.Field{
				{Type: "Tagged", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "id", Category: "numeric", Optionals: []bool{false}},
				{Type: "Tagged", FieldType: "StringField", ParquetType: "StringType", TypeName: "string", FieldNames: []string{"Name"}, ColumnName: "name", Category: "string", Optionals: []bool{false}},
			},
		},
		{
			name: "omit tag",
			typ:  "IgnoreMe",
			expected: []parse.Field{
				{Type: "IgnoreMe", FieldType: "Int32Field", ParquetType: "Int32Type", TypeName: "int32", FieldNames: []string{"ID"}, ColumnName: "id", Category: "numeric", Optionals: []bool{false}},
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
