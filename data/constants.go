package data

// types of screens
const (
	GENESIS = "GENESIS"
	OPEN    = "OPEN"
	CLOSED  = "CLOSED"
	END     = "END"
)

const (
	STATE_FILE = "_state.json"
)

var (
	DATA_FILE      = "data.json"
	DATA_DIRECTORY = "data"
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
	PAYMENT_METHOD = "payment_method"

	//	AIRTIME
	AIRTIME                     = "airtime"
	AIRTIME_OTHER_NUMBER_SELECT = "airtime_other_number_select"
	AIRTIME_OTHER_NUMBER        = "airtime_other_number"
	AIRTIME_AMOUNT              = "airtime_amount"
)

const (
	INT              = "INT"
	MIN              = "MIN"
	MAX              = "MAX"
	PHONE            = "PHONE"
	DISALLOW_CURRENT = "DISALLOW_CURRENT"
	MPESA_NUMBER     = "MPESA_NUMBER"
)

const (
	MPESA   = "MPESA"
	VOUCHER = "VOUCHER"
)
