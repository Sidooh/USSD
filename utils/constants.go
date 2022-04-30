package utils

var (
	DATA_FILE      = "data.json"
	DATA_DIRECTORY = "data"
)

// types of screens
const (
	GENESIS = "GENESIS"
	OPEN    = "OPEN"
	CLOSED  = "CLOSED"
	END     = "END"
)

// file that stores session state
const (
	STATE_FILE = "_state.json"
)

// screen keys
const (
	MAIN_MENU = "main_menu"

	// static
	ABOUT       = "about"
	COMING_SOON = "coming_soon"
	BETA        = "beta"
	CANCEL      = "cancel"

	// global
	PAYMENT_METHOD             = "payment_method"
	PAYMENT_CONFIRMATION       = "payment_confirmation"
	PAYMENT_END                = "payment_end"
	PAYMENT_OTHER_NUMBER_MPESA = "payment_other_number_mpesa"

	PIN_BLOCKED = "pin_blocked"

	// AIRTIME
	AIRTIME                     = "airtime"
	AIRTIME_OTHER_NUMBER_SELECT = "airtime_other_number_select"
	AIRTIME_OTHER_NUMBER        = "airtime_other_number"
	AIRTIME_AMOUNT              = "airtime_amount"

	// PAY
	PAY = "pay"

	// UTILITY
	UTILITY                = "utility"
	UTILITY_ACCOUNT_SELECT = "utility_account_select"
	UTILITY_OTHER_ACCOUNT  = "utility_other_account"
	UTILITY_AMOUNT         = "utility_amount"

	//TODO: Improve constants naming or find way to do enums to avoid conflicts
	// VOUCHER
	PAY_VOUCHER           = "voucher"
	VOUCHER_AMOUNT        = "voucher_amount"
	VOUCHER_OTHER_ACCOUNT = "voucher_other_account"
)

// validation keys
const (
	INT              = "INT"
	MIN              = "MIN"
	MAX              = "MAX"
	PHONE            = "PHONE"
	DISALLOW_CURRENT = "DISALLOW_CURRENT"
	PIN              = "PIN"
	PIN_LENGTH       = "PIN_LENGTH"
	UTILITY_AMOUNTS  = "UTILITY_AMOUNTS"
	EXISTING_ACCOUNT = "EXISTING_ACCOUNT"
)

// payment methods
const (
	MPESA   = "MPESA"
	VOUCHER = "VOUCHER"
)

// telco providers
const (
	SAFARICOM = "SAFARICOM"
	AIRTEL    = "AIRTEL"
	TELKOM    = "TELKOM"
	EQUITEL   = "EQUITEL"
	FAIBA     = "FAIBA"
)

// To confirm, use https://regoio.herokuapp.com/ & https://regexr.com/
const (
	SAFARICOM_REGEX = `^(\+?254|0)?((7([0129][0-9]|4[0123568]|5[789]|6[89])|(1[1][0-5]))[0-9]{6})$`
	AIRTEL_REGEX    = `^(\+?254|0)?(((7(3[0-9]|5[0-6]|6[27]|8[0-9]))|(1[0][0-6]))[0-9]{6})$`
	TELKOM_REGEX    = `^(\+?254|0)?(77[0-9]{7})$`
	EQUITEL_REGEX   = `^(\+?254|0)?(76[3-6][0-9]{6})$`
	FAIBA_REGEX     = `^(\+?254|0)?(747[0-9]{6})$`
)

// utility providers
const (
	KPLC_PREPAID  = "KPLC_PREPAID"
	KPLC_POSTPAID = "KPLC_POSTPAID"
	NAIROBI_WTR   = "NAIROBI_WTR"
	DSTV          = "DSTV"
	ZUKU          = "ZUKU"
	GOTV          = "GOTV"
	STARTIMES     = "STARTIMES"
)
