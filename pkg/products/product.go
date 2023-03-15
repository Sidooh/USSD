package products

import (
	"USSD.sidooh/pkg/data"
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

type ProductI interface {
	Initialize(vars map[string]string, screen *data.Screen)
	Process(input string)
	finalize()
}

type Product struct {
	vars       map[string]string
	screen     *data.Screen
	productRep string
}

func (p *Product) Initialize(vars map[string]string, screen *data.Screen) {
	logger.UssdLog.Println(" - PRODUCT: initialize")
	p.screen = screen
	p.vars = vars
}

func (p *Product) Process(input string) {
	logger.UssdLog.Println(" - PRODUCT: Process")

	switch p.screen.Key {
	case utils.PAYMENT_METHOD:
		p.setPaymentMethodText(input)
		break
	case utils.PAYMENT_OTHER_NUMBER_MPESA:
		p.vars["{mpesa_number}"], _ = utils.FormatPhone(input)
		p.vars["{payment_method_text}"] = utils.MPESA + " " + p.vars["{mpesa_number}"]
		break
	}
}

func (p *Product) finalize() {
	logger.UssdLog.Println(" - PRODUCT: finalize")
}

func (p *Product) setPaymentMethods(input string) {
	amount, _ := strconv.Atoi(input)
	voucherBalance, _ := strconv.ParseFloat(p.vars["{voucher_balance}"], 32)

	// Move user to top up flow if balance is not enough
	if int(voucherBalance) < amount {
		//TODO: Move user to top up flow
		p.screen.Next.Options[2].NextKey = utils.VOUCHER_BALANCE_INSUFFICIENT
	}

	// Delete voucher option if buying voucher for self
	if p.productRep == "voucher" && p.vars["{number}"] == p.vars["{phone}"] {
		delete(p.screen.Next.Options, 2)
	}

	hasPin := p.checkHasPin()
	if !hasPin {
		if p.productRep == "subscription" && p.screen.Key == utils.PAYMENT_METHOD {
			if _, ok := p.screen.Options[2]; ok {
				p.screen.Options[2].NextKey = utils.PIN_NOT_SET
				return
			}
		}

		if _, ok := p.screen.Next.Options[2]; ok {
			p.screen.Next.Options[2].NextKey = utils.PIN_NOT_SET
		}
		return
	}
}

func (p *Product) setPaymentMethodText(input string) {
	switch input {
	case "1":
		p.vars["{payment_method}"] = utils.MPESA
		p.vars["{payment_method_text}"] = utils.MPESA + " " + p.vars["{phone}"]
		p.vars["{payment_method_instruction}"] = "PLEASE ENTER MPESA PIN when prompted"

	case "2":
		p.vars["{payment_method}"] = utils.VOUCHER
		p.vars["{payment_method_text}"] = utils.VOUCHER + "(KES" + p.vars["{voucher_balance}"] + ")"
		p.vars["{payment_method_instruction}"] = fmt.Sprintf("Your %s will be debited automatically", p.vars["{payment_method_text}"])

		//next := p.screen.Next
		//for {
		//	if next.Key == utils.PAYMENT_CONFIRMATION {
		//		delete(next.Options, 3)
		//		break
		//	}
		//	if next.Next == nil {
		//		break
		//	}
		//	next = next.Next
		//}
	}
}

func (p *Product) checkHasPin() bool {
	accountId := p.vars["{account_id}"]

	// Check if user already has_pin in state else fetch from service
	var hasPin bool
	err := json.Unmarshal([]byte(p.vars["{has_pin}"]), &hasPin)
	if err != nil {
		hasPin = service.CheckHasPin(accountId)
		stringVars, _ := json.Marshal(hasPin)
		p.vars["{has_pin}"] = string(stringVars)
	}

	return hasPin
}
