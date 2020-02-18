package schema

import (
	"go/ast"
	"go/types"

	"github.com/bflad/tfproviderlint/helper/astutils"
	tfschema "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	SchemaFieldAtLeastOneOf     = `AtLeastOneOf`
	SchemaFieldComputed         = `Computed`
	SchemaFieldComputedWhen     = `ComputedWhen`
	SchemaFieldConfigMode       = `ConfigMode`
	SchemaFieldConflictsWith    = `ConflictsWith`
	SchemaFieldDefault          = `Default`
	SchemaFieldDefaultFunc      = `DefaultFunc`
	SchemaFieldDeprecated       = `Deprecated`
	SchemaFieldDescription      = `Description`
	SchemaFieldDiffSuppressFunc = `DiffSuppressFunc`
	SchemaFieldElem             = `Elem`
	SchemaFieldExactlyOneOf     = `ExactlyOneOf`
	SchemaFieldForceNew         = `ForceNew`
	SchemaFieldInputDefault     = `InputDefault`
	SchemaFieldMaxItems         = `MaxItems`
	SchemaFieldMinItems         = `MinItems`
	SchemaFieldOptional         = `Optional`
	SchemaFieldPromoteSingle    = `PromoteSingle`
	SchemaFieldRemoved          = `Removed`
	SchemaFieldRequired         = `Required`
	SchemaFieldSensitive        = `Sensitive`
	SchemaFieldSet              = `Set`
	SchemaFieldStateFunc        = `StateFunc`
	SchemaFieldType             = `Type`
	SchemaFieldValidateFunc     = `ValidateFunc`

	SchemaValueTypeBool   = `TypeBool`
	SchemaValueTypeFloat  = `TypeFloat`
	SchemaValueTypeInt    = `TypeInt`
	SchemaValueTypeList   = `TypeList`
	SchemaValueTypeMap    = `TypeMap`
	SchemaValueTypeSet    = `TypeSet`
	SchemaValueTypeString = `TypeString`

	TypeNameSchema    = `Schema`
	TypeNameSet       = `Set`
	TypeNameValueType = `ValueType`
)

// SchemaInfo represents all gathered Schema data for easier access
type SchemaInfo struct {
	AstCompositeLit *ast.CompositeLit
	Fields          map[string]*ast.KeyValueExpr
	Schema          *tfschema.Schema
	SchemaValueType string
	TypesInfo       *types.Info
}

// NewSchemaInfo instantiates a SchemaInfo
func NewSchemaInfo(cl *ast.CompositeLit, info *types.Info) *SchemaInfo {
	result := &SchemaInfo{
		AstCompositeLit: cl,
		Fields:          astutils.CompositeLitFields(cl),
		Schema:          &tfschema.Schema{},
		SchemaValueType: typeSchemaType(cl, info),
		TypesInfo:       info,
	}

	if kvExpr := result.Fields[SchemaFieldAtLeastOneOf]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.AtLeastOneOf = []string{}
	}

	if kvExpr := result.Fields[SchemaFieldComputed]; kvExpr != nil && astutils.ExprBoolValue(kvExpr.Value) != nil {
		result.Schema.Computed = *astutils.ExprBoolValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldConflictsWith]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.ConflictsWith = []string{}
	}

	if kvExpr := result.Fields[SchemaFieldDefault]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.Default = func() (interface{}, error) { return nil, nil }
	}

	if kvExpr := result.Fields[SchemaFieldDefaultFunc]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.DefaultFunc = func() (interface{}, error) { return nil, nil }
	}

	if kvExpr := result.Fields[SchemaFieldDescription]; kvExpr != nil && astutils.ExprStringValue(kvExpr.Value) != nil {
		result.Schema.Description = *astutils.ExprStringValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldExactlyOneOf]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.ExactlyOneOf = []string{}
	}

	if kvExpr := result.Fields[SchemaFieldDiffSuppressFunc]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.DiffSuppressFunc = func(k, old, new string, d *tfschema.ResourceData) bool { return false }
	}

	if kvExpr := result.Fields[SchemaFieldForceNew]; kvExpr != nil && astutils.ExprBoolValue(kvExpr.Value) != nil {
		result.Schema.ForceNew = *astutils.ExprBoolValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldInputDefault]; kvExpr != nil && astutils.ExprStringValue(kvExpr.Value) != nil {
		result.Schema.InputDefault = *astutils.ExprStringValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldMaxItems]; kvExpr != nil && astutils.ExprIntValue(kvExpr.Value) != nil {
		result.Schema.MaxItems = *astutils.ExprIntValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldMinItems]; kvExpr != nil && astutils.ExprIntValue(kvExpr.Value) != nil {
		result.Schema.MinItems = *astutils.ExprIntValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldOptional]; kvExpr != nil && astutils.ExprBoolValue(kvExpr.Value) != nil {
		result.Schema.Optional = *astutils.ExprBoolValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldPromoteSingle]; kvExpr != nil && astutils.ExprBoolValue(kvExpr.Value) != nil {
		result.Schema.PromoteSingle = *astutils.ExprBoolValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldRequired]; kvExpr != nil && astutils.ExprBoolValue(kvExpr.Value) != nil {
		result.Schema.Required = *astutils.ExprBoolValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldSensitive]; kvExpr != nil && astutils.ExprBoolValue(kvExpr.Value) != nil {
		result.Schema.Sensitive = *astutils.ExprBoolValue(kvExpr.Value)
	}

	if kvExpr := result.Fields[SchemaFieldStateFunc]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.StateFunc = func(interface{}) string { return "" }
	}

	if kvExpr := result.Fields[SchemaFieldValidateFunc]; kvExpr != nil && astutils.ExprValue(kvExpr.Value) != nil {
		result.Schema.ValidateFunc = func(interface{}, string) ([]string, []error) { return nil, nil }
	}

	return result
}

// DeclaresField returns true if the field name is present in the AST
func (info *SchemaInfo) DeclaresField(fieldName string) bool {
	return info.Fields[fieldName] != nil
}

// DeclaresBoolFieldWithZeroValue returns true if the field name is present and is false
func (info *SchemaInfo) DeclaresBoolFieldWithZeroValue(fieldName string) bool {
	kvExpr := info.Fields[fieldName]

	// Field not declared
	if kvExpr == nil {
		return false
	}

	valuePtr := astutils.ExprBoolValue(kvExpr.Value)

	// Value not readable
	if valuePtr == nil {
		return false
	}

	return !*valuePtr
}

// IsType returns true if the given input is equal to the Type
func (info *SchemaInfo) IsType(valueType string) bool {
	return info.SchemaValueType == valueType
}

// IsOneOfTypes returns true if one of the given input is equal to the Type
func (info *SchemaInfo) IsOneOfTypes(valueTypes ...string) bool {
	for _, valueType := range valueTypes {
		if info.SchemaValueType == valueType {
			return true
		}
	}

	return false
}

// GetSchemaMapAttributeNames returns all attribute names held in a map[string]*schema.Schema
func GetSchemaMapAttributeNames(cl *ast.CompositeLit) []ast.Expr {
	var result []ast.Expr

	for _, elt := range cl.Elts {
		switch v := elt.(type) {
		case *ast.KeyValueExpr:
			result = append(result, v.Key)
		}
	}

	return result
}

// GetSchemaMapSchemas returns all Schema held in a map[string]*schema.Schema
func GetSchemaMapSchemas(cl *ast.CompositeLit) []*ast.CompositeLit {
	var result []*ast.CompositeLit

	for _, elt := range cl.Elts {
		switch v := elt.(type) {
		case *ast.KeyValueExpr:
			switch v := v.Value.(type) {
			case *ast.CompositeLit:
				result = append(result, v)
			}
		}
	}

	return result
}

// IsTypeSchema returns if the type is Schema from the helper/schema package
func IsTypeSchema(t types.Type) bool {
	switch t := t.(type) {
	case *types.Named:
		return IsNamedType(t, TypeNameSchema)
	case *types.Pointer:
		return IsTypeSchema(t.Elem())
	default:
		return false
	}
}

// IsValueType returns if the Schema field Type matches
func IsValueType(e ast.Expr, info *types.Info) bool {
	switch e := e.(type) {
	case *ast.SelectorExpr:
		switch t := info.TypeOf(e).(type) {
		case *types.Named:
			return IsNamedType(t, TypeNameValueType)
		default:
			return false
		}
	default:
		return false
	}
}

// IsTypeSet returns if the type is Set from the helper/schema package
// Use IsTypeSchemaFieldType for verifying Type: schema.TypeSet ValueType
func IsTypeSet(t types.Type) bool {
	switch t := t.(type) {
	case *types.Named:
		return IsNamedType(t, TypeNameSet)
	case *types.Pointer:
		return IsTypeSet(t.Elem())
	default:
		return false
	}
}

// IsMapStringSchema returns if the type is map[string]*Schema from the helper/schema package
func IsMapStringSchema(cl *ast.CompositeLit, info *types.Info) bool {
	switch v := cl.Type.(type) {
	case *ast.MapType:
		switch k := v.Key.(type) {
		case *ast.Ident:
			if k.Name != "string" {
				return false
			}
		}

		return IsTypeSchema(info.TypeOf(v.Value))
	}

	return false
}

// typeSchemaType extracts the string representation of a Schema Type value
func typeSchemaType(schema *ast.CompositeLit, info *types.Info) string {
	kvExpr := astutils.CompositeLitField(schema, SchemaFieldType)

	if kvExpr == nil {
		return ""
	}

	if !IsValueType(kvExpr.Value, info) {
		return ""
	}

	return valueTypeString(kvExpr.Value)
}

// valueTypeString extracts the string representation of a Schema ValueType
func valueTypeString(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.SelectorExpr:
		return e.Sel.Name
	default:
		return ""
	}
}
