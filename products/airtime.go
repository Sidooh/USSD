package products

import (
	"USSD/data"
	"fmt"
)

type Airtime struct {
	Product
	screen *data.Screen
}

func (a *Airtime) Process(session string, screen *data.Screen) {
	fmt.Println("AIRTIME: process")
	a.initialize(session, screen)

	a.ProcessScreen()

	a.finalize()

}

func (a *Airtime) ProcessScreen() {
	fmt.Println("AIRTIME: process screen")

}
