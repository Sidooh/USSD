package products

import (
	"USSD/data"
	"fmt"
)

type Airtime struct {
	Product
	productRep string
}

func (a *Airtime) Process(input string) {
	fmt.Println("AIRTIME: process")
	a.productRep = "AIRTIME"

	a.processScreen(input)
}

func (a *Airtime) processScreen(input string) {
	fmt.Println("AIRTIME: process screen", input)

	switch a.screen.Key {
	case data.AIRTIME:
		fmt.Println("airtime selected")
		break
	case data.AIRTIME_OTHER_NUMBER_SELECT:
		fmt.Println("airtime other number select selected")
		break
	case data.AIRTIME_OTHER_NUMBER:
		fmt.Println("airtime other number selected")
		break
	case data.AIRTIME_AMOUNT:
		fmt.Println("airtime amount selected")
		break
	}
}
