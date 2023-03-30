package products

import (
	"USSD.sidooh/pkg/logger"
	"USSD.sidooh/pkg/service"
	"USSD.sidooh/pkg/service/client"
	"USSD.sidooh/utils"
	"strconv"
)

type Merchant struct {
	Pay
}

func (m *Merchant) Process(input string) {
	logger.UssdLog.Println(" -- PAY_MERCHANT: process", m.screen.Key, input)
	m.productRep = "pay_merchant"

	m.Pay.Process(input)
	m.processScreen(input)
	m.finalize()
}

func (m *Merchant) processScreen(input string) {
	switch m.screen.Key {
	case utils.MERCHANT:
		if input == "1" {
			m.vars["{merchant_type}"] = utils.MPESA_PAY_BILL
		} else {
			m.vars["{merchant_type}"] = utils.MPESA_BUY_GOODS
		}
		m.vars["{product}"] = "to Merchant"

	case utils.MERCHANT_AMOUNT:
		m.vars["{amount}"] = input

		m.setPaymentMethods(input)

		m.getCharge(input)
		m.setChargeText()

	}
}

func (m *Merchant) setChargeText() {
	charge := ""

	charge = "\nSave: KES" + m.vars["{merchant_fee}"]

	m.vars["{payment_charge_text}"] = charge
}

func (m *Merchant) finalize() {
	logger.UssdLog.Println(" -- PAY_MERCHANT: finalize", m.screen.Next.Type)

	if m.screen.Next.Type == utils.END {
		accountId, _ := strconv.Atoi(m.vars["{account_id}"])
		amount, _ := strconv.Atoi(m.vars["{amount}"])
		method := m.vars["{payment_method}"]
		merchantType := m.vars["{merchant_type}"]
		businessNumber := m.vars["{merchant_number}"]

		if accountId == 0 {
			logger.UssdLog.Println(" -- MERCHANT: creating acc")

			//	TODO: Fix nil value for invite code
			account, err := service.CreateAccount(m.vars["{phone}"], nil)
			if err != nil {
				// TODO: Send message to user
				logger.UssdLog.Error(err)
			}

			accountId = account.Id
		}

		request := client.MerchantPurchaseRequest{
			PurchaseRequest: client.PurchaseRequest{
				Initiator: utils.CONSUMER,
				Amount:    amount,
				Method:    method,
				AccountId: accountId,
			},
			MerchantType:   merchantType,
			BusinessNumber: businessNumber,
		}

		if merchantType == utils.MPESA_PAY_BILL {
			request.AccountNumber = m.vars["{merchant_account}"]
		}

		if _, ok := m.vars["{mpesa_number}"]; ok {
			request.DebitAccount = m.vars["{mpesa_number}"]
		}

		logger.UssdLog.Println(" -- PAY_MERCHANT: payment", request)

		service.PayMerchant(request)
	}
}

func (m *Merchant) getCharge(input string) {
	amount, _ := strconv.Atoi(input)
	fee := 0

	if m.vars["{merchant_type}"] == utils.MPESA_PAY_BILL {
		fee = service.GetPaybillCharge(amount)
	} else {
		fee = service.GetBuyGoodsCharge(amount)
	}

	m.vars["{merchant_fee}"] = strconv.Itoa(fee)
}
