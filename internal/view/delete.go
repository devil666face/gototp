package view

import (
	"fmt"
)

func (v *View) delete() error {
	id, err := v.SelectCode("❌ delete key ")
	if err != nil {
		return err
	}
	key := v.gototp.Data.Keystore.Keys[id]
	if err := v.gototp.Data.Keystore.Delete(id); err != nil {
		return err
	}
	fmt.Printf("❌ %s - deleted\r\n", key.Name)
	return nil
}
