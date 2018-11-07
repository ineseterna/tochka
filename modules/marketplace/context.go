package marketplace

import (
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

type Context struct {
	*util.Context
	// Localization
	Localization Localization `json:"-"`
	// General
	CanEdit         bool   `json:"can_edit,omitempty"`
	CaptchaId       string `json:"captcha_id,omitempty"`
	Error           string `json:"error,omitempty"`
	CanPostWarnings bool   `json:"-"`
	// Misc
	Pgp                 string                `json:"pgp,omitempty"`
	UserSettingsHistory []UserSettingsHistory `json:"user_settings_history,omitempty"`
	Language            string                `json:"language,omitempty"`
	// Paging & sorting
	SelectedPage  int    `json:"selected_page,omitempty"`
	Pages         []int  `json:"-,omitempty"`
	Page          int    `json:"page,omitempty"`
	NumberOfPages int    `json:"number_of_pages,omitempty"`
	Query         string `json:"query,omitempty"`
	SortBy        string `json:"sort_by,omitempty"`
	// Static Pages
	StaticPage  StaticPage   `json:"-,omitempty"`
	StaticPages []StaticPage `json:"-,omitempty"`
	// Messageboard
	IsPrivateMessage bool `json:"is_private_message,omitempty"`
	// ReadOnlyThread   bool
	ItemCategories []ItemCategory `json:"item_categories,omitempty"`
	ItemCategory   ItemCategory   `json:"-,omitempty"`
	// Menu
	Categories          []Category         `json:"-,omitempty"`
	Cities              map[string]int     `json:"-,omitempty"`
	City                string             `json:"city,omitempty"`
	GeoCities           []City             `json:"geo_cities,omitempty"`
	CityMetroStations   []CityMetroStation `json:"metro_stations,omitempty"`
	CityID              int                `json:"city_id,omitempty"`
	Countries           []Country          `json:"countries,omitempty"`
	Quantity            int                `json:"quantity,omitempty"`
	SelectedPackageType string             `json:"selected_package_type,omitempty"`
	SelectedSection     string             `json:"-,omitempty"`
	SelectedSectionID   int                `json:"-,omitempty"`
	SelectedStatus      string             `json:"selected_status,omitempty"`
	ShippingFrom        string             `json:"shipping_from,omitempty"`
	ShippingFromList    []string           `json:"shipping_from_list,omitempty"`
	ShippingTo          string             `json:"shipping_to,omitempty"`
	ShippingToList      []string           `json:"shipping_to_list,omitempty"`
	Account             string             `json:"account,omitempty"`
	Currency            string             `json:"currency,omitempty"`
	// Categories
	Category       string `json:"category,omitempty"`
	SubCategory    string `json:"sub_category,omitempty"`
	SubSubCategory string `json:"sub_sub_category,omitempty"`
	CategoryID     int    `json:"category_id,omitempty"`
	// Items page
	GroupPackages                        []GroupPackage                       `json:"-,omitempty"`
	GroupPackagesByTypeOriginDestination map[GroupedPackageKey][]GroupPackage `json:"-,omitempty"`
	GroupAvailability                    *GroupPackage                        `json:"group_package,omitempty"`
	NumberOfItems                        int                                  `json:"number_of_items,omitempty"`
	PackageCurrency                      string                               `json:"package_currency,omitempty"`
	PackagePrice                         string                               `json:"package_price,omitempty"`
	// Transactions page
	PendingCount   int `json:"pending_count,omitempty"`
	FailedCount    int `json:"failed_count,omitempty"`
	ReleasedCount  int `json:"released_count,omitempty"`
	CompletedCount int `json:"completed_count,omitempty"`
	AllCount       int `json:"all_count,omitempty"`
	// Models
	ExtendedUsers            []ExtendedUser            `json:"-,omitempty"`
	Invitation               Invitation                `json:"-,omitempty"`
	Invitations              []Invitation              `json:"-,omitempty"`
	Item                     Item                      `json:"-,omitempty"`
	Items                    Items                     `json:"-,omitempty"`
	Package                  Package                   `json:"-,omitempty"`
	Packages                 Packages                  `json:"-,omitempty"`
	Thread                   Thread                    `json:"-,omitempty"`
	Threads                  []Thread                  `json:"-,omitempty"`
	Transaction              Transaction               `json:"-,omitempty"`
	Transactions             []Transaction             `json:"-,omitempty"`
	MessageboardSections     []MessageboardSection     `json:"-,omitempty"`
	ViewMessageboardSections []ViewMessageboardSection `json:"-,omitempty"`
	ViewMessageboardSection  ViewMessageboardSection   `json:"-,omitempty"`
	MessageboardSection      MessageboardSection       `json:"-,omitempty"`
	RatingReview             RatingReview              `json:"-,omitempty"`
	// View Models
	ViewCurrentTransactionStatuses []ViewCurrentTransactionStatus `json:"transaction_statuses,omitempty"`
	ViewExtendedUsers              []ViewExtendedUser             `json:"-"`
	ViewFeedItems                  []ViewFeedItem                 `json:"-"`
	ViewInvitation                 ViewInvitation                 `json:"-"`
	ViewItem                       *ViewItem                      `json:"item,omitempty"`
	ViewItems                      []ViewItem                     `json:"items,omitempty"`
	ViewMessage                    ViewMessage                    `json:"-"`
	ViewMessages                   []ViewMessage                  `json:"-"`
	ViewPackage                    ViewPackage                    `json:"-"`
	ViewPackages                   []ViewPackage                  `json:"-"`
	ViewSeller                     *ViewUser                      `json:"vendor,omitempty"`
	ViewSellers                    []ViewUser                     `json:"-"`
	ViewThread                     *ViewThread                    `json:"thread,omitempty"`
	ViewThreads                    []ViewThread                   `json:"-"`
	ViewPrivateThreads             []ViewPrivateThread            `json:"-"`
	ViewTransaction                *ViewTransaction               `json:"transaction,omitempty"`
	ViewTransactions               []ViewTransaction              `json:"transactions,omitempty"`
	ViewUser                       ViewUser                       `json:"-"`
	ViewUsers                      []ViewUser                     `json:"-"`
	ViewUserWarnings               []ViewUserWarning              `json:"-"`
	// Stats
	NumberOfDailyTransactions     int `json:"-"`
	NumberOfMonthlyTransactions   int `json:"-"`
	NumberOfPrivateMessages       int `json:"-"`
	NumberOfSupportMessages       int `json:"-"`
	NumberOfTransactions          int `json:"-"`
	NumberOfUnreadPrivateMessages int `json:"-"`
	NumberOfUnreadSupportMessages int `json:"-"`
	NumberOfWeeklyTransactions    int `json:"-"`
	NumberOfDisputes              int `json:"-"`
	// Admin Stats
	NumberOfUsers int `json:"-"`
	// --- Vendor Statistics ---
	NumberOfVendors       int `json:"-"`
	NumberOfFreeVendors   int `json:"-"`
	NumberOfGoldVendors   int `json:"-"`
	NumberOfSilverVendors int `json:"-"`
	NumberOfBronzeVendors int `json:"-"`
	// --- User Statistics ---
	NumberOfNewUsers           int         `json:"-"`
	NumberOfActiveUsers        int         `json:"-"`
	NumberOfWeeklyActiveUsers  int         `json:"-"`
	NumberOfOnlineUsers        int         `json:"-"`
	NumberOfMonthlyActiveUsers int         `json:"-"`
	NumberOfInvitedUsers       int         `json:"-"`
	StatsItems                 []StatsItem `json:"-"`
	// User Stats
	StaffStats StaffStats `json:"-"`
	// Auth
	SecretText string `json:"secret_text,omitempty"`
	InviteCode string `json:"invite_code,omitempty"`
	// Bitcoin Wallets
	UserBitcoinBalance       *apis.BTCWalletBalance    `json:"btc_balance"`
	UserBitcoinWallets       UserBitcoinWallets        `json:"-"`
	UserBitcoinWallet        *UserBitcoinWallet        `json:"btc_wallet"`
	UserBitcoinWalletActions []UserBitcoinWalletAction `json:"-"`
	// Ethereum Wallets
	UserEthereumBalance       *apis.ETHWalletBalance     `json:"eth_balance"`
	UserEthereumWallets       UserEthereumWallets        `json:"-"`
	UserEthereumWallet        *UserEthereumWallet        `json:"eth_wallet"`
	UserEthereumWalletActions []UserEthereumWalletAction `json:"-"`
	// Bitcoin Cash Wallets
	UserBitcoinCashBalance       *apis.BCHWalletBalance        `json:"bch_balance"`
	UserBitcoinCashWallets       UserBitcoinCashWallets        `json:"-"`
	UserBitcoinCashWallet        *UserBitcoinCashWallet        `json:"bch_wallet"`
	UserBitcoinCashWalletActions []UserBitcoinCashWalletAction `json:"-"`
	// Referrals
	ReferralPayments []ReferralPayment `json:"-"`
	//Dispute
	Dispute      Dispute      `json:"-"`
	Disputes     []Dispute    `json:"-"`
	DisputeClaim DisputeClaim `json:"-"`
	// Deposit
	Deposit         *Deposit        `json:"-"`
	Deposits        Deposits        `json:"-"`
	DepositsSummary DepositsSummary `json:"-"`
	// Support
	ViewMessageboardThreads []ViewMessageboardThread `json:"-"`
	ViewSupportTicket       ViewSupportTicket        `json:"-"`
	ViewSupportTickets      []ViewSupportTicket      `json:"-"`
	// New Items List page
	ViewAvailableItems []ViewAvailableItem `json:"available_items,omitempty"`
	ViewVendors        []ViewVendor        `json:"-"`
	// Currency Rates
	CurrencyRates map[string]map[string]float64 `json:"-"`
	USDBTCRate    float64                       `json:"-"`
	// Wallet page
	BTCFee           float64                `json:"btc_fee,omitempty"`
	BCHFee           float64                `json:"bch_fee,omitempty"`
	Amount           float64                `json:"amount,omitempty"`
	Address          string                 `json:"address,omitempty"`
	Description      string                 `json:"description,omitempty"`
	BTCPaymentResult *apis.BTCPaymentResult `json:"btc_payment_result,omitempty"`
	BCHPaymentResult *apis.BCHPaymentResult `json:"bch_payment_result,omitempty""`
	ETHPaymentResult *apis.ETHPaymentResult `json:"eth_payment_results,omitempty""`
	// Membership Plans
	PriceBTC float64
	PriceETH float64
	PriceBCH float64
	PriceUSD float64
	// Advertising
	Advertisings    []Advertising `json:"-"`
	AdvertisingCost float64       `json:"-"`
	// ApiSession
	APISession *APISession `json:"api_session,omitempty"`
	// CSRF Token
	CSRFToken string `json:"-"`
	SiteName  string `json:"-"`
	SiteURL   string `json:"-"`
	// Messageboard Stats
	MessageboardSummaryStats MessageboardStats
	MessageboardDailyStats   MessageboardStats
	// Performance metrics
	RenderTime int64
}
