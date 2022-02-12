package products

import (
	"USSD/data"
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
	fmt.Println("PRODUCT: initialize")
	p.screen = screen
	p.vars = vars
}

func (p *Product) Process(input string) {
	fmt.Println("PRODUCT: Process")
}

func (p *Product) finalize() {
	fmt.Println("PRODUCT: finalize")

	//	Final checks
}
