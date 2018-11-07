package marketplace

import (
	"errors"
	"html/template"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jinzhu/gorm"
	"github.com/o1egl/govatar"
	"github.com/russross/blackfriday"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/apis"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

var (
	onlineDuration, _ = time.ParseDuration(MARKETPLACE_SETTINGS.OnlineDuration)
)

/*
	Models
*/

type User struct {
	// Base shit
	Uuid           string `json:"-" gorm:"primary_key"`
	Username       string `json:"username" gorm:"type:varchar(16);unique"`
	PassphraseHash string `json:"-"`

	// Login dates
	RegistrationDate time.Time  `json:"registration_date" gorm:"index"`
	LastLoginDate    *time.Time `json:"last_login_date" gorm:"index"`

	// Settings
	Language string `json:"-" gorm:"index"`
	Currency string `json:"-" gorm:"index"`

	// Profile
	Bitcoin                  string `json:"-"`
	BitcoinMultisigPublicKey string `json:"-"`
	Ethereum                 string `json:"-"`
	BitcoinCash              string `json:"-"`

	// Contacts
	Bitmessage      string `json:"bitmessage"`
	Tox             string `json:"tox"`
	Email           string `json:"email"`
	Pgp             string `json:"pgp" sql:"type:varchar(8192);"`
	Description     string `json:"description" sql:"size:140"`
	LongDescription string `json:"long_description" sql:"size:2048"`
	InviteCode      string `json:"-" gorm:"unique"`

	// Misc settings
	TwoFactorAuthentication bool `json:"-" gorm:"index"`
	HasTopBanner            bool `json:"-" gorm:"index"`
	Banned                  bool `json:"-" gorm:"index"`
	PossibleScammer         bool `json:"is_possible_scammer" gorm:"index"`
	VacationMode            bool `json:"vacation_mode" gorm:"index"`
	HasAvatar               bool `json:"-" gorm:"index"`

	// Relations
	Items                  []Item                  `json:"-"`
	RatingReviews          []RatingReview          `json:"-"`
	ShippingOptions        []ShippingOption        `json:"-"`
	UserBitcoinWallets     []UserBitcoinWallet     `json:"-"`
	UserEthereumWallets    []UserEthereumWallet    `json:"-"`
	UserBitcoinCashWallets []UserBitcoinCashWallet `json:"-"`
	UserWarnings           []UserWarning           `json:"-"`
	Deposits               []Deposit               `json:"-"`

	// Roles
	IsSeller        bool `json:"is_seller" gorm:"index"`
	IsTrustedSeller bool `json:"is_trustedseller" gorm:"index"`
	IsTester        bool `json:"-" gorm:"index"`
	IsModerator     bool `json:"is_moderator" gorm:"index"`
	IsAdmin         bool `json:"is_admin" gorm:"index"`
	IsStaff         bool `json:"is_staff" gorm:"index"`

	// User Verification v2
	IsSigned                 bool      `json:"is_signed" gorm:"index"`
	SignedAgreementText      string    `json:"signed_agreement_text" sql:"type:varchar(8192);"`
	AgreementDate            time.Time `json:"agreement_date"`
	HasRequestedVerification bool      `json:"-" gorm:"index"`

	// Integrations
	MattermostUsername string `json:"-"`

	// Index
	InviterUuid string `json:"-" gorm:"index"`

	// Support User UUID
	SupporterUuid string `json:"-" gorm:"index"`
	TrusteeUuid   string `json:"-" gorm:"index"`

	// Membership Plans
	IsFreeAccount   bool `json:"is_free_account" gorm:"index"`
	IsGoldAccount   bool `json:"is_gold_account" gorm:"index"`
	IsSilverAccount bool `json:"is_silver_account" gorm:"index"`
	IsBronzeAccount bool `json:"is_bronze_account" gorm:"index"`

	// ORM timestamps
	CreatedAt *time.Time `json:"-" gorm:"index"`
	UpdatedAt *time.Time `json:"-" gorm:"index"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

// ExtendedUser is an optimization to retrieve User and associated information from db without extra requests
type ExtendedUser struct {
	User

	InviterUsername     string
	SupportUserUsername string
	InviterCount        string

	BitcoinBalance                float64
	BitcoinUnconfirmedBalance     float64
	BitcoinCashBalance            float64
	BitcoinCashUnconfirmedBalance float64
	EthereumBalance               float64

	RedWarningCount    int
	GreenWarningCount  int
	YellowWarningCount int
}

type Users []User

type ExtendedUsers []ExtendedUser

type UserSettingsHistory struct {
	gorm.Model

	UserUuid string    `json:"user_uuid"`
	Datetime time.Time `json:"datetime"`
	Action   string    `json:"action" sql:"size:1024"`
	Type     string    `json:"string"`
	User     User
}

/*
	Model Methods
*/

func (u User) CheckPassphrase(passphrase string) bool {
	return u.PassphraseHash == util.PasswordHashV1(u.Username, passphrase)
}

func (u User) String() string {
	return u.Username
}

func (u User) NumberOfInvitedUsers() int {
	var count int
	database.
		Model(&User{}).
		Where(&User{InviterUuid: u.Uuid}).
		Count(&count)
	return count
}

// func (u User) Balance() float64

func (u User) Validate() error {
	if !usernameRegexp.MatchString(u.Username) {
		return errors.New("Username is not valid")
	}
	if u.Bitmessage != "" && !bitmessageRegexp.MatchString(u.Bitmessage) {
		return errors.New("Invalid Bitmessage")
	}
	if u.Pgp != "" {
		fromBlock, err := armor.Decode(strings.NewReader(u.Pgp))
		if err != nil || fromBlock.Type != openpgp.PublicKeyType {
			return errors.New("Invalid PGP Key")
		}
	}
	if u.Bitcoin != "" && !bitcoinRegexp.MatchString(u.Bitcoin) {
		return errors.New("Ivalid BTC address")
	}
	if u.Ethereum != "" && !ethereumRegexp.MatchString(u.Ethereum) {
		return errors.New("Ivalid ETH address")
	}
	if u.TwoFactorAuthentication && u.Pgp == "" {
		return errors.New("Fill PGP for 2FA")
	}
	if u.Email != "" && !emailRegexp.MatchString(u.Email) {
		return errors.New("Invalid email/xmpp address")
	}
	if u.InviteCode == "" {
		u.InviteCode = util.GenerateUuid()
	}
	return nil
}

func (u User) PGPFingerprint() string {
	return util.Fingerprint(u.Pgp)
}

func (s User) Score() float64 {
	score := float64(0.0)
	for _, r := range s.RatingReviews {
		score += float64(r.SellerScore)
	}
	if len(s.RatingReviews) > 0 {
		score /= float64(len(s.RatingReviews))
	}
	return math.Ceil(score*100) / float64(100.0)
}

func (s User) AllowedToSell() bool {
	return s.IsSeller && (s.IsGoldAccount || s.IsSilverAccount || s.IsBronzeAccount || s.IsFreeAccount)
}

/*
	Relations
*/
func (u User) Iniviter() *User {
	if u.InviterUuid == "" {
		return nil
	}
	user, _ := FindUserByUuid(u.InviterUuid)
	return user
}

func (u User) Fingerprint() string {
	return util.Fingerprint(u.Pgp)
}

func (u User) FindUserBitcoinWallets() UserBitcoinWallets {
	var uw UserBitcoinWallets

	database.
		Where(&UserBitcoinWallet{UserUuid: u.Uuid}).
		Order("created_at DESC").
		Find(&uw)

	if len(uw) == 0 {
		userWallet, err := CreateBitcoinWallet(u)
		if err != nil {
			return uw
		}

		return UserBitcoinWallets{
			*userWallet,
		}
	}

	return UserBitcoinWallets(uw)
}

func (u User) FindUserBitcoinCashWallets() UserBitcoinCashWallets {
	var uw UserBitcoinCashWallets

	database.
		Where(&UserBitcoinCashWallet{UserUuid: u.Uuid}).
		Order("created_at DESC").
		Find(&uw)

	if len(uw) == 0 {
		userWallet, err := CreateBitcoinCashWallet(u)
		if err != nil {
			return uw
		}

		return UserBitcoinCashWallets{
			*userWallet,
		}
	}

	return UserBitcoinCashWallets(uw)
}

func (u User) FindRecentBitcoinWallet() UserBitcoinWallet {
	return u.FindUserBitcoinWallets()[0]
}

func (u User) FindRecentBitcoinCashWallet() UserBitcoinCashWallet {
	return u.FindUserBitcoinCashWallets()[0]
}

func (u User) FindUserEthereumWallets() UserEthereumWallets {
	var uw UserEthereumWallets

	database.
		Where(&UserEthereumWallet{UserUuid: u.Uuid}).
		Order("created_at DESC").
		Find(&uw)

	if len(uw) == 0 {
		userWallet, err := CreateEthereumWallet(u)
		if err != nil {
			return uw
		}
		return UserEthereumWallets{*userWallet}
	}

	return UserEthereumWallets(uw)
}

func (u User) FindRecentEthereumWallet() UserEthereumWallet {
	return u.FindUserEthereumWallets()[0]
}

/*
	Queries
*/

func SettingsChangeHistoryByUser(uuid string) []UserSettingsHistory {
	var history []UserSettingsHistory
	database.
		Where(&UserSettingsHistory{UserUuid: uuid}).
		Find(&history)
	return history
}

func CreateUser(username string, passphrase string) (*User, error) {
	user, _ := FindUserByUsername(username)
	if user != nil {
		return nil, errors.New("Username is not unique")
	}
	invitation, _ := FindInvitationByUsername(username)
	if invitation != nil {
		return nil, errors.New("Username is reserved")
	}

	now := time.Now()

	user = &User{
		Uuid:             util.GenerateUuid(),
		InviteCode:       util.GenerateUuid(),
		PassphraseHash:   util.PasswordHashV1(username, passphrase),
		Username:         username,
		RegistrationDate: time.Now(),
		Currency:         "USD",
		Language:         "en",
		LastLoginDate:    &now,
	}
	validationError := user.Validate()
	if validationError != nil {
		return nil, validationError
	}
	user.Save()

	return user, nil
}

func GetExtendedUsersPage(page, pageSize int, orderBy, query string) []ExtendedUser {

	if orderBy == "last_login" {
		orderBy = "last_login_date desc nulls last"
	} else if orderBy == "invited_users" {
		orderBy = "inviter_count desc nulls last"
	} else if orderBy == "balance" {
		orderBy = "bitcoin_balance desc nulls last, bitcoin_cash_balance desc nulls last, ethereum_balance desc nulls last"
	} else {
		orderBy = "registration_date desc nulls last"
	}

	var users []ExtendedUser
	qry := database.
		Table("v_users").
		Model(&ExtendedUser{}).
		Order(orderBy).
		Offset(page * pageSize).
		Limit(pageSize)

	if query != "" {
		qry = qry.Where("username LIKE ?", "%"+query+"%")
	}

	qry.Find(&users)

	return users
}

func GetInvitedUserPage(inviterUuid string, page, pageSize int) []ExtendedUser {
	orderBy := "registration_date desc nulls last"
	var users []ExtendedUser
	qry := database.
		Table("v_users").
		Where("inviter_uuid = ?", inviterUuid).
		Model(&ExtendedUser{}).
		Order(orderBy).
		Offset(page * pageSize).
		Limit(pageSize)
	qry.Find(&users)
	return users
}

func CountUsers(dt *time.Time) int {
	var count int
	q := database.Table("users")
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountVendors(dt *time.Time) int {
	var count int
	q := database.Table("v_users").Where("is_seller=?", true)
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountFreeVendors(dt *time.Time) int {
	var count int
	q := database.Table("v_users").Where("is_seller=true and is_free_account=true")
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountGoldVendors(dt *time.Time) int {
	var count int
	q := database.Table("v_users").Where("is_gold_account=true and is_seller=true")
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountSilverVendors(dt *time.Time) int {
	var count int
	q := database.Table("v_users").Where("is_silver_account=true and is_seller=true")
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountBronzeVendors(dt *time.Time) int {
	var count int
	q := database.Table("v_users").Where("is_bronze_account=true and is_seller=true")
	if dt != nil {
		q = q.Where("registration_date < ?", dt)
	}
	q.Count(&count)
	return count
}

func CountUsersRegistredAfter(dt time.Time) int {
	var count int
	database.Table("users").Where("registration_date > ?", dt).Count(&count)
	return count
}

func CountUsersActiveAfter(dt time.Time) int {
	var count int
	database.Table("users").Where("last_login_date > ?", dt).Count(&count)
	return count
}

func FindUserByUuid(uuid string) (*User, error) {
	var user User

	err := database.
		Preload("Items").
		Preload("Deposits").
		Preload("ShippingOptions").
		First(&user, "uuid = ?", uuid).
		Error

	if err != nil {
		return nil, err
	}

	return &user, err
}

func FindUserByUsername(username string) (*User, error) {
	var user User
	err := database.
		Preload("Items").
		Preload("Deposits").
		Preload("ShippingOptions").
		First(&user, "username = ?", username).
		Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func FindValidVerificationRequests() (Vendors, error) {
	var users Vendors
	err := database.
		Table("v_users").
		Joins("JOIN v_verification_threads on v_verification_threads.uuid='store-verification-'||v_users.uuid").
		Where("v_verification_threads.last_message_by_staff=?", false).
		Where("is_seller = ?", true).
		Where("pgp <> ?", "").
		Where("(is_gold_account=true OR is_silver_account=true OR is_bronze_account=true OR is_free_account)").
		Where("is_signed=?", true).
		Where("is_trusted_seller=?", false).
		Order("has_requested_verification DESC NULLS LAST, registration_date DESC").
		Find(&users).Error

	return users, err
}

func FindUserByInviteCode(code string) (*User, error) {
	var user User
	err := database.
		First(&user, "invite_code = ?", code).
		Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

/*
	Staff
*/

func FindStaffMemebers() ([]ExtendedUser, error) {
	var users []ExtendedUser
	err := database.
		Table("v_users").
		Where("is_staff=true OR is_admin=true").
		Order("is_admin DESC, is_staff DESC, last_login_date ASC").
		Find(&users).
		Error
	return users, err
}

/*
	Model Methods
*/

func (u User) Save() error {
	err := u.Validate()
	if err != nil {
		return err
	}
	return u.SaveToDatabase()
}

func (u User) SaveToDatabase() error {
	if existing, _ := FindUserByUuid(u.Uuid); existing == nil {
		return database.Create(&u).Error
	}
	return database.Save(&u).Error
}

func (u User) Remove() error {
	return database.Delete(&u).Error
}

func (u UserSettingsHistory) Save() error {
	return database.Save(&u).Error
}

func (u User) GenerateAvatar() error {
	if u.HasAvatar {
		return nil
	}

	err := govatar.GenerateFileFromUsername(govatar.MALE, u.Username, "data/images/"+u.Uuid+"_av.jpeg")
	if err != nil {
		return err
	}

	u.HasAvatar = true
	return u.Save()
}

/*
	View Models
*/

type ViewUser struct {
	*User
	ViewItems            []ViewItem         `json:"-"`
	ItemCategories       []ItemCategory     `json:"-"`
	ViewRatingReviews    []ViewRatingReview `json:"-"`
	InviterUsername      string             `json:"-"`
	RegistrationDateStr  string             `json:"-"`
	LastLoginDateStr     string             `json:"-"`
	NumberOfItems        int                `json:"-"`
	Score                string             `json:"-"`
	ScoreFloat           float64            `json:"-"`
	IsOnline             bool               `json:"-"`
	MultisigEnabled      bool               `json:"-"`
	LongDescriptionHTML  template.HTML      `json:"-"`
	NumberOfSales        int                `json:"-"`
	NumberOfSalesStr     string             `json:"-"`
	SalesVolume          float64            `json:"-"`
	SalesVolumeStr       string             `json:"-"`
	NumberOfPurchases    int                `json:"-"`
	PurchaseVolume       float64            `json:"-"`
	RequiredVerification bool               `json:"-"`
	AllowedToSell        bool
	// Bitcoin Wallets `json:"-"`
	BitcoinBalance       apis.BTCWalletBalance     `json:"-"`
	BitcoinWallets       UserBitcoinWallets        `json:"-"`
	BitcoinWallet        UserBitcoinWallet         `json:"-"`
	BitcoinWalletActions []UserBitcoinWalletAction `json:"-"`
	// Ethereum Wallets
	EthereumBalance       apis.ETHWalletBalance      `json:"-"`
	EthereumWallets       UserEthereumWallets        `json:"-"`
	EthereumWallet        UserEthereumWallet         `json:"-"`
	EthereumWalletActions []UserEthereumWalletAction `json:"-"`
	// Bitcoin Cash Wallets
	BitcoinCashBalance       apis.BCHWalletBalance         `json:"-"`
	BitcoinCashWallets       UserBitcoinCashWallets        `json:"-"`
	BitcoinCashWallet        UserBitcoinCashWallet         `json:"-"`
	BitcoinCashWalletActions []UserBitcoinCashWalletAction `json:"-"`
	// Warnings
	NumberOfRedWarnings    int `json:"number_of_red_warnings"`
	NumberOfYellowWarnings int `json:"number_of_yellow_warnings"`
	NumberOfGreenWarnings  int `json:"number_of_green_warnings"`
	// Deposits
	HasDeposits     bool            `json:"has_deposits"`
	DepositsSummary DepositsSummary `json:"deposits_summary"`
	// Levels
	StaffStats  StaffStats
	VendorStats VendorStats
	StaffLevel  int `json:"vendor_level"`
	VendorLevel int `json:"vendor_level"`
	BuyerLevel  int `json:"buyer_level"`
}

type ViewUsers []ViewUser

type ViewExtendedUser struct {
	*ExtendedUser
	Balance             float64
	InviterUsername     string
	RegistrationDateStr string
	LastLoginDateStr    string
	IsOnline            bool
	// Levels
	StaffStats  StaffStats
	VendorStats VendorStats
	StaffLevel  int `json:"vendor_level"`
	VendorLevel int `json:"vendor_level"`
	BuyerLevel  int `json:"buyer_level"`
}

func (u ExtendedUser) ViewExtendedUser(lang string) ViewExtendedUser {

	user := ViewExtendedUser{
		ExtendedUser:        &u,
		RegistrationDateStr: util.HumanizeTime(u.RegistrationDate, lang),
	}

	if u.LastLoginDate != nil {

		if u.IsAdmin {
			dt := time.Now()
			u.LastLoginDate = &dt
		}

		user.LastLoginDateStr = util.HumanizeTime(*u.LastLoginDate, lang)
		user.IsOnline = u.LastLoginDate.After(time.Now().Add(-onlineDuration))
	}

	if u.IsStaff || u.IsAdmin {
		user.StaffStats = CacheGetStaffStats(u.Uuid)
		user.StaffLevel = CalculateStaffLevel(user.StaffStats)

	} else if u.IsSeller {
		user.VendorStats = CacheGetVendorStats(u.Uuid)
		user.VendorLevel = CalculateVendorLevel(user.VendorStats, user.RegistrationDate)
	} else {
		// buyer level code
		user.BuyerLevel = 1
	}

	return user
}

func (users ExtendedUsers) ViewExtendedUsers(lang string) []ViewExtendedUser {
	viewUsers := []ViewExtendedUser{}
	for _, user := range users {
		viewUsers = append(viewUsers, user.ViewExtendedUser(lang))
	}
	return viewUsers
}

func (s User) ViewUser(lang string) ViewUser {

	var lastLoginDate string = "?"
	if s.LastLoginDate != nil {

		if s.IsAdmin {
			dt := time.Now()
			s.LastLoginDate = &dt
		}

		if lang == "ru" {
			lastLoginDate = util.HumanizeTimeRU(*s.LastLoginDate)
		} else {
			lastLoginDate = humanize.Time(*s.LastLoginDate)
		}
	}

	score := s.Score()

	viewSeller := ViewUser{
		LastLoginDateStr:    lastLoginDate,
		LongDescriptionHTML: template.HTML(userHtmlPolicy.Sanitize(string(blackfriday.MarkdownCommon([]byte(s.LongDescription))))),
		MultisigEnabled:     (s.BitcoinMultisigPublicKey != "" && s.AllowedToSell()),
		RegistrationDateStr: humanize.Time(s.RegistrationDate),
		Score:               humanize.Ftoa(score),
		ScoreFloat:          score,
		AllowedToSell:       s.AllowedToSell(),
		User:                &s,
	}

	if s.IsStaff || s.IsAdmin {
		viewSeller.StaffStats = CacheGetStaffStats(s.Uuid)
		viewSeller.StaffLevel = CalculateStaffLevel(viewSeller.StaffStats)

	} else if s.IsSeller {
		viewSeller.VendorStats = CacheGetVendorStats(s.Uuid)
		viewSeller.VendorLevel = CalculateVendorLevel(viewSeller.VendorStats, viewSeller.RegistrationDate)
	} else {
		// buyer level code
		viewSeller.BuyerLevel = 1
	}

	if s.LastLoginDate != nil {
		viewSeller.IsOnline = s.LastLoginDate.After(time.Now().Add(-onlineDuration))
	}

	for _, item := range s.Items {
		nc := item.ItemCategory
		itemExists := false
		for _, ec := range viewSeller.ItemCategories {
			if ec.ID == nc.ID {
				itemExists = true
				break
			}
		}
		if !itemExists {
			viewSeller.ItemCategories = append(viewSeller.ItemCategories, nc)
		}
	}

	for _, review := range viewSeller.RatingReviews {
		if review.SellerReview == "" {
			continue
		}
		viewSeller.ViewRatingReviews = append(viewSeller.ViewRatingReviews, review.ViewRatingReview(lang))
	}

	sortFunc := func(i, j int) bool {
		return viewSeller.ViewRatingReviews[i].CreatedAt.After(*viewSeller.ViewRatingReviews[j].CreatedAt)
	}
	sort.Slice(viewSeller.ViewRatingReviews, sortFunc)

	if lang == "ru" {
		viewSeller.RegistrationDateStr = util.HumanizeTimeRU(s.RegistrationDate)
	}

	inviter := s.Iniviter()
	if inviter != nil {
		viewSeller.InviterUsername = inviter.Username
	}

	if viewSeller.IsSeller && !viewSeller.IsTrustedSeller {
		viewSeller.RequiredVerification = (viewSeller.Pgp == "" || !viewSeller.IsSigned || !s.AllowedToSell())
	}

	for _, uw := range s.UserWarnings {
		if !uw.HasExpired() && uw.IsApproved {
			switch uw.Severety {
			case "RED":
				viewSeller.NumberOfRedWarnings += 1

			case "YELLOW":
				viewSeller.NumberOfYellowWarnings += 1

			case "GREEN":
				viewSeller.NumberOfGreenWarnings += 1
			}
		}
	}

	viewSeller.HasDeposits = len(s.Deposits) > 0
	viewSeller.DepositsSummary = Deposits(s.Deposits).DepositsSummary()

	return viewSeller
}

func (users Users) ViewUsers(lang string) []ViewUser {
	viewUsers := []ViewUser{}
	for _, user := range users {
		viewUsers = append(viewUsers, user.ViewUser(lang))
	}
	return viewUsers
}

/*
	Cache
*/

func CacheGetUserUuid(username string) string {
	cUuid, _ := cache1h.Get(username)
	if cUuid == nil {
		user, err := FindUserByUsername(username)
		if err != nil {
			return ""
		}
		cache1h.Set(username, user.Uuid)
		return user.Uuid
	}

	return cUuid.(string)
}

/*
	Database Views
*/

func setupUserViews() {
	database.Exec("DROP VIEW IF EXISTS v_users CASCADE;")
	database.Exec(`
		CREATE VIEW v_users AS (
             select
				u1.*,
				u2.username as inviter_username,
				(select count(*) from users u3 where u3.inviter_uuid=u1.uuid) as inviter_count,
				(select count(*) from user_warnings uw where uw.user_uuid=u1.uuid and severety='RED' and is_approved=true and deleted_at is NULL) as red_warning_count,
				(select count(*) from user_warnings uw where uw.user_uuid=u1.uuid and severety='GREEN' and is_approved=true and deleted_at is NULL) as green_warning_count,
				(select count(*) from user_warnings uw where uw.user_uuid=u1.uuid and severety='YELLOW' and is_approved=true and deleted_at is NULL) as yellow_warning_count,
				COALESCE(vubb.balance, 0) as bitcoin_balance,
				COALESCE(vubb.unconfirmed_balance, 0) as bitcoin_unconfirmed_balance,
				COALESCE(vueb.balance, 0) as ethereum_balance,
				COALESCE(vubcb.balance, 0) as bitcoin_cash_balance,
				COALESCE(vubcb.unconfirmed_balance, 0) as bitcoin_cash_unconfirmed_balance
     		from
				users u1
     		left outer join
				users u2 on u2.uuid=u1.inviter_uuid
     		left outer join
				v_user_bitcoin_wallet_balances vubb on vubb.user_uuid=u1.uuid
     		left outer join
				v_user_ethereum_wallet_balances vueb on vueb.user_uuid=u1.uuid
     		left outer join
				v_user_bitcoin_cash_wallet_balances vubcb on vubcb.user_uuid=u1.uuid
    );
    `)
}

func refreshUsesMaterializedView() {
	database.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY vm_users ;")
}
