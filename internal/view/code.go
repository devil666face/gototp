package view

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/huh/spinner"
)

func (v *View) code() error {
	ctx, cancel := context.WithCancel(context.Background())
	id, err := v.SelectCode("🔑 generate code ")
	if err != nil {
		cancel()
		return err
	}
	key := v.gototp.Data.Keystore.Keys[id]
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				otpcode, err := key.GenCode()
				if err != nil {
					ErrorFunc(err)
					cancel()
					return
				}
				fmt.Printf("🔑 %s - %s\r\n", key.Name, otpcode)
				time.Sleep(time.Second * time.Duration(key.Period))
			}
		}
	}()
	spinner.New().Title(" updating").Context(ctx).Type(spinner.Monkey).Run()
	cancel()
	return nil
}
