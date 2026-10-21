package palindrome

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"racecar", true},
		{"Racecar", true},
		{"A man, a plan, a canal: Panama!", true},
		{"No 'x' in Nixon", true},
		{"Was it a car or a cat I saw?", true},
		{"Gopher", false},
		{"ab", false},
		{"été", true},
		{"éte", false},
		{"Ésope reste ici et se repose", false}, // é ≠ e
		{"上海自来水来自海上", true},
		{"世界", false},
		{"12321", true},
		{"123", false},
		{"1a2", false},
		{"a", true},
		{"", true},
		{"?!", true},
		{"🙂 a 🙃", true}, // emoji are neither letters nor digits
	}
	for _, tt := range tests {
		if got := IsPalindrome(tt.s); got != tt.want {
			t.Errorf("IsPalindrome(%q) = %t, want %t", tt.s, got, tt.want)
		}
	}
}
