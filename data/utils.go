package data

// TODO: Finalize phone utils
//func getPhoneProvider(input string) (string, error) {
//	safaricomReg, _ := regexp.Compile("/^(?:254|\+254|0)?((?:7(?:[0129][0-9]|4[0123568]|5[789]|6[89])|(1([1][0-5])))[0-9]{6})$/")
//	airtelReg, _ := regexp.Compile("/^(?:254|\+254|0)?((?:(7(?:(3[0-9])|(5[0-6])|(6[27])|(8[0-9])))|(1([0][0-6])))[0-9]{6})$/")
//	telkomReg, _ := regexp.Compile("/^(?:254|\+254|0)?(7(7[0-9])[0-9]{6})$/")
//	equitelReg, _ := regexp.Compile("/^(?:254|\+254|0)?(7(6[3-6])[0-9]{6})$/")
//	faibaReg, _ := regexp.Compile("/^(?:254|\+254|0)?(747[0-9]{6})$/")
//
//	switch true {
//	case safaricomReg.Match([]byte(input)):
//		return "SAFARICOM", nil
//	case airtelReg.Match([]byte(input)):
//		return "AIRTEL", nil
//	case telkomReg.Match([]byte(input)):
//		return "TELKOM", nil
//	case equitelReg.Match([]byte(input)):
//		return "EQUITEL", nil
//	case faibaReg.Match([]byte(input)):
//		return "FAIBA", nil
//	}
//
//	return "", errors.New("Phone does not seem to be supported");
//}
