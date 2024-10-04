package view

import (
	"fmt"
	"gototp/internal/models"

	"github.com/charmbracelet/huh"
)

func (v *View) add() error {
	var (
		input       = models.Input{}
		suggestions = v.gototp.Data.Keystore.Suggestions()
	)

	form := huh.NewForm(
		huh.NewGroup(
			Input(
				"name ", "",
				strValidator,
				&input.Name,
				suggestions...,
			),
			Input(
				"update period ", "30",
				numValidator,
				&input.Period,
				"30",
			),
			Select(
				"set digits ",
				digitopts,
				&input.Digit,
			),
			Select(
				"algorithm ",
				algorithmops,
				&input.Algorithm,
			),
			Input(
				"secret ", "",
				strValidator,
				&input.Secret,
			),
		).WithTheme(base16),
	)
	if err := form.Run(); err != nil {
		return err
	}
	v.gototp.Data.Keystore.Add(input)
	if err := v.gototp.Save(); err != nil {
		return err
	}
	fmt.Printf("✅ %s - created\r\n", input.Name)
	return nil
}
