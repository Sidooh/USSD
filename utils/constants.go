package utils

var (
	DATA_FILE      = "data.json"
	DATA_DIRECTORY = "pkg/data"
)

// types of screens
const (
	GENESIS = "GENESIS"
	OPEN    = "OPEN"
	CLOSED  = "CLOSED"
	END     = "END"
)

// screen keys
const (
	MAIN_MENU        = "main_menu"
	INVITE_CODE      = "invite_code"
	INACTIVE_ACCOUNT = "inactive_account"

	// static
	ABOUT                        = "about"
	COMING_SOON                  = "coming_soon"
	CANCEL                       = "cancel"
	PIN_NOT_SET                  = "pin_not_set"
	SECURITY_QUESTIONS_NOT_SET   = "security_questions_not_set"
	VOUCHER_BALANCE_INSUFFICIENT = "voucher_balance_insufficient"

	NOT_TRANSACTED = "not_transacted"

	// global
	PAYMENT_METHOD             = "payment_method"
	PAYMENT_OTHER_NUMBER_MPESA = "payment_other_number_mpesa"
	PAYMENT_PIN_CONFIRMATION   = "payment_pin_confirmation"
	PAYMENT_CONFIRMATION       = "payment_confirmation"

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

	// MERCHANT
	MERCHANT            = "merchant"
	MERCHANT_CONSENT    = "merchant_consent"
	MERCHANT_FIRST_NAME = "merchant_first_name"
	MERCHANT_LAST_NAME  = "merchant_last_name"
	MERCHANT_ID_NUMBER  = "merchant_id_number"
	MERCHANT_TERMS      = "merchant_terms"

	MERCHANT_KYB              = "merchant_kyb"
	MERCHANT_LOCATION         = "merchant_location"
	MERCHANT_COUNTY           = "merchant_county"
	MERCHANT_SUB_COUNTY       = "merchant_sub_county"
	MERCHANT_WARD             = "merchant_ward"
	MERCHANT_LANDMARK         = "merchant_landmark"
	MERCHANT_LANDMARK_OTHER   = "merchant_landmark_other"
	MERCHANT_KYB_CONFIRMATION = "merchant_kyb_confirmation"

	MERCHANT_FLOAT_AGENT  = "merchant_float_agent"
	MERCHANT_FLOAT_STORE  = "merchant_float_store"
	MERCHANT_FLOAT_AMOUNT = "merchant_float_amount"

	MERCHANT_PAY_BILL         = "merchant_pay_bill"
	MERCHANT_PAY_BILL_ACCOUNT = "merchant_pay_bill_account"
	MERCHANT_BUY_GOODS        = "merchant_buy_goods"
	MERCHANT_AMOUNT           = "merchant_amount"

	// SAVE
	SAVE = "save"

	// INVITE
	INVITE     = "invite"
	INVITE_PIN = "invite_pin"
	INVITE_END = "invite_end"

	// SUBSCRIPTION
	SUBSCRIPTION_RENEW              = "subscription_renew"
	SUBSCRIPTION_ACTIVE             = "subscription_active"
	SUBSCRIPTION_REGISTER           = "subscription_register"
	SUBSCRIPTION_SUBSCRIBER_NAME    = "subscription_subscriber_name"
	SUBSCRIPTION_SUBSCRIBER_CONFIRM = "subscription_subscriber_confirm"

	// ACCOUNT
	ACCOUNT          = "account"
	ACCOUNT_PROFILE  = "account_profile"
	ACCOUNT_BALANCES = "account_balances"
	ACCOUNT_WITHDRAW = "account_withdraw"

	PROFILE_SECURITY                   = "profile_security"
	PROFILE_NAME                       = "profile_name"
	PROFILE_NEW_PIN                    = "profile_new_pin"
	PROFILE_NEW_PIN_CONFIRM            = "profile_new_pin_confirm"
	PROFILE_CHANGE_PIN_METHODS         = "profile_change_pin_methods"
	PROFILE_CHANGE_PIN_QUESTION        = "profile_change_pin_question"
	PROFILE_SECURITY_QUESTIONS_PIN     = "profile_security_questions_pin"
	PROFILE_SECURITY_QUESTIONS_OPTIONS = "profile_security_questions_options"
	PROFILE_SECURITY_QUESTIONS_ANSWER  = "profile_security_questions_answer"
	PROFILE_UPDATE_NAME                = "profile_update_name"
	PROFILE_UPDATE_END                 = "profile_update_end"
	PROFILE_SECURITY_QUESTIONS_END     = "profile_security_questions_end"

	HAS_SECURITY_QUESTIONS = "has_security_questions"

	ACCOUNT_WITHDRAW_PIN        = "account_withdraw_pin"
	WITHDRAW_DESTINATION        = "withdraw_destination"
	WITHDRAW_MPESA              = "withdraw_mpesa"
	WITHDRAW_OTHER_NUMBER_MPESA = "withdraw_other_number_mpesa"
	WITHDRAW_CONFIRM_PIN        = "withdraw_confirm_pin"
)

// validation keys
const (
	INT              = "INT"
	MIN              = "MIN"
	MAX              = "MAX"
	ALPHANUM         = "ALPHANUM"
	MAX_CHARS        = "MAX_CHARS"
	LENGTH           = "LENGTH"
	PHONE            = "PHONE"
	DISALLOW_CURRENT = "DISALLOW_CURRENT"
	PIN              = "PIN"
	PIN_LENGTH       = "PIN_LENGTH"
	PIN_CONFIRMED    = "PIN_CONFIRMED"
	UTILITY_AMOUNTS  = "UTILITY_AMOUNTS"
	EXISTING_ACCOUNT = "EXISTING_ACCOUNT"

	NOT_INVITED_OR_EXISTING_ACCOUNT = "NOT_INVITED_OR_EXISTING_ACCOUNT"

	NAME = "NAME"

	SECURITY_QUESTION      = "SECURITY_QUESTION"
	WITHDRAW_LIMITS        = "WITHDRAW_LIMITS"
	INVITE_CODE_VALIDATION = "INVITE_CODE_VALIDATION"
)

// initiator types
const (
	CONSUMER = "CONSUMER"
)

// payment methods
const (
	MPESA   = "MPESA"
	VOUCHER = "VOUCHER"
	FLOAT   = "FLOAT"
	BANK    = "BANK"
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

// states
const (
	ACTIVE  = "ACTIVE"
	EXPIRED = "EXPIRED"
)

// merchant types
const (
	MPESA_PAY_BILL  = "MPESA_PAY_BILL"
	MPESA_BUY_GOODS = "MPESA_BUY_GOODS"
)

const (
	DEFAULT_SESSIONS_LIMIT = 300
)
