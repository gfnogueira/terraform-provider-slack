package provider

import "testing"

func TestProvider(t *testing.T) {
	p := Provider()
	if p == nil {
		t.Fatal("provider returned nil")
	}
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if _, ok := p.Schema["token"]; !ok {
		t.Errorf("token attribute missing in provider schema")
	}
}