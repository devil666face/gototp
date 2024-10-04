package view

import (
	"fmt"
	"gototp/internal/models"

	"github.com/charmbracelet/huh"
)

func (v *View) change() error {
	id, err := v.SelectCode("🔄 change key ")
	if err != nil {
		return err
	}
	key := v.gototp.Data.Keystore.Keys[id]
	var input = models.Input{}

	form := huh.NewForm(
		huh.NewGroup(
			Input(
				"name ", key.Name,
				strValidator,
				&input.Name,
				fmt.Sprint(key.Name),
			),
			Input(
				"update period ", fmt.Sprint(key.Period),
				numValidator,
				&input.Period,
				fmt.Sprint(key.Period),
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
				"secret ", fmt.Sprint(key.Secret),
				strValidator,
				&input.Secret,
				fmt.Sprint(key.Secret),
			),
		).WithTheme(base16),
	)
	if err := form.Run(); err != nil {
		return err
	}
	v.gototp.Data.Keystore.Keys[id] = models.NewKey(
		input.Name,
		input.Period,
		input.Digit,
		input.Algorithm,
		input.Secret,
	)
	if err := v.gototp.Save(); err != nil {
		return err
	}
	fmt.Printf("🔄 %s - changed\r\n", input.Name)
	return nil
}
