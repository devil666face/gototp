package view

import (
	"errors"
	"fmt"
	"gototp/internal/gototp"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"

	"github.com/pquerna/otp"
)

var (
	strValidator = func(in string) error {
		if in == "" {
			return fmt.Errorf("⚠️  value is empty")
		}
		return nil
	}
	numValidator = func(in string) error {
		if _, err := strconv.Atoi(in); err != nil {
			return fmt.Errorf("⚠️  value must be digit")
		}
		return nil
	}
)

type View struct {
	gototp *gototp.Gototp
}

func New() *View {
	var passphrase string
	form := InputForm(
		"*️⃣  enter your secret", "",
		strValidator,
		&passphrase,
	)
	if err := form.Run(); err != nil {
		if !errors.Is(err, huh.ErrUserAborted) {
			ErrorFunc(err)
		}
		os.Exit(0)
	}

	_gototp, err := gototp.New(passphrase)
	if err != nil {
		ErrorFunc(err)
		os.Exit(1)
	}

	v := View{
		gototp: _gototp,
	}
	return &v
}

const (
	show   = "#️⃣  show"
	code   = "🔑 code"
	add    = "🆕 add"
	delete = "❌ delete"
)

var mainopts = []huh.Option[string]{
	huh.NewOption[string](show, show),
	huh.NewOption[string](code, code),
	huh.NewOption[string](add, add),
	huh.NewOption[string](delete, delete),
}

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

func (v *View) SelectCode(title string) (int, error) {
	var (
		strid string
		opts  []huh.Option[string]
	)
	if len(v.gototp.Data.Keystore.Keys) == 0 {
		return -1, fmt.Errorf("keys not created")
	}
	for id, f := range v.gototp.Data.Keystore.Keys {
		opts = append(opts, huh.NewOption[string](f.Name, fmt.Sprint(id)))
	}
	form := SelectForm(title, opts, &strid)
	if err := form.Run(); err != nil {
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
		var action string
		form := SelectForm("🔐 gototp ", mainopts, &action)
		if err := form.Run(); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return
			}
			ErrorFunc(err)
			continue
		}
		switch action {
		case show:
			if err := v.show(); err != nil {
				ErrorFunc(err)
				continue
			}
		case code:
			if err := v.code(); err != nil {
				ErrorFunc(err)
				continue
			}
		case add:
			if err := v.add(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				ErrorFunc(err)

			}
		case delete:
			if err := v.delete(); err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					continue
				}
				ErrorFunc(err)
			}
		}
	}
}
