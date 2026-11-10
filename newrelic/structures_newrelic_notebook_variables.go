package newrelic

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func notebookVariableSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"from": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"list", "query", "input"}, false),
				Description:  `Variable type. One of "list", "query", or "input".`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier used in {{name}} template syntax inside widget queries.",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Display name shown in the variables bar. Use an empty string for no label.",
			},
			"as": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "string",
				ValidateFunc: validation.StringInSlice([]string{"string", "number", "identifier"}, false),
				Description:  `Output format for the variable value when substituted in queries. One of "string", "number", or "identifier". Defaults to "string".`,
			},
			"hidden": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, hides the variable from the variables bar entirely.",
			},
			"omit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `When true, the variable value is not injected into queries (equivalent to the "Include variable" toggle off in the UI).`,
			},
			"defaults": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Default values pre-selected when the notebook loads. Omit this block entirely when no default is desired. Use ["*"] for "Select all" on multi-value list or query variables.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// ── List-specific ────────────────────────────────────────────────────────
			"item": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Static choice items for a list variable. Required when from = "list".`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Item value. Always a string even when as = "number"; conversion happens at query substitution time.`,
						},
					},
				},
			},
			// ── Query-specific ───────────────────────────────────────────────────────
			"nrql_query": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `NRQL query that dynamically populates variable choices. Required when from = "query".`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The NRQL query string. May reference other variables with {{varName}} syntax.",
						},
						"account_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Account IDs to run the query against.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"ignore_time_range": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `When true (default), the variable query ignores the notebook time picker and uses the query's own time range. Only valid when from = "query".`,
			},
			// ── List + Query ─────────────────────────────────────────────────────────
			"multi": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Enables multi-value selection so users can pick more than one option. Only valid when from = "list" or "query".`,
			},
			"defer_execution": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Shows an Apply button when multi-selecting instead of applying each selection immediately. Only valid when from = "list" or "query". Defaults to true.`,
			},
		},
	}
}

// expandNotebookVariables converts `variable` blocks from ResourceData into a
// []interface{} ready for JSON marshaling as the blob's top-level "variables" key.
func expandNotebookVariables(d *schema.ResourceData) ([]interface{}, error) {
	raw, ok := d.GetOk("variable")
	if !ok {
		return nil, nil
	}

	blocks := raw.([]interface{})
	vars := make([]interface{}, 0, len(blocks))

	for i, b := range blocks {
		block := b.(map[string]interface{})
		from, _ := block["from"].(string)

		v := map[string]interface{}{
			"from":   from,
			"name":   block["name"].(string),
			"label":  block["label"].(string),
			"as":     block["as"].(string),
			"omit":   block["omit"].(bool),
			"hidden": block["hidden"].(bool),
		}

		// defaults: omit the key entirely when no defaults are configured
		if defs, ok := block["defaults"].([]interface{}); ok && len(defs) > 0 {
			strs := make([]string, len(defs))
			for j, d := range defs {
				strs[j], _ = d.(string)
			}
			v["defaults"] = strs
		}

		switch from {
		case "list":
			items, err := expandNotebookVarItems(block, i)
			if err != nil {
				return nil, err
			}
			v["items"] = items
			v["multi"] = block["multi"].(bool)
			v["defer"] = block["defer_execution"].(bool)

		case "query":
			queries, err := expandNotebookVarNRQLQueries(block, i)
			if err != nil {
				return nil, err
			}
			v["nrqlQueries"] = queries
			v["options"] = map[string]interface{}{
				"ignoreTimeRange": block["ignore_time_range"].(bool),
			}
			v["multi"] = block["multi"].(bool)
			v["defer"] = block["defer_execution"].(bool)

		case "input":
			v["type"] = "text"
		}

		vars = append(vars, v)
	}

	return vars, nil
}

func expandNotebookVarItems(block map[string]interface{}, varIdx int) ([]map[string]interface{}, error) {
	raw, _ := block["item"].([]interface{})
	if len(raw) == 0 {
		return nil, fmt.Errorf(`variable[%d]: from = "list" requires at least one item block`, varIdx)
	}
	items := make([]map[string]interface{}, len(raw))
	for i, r := range raw {
		m := r.(map[string]interface{})
		items[i] = map[string]interface{}{"value": m["value"].(string)}
	}
	return items, nil
}

func expandNotebookVarNRQLQueries(block map[string]interface{}, varIdx int) ([]map[string]interface{}, error) {
	raw, _ := block["nrql_query"].([]interface{})
	if len(raw) == 0 {
		return nil, fmt.Errorf(`variable[%d]: from = "query" requires at least one nrql_query block`, varIdx)
	}
	queries := make([]map[string]interface{}, len(raw))
	for i, r := range raw {
		m := r.(map[string]interface{})
		rawIDs, _ := m["account_ids"].([]interface{})
		if len(rawIDs) == 0 {
			return nil, fmt.Errorf("variable[%d].nrql_query[%d]: account_ids must contain at least one account ID", varIdx, i)
		}
		ids := make([]int, len(rawIDs))
		for j, id := range rawIDs {
			ids[j] = id.(int)
		}
		q, _ := m["query"].(string)
		if q == "" {
			return nil, fmt.Errorf("variable[%d].nrql_query[%d]: query must not be empty", varIdx, i)
		}
		queries[i] = map[string]interface{}{
			"query":      q,
			"accountIds": ids,
		}
	}
	return queries, nil
}

// flattenNotebookVariables reads the "variables" array from a decoded blob
// document and writes it into the `variable` blocks in state.
func flattenNotebookVariables(vars []interface{}, d *schema.ResourceData) error {
	if len(vars) == 0 {
		return d.Set("variable", []interface{}{})
	}

	blocks := make([]interface{}, 0, len(vars))
	for _, raw := range vars {
		v, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		from, _ := v["from"].(string)

		block := map[string]interface{}{
			"from":   from,
			"name":   stringOrEmptyNb(v["name"]),
			"label":  stringOrEmptyNb(v["label"]),
			"as":     stringOrEmptyNb(v["as"]),
			"hidden": boolOrFalseNb(v["hidden"]),
			"omit":   boolOrFalseNb(v["omit"]),
		}

		// defaults: absent in blob → empty list in state
		if defs, ok := v["defaults"].([]interface{}); ok && len(defs) > 0 {
			strs := make([]interface{}, len(defs))
			for i, d := range defs {
				strs[i], _ = d.(string)
			}
			block["defaults"] = strs
		} else {
			block["defaults"] = []interface{}{}
		}

		switch from {
		case "list":
			block["item"] = flattenNotebookVarItems(v["items"])
			block["multi"] = boolOrFalseNb(v["multi"])
			block["defer_execution"] = boolOrTrueNb(v["defer"])
			block["nrql_query"] = []interface{}{}
			block["ignore_time_range"] = true

		case "query":
			block["nrql_query"] = flattenNotebookVarNRQLQueries(v["nrqlQueries"])
			block["multi"] = boolOrFalseNb(v["multi"])
			block["defer_execution"] = boolOrTrueNb(v["defer"])
			block["ignore_time_range"] = flattenIgnoreTimeRange(v["options"])
			block["item"] = []interface{}{}

		case "input":
			block["item"] = []interface{}{}
			block["nrql_query"] = []interface{}{}
			block["multi"] = false
			block["defer_execution"] = true
			block["ignore_time_range"] = true
		}

		blocks = append(blocks, block)
	}

	return d.Set("variable", blocks)
}

func flattenNotebookVarItems(raw interface{}) []interface{} {
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(arr))
	for _, r := range arr {
		m, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"value": stringOrEmptyNb(m["value"]),
		})
	}
	return out
}

func flattenNotebookVarNRQLQueries(raw interface{}) []interface{} {
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(arr))
	for _, r := range arr {
		m, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		rawIDs, _ := m["accountIds"].([]interface{})
		ids := make([]interface{}, len(rawIDs))
		for i, id := range rawIDs {
			switch v := id.(type) {
			case float64:
				ids[i] = int(v)
			case int:
				ids[i] = v
			default:
				ids[i] = 0
			}
		}
		out = append(out, map[string]interface{}{
			"query":       stringOrEmptyNb(m["query"]),
			"account_ids": ids,
		})
	}
	return out
}

func flattenIgnoreTimeRange(opts interface{}) bool {
	m, ok := opts.(map[string]interface{})
	if !ok {
		return true
	}
	v, ok := m["ignoreTimeRange"].(bool)
	if !ok {
		return true
	}
	return v
}

// validateNotebookVariables validates variable maps before writing to the blob API.
// Returns plain error strings for pooling with other validation results.
func validateNotebookVariables(vars []map[string]interface{}) []string {
	var errs []string
	validFrom := map[string]bool{"list": true, "query": true, "input": true}
	validAs := map[string]bool{"string": true, "number": true, "identifier": true}

	for i, v := range vars {
		path := fmt.Sprintf("variable[%d]", i)

		name, _ := v["name"].(string)
		if name == "" {
			errs = append(errs, fmt.Sprintf("%s: name must not be empty", path))
		}

		from, _ := v["from"].(string)
		if !validFrom[from] {
			errs = append(errs, fmt.Sprintf(`%s: from must be one of "list", "query", "input" (lowercase), got %q`, path, from))
			continue
		}

		if as, _ := v["as"].(string); as != "" && !validAs[as] {
			errs = append(errs, fmt.Sprintf(`%s: as must be one of "string", "number", "identifier" (lowercase), got %q`, path, as))
		}

		switch from {
		case "list":
			items, _ := v["items"].([]interface{})
			if len(items) == 0 {
				errs = append(errs, fmt.Sprintf(`%s: from = "list" requires at least one item`, path))
			}
			if _, hasNRQL := v["nrqlQueries"]; hasNRQL {
				errs = append(errs, fmt.Sprintf(`%s: nrqlQueries must not be set on a list variable`, path))
			}

		case "query":
			queries, _ := v["nrqlQueries"].([]interface{})
			if len(queries) == 0 {
				errs = append(errs, fmt.Sprintf(`%s: from = "query" requires at least one nrqlQueries entry`, path))
			} else {
				for qi, q := range queries {
					qm, ok := q.(map[string]interface{})
					if !ok {
						errs = append(errs, fmt.Sprintf(`%s.nrqlQueries[%d]: must be an object with query and accountIds`, path, qi))
						continue
					}
					if qs, _ := qm["query"].(string); qs == "" {
						errs = append(errs, fmt.Sprintf(`%s.nrqlQueries[%d]: query must not be empty`, path, qi))
					}
					ids, _ := qm["accountIds"].([]interface{})
					if len(ids) == 0 {
						errs = append(errs, fmt.Sprintf(`%s.nrqlQueries[%d]: accountIds must contain at least one account ID`, path, qi))
					}
				}
			}
			if _, hasItems := v["items"]; hasItems {
				errs = append(errs, fmt.Sprintf(`%s: items must not be set on a query variable`, path))
			}

		case "input":
			if _, hasItems := v["items"]; hasItems {
				errs = append(errs, fmt.Sprintf(`%s: items must not be set on an input variable`, path))
			}
			if _, hasNRQL := v["nrqlQueries"]; hasNRQL {
				errs = append(errs, fmt.Sprintf(`%s: nrqlQueries must not be set on an input variable`, path))
			}
		}

		if defs, ok := v["defaults"]; ok && defs != nil {
			defsArr, isArr := defs.([]interface{})
			if !isArr {
				errs = append(errs, fmt.Sprintf(`%s: defaults must be an array of strings`, path))
			} else {
				if from == "input" && len(defsArr) > 1 {
					errs = append(errs, fmt.Sprintf(`%s: input variables support at most one default value`, path))
				}
				for di, d := range defsArr {
					if _, ok := d.(string); !ok {
						errs = append(errs, fmt.Sprintf(`%s.defaults[%d]: must be a string`, path, di))
					}
				}
			}
		}
	}

	return errs
}

// mergeVariablesIntoBlob injects or removes the "variables" key in a decoded
// blob document map. Pass nil vars to remove the key (no variables configured).
func mergeVariablesIntoBlob(doc map[string]interface{}, vars []interface{}) map[string]interface{} {
	if len(vars) == 0 {
		delete(doc, "variables")
	} else {
		doc["variables"] = vars
	}
	return doc
}

// extractVariablesFromBlob parses a raw JSON blob string and returns the
// top-level "variables" array. Returns nil when the key is absent.
func extractVariablesFromBlob(raw string) ([]interface{}, error) {
	var doc map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil, fmt.Errorf("invalid blob JSON: %w", err)
	}
	v, ok := doc["variables"]
	if !ok || v == nil {
		return nil, nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf(`blob "variables" field is not an array`)
	}
	return arr, nil
}

// ── Type coercion helpers ─────────────────────────────────────────────────────
// Suffixed with "Nb" to avoid collision with any existing helpers in the package.

func stringOrEmptyNb(v interface{}) string {
	s, _ := v.(string)
	return s
}

func boolOrFalseNb(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

func boolOrTrueNb(v interface{}) bool {
	if v == nil {
		return true
	}
	b, ok := v.(bool)
	if !ok {
		return true
	}
	return b
}
