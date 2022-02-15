package products

import (
	"USSD/data"
	"fmt"
	"strings"
)

type Airtime struct {
	Product
	productRep string
}

func (a *Airtime) Process(input string) {
	fmt.Println("\t -- AIRTIME: process")
	a.productRep = "airtime"

	a.processScreen(input)
}

func (a *Airtime) processScreen(input string) {
	fmt.Println("\t -- AIRTIME: process screen", input)
	fmt.Println("\t --> selected: ", a.screen.Key)

	switch a.screen.Key {
	case data.AIRTIME:
		a.vars["{product}"] = a.productRep
		a.vars["{number}"] = a.vars["{phone}"]
		break
	case data.AIRTIME_OTHER_NUMBER_SELECT:
		break
	case data.AIRTIME_OTHER_NUMBER:
		break
	case data.AIRTIME_AMOUNT:
		a.vars["{amount}"] = input
		a.screen.Next.Title = "Buy" + strings.TrimPrefix(a.screen.Next.Title, "Pay")
		break
	}
}
