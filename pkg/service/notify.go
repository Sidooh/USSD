package service

func GetUSSDBalance() (int, error) {
	balance, err := notifyClient.GetUSSDBalance()
	if err != nil {
		return 0, err
	}

	return balance, nil
}
