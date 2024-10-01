package view

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
)

func (v *View) show() error {
	var (
		sb strings.Builder
		w  = tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
	)
	if len(v.gototp.Data.Keystore.Keys) == 0 {
		return fmt.Errorf("keys not created")
	}
	fmt.Fprintf(w, "#\t%s\t%s\t%s\t%s\t%s", "name", "period", "digits", "algorithm", "otp")
	for i, k := range v.gototp.Data.Keystore.Keys {
		code, err := k.GenCode()
		if err != nil {
			code = "❌"
		}
		fmt.Fprintf(
			w,
			"\n%d\t%s\t%d\t%s\t%s\t%s",
			i+1,
			k.Name,
			k.Period,
			k.Digits,
			k.Algorithm,
			code,
		)
	}
	w.Flush()
	fmt.Println(
		lipgloss.NewStyle().
			Padding(0, 1).
			Render(sb.String()),
	)
	return nil
}
