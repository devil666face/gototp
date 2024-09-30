package view

import (
	"errors"
	"fmt"
	"gototp/internal/gototp"
	"gototp/internal/models"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/pquerna/otp"
)

type View struct {
	gototp *gototp.Gototp
}

func New(_gototp *gototp.Gototp) *View {
	v := View{
		gototp: _gototp,
	}
	return &v
}

var (
	base16 *huh.Theme = huh.ThemeBase16()
)

var ErrorFunc func(err error) = func(err error) {
	fmt.Printf(err.Error() + "\r\n")
}

func RunSelect(
	title string,
	opts []huh.Option[string],
	value *string,
) error {
	s := huh.NewSelect[string]().
		Title(title).
		Options(opts...).
		Value(value).
		WithTheme(base16)
	return huh.NewForm(huh.NewGroup(s)).Run()
}

const (
	show   = "show"
	code   = "code"
	add    = "add"
	delete = "delete"
)

var (
	algorithmops = huh.NewOptions[string](
		fmt.Sprint(otp.AlgorithmSHA1),
		fmt.Sprint(otp.AlgorithmSHA256),
		fmt.Sprint(otp.AlgorithmSHA512),
		fmt.Sprint(otp.AlgorithmMD5),
	)
	digitopts = huh.NewOptions[string](
		fmt.Sprint(otp.DigitsSix),
		fmt.Sprint(otp.DigitsEight),
	)
)

var (
	strValidator = func(in string) error {
		if in == "" {
			return fmt.Errorf("value is empty")
		}
		return nil
	}
	numValidator = func(in string) error {
		if _, err := strconv.Atoi(in); err != nil {
			return fmt.Errorf("value must be digit")
		}
		return nil
	}
)

func (v *View) SelectCode(title string) (int, error) {
	var (
		strid string
		opts  []huh.Option[string]
	)
	for id, f := range v.gototp.Data.Keystore.Keys {
		opts = append(opts, huh.NewOption[string](f.Name, fmt.Sprint(id)))
	}
	if err := RunSelect(title, opts, &strid); err != nil {
		return -1, err
	}
	id, err := strconv.Atoi(strid)
	if err != nil {
		return -1, fmt.Errorf("failed to get id")
	}
	return id, nil
}

func (v *View) Run() {
	for {
		var (
			action string
			opts   = []huh.Option[string]{
				huh.NewOption[string](show, show),
				huh.NewOption[string](code, code),
				huh.NewOption[string](add, add),
				huh.NewOption[string](delete, delete),
			}
		)
		if err := RunSelect("select action ", opts, &action); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return
			}
			ErrorFunc(err)
			continue
		}
		switch action {
		case show:
			var (
				sb strings.Builder
				w  = tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
			)
			if len(v.gototp.Data.Keystore.Keys) == 0 {
				ErrorFunc(fmt.Errorf("keys not created"))
				continue
			}
			fmt.Fprintf(w, "#\t%s\t%s\t%s\t%s", "name", "period", "digits", "algorithm")
			for i, v := range v.gototp.Data.Keystore.Keys {
				fmt.Fprintf(
					w,
					"\n%d\t%s\t%d\t%s\t%s",
					i+1,
					v.Name,
					v.Period,
					v.Digits,
					v.Algorithm,
				)
			}
			w.Flush()
			fmt.Println(
				lipgloss.NewStyle().
					Padding(0, 1).
					Render(sb.String()),
			)
		case code:
			id, err := v.SelectCode("generate code ")
			if err != nil {
				ErrorFunc(err)
				continue
			}
			key := v.gototp.Data.Keystore.Keys[id]
			otpcode, err := key.GenCode()
			if err != nil {
				ErrorFunc(err)
				continue
			}
			fmt.Println(otpcode)
		case add:
			var (
				input = models.Input{}
			)
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("name ").
						Prompt("> ").
						Validate(strValidator).
						Value(&input.Name).
						Placeholder(""),
					huh.NewInput().
						Title("update period ").
						Prompt("> ").
						Validate(numValidator).
						Value(&input.Period).
						Placeholder("30"),
					huh.NewSelect[string]().
						Title("set digits ").
						Options(digitopts...).
						Value(&input.Digit),
					huh.NewSelect[string]().
						Title("algorithm ").
						Options(algorithmops...).
						Value(&input.Algorithm),
					huh.NewInput().
						Title("secret ").
						Prompt("> ").
						Validate(strValidator).
						Value(&input.Secret).
						Placeholder(""),
				).WithTheme(base16),
			)
			if err := form.Run(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				ErrorFunc(err)
				continue
			}
			v.gototp.Data.Keystore.Add(input)
			if err := v.gototp.Save(); err != nil {
				ErrorFunc(err)
				continue
			}
			fmt.Printf("%s - created\r\n", input.Name)
		case delete:
		}
	}

}
