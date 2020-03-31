package nerdgraph

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ResolveSchemaTypes(schema Schema, typeNames []string) (map[string]string, error) {
	typeKeeper := make(map[string]string)

	for _, typeName := range typeNames {
		typeGenResult, err := TypeGen(schema, typeName)
		if err != nil {
			log.Errorf("error while generating type %s: %s", typeName, err)
		}

		for k, v := range typeGenResult {
			typeKeeper[k] = v
		}
	}

	return typeKeeper, nil
}

func handleEnumType(schema Schema, t SchemaType) map[string]string {
	types := make(map[string]string)

	// output collects each line of a struct type
	output := []string{}

	// Add a comment for golint to ignore
	output = append(output, "")

	output = append(output, "// nolint:golint")
	output = append(output, fmt.Sprintf("type %s string ", t.Name))
	output = append(output, "")

	output = append(output, "const (")
	for _, v := range t.EnumValues {

		if v.Description != "" {
			output = append(output, fmt.Sprintf("\t /* %s */", v.Description))
		}

		output = append(output, fmt.Sprintf("\t%s %s = \"%s\" // nolint:golint", v.Name, t.Name, v.Name))
	}

	output = append(output, ")")
	output = append(output, "")

	types[t.Name] = strings.Join(output, "\n")

	return types
}

func kindTree(t SchemaTypeRef) []string {
	tree := []string{}

	if t.Kind != "" {
		tree = append(tree, t.Kind)
	}

	if t.OfType.Kind != "" {
		tree = append(tree, t.OfType.Kind)
	}

	if t.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.Kind)
	}

	if t.OfType.OfType.OfType.OfType.OfType.OfType.Kind != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.OfType.Kind)
	}

	return tree
}

func nameTree(t SchemaTypeRef) []string {
	tree := []string{}

	if t.Name != "" {
		tree = append(tree, t.Name)
	}

	if t.OfType.Name != "" {
		tree = append(tree, t.OfType.Name)
	}

	if t.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.Name)
	}

	if t.OfType.OfType.OfType.OfType.OfType.OfType.Name != "" {
		tree = append(tree, t.OfType.OfType.OfType.OfType.OfType.OfType.Name)
	}

	return tree
}

func removeNonNullValues(tree []string) []string {
	a := []string{}

	for _, x := range tree {
		if x != "NON_NULL" {
			a = append(a, x)
		}
	}

	return a
}

// fieldTypeFromTypeRef resolves the given SchemaInputValue into a field name to use on a go struct.
func fieldTypeFromTypeRef(t SchemaTypeRef) (string, bool, error) {

	switch t := nameTree(t)[0]; t {
	case "String":
		return "string", false, nil
	case "Int":
		return "int", false, nil
	case "Boolean":
		return "bool", false, nil
	case "Float":
		return "float64", false, nil
	case "ID":
		// ID is a nested object, but behaves like an integer.  This may be true of other SCALAR types as well, so logic here could potentially be moved.
		return "int", false, nil
	case "":
		return "", true, fmt.Errorf("empty field name: %+v", t)
	default:
		return t, true, nil
	}
}

// handleObjectType will operate on a SchemaType who's Kind is OBJECT or INPUT_OBJECT.
func handleObjectType(schema Schema, t SchemaType) map[string]string {
	types := make(map[string]string)
	var err error
	recurse := false

	// output collects each line of a struct type
	output := []string{}

	// Add a comment for golint to ignore
	output = append(output, "// nolint:golint")

	output = append(output, fmt.Sprintf("type %s struct {", t.Name))

	// Fill in the struct fields for an input type
	for _, f := range t.InputFields {
		var fieldType string

		log.Debugf("handling kind %s: %+v\n\n", f.Type.Kind, f)
		fieldType, recurse, err = fieldTypeFromTypeRef(f.Type)
		if err != nil {
			// If we have an error, then we don't know how to handle the type to
			// determine the field name.  This indicates that
			log.Errorf("error resolving first non-empty name from field: %s: %s", f, err)
		}

		if recurse {
			// The name of the nested sub-type.  We take the first value here as the root name for the nested type.
			subTName := nameTree(f.Type)[0]

			subT, err := typeByName(schema, subTName)
			if err != nil {
				log.Warnf("non_null: unhandled type: %+v\n", f)
				continue
			}

			// Determnine if we need to resolve the sub type, or if it already
			// exists in the map.
			if _, ok := types[subT.Name]; !ok {
				result, err := TypeGen(schema, subT.Name)
				if err != nil {
					log.Errorf("ERROR while resolving sub type %s: %s\n", subT.Name, err)
				}

				log.Debugf("resolved type result:\n%+v\n", result)

				for k, v := range result {
					if _, ok := types[k]; !ok {
						types[k] = v
					}
				}
			}

			fieldType = subT.Name
		}

		var fieldName string

		if f.Name == "ids" {
			// special case to avoid the struct field Ids, and prefer IDs instead
			fieldName = "IDs"
		} else {
			fieldName = strings.Title(f.Name)
		}

		// The prefix is used to ensure that we handle LIST or slices correctly.
		fieldTypePrefix := ""

		if removeNonNullValues(kindTree(f.Type))[0] == "LIST" {
			fieldTypePrefix = "[]"
		}

		// Include some documentation
		if f.Description != "" {
			output = append(output, fmt.Sprintf("\t /* %s */", f.Description))
		}

		var fieldTags string
		if f.Name == "id" {
			fieldTags = fmt.Sprintf("`json:\"%s,string\"`", f.Name)
		} else {
			fieldTags = fmt.Sprintf("`json:\"%s\"`", f.Name)
		}

		output = append(output, fmt.Sprintf("\t %s %s%s %s", fieldName, fieldTypePrefix, fieldType, fieldTags))
		output = append(output, "")
	}

	for _, f := range t.Fields {
		log.Warnf("Fields: %+v", f)
		output = append(output, "")
		output = append(output, lineForField(f.Name, f.Description, f.Type)...)
	}

	// Close the struct
	output = append(output, "}\n")
	types[t.Name] = strings.Join(output, "\n")

	return types
}

func lineForField(name string, description string, typeRef SchemaTypeRef) []string {
	var output []string
	var fieldName string

	log.Infof("handling kind %s: %+v", typeRef.Kind, typeRef)
	fieldType, _, err := fieldTypeFromTypeRef(typeRef)
	if err != nil {
		// If we have an error, then we don't know how to handle the type to
		// determine the field name.  This indicates that
		log.Errorf("error resolving first non-empty name from field: %s: %s", typeRef, err)
	}

	// TODO determine if we need to handle the subT name as is in the
	// handleObjectType() method above.  If we do indded have a reason to handle
	// it, this code matches pretty close to what is above.  With a little love,
	// we could DRY this up a bit.

	if name == "ids" {
		// special case to avoid the struct field Ids, and prefer IDs instead
		fieldName = "IDs"
	} else if name == "id" {
		fieldName = "ID"
	} else if name == "accountId" {
		fieldName = "AccountID"
	} else {
		fieldName = strings.Title(name)
	}

	fieldTypePrefix := ""

	if removeNonNullValues(kindTree(typeRef))[0] == "LIST" {
		fieldTypePrefix = "[]"
	}

	// Include some documentation
	if description != "" {
		output = append(output, fmt.Sprintf("\t /* %s */", description))
	}

	var fieldTags string
	if name == "id" {
		fieldTags = fmt.Sprintf("`json:\"%s,string\"`", name)
	} else {
		fieldTags = fmt.Sprintf("`json:\"%s\"`", name)
	}

	output = append(output, fmt.Sprintf("\t %s %s%s %s", fieldName, fieldTypePrefix, fieldType, fieldTags))

	return output
}

// TypeGen is the mother type generator.
func TypeGen(schema Schema, typeName string) (map[string]string, error) {

	// The total known types.  Keyed by the typeName, and valued as the string
	// output that one would write to a file where Go structs are kept.
	types := make(map[string]string)

	t, err := typeByName(schema, typeName)
	if err != nil {
		log.Error(err)
	}

	log.Infof("starting on %s: %+v", typeName, t.Kind)

	// To store the results from the single
	results := make(map[string]string)

	if t.Kind == "INPUT_OBJECT" || t.Kind == "OBJECT" {
		results = handleObjectType(schema, *t)
	} else if t.Kind == "ENUM" {
		results = handleEnumType(schema, *t)
	} else {
		log.Warnf("WARN: unhandled object Kind: %s\n", t.Kind)
	}

	for k, v := range results {
		types[k] = v
	}

	// return strings.Join(output, "\n"), nil
	return types, nil
}

func typeByName(schema Schema, typeName string) (*SchemaType, error) {
	log.Debugf("looking for typeName: %s", typeName)

	for _, t := range schema.Types {
		if t.Name == typeName {
			return t, nil
		}
	}

	return nil, fmt.Errorf("type by name %s not found", typeName)
}
