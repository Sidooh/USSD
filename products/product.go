package products

import (
	"USSD/data"
	"fmt"
)

type ProductI interface {
	Process(screen *data.Screen)
	ProcessScreen()
}

type Product struct {
	vars   map[string]string
	screen *data.Screen
}

func (p *Product) Process(screen *data.Screen) {
	fmt.Println("PRODUCT: Process")

	p.initialize(screen)
	p.finalize()
}

func (p *Product) initialize(screen *data.Screen) {
	fmt.Println("PRODUCT: initialize")
	p.screen = screen

	p.retrieveState()
}

func (p *Product) finalize() {
	fmt.Println("PRODUCT: finalize")

	p.translateScreens()

	p.saveState()

	//	Final checks
}

func (p *Product) retrieveState() {

}

func (p *Product) translateScreens() {

}

func (p *Product) saveState() {

}

func (p *Product) ProcessScreen() {
	fmt.Println("PRODUCT: process previous")

}
