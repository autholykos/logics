package common

import "github.com/manifoldco/promptui"

func YNPrompt(desc string) bool {
	prompt := promptui.Select{
		Label: desc,
		Items: []string{"Nay", "Yay"},
	}
	_, yayOrNay, e := prompt.Run()
	if e != nil {
		panic(e)
	}

	return yayOrNay == "Yay"
}
