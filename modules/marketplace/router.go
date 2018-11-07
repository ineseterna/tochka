package marketplace

import (
	"github.com/gocraft/web"
)

func ConfigureRouter(router *web.Router) *web.Router {

	// Middleware
	router.Middleware((*Context).LoggerMiddleware)

	// router.Middleware((*Context).BotCheckMiddleware)
	router.Middleware((*Context).SecurityHeadersMiddleware)
	router.Middleware((*Context).AuthMiddleware)
	router.Middleware((*Context).GlobalsMiddleware)
	router.Middleware((*Context).CurrencyMiddleware)

	// Images
	router.Get("/item-image/:item", (*Context).ItemImage)
	router.Get("/user-avatar/:user", (*Context).UserAvatar)

	// Index
	router.Get("/", (*Context).Index)

	// Captcha
	router.Get("/captcha/:captcha_id", (*Context).ViewCaptchaImage)

	// Bot Check
	router.Get("bot-check", (*Context).BotCheckGET)
	router.Post("bot-check", (*Context).BotCheckPOST)

	// ---
	// API
	// ---

	apiRouter := router.Subrouter(Context{}, "/api")
	apiRouter.Get("/auth/login", (*Context).ViewAPILoginRegisterGET)
	apiRouter.Post("/auth/login", (*Context).ViewAPILoginPOST)
	apiRouter.Get("/auth/register", (*Context).ViewAPILoginRegisterGET)
	apiRouter.Post("/auth/register", (*Context).ViewAPIRegisterPOST)
	apiRouter.Get("/serp", (*Context).ViewAPISERP)

	// Store
	apiStoreRouter := apiRouter.Subrouter(Context{}, "/user/:store")
	apiStoreRouter.Middleware((*Context).VendorMiddleware)
	apiStoreRouter.Get("/", (*Context).ViewAPIShowStore)

	// Store item
	apiStoreItemRouter := apiStoreRouter.Subrouter(Context{}, "/item/:item")
	apiStoreItemRouter.Middleware((*Context).VendorItemMiddleware)
	apiStoreItemRouter.Get("/", (*Context).ViewAPIShowItem)

	// Reserve item
	apiStoreItemRouter.Post("/package/:hash/book", (*Context).ViewAPIBookPackage)

	// Transactions
	apiTransactionRouter := apiRouter.Subrouter(Context{}, "/payments")
	apiTransactionRouter.Get("/", (*Context).ViewAPIListCurrentTransactionStatuses)

	singleApiTransactionRouter := apiTransactionRouter.Subrouter(Context{}, "/:transaction")
	singleApiTransactionRouter.Middleware((*Context).TransactionMiddleware)
	singleApiTransactionRouter.Get("/", (*Context).ViewAPIShowTransactionGET)
	singleApiTransactionRouter.Post("/", (*Context).ViewAPIShowTransactionPOST)

	// Wallets
	apiWalletsMiddleware := apiRouter.Subrouter(Context{}, "/wallet")
	apiWalletsMiddleware.Middleware((*Context).BitcoinWalletMiddleware)
	apiWalletsMiddleware.Middleware((*Context).BitcoinCashWalletMiddleware)
	apiWalletsMiddleware.Middleware((*Context).EthereumWalletMiddleware)
	apiWalletsMiddleware.Middleware((*Context).WalletsMiddleware)

	apiWalletsMiddleware.Get("/", (*Context).ViewAPIWallet)

	// --------------
	// HTTP Interface
	// --------------

	loggedInRouter := router.Subrouter(Context{}, "/")
	loggedInRouter.Middleware((*Context).MessageStatsMiddleware)
	loggedInRouter.Middleware((*Context).TransactionStatsMiddleware)
	loggedInRouter.Middleware((*Context).DisputeStatsMiddleware)
	loggedInRouter.Middleware((*Context).WalletsMiddleware)
	loggedInRouter.Middleware((*Context).AdvertisingMiddleware)

	// Static
	staticRouter := loggedInRouter.Subrouter(Context{}, "/help")
	staticRouter.Get("/", (*Context).Help)
	staticRouter.Get("/:filename", (*Context).HelpItem)

	// SERP
	loggedInRouter.Get("/marketplace/", (*Context).ListAvailableItems)
	loggedInRouter.Get("/marketplace/:package_type", (*Context).ListAvailableItems)
	loggedInRouter.Get("/vendors/", (*Context).ListAvailableVendors)
	loggedInRouter.Get("/vendors/:package_type", (*Context).ListAvailableVendors)

	// Auth
	authRouter := loggedInRouter.Subrouter(Context{}, "/auth")
	authRouter.Get("/login", (*Context).LoginGET)
	authRouter.Post("/login", (*Context).LoginPOST)
	authRouter.Get("/recover", (*Context).ViewRecoverGET)
	authRouter.Post("/recover", (*Context).ViewRecoverPOST)
	authRouter.Get("/register", (*Context).RegisterGET)
	authRouter.Post("/register", (*Context).RegisterPOST)
	authRouter.Get("/register/:invite_code", (*Context).RegisterGET)
	authRouter.Post("/register/:invite_code", (*Context).RegisterPOST)
	authRouter.Post("/logout", (*Context).Logout)

	// Main Page
	shoutboxRouter := loggedInRouter.Subrouter(Context{}, "/shoutbox")
	shoutboxRouter.Get("/", (*Context).ViewShoutboxGET)
	shoutboxRouter.Post("/", (*Context).ViewShoutboxPOST)

	// Referral Admin
	referralAdminRouter := loggedInRouter.Subrouter(Context{}, "/referrals/admin")
	referralAdminRouter.Middleware((*Context).AdminMiddleware)
	referralAdminRouter.Get("/", (*Context).ViewAdminListReferralPayments)

	// Auth Admin
	authAdminRouter := authRouter.Subrouter(Context{}, "/admin")
	authAdminRouter.Middleware((*Context).AdminMiddleware)
	authAdminRouter.Get("/users", (*Context).AdminUsers)
	authAdminRouter.Post("/users/:user/login", (*Context).LoginAsUser)
	authAdminRouter.Post("/users/:user/ban", (*Context).BanUser)
	authAdminRouter.Post("/users/:user/staff", (*Context).GrantStaff)
	authAdminRouter.Post("/users/:user/possible_scammer", (*Context).MarkPossibleScammer)
	authAdminRouter.Post("/users/:user/trusted_seller", (*Context).GrantTrustedSeller)
	authAdminRouter.Post("/users/:user/seller", (*Context).GrantSeller)
	authAdminRouter.Post("/users/:user/tester", (*Context).GrantTester)
	authAdminRouter.Get("/reviews", (*Context).AdminReviews)

	// Staff
	staffRouter := loggedInRouter.Subrouter(Context{}, "/staff")
	staffRouter.Middleware((*Context).StaffMiddleware)

	staffRouter.Get("", (*Context).ViewStaff)
	staffRouter.Get("/staff", (*Context).ViewStaffListStaff)
	staffRouter.Get("/users", (*Context).ViewStaffListSupportTickets)
	staffRouter.Get("/advertising", (*Context).ViewStaffAdvertisings)
	staffRouter.Get("/warnings", (*Context).ViewStaffListWarnings)
	staffRouter.Get("/messages", (*Context).ViewStaffListReportedMessages)
	staffRouter.Get("/messages/:message_uuid", (*Context).ViewStaffShowReportedMessage)
	staffRouter.Get("/items", (*Context).ViewStaffListItems)
	staffRouter.Get("/disputes", (*Context).ViewStaffListDisputes)
	staffRouter.Get("/stats", (*Context).ViewStaffStats)

	// Staff - Stats
	staffRouter.Get("/stats/users.png", (*Context).ViewStatsNumberOfUsersGraph)
	staffRouter.Get("/stats/vendors.png", (*Context).ViewStatsNumberOfVendorsGraph)
	staffRouter.Get("/stats/trade-amount-btc.png", (*Context).ViewStatsBTCTradeAmountGraph)
	staffRouter.Get("/stats/trade-amount-eth.png", (*Context).ViewStatsETHTradeAmountGraph)
	staffRouter.Get("/stats/trade-amount-bch.png", (*Context).ViewStatsBCHTradeAmountGraph)
	staffRouter.Get("/stats/transactions.png", (*Context).ViewStatsNumberOfTransactionsGraph)

	// Staff - Vendors
	staffRouter.Get("/vendors", (*Context).ViewStaffListVendors)
	staffRouter.Get("/vendors/:username", (*Context).ViewStaffVendorVerificationThreadGET)
	staffRouter.Post("/vendors/:username", (*Context).ViewStaffVendorVerificationThreadPOST)

	// Staff - User
	staffRouter.Get("/users/:username/payments", (*Context).ViewStaffUserPayments)
	staffRouter.Get("/users/:username/finance", (*Context).ViewStaffUserFinance)
	staffRouter.Get("/users/:username/tickets", (*Context).ViewStaffUserTickets)
	staffRouter.Get("/users/:username/tickets/:id", (*Context).ShowSupportTicketGET)
	staffRouter.Post("/users/:username/tickets/:id", (*Context).ShowSupportTicketPOST)
	staffRouter.Get("/users/:username/actions", (*Context).ViewStaffUserAdminActions)

	// Staff - User - Admin Actions
	staffRouter.Post("/users/:user/ban", (*Context).BanUser)
	staffRouter.Post("/users/:user/staff", (*Context).GrantStaff)
	staffRouter.Post("/users/:user/possible_scammer", (*Context).MarkPossibleScammer)
	staffRouter.Post("/users/:user/seller", (*Context).GrantSeller)
	staffRouter.Post("/users/:user/trusted", (*Context).GrantTrustedSeller)
	// -- account types
	staffRouter.Post("/users/:user/gold", (*Context).GrantGoldAccount)
	staffRouter.Post("/users/:user/silver", (*Context).GrantSilverAccount)
	staffRouter.Post("/users/:user/bronze", (*Context).GrantBronzeAccount)
	staffRouter.Post("/users/:user/free", (*Context).GrantFreeAccount)

	// Staff - CRUD
	staffRouter.Get("/news", (*Context).ViewStaffEditNewsGET)
	staffRouter.Post("/news", (*Context).ViewStaffEditNewsPOST)
	staffRouter.Get("/item_categories", (*Context).ViewStaffCategories)
	staffRouter.Get("/item_categories/:id", (*Context).ViewStaffCategoriesEdit)
	staffRouter.Post("/item_categories/:id", (*Context).ViewStaffCategoriesEditPOST)
	staffRouter.Post("/item_categories/:id/delete", (*Context).ViewStaffCategoriesDelete)

	// Warnings
	loggedInRouter.Get("/user/active_reservation", (*Context).ActiveReservation)
	// loggedInRouter.Get("/user/self_purchase", (*Context).SelfPurchase)

	// Store
	storeRouter := loggedInRouter.Subrouter(Context{}, "/user/:store")
	storeRouter.Middleware((*Context).VendorMiddleware)
	storeRouter.Get("/items", (*Context).ViewShowStore)
	storeRouter.Get("/", (*Context).ViewAboutStore)
	storeRouter.Get("/reviews", (*Context).ViewStoreReviews)
	storeRouter.Get("/warnings", (*Context).ViewStoreWarningsGET)
	storeRouter.Post("/warnings", (*Context).ViewStoreWarningsPOST)
	storeRouter.Post("/warnings/:uuid", (*Context).ViewStoreWarningUpdateStatusPOST)

	// Store item
	storeItemRouter := storeRouter.Subrouter(Context{}, "/item/:item")
	storeItemRouter.Middleware((*Context).VendorItemMiddleware)
	storeItemRouter.Get("/", (*Context).ViewShowItem)
	storeItemRouter.Get("/package/:hash/book", (*Context).PreBookPackage)
	storeItemRouter.Post("/package/:hash/book", (*Context).BookPackage)

	// Board
	boardRouter := loggedInRouter.Subrouter(Context{}, "/board")
	boardRouter.Get("/", (*Context).ViewListMessageboardSections)
	boardRouter.Get("/:section_id", (*Context).ViewListThreads)
	boardRouter.Get("/message/:uuid/image", (*Context).MessageImage)
	boardRouter.Get("/thread/new", (*Context).EditThread)
	boardRouter.Post("/thread/new", (*Context).EditThreadPOST)
	boardRouter.Get("/thread/:uuid", (*Context).ShowThread)
	boardRouter.Get("/thread/:uuid/edit", (*Context).EditThread)
	boardRouter.Post("/thread/:uuid/edit", (*Context).EditThreadPOST)
	boardRouter.Post("/:uuid/delete", (*Context).DeleteThread)
	boardRouter.Post("/:uuid/report", (*Context).ViewMessageReportPOST)
	boardRouter.Post("/thread/:uuid", (*Context).ReplyToThread)

	// Messages
	messagesRouter := loggedInRouter.Subrouter(Context{}, "/messages")
	messagesRouter.Middleware((*Context).MessagesMiddleware)
	messagesRouter.Get("/", (*Context).ViewListPrivateMessages)
	messagesRouter.Get("/:username", (*Context).ViewShowPrivateMessageGET)
	messagesRouter.Post("/:username", (*Context).ViewShowPrivateMessagePOST)

	messagesAdminRouter := messagesRouter.Subrouter(Context{}, "/admin")
	messagesAdminRouter.Middleware((*Context).AdminMiddleware)
	messagesAdminRouter.Get("/:uuid", (*Context).AdminMessagesShow)

	// Package
	packagesRouter := loggedInRouter.Subrouter(Context{}, "/packages")
	packagesRouter.Get("/", (*Context).ListPackages)
	singlePackageRouter := packagesRouter.Subrouter(Context{}, "/:package")
	singlePackageRouter.Middleware((*Context).PackageMiddleware)
	singlePackageRouter.Get("/image", (*Context).PackageImage)

	// Profile
	loggedInRouter.Get("/profile", (*Context).ProfileGET)
	loggedInRouter.Post("/profile", (*Context).ProfilePOST)
	loggedInRouter.Get("/referrals", (*Context).Referrals)
	loggedInRouter.Get("/verification/encryption", (*Context).ViewVerificationEncryptionGET)
	loggedInRouter.Get("/verification/agreement", (*Context).ViewVerificationAgreementGET)
	loggedInRouter.Post("/verification/agreement", (*Context).ViewVerificationAgreementPOST)
	loggedInRouter.Get("/verification/plan", (*Context).ViewVerificationPlanGET)
	loggedInRouter.Get("/verification/plan/:account", (*Context).ViewVerificationPlanPurchaseGET)
	loggedInRouter.Post("/verification/plan/:account", (*Context).ViewVerificationPlanPurchasePOST)
	loggedInRouter.Get("/free_restrictions", (*Context).FreeRestrictions)
	loggedInRouter.Get("/verification/thread", (*Context).ViewVerificationThreadGET)
	loggedInRouter.Post("/verification/thread", (*Context).ViewVerificationThreadPOST)
	loggedInRouter.Get("/settings/pgp/step1", (*Context).SetupPGPViaDecryptionStep1GET)
	loggedInRouter.Post("/settings/pgp/step1", (*Context).SetupPGPViaDecryptionStep1POST)
	loggedInRouter.Post("/settings/pgp/step2", (*Context).SetupPGPViaDecryptionStep2POST)

	loggedInRouter.Get("/settings/currency/:currency", (*Context).SetCurrency)
	loggedInRouter.Get("/settings/language/:language", (*Context).SetLanguage)

	// Profile
	loggedInRouter.Post("/shipping", (*Context).SaveShippingOption)
	loggedInRouter.Post("/shipping/delete", (*Context).DeleteShippingOption)

	// Support
	supportRouter := loggedInRouter.Subrouter(Context{}, "/support")
	supportRouter.Get("/", (*Context).ViewListSupportTickets)
	supportRouter.Get("/new", (*Context).ShowSupportTicketGET)
	supportRouter.Post("/new", (*Context).ShowSupportTicketPOST)
	supportRouter.Get("/:id", (*Context).ShowSupportTicketGET)
	supportRouter.Post("/:id", (*Context).ShowSupportTicketPOST)
	supportRouter.Post("/:id/status", (*Context).ViewUpdateTicketStatusPOST)

	// Feed
	feedRouter := loggedInRouter.Subrouter(Context{}, "/feed")
	feedRouter.Get("/", (*Context).ShowFeed)

	// Store CMS
	sellerRouter := loggedInRouter.Subrouter(Context{}, "/seller/:store")
	sellerRouter.Middleware((*Context).SellerMiddleware)

	// Advertisings
	sellerRouter.Get("/advertisings", (*Context).EditAdvertisings)
	sellerRouter.Post("/advertisings", (*Context).AddAdvertisingsPOST)

	// Items CMS
	itemRouter := sellerRouter.Subrouter(Context{}, "/item/:item")
	itemRouter.Middleware((*Context).SellerItemMiddleware)
	itemRouter.Get("/edit", (*Context).EditItem)
	itemRouter.Post("/edit", (*Context).SaveItem)
	itemRouter.Post("/delete", (*Context).DeleteItem)

	// Package CMS
	packageRouter := itemRouter.Subrouter(Context{}, "/package/:package")
	packageRouter.Middleware((*Context).SellerItemPackageMiddleware)
	packageRouter.Get("/edit", (*Context).EditPackageStep1)
	packageRouter.Post("/edit", (*Context).SavePackage)
	packageRouter.Post("/duplicate", (*Context).DuplicatePackage)
	packageRouter.Post("/delete", (*Context).DeletePackage)

	// Transactions
	transactionRouter := loggedInRouter.Subrouter(Context{}, "/payments")
	transactionRouter.Get("/", (*Context).ListCurrentTransactionStatuses)

	singleTransactionRouter := transactionRouter.Subrouter(Context{}, "/:transaction")
	singleTransactionRouter.Middleware((*Context).TransactionMiddleware)
	singleTransactionRouter.Get("/", (*Context).ShowTransaction)
	singleTransactionRouter.Get("/image", (*Context).TransactionImage)
	singleTransactionRouter.Post("/", (*Context).ShowTransactionPOST)
	singleTransactionRouter.Post("/shipping", (*Context).SetTransactionShippingStatus)
	singleTransactionRouter.Post("/release", (*Context).ReleaseTransaction)
	singleTransactionRouter.Post("/cancel", (*Context).CancelTransaction)
	singleTransactionRouter.Post("/complete", (*Context).CompleteTransactionPOST)

	//Disputes
	disputeRouter := loggedInRouter.Subrouter(Context{}, "/dispute")
	disputeRouter.Get("/", (*Context).ViewDisputeList)
	disputeRouter.Post("/start/:transaction_uuid", (*Context).ViewDisputeStart)

	concreteDisputeRouter := disputeRouter.Subrouter(Context{}, "/:dispute_uuid")
	concreteDisputeRouter.Middleware((*Context).DisputeMiddleware)
	concreteDisputeRouter.Get("/", (*Context).ViewDispute)
	concreteDisputeRouter.Post("/status", (*Context).ViewDisputeSetStatus)
	concreteDisputeRouter.Post("/claim", (*Context).ViewDisputeClaimAdd)

	disputeClaimRouter := concreteDisputeRouter.Subrouter(Context{}, "/:dispute_claim_id")
	disputeClaimRouter.Middleware((*Context).DisputeClaimMiddleware)
	disputeClaimRouter.Get("/", (*Context).ViewDisputeClaimGET)
	disputeClaimRouter.Post("/", (*Context).ViewDisputeClaimPOST)

	disputeAdminROuter := disputeRouter.Subrouter(Context{}, "/admin")
	disputeAdminROuter.Middleware((*Context).AdminMiddleware)
	disputeAdminROuter.Get("/list", (*Context).AdminDisputeList)

	// Wallet
	walletRouter := loggedInRouter.Subrouter(Context{}, "/wallet")

	// Bitcoin Wallet
	bitcoinWalletRouter := walletRouter.Subrouter(Context{}, "/bitcoin")
	bitcoinWalletRouter.Middleware((*Context).BitcoinWalletMiddleware)
	bitcoinWalletRouter.Get("/receive", (*Context).BitcoinWalletRecieve)
	bitcoinWalletRouter.Get("/send", (*Context).BitcoinWalletSendStep1GET)
	bitcoinWalletRouter.Post("/send", (*Context).BitcoinWalletSendPOST)
	bitcoinWalletRouter.Get("/:address/image", (*Context).BitcoinWalletImage)
	bitcoinWalletRouter.Get("/actions", (*Context).BitcoinWalletActions)

	// Bitcoin Cash Wallet
	bitcoinCashWalletRouter := walletRouter.Subrouter(Context{}, "/bitcoin_cash")
	bitcoinCashWalletRouter.Middleware((*Context).BitcoinCashWalletMiddleware)
	bitcoinCashWalletRouter.Get("/receive", (*Context).BitcoinCashWalletRecieve)
	bitcoinCashWalletRouter.Get("/send", (*Context).BitcoinCashWalletSendStep1GET)
	bitcoinCashWalletRouter.Post("/send", (*Context).BitcoinCashWalletSendPOST)
	bitcoinCashWalletRouter.Post("/send_to_btc", (*Context).BitcoinCashWalletSendToBitcoinPOST)
	bitcoinCashWalletRouter.Get("/:address/image", (*Context).BitcoinCashWalletImage)
	bitcoinCashWalletRouter.Get("/actions", (*Context).BitcoinCashWalletActions)

	// Ethereum Wallet
	ethereumWalletRouter := walletRouter.Subrouter(Context{}, "/ethereum")
	ethereumWalletRouter.Middleware((*Context).EthereumWalletMiddleware)
	ethereumWalletRouter.Get("/receive", (*Context).EthereumWalletRecieve)
	ethereumWalletRouter.Get("/send", (*Context).EthereumWalletSendGET)
	ethereumWalletRouter.Post("/send", (*Context).EthereumWalletSendPOST)
	ethereumWalletRouter.Get("/:address/image", (*Context).EthereumWalletImage)
	ethereumWalletRouter.Get("/actions", (*Context).EthereumWalletActions)

	// Deposit
	depositRouter := loggedInRouter.Subrouter(Context{}, "/deposit")
	depositRouter.Get("/", (*Context).ViewShowDeposit)
	depositRouter.Get("/add", (*Context).ViewShowDepositAddGET)
	depositRouter.Post("/add", (*Context).ViewShowDepositAddPOST)

	// Transactions Admin
	transactionAdminRouter := transactionRouter.Subrouter(Context{}, "/admin")
	transactionAdminRouter.Middleware((*Context).StaffMiddleware)
	transactionAdminRouter.Get("/list", (*Context).AdminListTransactions)
	transactionAdminRouter.Post("/:transaction/cancel", (*Context).AdminTransactionCancel)
	transactionAdminRouter.Post("/:transaction/fail", (*Context).AdminTransactionFail)
	transactionAdminRouter.Post("/:transaction/pending", (*Context).AdminTransactionMakePending)
	transactionAdminRouter.Post("/:transaction/complete", (*Context).AdminTransactionComplete)
	transactionAdminRouter.Post("/:transaction/release", (*Context).AdminTransactionRelease)
	transactionAdminRouter.Post("/:transaction/freeze", (*Context).AdminTransactionFreeze)
	transactionAdminRouter.Post("/:transaction/update", (*Context).AdminTransactionUpdateStatus)

	// Messageboard sections Admin
	messageboardSectionsAdmin := loggedInRouter.Subrouter(Context{}, "/messageboard_sections/admin")
	messageboardSectionsAdmin.Middleware((*Context).AdminMiddleware)

	messageboardSectionsAdmin.Get("/", (*Context).AdminMessageboardSections)
	messageboardSectionsAdmin.Get("/:id", (*Context).AdminMessageboardSectionsEdit)
	messageboardSectionsAdmin.Post("/:id", (*Context).AdminMessageboardSectionsEditPOST)
	messageboardSectionsAdmin.Post("/:id/delete", (*Context).AdminMessageboardSectionsDelete)

	return router
}
