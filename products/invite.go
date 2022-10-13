package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
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
		_, err := service.CreateInvite(accountId, number)
		if err != nil {
			i.screen.Next.Title = "Sorry. We failed to process your invite, please try again later."
			return
		}

		name := i.vars["{full_name}"] + " " + i.vars["{phone}"]
		code := "*384*99#"

		message := fmt.Sprintf("Hi, %s has invited you to try out Sidooh, "+
			"a digital platform that gives you loyalty points on every item you purchase and pay for through "+
			"the platform. After which, the earned loyalty points are automatically saved and then invested in "+
			"secure financial assets like Treasury Bills & Bonds so as to generate extra income for you.\n"+
			"Dial %s NOW for FREE on your Safaricom line to buy airtime & start investing using your points.", name, code)
		request := client.NotificationRequest{
			Channel:     "SMS",
			Destination: []string{number},
			EventType:   "REFERRAL_INVITE", //TODO: Change notify referral types to invite
			Content:     message,
		}

		service.Notify(&request)
	}
}
