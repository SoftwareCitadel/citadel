package auth

import "citadel/cmd/citadel/util"

func Logout() error {
	if err := util.RemoveConfigFile(); err != nil {
		return err
	}

	return nil
}
