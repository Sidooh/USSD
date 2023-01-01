package main

//
//import (
//	"USSD.sidooh/utils"
//)
//
//var openScreenSampleInputs = map[string][]string{
//	utils.AIRTIME_AMOUNT:       {"10", "20", "10000"},
//	utils.AIRTIME_OTHER_NUMBER: {"10", "701100100"},
//}
//
//var openEndScreens = map[string]string{
//	"about":     "string",
//	"refer_pin": "string",
//}
//
//type testPath struct {
//	inputs   []string
//	expected string
//}
//
////func TestProcess(t *testing.T) {
////	err := os.Setenv("DB_DSN", "file:test.db?cache=shared&mode=memory")
////	err = os.Setenv("APP_ENV", "TEST")
////	if err != nil {
////		return
////	}
////
////	initUssd()
////	// TODO: Finalize option 1 : fully automated
////	//var depth = 3
////
////	//fmt.Println(call("a", "", depth))
////
////	// TODO: Flesh out option 2 : somewhat automated
////	paths := map[string][]string{
////		"about":                        {"", "1"},
////		"airtime_self_20_mpesa_accept": {"", "2", "1", "20", "1", "1"},
////		"pay_utility_tokens_existing-acc_100_mpesa_accept": {"", "3", "4", "1", "100", "1", "1"},
////	}
////
////	//paths := map[string][]testPath{}
////	//{
////	//	"about": {
////	//		inputs:   []string{"", "1"},
////	//		expected: "sidooh is a ",
////	//	},
////	//}
////
////	for path, inputs := range paths {
////		for _, input := range inputs {
////			Process("test", "714611696", path, input)
////		}
////	}
////
////}
//
////func call(session, input string, depth int) error {
////	if depth == 0 {
////		fmt.Println("Max Depth reached, exiting...")
////		return nil
////	}
////	fmt.Println("Session ", session, "-", depth)
////
////	stateData := Process(session, input)
////	fmt.Println(stateData.ScreenPath.Key)
////
////	if stateData.ScreenPath.Type == data.GENESIS {
////		v := "2"
////		newSession := session + "*" + v
////		copy(session, newSession)
////		fmt.Println("copied file, call recursive")
////		return call(newSession, v, depth-1)
////
////	} else if stateData.ScreenPath.Type == data.CLOSED {
////		for i := range stateData.ScreenPath.Options {
////
////			v := strconv.Itoa(i)
////
////			newSession := session + "*" + v
////			copy(session, newSession)
////			fmt.Println("copied file, call recursive")
////			return call(newSession, v, depth-1)
////
////		}
////	} else if stateData.ScreenPath.Type == data.OPEN {
////		if _, ok := openEndScreens[stateData.ScreenPath.Key]; ok {
////			return nil
////		}
////
////		inputs := openScreenSampleInputs[stateData.ScreenPath.Key]
////
////		for _, s := range inputs {
////			newSession := session + "*" + s
////			copy(session, newSession)
////			return call(newSession, s, depth-1)
////		}
////
////	} else if stateData.ScreenPath.Type == data.END {
////		//fmt.Println("ENDed flow")
////		return nil
////	}
////
////	return errors.New("failure")
////}
////
////func duplicate(name, new string) {
////	file, err := data.ReadFile(name + data.STATE_FILE)
////	if err != nil {
////		panic(err)
////	}
////
////	err = data.WriteFile(file, new+data.STATE_FILE)
////	if err != nil {
////		panic(err)
////	}
////}
////
////func copy(src, dst string) (int64, error) {
////	src += data.STATE_FILE
////	dst += data.STATE_FILE
////
////	wd, err := os.Getwd()
////	sourceFileStat, err := os.Stat(filepath.Join(wd, data.DATA_DIRECTORY, src))
////	if err != nil {
////		return 0, err
////	}
////
////	if !sourceFileStat.Mode().IsRegular() {
////		return 0, fmt.Errorf("%s is not a regular file", src)
////	}
////
////	source, err := os.Open(filepath.Join(wd, data.DATA_DIRECTORY, src))
////	if err != nil {
////		return 0, err
////	}
////	defer source.Close()
////
////	destination, err := os.Create(filepath.Join(wd, data.DATA_DIRECTORY, dst))
////	if err != nil {
////		return 0, err
////	}
////	defer destination.Close()
////	nBytes, err := io.Copy(destination, source)
////	return nBytes, err
////}
