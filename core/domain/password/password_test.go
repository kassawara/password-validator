package password

import (
	"testing"
)

func TestPasswordValidation(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		isValid bool
	}{
		{"valid password", "Abcdef1!2", true},
		{"valid with spaces", "  Abc def1!2  ", true},
		{"too short", "Abc1!2", false},
		{"too short with spaces", "A b c 1 ! 2", false},
		{"no digit", "Abcdefgh!", false},
		{"no digit, spaces", "Abc def gh !", false},
		{"no lower", "ABCDEFG1!", false},
		{"no lower, spaces", "A B C D E F G 1 !", false},
		{"no upper", "abcdefg1!", false},
		{"no upper, spaces", "a b c d e f g 1 !", false},
		{"no special", "Abcdefg12", false},
		{"no special, spaces", "A b c d e f g 1 2", false},
		{"repeated char", "Abcdef1!A", false},
		{"repeated char, spaces", "A b c d e f 1 ! A", false},
		{"all rules fail", "abc", false},
		{"all rules fail, spaces", "a b c", false},
		{"empty", "", false},
		{"only spaces", "         ", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, _ := New(func(p *Password) { p.password = tc.value })
			if p.isValid != tc.isValid {
				t.Errorf("expected valid=%v, got %v for password '%s'", tc.isValid, p.isValid, tc.value)
			}
		})
	}
}
