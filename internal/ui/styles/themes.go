package styles

import 

type Palette struct {
	Name   string
	GradA  string
	GradB  string
	Accent string
	Ok     string
	Warn   string
	Bad    string
	Dim    string
	Text   string
}

var Themes = []Palette{
	{Name: "crush", GradA: "#ff5ff0", GradB: "#6b50ff", Accent: "#c86bff", Ok: "#5fffaf", Warn: "#ffd75f", Bad: "#ff5f87", Dim: "#6f6f8f", Text: "#e6e6f2"},
	{Name: "ocean", GradA: "#00e5ff", GradB: "#2962ff", Accent: "#4fc3f7", Ok: "#69f0ae", Warn: "#ffe082", Bad: "#ff8a80", Dim: "#607d9b", Text: "#e3f2fd"},
	{Name: "sunset", GradA: "#ffb347", GradB: "#ff2f92", Accent: "#ff8a65", Ok: "#a5d6a7", Warn: "#ffe57f", Bad: "#ff5252", Dim: "#8d6e63", Text: "#fff3e0"},
	{Name: "matrix", GradA: "#d4ff00", GradB: "#00c853", Accent: "#69f0ae", Ok: "#00e676", Warn: "#ffee58", Bad: "#ff5252", Dim: "#3e7a4f", Text: "#d7ffd9"},
	{Name: "dracula", GradA: "#bd93f9", GradB: "#ff79c6", Accent: "#8be9fd", Ok: "#50fa7b", Warn: "#f1fa8c", Bad: "#ff5555", Dim: "#6272a4", Text: "#f8f8f2"},
	{Name: "aurora", GradA: "#00ffa3", GradB: "#a259ff", Accent: "#7cf7d4", Ok: "#5fffaf", Warn: "#ffd75f", Bad: "#ff6b9d", Dim: "#5f6f8f", Text: "#e8fff6"},
}

var currentThemeIdx = 0

// CycleTheme switches to the next theme and returns its name
func CycleTheme() string {
	currentThemeIdx = (currentThemeIdx + 1) % len(Themes)
	ApplyTheme(Themes[currentThemeIdx].Name)
	return Themes[currentThemeIdx].Name
}

func GetCurrentPalette() Palette {
	return Themes[currentThemeIdx]
}
