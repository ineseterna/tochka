package marketplace

import (
	"fmt"
	"math"
	"sort"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type AvaiableItem struct {
	ItemUuid string `json:"item_uuid"`

	// Vendor details
	VendorUuid                      string    `json:"vendor_uuid"`
	VendorUsername                  string    `json:"vendor_username"`
	VendorDescription               string    `json:"vendor_description"`
	VendorLanguage                  string    `json:"vendor_language"`
	IsTrustedSeller                 bool      `json:"vendor_is_trusted"`
	IsSigned                        bool      `json:"vendor_is_trusted"`
	LastLoginDate                   time.Time `json:"vendor_last_login_date"`
	RegistrationDate                time.Time `json:"vendor_registration_date"`
	BitcoinMultisigPublicKeyEnabled bool      `json:"-"`
	NumberOfRedWarnings             int       `json:"number_of_red_warnings"`
	NumberOfGreenWarnings           int       `json:"number_of_green_warnings"`
	NumberOfYellowWarnings          int       `json:"number_of_yellow_warnings"`
	ReviewedByUserUuid              string    `json:"reviewed_by_user_uuid"`
	// Vendor account type
	VendorIsFreeAccount   bool `json:"vendor_is_free_account"`
	VendorIsGoldAccount   bool `json:"vendor_is_gold_account"`
	VendorIsSilverAccount bool `json:"vendor_is_silver_account"`
	VendorIsBronzeAccount bool `json:"vendor_is_bronze_account"`
	// Item details
	Type                       string    `json:"type"`
	ItemCreatedAt              time.Time `json:"item_created_at"`
	ItemName                   string    `json:"item_name"`
	ItemDescription            string    `json:"item_description"`
	ItemCategoryId             int       `json:"item_category_id"`
	ParentItemCategoryId       int       `json:"item_parent_category_id"`
	ParentParentItemCategoryId int       `json:"item_parent_parent_category_id"`
	// Price details
	MinPrice float64 `json:"-"`
	MaxPrice float64 `json:"-"`
	Currency string  `json:"-"`
	// Scores
	VendorScore      float64 `json:"vendor_score"`
	VendorScoreCount int     `json:"vendor_score_count"`
	ItemScore        float64 `json:"item_score"`
	ItemScoreCount   int     `json:"item_score_count"`

	VendorBitcoinTxNumber float64 `json:"-"`
	VendorBitcoinTxVolume float64 `json:"-"`
	ItemBitcoinTxNumber   float64 `json:"-"`
	ItemBitcoinTxVolume   float64 `json:"-"`

	VendorBitcoinCashTxNumber float64 `json:"-"`
	VendorBitcoinCashTxVolume float64 `json:"-"`
	ItemBitcoinCashTxNumber   float64 `json:"-"`
	ItemBitcoinCashTxVolume   float64 `json:"-"`

	VendorEthereumTxNumber float64 `json:"-"`
	VendorEthereumTxVolume float64 `json:"-"`
	ItemEthereumTxNumber   float64 `json:"-"`
	ItemEthereumTxVolume   float64 `json:"-"`

	CountryNameEnShippingFrom string `json:"country_shipping_from"`
	CountryNameEnShippingTo   string `json:"country_shipping_to"`
	DropCityId                int    `json:"geoname_id"`

	Price map[string][2]float64 `json:"price"`

	VendorLevel int

	GeoCity        City    `gorm:"ForeignKey:DropCityId" json:"-"`
	GeoCountryFrom Country `gorm:"ForeignKey:CountryNameEnShippingFrom" json:"-"`
	GeoCountryTo   Country `gorm:"ForeignKey:CountryNameEnShippingTo" json:"-"`
}

type Vendor struct {
	Username         string
	LastLoginDate    time.Time
	RegistrationDate time.Time

	BitcoinTxNumber     float64 `json:"btc_tx_number"`
	BitcoinTxVolume     float64 `json:"btc_tx_volume"`
	BitcoinCashTxNumber float64 `json:"bch_tx_number"`
	BitcoinCashTxVolume float64 `json:"bch_tx_volume"`
	EthereumTxNumber    float64 `json:"eth_tx_number"`
	EthereumTxVolume    float64 `json:"eth_tx_volume"`

	VendorScore float64 `json:"vendor_score"`
	VendorLevel int     `json:"vendor_level"`

	VendorDescription string `json:"vendor_description"`
	Language          string `json:"language"`

	IsFreeAccount            bool `json:"vendor_is_free_account"`
	IsGoldAccount            bool `json:"vendor_is_gold_account"`
	IsSilverAccount          bool `json:"vendor_is_silver_account"`
	IsBronzeAccount          bool `json:"vendor_is_bronze_account"`
	IsTrustedSeller          bool `json:"vendor_is_trusted"`
	HasRequestedVerification bool `json:"-"`
	IsSigned                 bool `json:"vendor_agreement_signed"`
}

type Vendors []Vendor

/*
	Model Methods
*/

func (ai AvaiableItem) USDVolume() float64 {
	return ai.ItemBitcoinTxVolume*GetCurrencyRate("BTC", "USD") +
		ai.ItemBitcoinCashTxVolume*GetCurrencyRate("BCH", "USD") +
		ai.ItemEthereumTxVolume*GetCurrencyRate("ETH", "USD")
}

/*
	Currency Rates
*/

func (ai AvaiableItem) GetPrice(currency string) [2]float64 {
	return [2]float64{
		ai.MinPrice / GetCurrencyRate(currency, ai.Currency),
		ai.MaxPrice / GetCurrencyRate(currency, ai.Currency),
	}
}

/*
	View Item
*/

type ViewAvailableItem struct {
	*AvaiableItem
	IsOnline            bool      `json:"vendor_is_online"`
	LastLoginDateStr    string    `json:"vendor_last_login_date"`
	RegistrationDateStr string    `json:"vendor_registration_date"`
	PriceRangeStr       [2]string `json:"price_range"`
	PriceStr            string    `json:"price"`
}

func (ai AvaiableItem) ViewAvailableItem(lang, currency string) ViewAvailableItem {

	price := ai.Price[currency]

	vai := ViewAvailableItem{
		AvaiableItem:        &ai,
		LastLoginDateStr:    util.HumanizeTime(ai.LastLoginDate, lang),
		PriceRangeStr:       [2]string{fmt.Sprintf("%d", int(math.Ceil(price[0]))), fmt.Sprintf("%d", int(math.Ceil(price[1])))},
		RegistrationDateStr: util.HumanizeTime(ai.RegistrationDate, lang),
	}

	if currency == "BTC" || currency == "ETH" || currency == "BCH" {
		vai.PriceRangeStr = [2]string{
			fmt.Sprintf("%f", price[0]),
			fmt.Sprintf("%f", price[1]),
		}
	}

	if price[0] == price[1] {
		vai.PriceStr = vai.PriceRangeStr[0]
	}

	return vai
}

type ViewVendor struct {
	*Vendor
	IsOnline            bool
	VendorScoreStr      string
	LastLoginDateStr    string
	RegistrationDateStr string
}

func (v Vendor) ViewVendor(lang string) ViewVendor {
	return ViewVendor{
		LastLoginDateStr:    util.HumanizeTime(v.LastLoginDate, lang),
		RegistrationDateStr: util.HumanizeTime(v.RegistrationDate, lang),
		Vendor:              &v,
	}
}

type AvailableItems []AvaiableItem

func (ais AvailableItems) ViewAvailableItems(lang, currency string) []ViewAvailableItem {
	var vais []ViewAvailableItem
	for _, ai := range ais {
		vai := ai.ViewAvailableItem(lang, currency)
		vais = append(vais, vai)
	}
	return vais
}

func (vs Vendors) ViewVendors(lang string) []ViewVendor {
	var vvs []ViewVendor
	for _, v := range vs {
		vv := v.ViewVendor(lang)
		vvs = append(vvs, vv)
	}
	return vvs
}

/*
	Collection Fields
*/

func (ais AvailableItems) Sort(sortyBy string) AvailableItems {

	var sortByFunc func(int, int) bool

	switch sortyBy {
	case "date_logged_in":
		sortByFunc = func(i, j int) bool {
			return ais[i].LastLoginDate.After(ais[j].LastLoginDate)
		}
	case "price":
		sortByFunc = func(i, j int) bool {

			if ais[i].VendorIsGoldAccount != ais[j].VendorIsGoldAccount { // GOLD
				return ais[i].VendorIsGoldAccount
			} else if ais[i].VendorIsSilverAccount != ais[j].VendorIsSilverAccount { // SILVER
				return ais[i].VendorIsSilverAccount
			} else if ais[i].VendorIsBronzeAccount != ais[j].VendorIsBronzeAccount { // BRONZE
				return ais[i].VendorIsBronzeAccount
			} else { // by price
				return ais[i].MinPrice < ais[j].MaxPrice
			}
		}
	case "popularity":
		sortByFunc = func(i, j int) bool {

			if ais[i].VendorIsGoldAccount != ais[j].VendorIsGoldAccount { // GOLD
				return ais[i].VendorIsGoldAccount
			} else if ais[i].VendorIsSilverAccount != ais[j].VendorIsSilverAccount { // SILVER
				return ais[i].VendorIsSilverAccount
			} else if ais[i].VendorIsBronzeAccount != ais[j].VendorIsBronzeAccount { // BRONZE
				return ais[i].VendorIsBronzeAccount
			} else { // by price
				return ais[i].USDVolume() > ais[j].USDVolume()
			}

		}
	case "date_added":
		sortByFunc = func(i, j int) bool {
			return ais[i].ItemCreatedAt.Before(ais[j].ItemCreatedAt)
		}
	case "rating":
		sortByFunc = func(i, j int) bool {
			return ais[i].ItemScore*float64(ais[i].VendorLevel) >
				ais[j].ItemScore*float64(ais[i].VendorLevel)

		}
	default:
		sortByFunc = func(i, j int) bool { return true }
	}

	sort.Slice(ais, sortByFunc)
	return ais
}

func (ais AvailableItems) Where(predicate func(AvaiableItem) bool) AvailableItems {
	newAis := AvailableItems{}
	for i, _ := range ais {
		if predicate(ais[i]) {
			newAis = append(newAis, ais[i])
		}
	}
	return newAis
}

func (ais AvailableItems) Filter(category, dropCityId int,
	packageType, query, to, from, accountType, vendorUuid string) AvailableItems {

	var searchResults []string
	if query != "" {
		searchResults = SearchItems(query)
	}

	categoryPredicate := func(ai AvaiableItem) bool {
		if category != 0 {
			return ai.ItemCategoryId == category || ai.ParentItemCategoryId == category || ai.ParentParentItemCategoryId == category
		}
		return true
	}

	accountTypePredicate := func(ai AvaiableItem) bool {
		switch accountType {
		case "gold":
			return ai.VendorIsGoldAccount
		case "silver":
			return ai.VendorIsSilverAccount
		case "bronze":
			return ai.VendorIsBronzeAccount
		case "free":
			return ai.VendorIsFreeAccount
		default:
			return true
		}
	}

	typePredicate := func(ai AvaiableItem) bool {
		if packageType != "" && packageType != "all" {
			return ai.Type == packageType
		}
		return true
	}

	queryPredicate := func(ai AvaiableItem) bool {
		if query != "" && !inSet(ai.ItemUuid, searchResults) {
			return false
		}
		return true
	}

	shipingToPredicate := func(ai AvaiableItem) bool {
		if to == "" {
			return true
		}
		return ai.CountryNameEnShippingTo == to
	}

	shipingFromPredicate := func(ai AvaiableItem) bool {
		if from == "" {
			return true
		}
		return ai.CountryNameEnShippingFrom == from
	}

	cityPredicate := func(ai AvaiableItem) bool {
		if dropCityId == 0 {
			return true
		}
		return ai.DropCityId == dropCityId
	}

	vendorPredicate := func(ai AvaiableItem) bool {
		if vendorUuid == "" {
			return true
		}
		return ai.VendorUuid == vendorUuid
	}

	filteredAvailableItems := ais.
		Where(typePredicate).
		Where(categoryPredicate).
		Where(queryPredicate).
		Where(shipingToPredicate).
		Where(shipingFromPredicate).
		Where(cityPredicate).
		Where(accountTypePredicate).
		Where(vendorPredicate)

	return filteredAvailableItems
}

func (ais AvailableItems) DropCitiesList() []City {
	locationMap := map[int]City{}
	for _, a := range ais {
		locationMap[a.DropCityId] = a.GeoCity
	}
	locations := []City{}
	for _, city := range locationMap {
		locations = append(locations, city)
	}
	return locations
}

func (ais AvailableItems) ShippingToList() []string {
	locationMap := map[string]bool{}
	for _, a := range ais {
		locationMap[a.CountryNameEnShippingTo] = true
	}
	locations := []string{}
	for l, _ := range locationMap {
		locations = append(locations, l)
	}
	sort.Strings(locations)
	return locations
}

func (ais AvailableItems) ShippingFromList() []string {
	locationMap := map[string]bool{}
	for _, a := range ais {
		locationMap[a.CountryNameEnShippingFrom] = true
	}
	locations := []string{}
	for l, _ := range locationMap {
		locations = append(locations, l)
	}
	sort.Strings(locations)
	return locations
}

func (ais AvailableItems) VendorList() Vendors {
	vendorsMap := map[string]Vendor{}
	for _, a := range ais {

		vendorLevel := CalculateVendorLevel(
			VendorStats{
				NumberOfReleasedTransactions: int(a.VendorBitcoinTxNumber + a.VendorBitcoinCashTxNumber + a.VendorEthereumTxNumber),
			},
			a.RegistrationDate,
		)

		vendorsMap[a.VendorUsername] = Vendor{
			// Basic info
			Username:          a.VendorUsername,
			LastLoginDate:     a.LastLoginDate,
			RegistrationDate:  a.RegistrationDate,
			VendorDescription: a.VendorDescription,
			IsTrustedSeller:   a.IsTrustedSeller,
			// Numerical characteristics
			VendorLevel: vendorLevel,
			VendorScore: a.VendorScore,
			// Account Type
			IsFreeAccount:   a.VendorIsFreeAccount,
			IsGoldAccount:   a.VendorIsGoldAccount,
			IsSilverAccount: a.VendorIsSilverAccount,
			IsBronzeAccount: a.VendorIsBronzeAccount,
			// Tx Stats
			BitcoinTxNumber:     a.VendorBitcoinTxNumber,
			BitcoinTxVolume:     a.VendorBitcoinTxVolume,
			BitcoinCashTxNumber: a.VendorBitcoinCashTxNumber,
			BitcoinCashTxVolume: a.VendorBitcoinCashTxVolume,
			EthereumTxNumber:    a.VendorEthereumTxNumber,
			EthereumTxVolume:    a.VendorEthereumTxVolume,
		}
	}

	vendors := []Vendor{}
	for _, v := range vendorsMap {
		vendors = append(vendors, v)
	}

	return Vendors(vendors)
}

func (vvs Vendors) Sort(sortyBy string) Vendors {

	var sortByFunc func(int, int) bool

	switch sortyBy {
	case "new":
		sortByFunc = func(i, j int) bool {
			return vvs[i].RegistrationDate.After(vvs[j].RegistrationDate)
		}
	case "date_logged_in":
		sortByFunc = func(i, j int) bool {
			return vvs[i].LastLoginDate.After(vvs[j].LastLoginDate)
		}
	case "popularity":
		sortByFunc = func(i, j int) bool {
			return int(vvs[i].VendorLevel) > int(vvs[j].VendorLevel)
		}
	case "date_added":
		sortByFunc = func(i, j int) bool {
			return vvs[i].RegistrationDate.After(vvs[j].RegistrationDate)
		}
	case "rating":
		sortByFunc = func(i, j int) bool {
			return vvs[i].VendorScore > vvs[j].VendorScore
		}
	default:
		sortByFunc = func(i, j int) bool { return true }
	}

	sort.Slice(vvs, sortByFunc)
	return vvs
}

/*
	Database Queries
*/

func FindAvailableItems(userUuid string) AvailableItems {

	items := []AvaiableItem{}

	if userUuid == "" {
		database.
			Table("vm_available_items").
			Preload("GeoCity").
			Find(&items)
	} else {
		database.
			Table("vm_available_items").
			Where("vendor_uuid=?", userUuid).
			Preload("GeoCity").
			Find(&items)
	}

	for i, _ := range items {
		items[i].Price = map[string][2]float64{
			"AUD": items[i].GetPrice("AUD"),
			"BCH": items[i].GetPrice("BCH"),
			"BTC": items[i].GetPrice("BTC"),
			"ETH": items[i].GetPrice("ETH"),
			"EUR": items[i].GetPrice("EUR"),
			"GBP": items[i].GetPrice("GBP"),
			"RUB": items[i].GetPrice("RUB"),
			"USD": items[i].GetPrice("USD"),
		}

		items[i].VendorLevel = CalculateVendorLevel(VendorStats{
			NumberOfReleasedTransactions: int(
				items[i].VendorBitcoinTxNumber +
					items[i].VendorBitcoinCashTxNumber +
					items[i].VendorEthereumTxNumber),
		}, items[i].RegistrationDate)
		items[i].ItemScore = float64(int(items[i].ItemScore*100)) / float64(100.0)
		items[i].VendorScore = float64(int(items[i].VendorScore*100)) / float64(100.0)
	}

	return AvailableItems(items)
}

/*
	Utils
*/

func inSet(s string, ss []string) bool {
	for _, i := range ss {
		if i == s {
			return true
		}
	}
	return false
}

/*
	Database Views
*/

func setupAvailableItemsView() {

	database.Exec("DROP VIEW IF EXISTS v_available_items CASCADE;")
	database.Exec(`
		CREATE VIEW v_available_items AS (
			SELECT * FROM
			(
				select 
					v_packages.item_uuid,
					v_packages.drop_city_id,
					v_packages.country_name_en_shipping_from,
					v_packages.country_name_en_shipping_to,
					v_packages.currency,
					min(v_packages.price) as min_price,
					max(v_packages.price) as max_price,
					users.uuid as vendor_uuid,
					users.description as vendor_description,
					users.username as vendor_username,
					users.language as vendor_language,
					users.is_free_account as vendor_is_free_account,
					users.is_gold_account as vendor_is_gold_account,
					users.is_silver_account as vendor_is_silver_account,
					users.is_bronze_account as vendor_is_bronze_account,
					users.is_signed,
					users.is_trusted_seller,
					users.last_login_date,
					users.registration_date,
					users.bitcoin_multisig_public_key != '' as bitcoin_multisig_public_key_enabled, 
					type,
					items.created_at as item_created_at, 
					items.name as item_name,
					items.description as item_description,
					items.item_category_id,
					items.reviewed_by_user_uuid,
					ic_parent.id as parent_item_category_id,
					ic_parent.parent_id as parent_parent_item_category_id,
					COALESCE(avg(r1.seller_score), 0) as vendor_score,
					COALESCE(count(r1.seller_score), 0) as vendor_score_count,
					COALESCE(avg(r2.item_score), 0) as item_score,
					COALESCE(count(r2.item_score), 0) as item_score_count,
					AVG(COALESCE(v_vendor_bitcoin_tx_stats.tx_number, 0)) as vendor_bitcoin_tx_number, 
					AVG(COALESCE(v_vendor_bitcoin_tx_stats.tx_volume, 0)) as vendor_bitcoin_tx_volume,
					AVG(COALESCE(v_vendor_bitcoin_cash_tx_stats.tx_number, 0)) as vendor_bitcoin_cash_tx_number, 
					AVG(COALESCE(v_vendor_bitcoin_cash_tx_stats.tx_volume, 0)) as vendor_bitcoin_cash_tx_volume,
					AVG(COALESCE(v_vendor_ethereum_tx_stats.tx_number, 0)) as vendor_ethereum_tx_number, 
					AVG(COALESCE(v_vendor_ethereum_tx_stats.tx_volume, 0)) as vendor_ethereum_tx_volume,
					AVG(COALESCE(v_item_bitcoin_tx_stats.tx_number, 0)) as item_bitcoin_tx_number,
					AVG(COALESCE(v_item_bitcoin_tx_stats.tx_volume, 0)) as item_bitcoin_tx_volume,
					AVG(COALESCE(v_item_bitcoin_cash_tx_stats.tx_number, 0)) as item_bitcoin_cash_tx_number,
					AVG(COALESCE(v_item_bitcoin_cash_tx_stats.tx_volume, 0)) as item_bitcoin_cash_tx_volume,
					AVG(COALESCE(v_item_ethereum_tx_stats.tx_number, 0)) as item_ethereum_tx_number,
					AVG(COALESCE(v_item_ethereum_tx_stats.tx_volume, 0)) as item_ethereum_tx_volume,
					(select count(*) from user_warnings uw where uw.user_uuid=users.uuid and severety='GREEN' and deleted_at IS NULL and is_approved=true and uw.created_at >= now() - interval '1 week') as number_of_green_warnings,
					(select count(*) from user_warnings uw where uw.user_uuid=users.uuid and severety='YELLOW' and deleted_at IS NULL and is_approved=true and uw.created_at >= now() - interval '2 weeks') as number_of_yellow_warnings,
					(select count(*) from user_warnings uw where uw.user_uuid=users.uuid and severety='RED' and deleted_at IS NULL and uw.created_at >= (now() - interval '4 weeks') and is_approved=true) as number_of_red_warnings
				from v_packages 
				join items on items.uuid = v_packages.item_uuid
				left join v_vendor_bitcoin_tx_stats on v_vendor_bitcoin_tx_stats.seller_uuid = items.user_uuid
				left join v_vendor_bitcoin_cash_tx_stats on v_vendor_bitcoin_cash_tx_stats.seller_uuid = items.user_uuid
				left join v_vendor_ethereum_tx_stats on v_vendor_ethereum_tx_stats.seller_uuid = items.user_uuid
				left join v_item_bitcoin_tx_stats on v_item_bitcoin_tx_stats.item_uuid = items.uuid
				left join v_item_bitcoin_cash_tx_stats on v_item_bitcoin_cash_tx_stats.item_uuid = items.uuid
				left join v_item_ethereum_tx_stats on v_item_ethereum_tx_stats.item_uuid = items.uuid
				join users on users.uuid = items.user_uuid
				left join item_categories ic on ic.id = items.item_category_id
				left join item_categories ic_parent on ic_parent.id = ic.parent_id
				left join rating_reviews r1 on r1.seller_uuid = items.user_uuid
				left join rating_reviews r2 on r2.item_uuid = v_packages.item_uuid
				WHERE 
					items.deleted_at IS NULL AND 
					users.banned=false and v_packages.deleted_at IS NULL AND
					items.reviewed_by_user_uuid <> ''
				group by 
					users.uuid,
					users.username, 
					users.is_gold_account,
					users.is_silver_account,
					users.is_bronze_account,
					users.is_free_account,
					users.is_signed,
					users.is_trusted_seller,
					users.last_login_date, 
					users.registration_date, 
					users.bitcoin_multisig_public_key,
					v_packages.item_uuid, 
					v_packages.drop_city_id,
					v_packages.country_name_en_shipping_from,
					v_packages.country_name_en_shipping_to,
					v_packages.currency,
					r1.seller_uuid, 
					r2.item_uuid, 
					type, 
					items.created_at,
					items.name, 
					items.description,
					items.item_category_id,
					items.reviewed_by_user_uuid,
					parent_item_category_id,
					parent_parent_item_category_id
		) ais
		WHERE number_of_red_warnings=0 and is_signed=true AND is_trusted_seller=true
	);`)

	database.Exec("DROP MATERIALIZED VIEW IF EXISTS vm_available_items CASCADE;")
	database.Exec(`
		CREATE MATERIALIZED VIEW vm_available_items AS (
			SELECT 
				* 
			FROM 
				v_available_items
	);`)

	database.Exec("CREATE UNIQUE INDEX idx_vm_available_items ON vm_available_items (item_uuid, drop_city_id,country_name_en_shipping_from,country_name_en_shipping_to,currency,min_price,max_price);")

}

func RefreshAvailableItemsMaterializedView() {
	database.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY vm_available_items;")
}
