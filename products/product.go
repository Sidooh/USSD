package products

import (
	"USSD.sidooh/data"
	"USSD.sidooh/logger"
)

type ProductI interface {
	Initialize(vars map[string]string, screen *data.Screen)
	Process(input string)
	finalize()
}

type Product struct {
	vars   map[string]string
	screen *data.Screen
}

func (p *Product) Initialize(vars map[string]string, screen *data.Screen) {
	logger.UssdLog.Println(" - PRODUCT: initialize")
	p.screen = screen
	p.vars = vars
}

func (p *Product) Process(input string) {
	logger.UssdLog.Println(" - PRODUCT: Process")
}

func (p *Product) finalize() {
	logger.UssdLog.Println(" - PRODUCT: finalize")
}
