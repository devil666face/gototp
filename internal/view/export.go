package view

import "fmt"

func (v *View) export() error {
	id, err := v.SelectCode("💾 export key ")
	if err != nil {
		return err
	}
	key := v.gototp.Data.Keystore.Keys[id]
	if err := v.gototp.SaveFile(&key); err != nil {
		return err
	}
	fmt.Printf("💾 key %s saved\r\n", key.Name)
	return nil
}
