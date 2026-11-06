package main

type theme struct {
	name       string
	ink        string
	onAccent   string
	muted      string
	edge       string
	cell       string
	shade      string
	wall       string
	accent     string
	accentSoft string
	cool       string
	coolSoft   string
	leaf       string
	leafSoft   string
	sectors    [6]string
}

func themes() []theme {
	return []theme{
		{
			name:       "light",
			ink:        "#1c1713",
			onAccent:   "#1c1713",
			muted:      "#8a8178",
			edge:       "#cfc7be",
			cell:       "#fdfbf9",
			shade:      "#e6e0d9",
			wall:       "#4a433d",
			accent:     "#e8883a",
			accentSoft: "#f9dcc0",
			cool:       "#2f7fb8",
			coolSoft:   "#cfe4f3",
			leaf:       "#4f9a3c",
			leafSoft:   "#d6ebcb",
			sectors:    [6]string{"#f9dcc0", "#f4e7b5", "#d6ebcb", "#c9e8e4", "#cfe4f3", "#e4d9f2"},
		},
		{
			name:       "dark",
			ink:        "#f0e9e2",
			onAccent:   "#1c1713",
			muted:      "#a39a90",
			edge:       "#5a524b",
			cell:       "#221d19",
			shade:      "#15110f",
			wall:       "#d8cfc5",
			accent:     "#e8883a",
			accentSoft: "#70451f",
			cool:       "#6bb4e3",
			coolSoft:   "#234a66",
			leaf:       "#7cc46a",
			leafSoft:   "#2f4f27",
			sectors:    [6]string{"#70451f", "#5c5121", "#2f4f27", "#1f4d48", "#234a66", "#44345e"},
		},
	}
}
