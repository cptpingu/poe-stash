package generate

import (
	"testing"
)

func mapEqual(a, b map[string][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok { //|| v != w {
			_, _ = v, w
			return false
		}
	}
	return true
}

// TestPoEMarkup tests parsing custom langage works.
func TestPoEMarkup(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<gemitem>{Vaal lightning orb}",
			expected: `<span class="PoEMarkup gemitem">Vaal lightning orb</span>`,
		},
		{
			input:    "<size:31>{Vaal lightning orb}",
			expected: `<span class="PoEMarkup" style="font-size:15.5px">Vaal lightning orb</span>`,
		},
		{
			input:    "<gemitem>{<size:31>{Vaal lightning orb}}",
			expected: `<span class="PoEMarkup gemitem"><span class="PoEMarkup" style="font-size:15.5px">Vaal lightning orb</span></span>`,
		},
		{
			input:    "<size:31>{<gemitem>{Vaal lightning orb}}",
			expected: `<span class="PoEMarkup" style="font-size:15.5px"><span class="PoEMarkup gemitem">Vaal lightning orb</span></span>`,
		},
		{
			input:    "before <gemitem>{Vaal lightning orb} after",
			expected: `before <span class="PoEMarkup gemitem">Vaal lightning orb</span> after`,
		},
		{
			input:    "before <gemitem>{Vaal lightning orb} after <rareItem>{Unique ring}",
			expected: `before <span class="PoEMarkup gemitem">Vaal lightning orb</span> after <span class="PoEMarkup rareItem">Unique ring</span>`,
		},
		{
			input:    "<gemitem>{Vaal <rareItem>{lightning} orb}",
			expected: `<span class="PoEMarkup gemitem">Vaal <span class="PoEMarkup rareItem">lightning</span> orb</span>`,
		},
		{
			input:    "",
			expected: ``,
		},
		{
			input:    "simple text",
			expected: `simple text`,
		},
		{
			input:    "single open <",
			expected: tokenErr,
		},
		{
			input:    "single open >",
			expected: `single open >`,
		},
		{
			input:    "single open <item>{",
			expected: tokenErr,
		},
		{
			input:    "single open <item>{{}",
			expected: tokenErr,
		},
		{
			input:    "single open <item>{}}",
			expected: `single open <span class="PoEMarkup item"></span>}`,
		},
		{
			input:    "<size:toto>{test}",
			expected: tokenErr,
		},
	}

	for _, current := range tests {
		res := ReplacePoEMarkup(current.input, false)
		if res != current.expected {
			t.Errorf("\n\texpected:\n%v\n\tbut got:\n%v\n", current.expected, res)
		}
	}
}
