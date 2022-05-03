package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/utils"
)

type Invite struct {
	Product
}

func (i *Invite) Process(input string) {
	logger.UssdLog.Println(" -- INVITE: process", i.screen.Key, input)
	i.productRep = "invite"

	i.Product.Process(input)
	i.processScreen(input)
	i.finalize()
}

func (i *Invite) processScreen(input string) {
	switch i.screen.Key {
	case utils.INVITE_PIN:
		// TODO: Get from env
		i.vars["{enrollment_time}"] = "48"

	case utils.INVITE:
		// TODO: Get from env
		i.vars["{number}"] = input
	}
}

func (i *Invite) finalize() {
	logger.UssdLog.Println(" -- INVITE: finalize", i.screen.Next.Type)

	if i.screen.Key == utils.INVITE {
		accountId, _ := i.vars["{account_id}"]
		number := i.vars["{number}"]

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		_, err := service.CreateInvite(accountId, number)
		if err != nil {
			i.screen.Next.Title = "Sorry. We failed to process your invite, please try again later."
		}
	}
}
