package service

import (
	"USSD.sidooh/pkg/service/client"
	"strconv"
)

func GetWithdrawalCharge(amount int) int {
	charges, err := paymentsClient.GetWithdrawalCharges()
	if err != nil {
		return 0
	}

	for _, charge := range charges {
		if charge.Min <= amount && amount <= charge.Max {
			return charge.Charge
		}
	}

	return 0
}

func GetPayBillCharge(amount int) int {
	charges, err := paymentsClient.GetPayBillCharges()
	if err != nil {
		return 0
	}

	for _, charge := range charges {
		if charge.Min <= amount && amount <= charge.Max {
			return charge.Charge
		}
	}

	return 0
}

func GetBuyGoodsCharge(amount int) int {
	charges, err := paymentsClient.GetBuyGoodsCharges()
	if err != nil {
		return 0
	}

	for _, charge := range charges {
		if charge.Min <= amount && amount <= charge.Max {
			return charge.Charge
		}
	}

	return 0
}

func GetMpesaFloatCharge(amount int) int {
	charges, err := paymentsClient.GetMpesaFloatCharges()
	if err != nil {
		return 0
	}

	for _, charge := range charges {
		if charge.Min <= amount && amount <= charge.Max {
			return charge.Charge
		}
	}

	return 0
}

func GetMpesaWithdrawalCharge(amount int) int {
	charges, err := paymentsClient.GetMpesaWithdrawalCharges()
	if err != nil {
		return 0
	}

	for _, charge := range charges {
		if charge.Min <= amount && amount <= charge.Max {
			return charge.Charge
		}
	}

	return 0
}

func GetMpesaCollectionCharge(amount int) int {
	charge, err := paymentsClient.GetMpesaCollectionCharge(strconv.Itoa(amount))
	if err != nil {
		return 0
	}

	return charge
}

func SearchMerchant(code string) (*client.Merchant, error) {
	return paymentsClient.SearchMerchant(code)
}
