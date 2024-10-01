package view

import "fmt"

func (v *View) importa() error {
	filename, err := v.SelectLocalFile("📥 import .gototp file", ".gototp")
	if err != nil {
		return err
	}
	key, err := v.gototp.LoadFile(filename)
	v.gototp.Data.Keystore.Keys = append(v.gototp.Data.Keystore.Keys, *key)
	if err := v.gototp.Save(); err != nil {
		return err
	}
	fmt.Printf("📥 key %s imported\r\n", key.Name)
	return nil
}
