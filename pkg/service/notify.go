package service

func GetUSSDBalance() (float64, error) {
	balance, err := notifyClient.GetUSSDBalance()
	if err != nil {
		return 0, err
	}

	return balance, nil
}
