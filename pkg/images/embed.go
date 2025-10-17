package images

import (
	_ "embed"
)

var (
	//go:embed cards_sheet.png
	CardsSheet []byte

	//go:embed table_green.png
	TableGreen []byte
)
