package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"fmt"
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

	//	TODO: Navigating back from screen causes a bug - sends request but display shows error
	case utils.INVITE:
		i.vars["{number}"], _ = utils.FormatPhone(input)
	}
}

func (i *Invite) finalize() {
	logger.UssdLog.Println(" -- INVITE: finalize", i.screen.Next.Type)

	if i.screen.Key == utils.INVITE {
		accountId, _ := i.vars["{account_id}"]
		number := i.vars["{number}"]

		// TODO: Make into goroutine if applicable
		// TODO: Should we check returned value? Or should we make it a void function?
		_, err := service.CreateInvite(accountId, number, "INVITED")
		if err != nil {
			i.screen.Next.Title = "Sorry. We failed to process your invite, please try again later."
			return
		}

		name := i.vars["{full_name}"] + " " + i.vars["{phone}"]
		code := "*384*99#"

		message := fmt.Sprintf("Hi, %s has invited you to try out Sidooh, "+
			"a cashback platform that rewards you with cash points on every item you buy or pay for via "+
			"the platform. The earned cash points are then automatically saved and invested in "+
			"secure financial assets such as Treasury Bills and Bonds so as to generate passive income for you.\n"+
			"Dial %s NOW on your Safaricom line for FREE to buy airtime and start saving and investing with your cash points.", name, code)
		request := client.NotificationRequest{
			Channel:     "SMS",
			Destination: []string{number},
			EventType:   "REFERRAL_INVITE", //TODO: Change notify referral types to invite
			Content:     message,
		}

		service.Notify(&request)
	}
}
