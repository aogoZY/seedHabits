package services

import "strings"

func ParamsValid(line string)bool{
	deleteBlanckLine := strings.Replace(line, " ","",-1)
	if deleteBlanckLine == ""{
		return false
	}
	return true
}

