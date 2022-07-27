package main

//func X() {
//	initUssd()
//
//	paths := map[string][]string{
//		// 1. ########## ABOUT
//		// ... > About
//		//"about": {"", "1"}, // --- valid
//		//
//		//############## ABOUT END
//
//		// 2. ########## AIRTIME
//		// ... > Airtime > self > amount > mpesa > final
//		//"airtime_self_20_mpesa_accept": {"", "2", "1", "20", "1", "1"}, // --- valid
//		//
//		// ... > Airtime > self > amount > other mpesa > final
//		//"airtime_self_20_other-mpesa_254714611696_accept": {"", "2", "1", "20", "1", "3", "254714611696", "1"}, // --- valid
//		//
//		// ... > Airtime > self > amount > voucher > final
//		//"airtime_self_20_voucher_pin_accept": {"", "2", "1", "31", "2", "1234", "1"}, // --- valid
//		//
//		// ... > Airtime > other > new phone > amount > payment > final
//		//"airtime_other_new-phone_20_mpesa_accept": {"", "2", "2", "780611696", "20", "1", "1"}, // --- valid
//		//
//		// ... > Airtime > other > phone > amount > payment > final
//		//"airtime_other_phone_20_mpesa_accept": {"", "2", "2", "1", "20", "1", "1"}, // --- valid
//		//
//		//	... > Extra paths
//		//"airtime_self_20_mpesa_cancel": {"", "2", "1", "20", "1", "2"}, // --- valid
//		//"airtime_self_20_voucher_invalid-pin_blocked": {"", "2", "1", "20", "2", "123123", "1231", "1232", "7667", "3245"}, // --- valid
//		//"airtime_other_existing-new-phone_20_mpesa_accept": {"", "2", "2", "9", "254780611696", "20", "1", "1"}, // --- valid
//		//"airtime_other_existing_20_mpesa_other_254714611696_accept": {"", "2", "2", "1", "20", "1", "3", "254110039317", "1"}, // --- valid
//		//
//		//"airtime_self_20_voucher_back": {"", "2", "1", "20", "2", "1234", "3", "0"}, // --- valid
//		//
//		// ############## AIRTIME END
//
//		// 3.1 ########## UTILITY
//		// ... > Pay > Utility > provider > select account > amount > payment > final
//		//"pay_utility_kplc_existing-acc_200_mpesa_accept": {"", "3", "1", "2", "1", "200", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Utility > provider > no account > amount > payment > final
//		//"pay_utility_dstv_new-acc_200_mpesa_accept": {"", "3", "1", "4", "1234567", "200", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Utility > provider > existing but new account > amount > payment > final
//		//"pay_utility_kplc_new-acc_200_mpesa_accept": {"", "3", "1", "2", "9", "1234567", "200", "1", "1"},
//		//
//		// ############## UTILITY END
//
//		// 3.2 ########## VOUCHER
//		// ... > Pay > Voucher > self > amount > payment > final
//		//"voucher_self_100_mpesa_accept": {"", "3", "2", "1", "100", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Voucher > other > account > amount > mpesa > final
//		//"voucher_other_phone_100_mpesa_accept": {"", "3", "2", "2", "110039317", "100", "1", "1"}, // --- valid
//		//
//		// ... > Pay > Voucher > other > account > amount > voucher > final
//		//"voucher_other_phone_100_voucher_accept": {"", "3", "2", "2", "110039317", "100", "2", "1234", "1"}, // --- valid
//		//
//		// ############## VOUCHER END
//
//		// 4 ########## SAVE
//		// ... > Save > Voucher > self > amount > payment > final
//		//"voucher_self_100_mpesa_accept": {"", "3", "2", "1", "100", "1", "1"}, // --- valid
//		//
//		// ... > Save > Voucher > other > account > amount > mpesa > final
//		//"voucher_other_phone_100_mpesa_accept": {"", "3", "2", "2", "110039317", "100", "1", "1"}, // --- valid
//		//
//		// ... > Save > Voucher > other > account > amount > voucher > final
//		//"voucher_other_phone_100_voucher_accept": {"", "3", "2", "2", "110039317", "100", "2", "1234", "1"}, // --- valid
//		//
//		// ############## SAVE END
//
//		// 5 ########## INVITE
//		// ... > Invite > Pin > phone > final
//		//"invite_pin_716611696_end": {"", "5", "1234", "716611696"},
//		//
//		// ... > Invite > Pin > phone [existing invite] > final
//		//"invite_pin_718611696_end": {"", "5", "1234", "718611696"}, // --- valid
//		//
//		// ... > Invite > Pin > phone [existing account] > final
//		//"invite_pin_110039317_end": {"", "5", "1234", "110039317"}, // --- valid
//		//
//		// ############## INVITE END
//
//		// 6 ########## SUBSCRIPTION
//		// ... > Subscription > info > name > confirm > payment > final
//		//"subscription_info_Dr-H_confirm_payment_end": {"", "6", "1", "1", "Dr H", "1", "2", "1234" /*, "1"*/},
//		//
//		// ... > Subscription > renew > payment > final
//		//"subscription_renew_payment_end": {"", "6", "1", "1", "1"},
//		//
//		// ############## SUBSCRIPTION END
//
//		// 7 ########## ACCOUNT
//		// ... > Account > Profile > view
//		//"account_profile_end": {"", "7", "1"},
//		//
//		// ... > Account > Profile > pin > set > name > new > confirm
//		//"account_profile_pin_new": {"", "7", "1", "1", "1", "Dr H", "1000", "1000"},
//		//
//		// ... > Account > Profile > pin > change > use pin > new > confirm
//		//"account_profile_pin_change_use_pin": {"", "7", "1", "1", "2", "1", "1234", "1000", "1001"},
//		//
//		// ... > Account > Profile > pin > change > use sec qns > new > confirm
//		//"account_profile_pin_change_use_qns": {"", "7", "1", "1", "2", "2", "Jack", "Summers", "Blue", "1000", "1000"},
//		//
//		// ... > Account > Profile > sec qn > option 1 > choice 1 > ...
//		//"account_profile_security_questions": {"", "7", "1", "1", "3", "1234", "2", "Blue", "1", "Jack", "1", "Dabber"},
//		//
//		// TODO: Update sec qns if I know pin.
//		//
//		// ... > Account > Profile > update > pin > name > end
//		//"account_profile_update_name": {"", "7", "1", "2", "1234", "Jack Dabbs"},
//		//
//		// ############## ACCOUNT END
//	}
//	x := time.Now()
//	for path, inputs := range paths {
//		for _, input := range inputs {
//			//254110039317
//			// TODO: Test with 7, 07, 2547, +2547... determine if mpesa validation will work for different scenarios
//			fmt.Println(processAndRespond("*384*99#", "254714611696", "254714611696"+path, input))
//			//time.Sleep(300 * time.Millisecond)
//			//fmt.Println(processAndRespond("*384*99#", "254110039317", "254110039317"+path, input))
//			//time.Sleep(200 * time.Millisecond)
//
//		}
//	}
//
//	fmt.Println(time.Since(x))
//}
