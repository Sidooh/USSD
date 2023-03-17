package service

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

func GetPaybillCharge(amount int) int {
	charges, err := paymentsClient.GetPaybillCharges()
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
