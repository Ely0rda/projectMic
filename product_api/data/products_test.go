package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "milo",
		Price: 7,
		SKU:   "abc-abc-",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
