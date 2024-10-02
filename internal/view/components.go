package view

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

var (
	base16 *huh.Theme = huh.ThemeBase16()
)

var ErrorFunc func(err error) = func(err error) {
	fmt.Printf("❗️" + err.Error() + "\r\n")
}

func Select(
	title string,
	opts []huh.Option[string],
	value *string,
) *huh.Select[string] {
	s := huh.NewSelect[string]().
		Title(title).
		Options(opts...).
		Value(value)
	return s
}

func SelectForm(
	title string,
	opts []huh.Option[string],
	value *string,
) *huh.Form {
	s := Select(
		title,
		opts,
		value,
	)
	return huh.NewForm(huh.NewGroup(s)).WithTheme(base16)
}

func Input(
	title string,
	placeholded string,
	validator func(string) error,
	value *string,
) *huh.Input {
	i := huh.NewInput().
		Title(title).
		Value(value).
		Prompt("> ").
		Validate(validator).
		Placeholder(placeholded)
	return i
}

func InputForm(
	title string,
	placeholded string,
	validator func(string) error,
	value *string,
) *huh.Form {
	i := Input(
		title,
		placeholded,
		validator,
		value,
	)
	return huh.NewForm(huh.NewGroup(i)).WithTheme(base16)
}

func PasswordForm(
	title string,
	placeholded string,
	validator func(string) error,
	value *string,
) *huh.Form {
	i := huh.NewInput().
		Title(title).
		Value(value).
		Prompt("> ").
		Validate(validator).
		Placeholder(placeholded).
		EchoMode(huh.EchoModePassword)
	return huh.NewForm(huh.NewGroup(i)).WithTheme(base16)
}
