package products

import (
	"USSD.sidooh/data"
	"fmt"
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
	fmt.Println("\t - PRODUCT: initialize")
	p.screen = screen
	p.vars = vars
}

func (p *Product) Process(input string) {
	fmt.Println("\t - PRODUCT: Process")
}

func (p *Product) finalize() {
	fmt.Println("\t - PRODUCT: finalize")

	//	Final checks
}
