package products

import (
	"USSD.sidooh/data"
	"USSD.sidooh/logger"
	"USSD.sidooh/utils"
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
		p.vars["{mpesa_number}"] = input
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

	// Delete voucher option if balance is not enough or buying voucher for self
	if int(voucherBalance) < amount {
		delete(p.screen.Next.Options, 2)
	} else if p.productRep == "Voucher" && p.vars["{number}"] == p.vars["{phone}"] {
		delete(p.screen.Next.Options, 2)
	}
}

func (p *Product) setPaymentMethodText(input string) {
	switch input {
	case "1":
		p.vars["{payment_method}"] = utils.MPESA
		p.vars["{payment_method_text}"] = utils.MPESA + " " + p.vars["{phone}"]
		p.vars["{payment_method_instruction}"] = "PLEASE ENTER MPESA PIN when prompted"
		break
	case "2":
		p.vars["{payment_method}"] = utils.VOUCHER
		p.vars["{payment_method_text}"] = utils.VOUCHER + "(KES" + p.vars["{voucher_balance}"] + ")"
		p.vars["{payment_method_instruction}"] = fmt.Sprintf("Your %s will be debited automatically", p.vars["{payment_method_text}"])

		delete(p.screen.Next.Next.Options, 3)
		break
	}
}
