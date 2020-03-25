package nerdgraph

import (
	"fmt"
	"strings"
)

func ResolveSchemaTypes(schema Schema, typeNames []string) (map[string]string, error) {
	typeKeeper := make(map[string]string)

	for _, typeName := range typeNames {
		typeGenResult, err := TypeGen(schema, typeName)
		if err != nil {
			fmt.Printf("ERROR while generating type %s: %s", typeName, err)
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
	output = append(output, fmt.Sprintf("type %s int ", t.Name))
	output = append(output, "")

	output = append(output, "const (")
	for i, v := range t.EnumValues {
		if i == 0 {
			output = append(output, fmt.Sprintf("\t%s = iota", v.Name))
		} else {
			output = append(output, fmt.Sprintf("\t%s", v.Name))
		}

	}
	output = append(output, ")")
	output = append(output, "")

	types[t.Name] = strings.Join(output, "\n")

	return types
}

// firstTypeName returns the first non-empty name in the type tree that is found.
func firstTypeName(f SchemaInputValue) string {
	if f.Type.Name != "" {
		return f.Type.Name
	} else if f.Type.OfType.Name != "" {
		return f.Type.OfType.Name
	} else if f.Type.OfType.OfType.Name != "" {
		return f.Type.OfType.OfType.Name
	} else if f.Type.OfType.OfType.OfType.Name != "" {
		return f.Type.OfType.OfType.OfType.Name
	} else if f.Type.OfType.OfType.OfType.OfType.Name != "" {
		return f.Type.OfType.OfType.OfType.OfType.Name
	} else if f.Type.OfType.OfType.OfType.OfType.OfType.Name != "" {
		return f.Type.OfType.OfType.OfType.OfType.OfType.Name
	} else if f.Type.OfType.OfType.OfType.OfType.OfType.OfType.Name != "" {
		return f.Type.OfType.OfType.OfType.OfType.OfType.OfType.Name
	}

	return ""
}

// fieldTypeFromTypeRef resolves the given SchemaInputValue into a field name to use on a go struct.
func fieldTypeFromTypeRef(f SchemaInputValue) (string, bool, error) {

	switch t := firstTypeName(f); t {
	case "String":
		return "string", false, nil
	case "Int":
		return "int", false, nil
	case "Boolean":
		return "bool", false, nil
	case "Float":
		return "float64", false, nil
	case "":
		return "", true, fmt.Errorf("empty field name: %+v", f)
	}

	return "", true, fmt.Errorf("need to handle Field f: %+v", f)
}

func handleInputType(schema Schema, t SchemaType) map[string]string {
	types := make(map[string]string)
	var err error
	recurse := false

	// output collects each line of a struct type
	output := []string{}

	output = append(output, fmt.Sprintf("type %s struct {", t.Name))

	// Fill in the struct fields for an input type
	for _, f := range t.InputFields {
		var fieldType string

		fmt.Printf("handling kind %s: %+v\n\n", f.Type.Kind, f)
		fieldType, recurse, err = fieldTypeFromTypeRef(f)
		if err != nil {
			// If we have an error, then we don't know how to handle the type to
			// determine the field name.  This indicates that
			fmt.Printf("error resolving first non-empty name from field: %s: %s", f, err)
		}

		if recurse {
			var subTName string

			// Determine where to start.  For NON_NULL type, the name will be
			// empty, so we start our search at the nested OfType name.
			if f.Type.Kind == "NON_NULL" {
				subTName = firstTypeName(f)
			} else {
				subTName = f.Type.Name
			}

			subT, err := typeByName(schema, subTName)
			if err != nil {
				fmt.Printf("non_null: unhandled type: %+v\n", f)
				continue
			}

			// Determnine if we need to resolve the sub type, or if it already
			// exists in the map.
			if _, ok := types[subT.Name]; !ok {
				result, err := TypeGen(schema, subT.Name)
				if err != nil {
					fmt.Printf("ERROR while resolving sub type %s: %s\n", subT.Name, err)
				}

				fmt.Printf("resolved type result:\n%+v\n", result)

				for k, v := range result {
					if _, ok := types[k]; !ok {
						types[k] = v
					}
				}
			}

			fieldType = subT.Name
		}

		fieldName := strings.Title(f.Name)

		// Include some documentation
		if f.Description != "" {
			output = append(output, fmt.Sprintf("\t /* %s */", f.Description))
		}

		fieldTags := fmt.Sprintf("`json:\"%s\"`", f.Name)

		output = append(output, fmt.Sprintf("\t %s %s %s", fieldName, fieldType, fieldTags))
		output = append(output, "")
	}

	for _, f := range t.EnumValues {
		fmt.Printf("\n\nEnums: %+v\n", f)
	}

	for _, f := range t.Fields {
		fmt.Printf("\n\nFields: %+v\n", f)
	}

	// Close the struct
	output = append(output, "}\n")
	types[t.Name] = strings.Join(output, "\n")

	return types
}

// TypeGen is the mother type generator.
func TypeGen(schema Schema, typeName string) (map[string]string, error) {

	// The total known types.  Keyed by the typeName, and valued as the string
	// output that one would write to a file where Go structs are kept.
	types := make(map[string]string)

	t, err := typeByName(schema, typeName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n\n\nSTARTING on %s\n%+v\n", typeName, t.Kind)

	// To store the results from the single
	results := make(map[string]string)

	if t.Kind == "INPUT_OBJECT" {
		results = handleInputType(schema, *t)
	} else if t.Kind == "ENUM" {
		results = handleEnumType(schema, *t)
		// } else if t.Kind == "OBJECT" {
		// TODO
	} else {
		fmt.Printf("WARN: unhandled object Kind: %s\n", t.Kind)
	}

	for k, v := range results {
		types[k] = v
	}

	// return strings.Join(output, "\n"), nil
	return types, nil
}

func typeByName(schema Schema, typeName string) (*SchemaType, error) {
	// fmt.Printf("looking for type %s\n", typeName)

	for _, t := range schema.Types {
		// fmt.Printf("checking name %s\n", t.Name)
		if t.Name == typeName {
			return t, nil
		}
	}

	return nil, fmt.Errorf("type by name %s not found", typeName)
}
