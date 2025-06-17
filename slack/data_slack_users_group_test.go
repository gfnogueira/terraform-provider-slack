package slack

import "testing"

func TestHashStringSlice(t *testing.T) {
	cases := []struct {
		name string
		data []string
		want string
	}{
		{"empty", []string{}, "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"single", []string{"a"}, "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8"},
		{"multi", []string{"a", "b"}, "da23614e02469a0d7c7bd1bdab5c9c474b1904dc"},
	}
	for _, tt := range cases {
		if got := hashStringSlice(tt.data); got != tt.want {
			t.Errorf("%s: expected %s, got %s", tt.name, tt.want, got)
		}
	}
}

func TestConvertStringSliceToInterface(t *testing.T) {
	in := []string{"a", "b"}
	got := convertStringSliceToInterface(in)
	if len(got) != len(in) {
		t.Fatalf("unexpected length: %d", len(got))
	}
	for i, v := range in {
		if got[i] != v {
			t.Errorf("index %d: expected %s, got %v", i, v, got[i])
		}
	}
}