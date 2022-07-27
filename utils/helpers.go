package utils

import "USSD.sidooh/service"

func GetPotentialEarnings(provider string, amount int, subscribed bool) float64 {
	// Get users earning ratio
	ratio := .6

	// Get ripples
	ripples := 6

	if subscribed {
		// Get subscribed users earning ratio
		ratio = 1.0

		// Get ripples (Subscribed users earn 100% pass-thru at the moment)
		ripples = 1
	}

	return calculateEarnings(provider, float64(amount), ratio, float64(ripples))
}

func calculateEarnings(provider string, amount float64, ratio float64, ripples float64) float64 {
	var discountAmount float64

	rate, err := service.GetEarningRate(provider)
	if err == nil && rate.Value != 0 {
		switch rate.Type {
		case "%":
			discountAmount = rate.Value * amount
		case "$":
			discountAmount = rate.Value
		}
	} else {
		discountAmount = float64(amount) * 0
	}

	return discountAmount * ratio / ripples
}

// TODO: create currency formatter function
// 	picks currency from env, parameters, amount, precision, thousand sep... override prefix..etc
