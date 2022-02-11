package products

import (
	"USSD/data"
	"fmt"
)

type ProductI interface {
	Process(session string, screen *data.Screen)
	ProcessScreen()
}

type Product struct {
	session string
	vars    map[string]string
	screen  *data.Screen
}

func (p *Product) Process(session string, screen *data.Screen) {
	fmt.Println("PRODUCT: Process")

	p.initialize(session, screen)
	p.finalize()
}

func (p *Product) initialize(session string, screen *data.Screen) {
	fmt.Println("PRODUCT: initialize")
	p.screen = screen
	p.session = session

	p.retrieveState()
}

func (p *Product) finalize() {
	fmt.Println("PRODUCT: finalize")

	p.translateScreens()

	p.saveState()

	//	Final checks
}

func (p *Product) retrieveState() {
	err := data.UnmarshalFromFile(p.vars, p.session+"_vars.json")
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Product) translateScreens() {

}

func (p *Product) saveState() {
	err := data.WriteFile(p, p.session+"_vars.json")
	if err != nil {
		panic(err)
	}
}

func (p *Product) ProcessScreen() {
	fmt.Println("PRODUCT: process previous")

}
