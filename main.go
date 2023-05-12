package main

import (
	_ "embed"

	"github.com/definenulls/komerco-chain/command/root"
	"github.com/definenulls/komerco-chain/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
