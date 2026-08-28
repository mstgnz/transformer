package scalar

import "testing"

func TestIsJSONNumber(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"0", true},
		{"42", true},
		{"-17", true},
		{"3.14", true},
		{"-0.5", true},
		{"1e10", true},
		{"1E-10", true},
		{"9007199254740993", true},
		{"12345678901234567890", true},
		{"01", false},   // JSON has no leading zeros
		{"+5", false},   // JSON has no leading plus
		{".5", false},   // JSON needs the integer part
		{"0x1F", false}, // hexadecimal is YAML, not JSON
		{"1_000", false},
		{"Inf", false},
		{"NaN", false},
		{"true", false}, // json.Valid accepts it; a number it is not
		{"null", false},
		{"not-a-number", false},
		{"12 34", false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := IsJSONNumber(tt.in); got != tt.want {
				t.Errorf("IsJSONNumber(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNumber(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		want  string
		wantK bool
	}{
		{"plain integer is untouched", "42", "42", true},
		{"float is untouched", "3.14", "3.14", true},
		{"exponent keeps its notation", "1e10", "1e10", true},
		{"big integer keeps every digit", "9007199254740993", "9007199254740993", true},
		{"beyond int64 keeps every digit", "12345678901234567890", "12345678901234567890", true},
		{"hexadecimal becomes decimal", "0x1F", "31", true},
		{"octal becomes decimal", "0o17", "15", true},
		{"underscores are dropped", "1_000", "1000", true},
		{"leading plus is dropped", "+5", "5", true},
		{"leading zeros are dropped", "007", "7", true},
		{"infinity is not a number", ".inf", ".inf", false},
		{"nan is not a number", ".nan", ".nan", false},
		{"text is not a number", "not-a-number", "not-a-number", false},
		{"empty is not a number", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Number(tt.in)
			if ok != tt.wantK {
				t.Fatalf("Number(%q) ok = %v, want %v", tt.in, ok, tt.wantK)
			}
			if got != tt.want {
				t.Errorf("Number(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNumber_ResultIsAlwaysWritable(t *testing.T) {
	// Whatever Number reports as a number must be something JSON can write,
	// otherwise the writers have to second-guess it.
	inputs := []string{"42", "3.14", "1e10", "0x1F", "0o17", "1_000", "+5", "007", "-0.5"}
	for _, in := range inputs {
		if got, ok := Number(in); ok && !IsJSONNumber(got) {
			t.Errorf("Number(%q) = %q, which JSON cannot write", in, got)
		}
	}
}

func TestIsInteger(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"42", true},
		{"-17", true},
		{"12345678901234567890", true},
		{"3.14", false},
		{"1e10", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := IsInteger(tt.in); got != tt.want {
				t.Errorf("IsInteger(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		in    string
		want  bool
		wantK bool
	}{
		{"true", true, true},
		{"false", false, true},
		{"TRUE", true, true},
		{"False", false, true},
		{"yes", true, true},
		{"no", false, true},
		{"on", true, true},
		{"off", false, true},
		{"1", true, true},
		{"0", false, true},
		{"maybe", false, false},
		{"", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, ok := Bool(tt.in)
			if ok != tt.wantK {
				t.Fatalf("Bool(%q) ok = %v, want %v", tt.in, ok, tt.wantK)
			}
			if got != tt.want {
				t.Errorf("Bool(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestNumber_WiderForms(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		want  string
		wantK bool
	}{
		{"hexadecimal past int64", "0xFFFFFFFFFFFFFFFF", "18446744073709551615", true},
		{"leading dot gains its zero", ".5", "0.5", true},
		{"negative leading dot", "-.5", "-0.5", true},
		{"binary", "0b1010", "10", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Number(tt.in)
			if ok != tt.wantK {
				t.Fatalf("Number(%q) ok = %v, want %v", tt.in, ok, tt.wantK)
			}
			if got != tt.want {
				t.Errorf("Number(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if ok && !IsJSONNumber(got) {
				t.Errorf("Number(%q) = %q, which JSON cannot write", tt.in, got)
			}
		})
	}
}
