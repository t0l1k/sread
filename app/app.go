package app

import "github.com/t0l1k/sread/ui"

func NewApp() *ui.Ui {
	a := ui.GetUi()
	a.SetupSettings(LoadSettings())
	a.SetupTheme(NewTheme())
	a.SetupScreen("Rapid Read")
	a.Push(NewRRLoadScene())
	return a
}

func NewTheme() *ui.Theme {
	t := ui.NewTheme()
	t.Set("bg", ui.Gray)
	t.Set("fg", ui.Black)
	t.Set("fg2", ui.Red)
	t.Set("button bg", ui.GreenYellow)
	t.Set("button fg", ui.Black)
	return &t
}

func LoadSettings() *ui.Preferences {
	s := ui.NewPreferences()
	s.Set("fullscreen", false)
	s.Set("default words per minute speed", 300)
	s.Set("step", 60)
	return &s
}
