package apm

// Label represents a New Relic label.
type Label struct {
	Key      string     `json:"key,omitempty"`
	Category string     `json:"category,omitempty"`
	Name     string     `json:"name,omitempty"`
	Links    LabelLinks `json:"links,omitempty"`
}

// LabelLinks represents external references on the Label.
type LabelLinks struct {
	Applications []int `json:"applications"`
	Servers      []int `json:"servers"`
}
