package products

import (
	"USSD.sidooh/logger"
	"USSD.sidooh/service"
	"USSD.sidooh/service/client"
	"USSD.sidooh/utils"
	"strconv"
)

type Merchant struct {
	Pay
}

func (m *Merchant) Process(input string) {
	logger.UssdLog.Println(" -- PAY_MERCHANT: process", m.screen.Key, input)
	m.productRep = "pay_merchant"

	m.Product.Process(input)
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

	case utils.MERCHANT_PAY_BILL, utils.MERCHANT_BUY_GOODS:
		m.vars["{merchant_number}"] = input

	case utils.MERCHANT_PAY_BILL_ACCOUNT:
		m.vars["{merchant_account}"] = input

	case utils.MERCHANT_AMOUNT:
		m.vars["{amount}"] = input
		m.setPaymentMethods(input)
		m.setProduct()

	}
}

func (m *Merchant) setProduct() {
	number := ""

	if m.vars["{merchant_type}"] == utils.MPESA_PAY_BILL {
		number = "Paybill " + m.vars["{merchant_number}"] + ", Account " + m.vars["{merchant_account}"]
	} else {
		number = "Till " + m.vars["{merchant_number}"]
	}

	m.vars["{number}"] = number
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

			account, err := service.CreateAccount(m.vars["{phone}"])
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

		logger.UssdLog.Println(" -- PAY_MERCHANT: payment", request)

		service.PayMerchant(request)
	}
}
