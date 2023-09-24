package products

import (
	"USSD.sidooh/pkg/logger"
)

type MerchantAccount struct {
	Merchant
}

func (f *MerchantAccount) Process(input string) {
	logger.UssdLog.Println(" -- ACCOUNT: process", f.screen.Key, input)
	f.productRep = "account"

	f.Product.Process(input)
	f.processScreen(input)
	f.finalize()
}

func (f *MerchantAccount) processScreen(input string) {
	switch f.screen.Key {

	}
}

func (f *MerchantAccount) finalize() {
	logger.UssdLog.Println(" -- ACCOUNT: finalize", f.screen.Next.Type)

}
