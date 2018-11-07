# Tochka Free Market REST API

List of methods

* Captcha
* Login
* Register
* SERP (list items)
* Vendor info
* Item info
* Transaction list
* Transaction info
* Reply to Transaction Chat
* Wallet list & Balances

## Captcha

### GET http://tochka3evlj3sxdv.onion/captcha/:captcha_id

Use captcha_id from both responses to display captcha: 

http://tochka3evlj3sxdv.onion/captcha/:captcha_id

## Login 

Endpoint: /api/auth/login

### GET http://tochka3evlj3sxdv.onion/api/auth/login

Example:
	
	> curl http://tochka3evlj3sxdv.onion/api/auth/login
	
	{"captcha_id":"jVnRGZ9sQHzaJVPiv9mY"}

Captcha id must be used to obtain captcha image on /captcha/:captcha_id route

### POST http://tochka3evlj3sxdv.onion/api/auth/login

Parameters:

* username
* passphrase
* captcha_id
* captcha

Example:
	
	> curl --data "username=test&passphrase=test&captcha_id=9mkx5dXi94QfSTmLbOTI&captcha=7005" http://tochka3evlj3sxdv.onion/api/auth/login
	
	{
		"api_session": {
			"token": "a6949a9fa5f24de272fbe1e0cab9cdf6",
			"end_date": "2018-03-23T22:33:47.259373345+03:00",
			"is_2fa_session": false,
			"is_2fa_completed": false
		}
	}

**All subsequent requests must include token get parameter**

## Login 2FA

### POST Response (1st Factor)

Parameters:

* username
* passphrase
* captcha_id
* captcha

Example:
	
	> curl --data "username=test&passphrase=test&captcha_id=9mkx5dXi94QfSTmLbOTI&captcha=7005" http://tochka3evlj3sxdv.onion/api/auth/login
	
	{
		"api_session": {
			"token": "a6949a9fa5f24de272fbe1e0cab9cdf6",
			"end_date": "2018-03-23T22:33:47.259373345+03:00",
			"is_2fa_session": true,
			"is_2fa_completed": false
		},
		"secret_text": "xxx"
	}

### POST Response (nd Factor)

Parameters:

* token
* decryptedmessage

Example:
	
	> curl --data "username=test_1&passphrase=test_1&token=a6949a9fa5f24de272fbe1e0cab9cdf6&decryptedmessage=xxx" http://tochka3evlj3sxdv.onion/api/auth/login
	
	{
		"api_session": {
			"token": "a6949a9fa5f24de272fbe1e0cab9cdf6",
			"end_date": "2018-03-23T22:33:47.259373345+03:00",
			"is_2fa_session": true,
			"is_2fa_completed": true
		}
	}

## Register

### GET Response

Example:
	
	> curl http://tochka3evlj3sxdv.onion/api/auth/register
	
	{"captcha_id":"jVnRGZ9sQHzaJVPiv9mY"}

Captcha id must be used to obtain captcha image on /captcha/:captcha_id route

### POST Response

Parameters:

* username
* passphrase_1
* passphrase_2
* captcha_id
* captcha

Example:
	
	> curl --data "username=test&passphrase_1=test&passphrase_2=test&role=seller&captcha_id=CumVFfRdnwEr7VPPbK6s&captcha=6207" http://tochka3evlj3sxdv.onion/auth/login?json=true
	
	{"captcha_id":"WxA3CXmhiUYs8CJnvSVe","error":"Invalid captcha"} 

## Items SERP Response

### GET http://tochka3evlj3sxdv.onion/api/serp?token=:token

Example 

	> curl http://tochka3evlj3sxdv.onion/api/serp?token=:token
	
	{
		"page": 1,
		"number_of_pages": 481,
		"sort_by": "popularity",
		"city": "All countries",
		"geo_cities": [
			{
				"geonameid": 0,
				"name": "",
				"country": "",
				"subcountry": ""
			},
			...
			{
				"geonameid": 6174041,
				"name": "Victoria",
				"country": "Canada",
				"subcountry": "British Columbia"
			}
		],
		"shipping_from_list": [
			"",
			"Afghanistan",
			"Aland Islands",
			...
			"United Kingdom",
			"United States",
			"Uruguay",
			"Worldwide"
		],
		"shipping_to_list": [
			"",
			"Afghanistan",
			"Andorra",
			...
			"United Arab Emirates",
			"United Kingdom",
			"United States",
			"Vatican",
			"Western Sahara",
			"Worldwide",
			"Yemen"
		],
		"account": "all",
		"available_items": [
			{
				"item_uuid": "cc2bfb67891d422e44415d94db75b492",
				"vendor_uuid": "f268db63a6054352705808cda6042e14",
				"vendor_username": "joshkingseller",
				"vendor_description": "Hello everyone I'm here to make u rich  ! No escrow on my product",
				"vendor_language": "en",
				"vendor_is_premium": true,
				"vendor_is_premium_plus": false,
				"vendor_is_trusted": false,
				"type": "digital",
				"item_created_at": "2017-08-18T05:50:22.34331+03:00",
				"item_name": "xxxx",
				"item_description": "xxxx",
				"item_category_id": 25,
				"item_parent_category_id": 23,
				"item_parent_parent_category_id": 0,
				"vendor_score": 4.57,
				"vendor_score_count": 14784,
				"item_score": 4.59,
				"item_score_count": 14784,
				"country_shipping_from": "Interwebs",
				"country_shipping_to": "Interwebs",
				"geoname_id": 0,
				"vendor_is_online": false,
				"vendor_last_login_date": "2 days ago",
				"vendor_registration_date": "8 months ago",
				"price_range": [
					"15",
					"50"
				],
				"price": "",
				"vendor_btc_tx_number": "100+",
				"vendor_btc_tx_volume": "0.5-1 BTC",
				"item_btc_tx_number": "50-100",
				"vendor_eth_tx_number": "10-20",
				"item_eth_tx_number": "5-10"
			}
		],
		"item_categories": [
			{
				"id": 1,
				"parent_id": 0,
				"icon": "",
				"price_en": "Drugs",
				"name_ru": "ÐŸÑÐ¸Ñ…Ð¾Ð°ÐºÑ‚Ð¸Ð²Ð½Ñ‹Ðµ Ð²ÐµÑ‰ÐµÑÑ‚Ð²Ð°",
				"name_de": "",
				"name_es": "",
				"name_fr": "",
				"name_rs": "",
				"name_tr": "",
				"item_count": 5394,
				"user_count": 1939,
				"subcategories": [
				...
			}
		]
	}

## Vendor Info

## GET http://127.0.0.1:8081/api/user/:username

Example
	
	> curl http://127.0.0.1:8081/api/user/:username

	{
		"items": [
			{
				"uuid": "8bf99758c71c4a6f64cac1d41b6f772b",
				"name": "xxxx",
				"description": "xxxx",
				"category_id": 9,
				"user_uuid": "91c76c945f644979449f3ed5cc19d878",
				"is_promoted": false,
				"number_of_sales": 0,
				"number_of_views": 0,
				"ReviewedByUserUuid": "",
				"ReviewedAt": null,
				"created_at": "2018-01-10T12:52:45.30687+03:00",
				"updated_at": "2018-02-07T02:05:50.162589+03:00",
				"deleted_at": null,
				"description_html": "\u003cp\u003exxxx\u003c/p\u003e\n",
				"short_description_html": "\u003cp\u003exxxx\u003c/p\u003e\n",
				"group_packages": [
					{
						"package_name": "3 pcs",
						"country_shipping_from": "united kingdom",
						"country_shipping_to": "worldwide",
						"drop_city_id": 0,
						"drop_city": {
							"geonameid": 0,
							"name": "",
							"country": "",
							"subcountry": ""
						},
						"price_aud": "24",
						"price_btc": "24",
						"price_bch": "24",
						"price_eth": "24",
						"price_eur": "24",
						"price_gbp": "24",
						"price_rub": "24",
						"price_usd": "24",
						"hash": "82f1739975",
						"item_name": "",
						"item_uuid": "",
						"premium": false,
						"score": 0,
						"type": "mail",
						"username": ""
					},
					...
				]
			},
			...
		],
		"vendor": {
			"uuid": "91c76c945f644979449f3ed5cc19d878",
			"username": "username",
			"registration_date": "2017-12-05T03:44:34.280726+03:00",
			"last_login_date": "2018-04-02T10:45:48.293078+03:00",
			"bitmessage": "",
			"tox": "",
			"email": "",
			"Pgp": "",
			"description": "",
			"long_description": "",
			"invite_code": "",
			"is_premium": true,
			"is_premium_plus": false,
			"is_possible_scammer": false,
			"vacation_mode": false,
			"is_seller": true,
			"is_trustedseller": false,
			"is_tester": false,
			"is_moderator": false,
			"is_admin": false,
			"is_staff": false
		}
	}

## Item Details

### GET http://127.0.0.1:8081/api/user/:username/item/:item_uuid 

Example 

	> curl http://127.0.0.1:8081/api/user/:username/item/:item_uuid 

	{
		"item": {
			"uuid": "8bf99758c71c4a6f64cac1d41b6f772b",
			"name": "xxxx",
			"description": "xxxx",
			"category_id": 9,
			"user_uuid": "91c76c945f644979449f3ed5cc19d878",
			"is_promoted": false,
			"number_of_sales": 0,
			"number_of_views": 0,
			"ReviewedByUserUuid": "",
			"ReviewedAt": null,
			"created_at": "2018-01-10T12:52:45.30687+03:00",
			"updated_at": "2018-02-07T02:05:50.162589+03:00",
			"deleted_at": null,
			"description_html": "\u003cp\u003exxxx\u003c/p\u003e\n",
			"short_description_html": "\u003cp\u003exxxx\u003c/p\u003e\n",
			"group_packages": [
				{
					"package_name": "3 pcs",
					"country_shipping_from": "united kingdom",
					"country_shipping_to": "worldwide",
					"drop_city_id": 0,
					"drop_city": {
						"geonameid": 0,
						"name": "",
						"country": "",
						"subcountry": ""
					},
					"price_aud": "24",
					"price_btc": "24",
					"price_bch": "24",
					"price_eth": "24",
					"price_eur": "24",
					"price_gbp": "24",
					"price_rub": "24",
					"price_usd": "24",
					"hash": "82f1739975",
					"item_name": "",
					"item_uuid": "",
					"premium": false,
					"score": 0,
					"type": "mail",
					"username": ""
				},
				...
			],
			"rating_reviews": [
				{
					"uuid": "034a83c1612e498975d9fba8bfef08db",
					"user_uuid": "0d69e311a31b4ec9483c3117e29eca27",
					"item_uuid": "8bf99758c71c4a6f64cac1d41b6f772b",
					"item_score": 5,
					"item_review": "Arrived. thanks again.",
					"seller_score": 5,
					"seller_review": "Arrived. thanks again.",
					"marketplace_score": 5,
					"marketplace_review": "",
					"created_at": "2018-02-28T14:33:19.07369+03:00"
				},
				...
			],
			"score": 5,
			"premium": false
		}
	}

## Package Reservation

### POST http://127.0.0.1:8081/api/user/:username/item/:item_uuid/package/:hash

Parameters:

* shipping_id
* quantity
* type [bitcoin,bitcoin_cash,ethereum]

Example 

	> curl http://127.0.0.1:8081/api/user/:username/item/:item_uuid/package/:hash --data "shipping_id=:shipping_id&quantity=:quantity&type=:type"

## Payments List Reservation

### GET http://127.0.0.1:8081/api/payments

Example 

	> curl http://127.0.0.1:8081/api/payments

	{
		"transaction_statuses": [
			{
				"amount": 0.3036861982905228,
				"buyer_username": "xxxx",
				"created_at": "2018-02-14T23:12:57.474391+03:00",
				"current_shipping_status": "DISPATCH PENDING",
				"current_status": "FAILED",
				"description": "xxx :: xxx x 1",
				"type": "ethereum",
				"updated_at": "2018-02-15T11:14:23.166919+03:00",
				"uuid": "0xf0971709d0bd2cec450c523d50f4b528802acfda",
				"vendor_username": "drugs-park",
				"updated_at_string": "15.02.2018 11:14",
				"created_at_string": "14.02.2018 23:12"
			},
			...
		]
	}

## Transaction List

### GET http://tochka3evlj3sxdv.onion/api/payments/?token=:token

Example

	> curl http://tochka3evlj3sxdv.onion/api/payments/?token=:token

	{
		"number_of_pages": 1,
		"transaction_statuses": [
			{
				"amount": 0.04958546550835019,
				"buyer_username": "test1230",
				"created_at": "2018-05-25T10:21:19.81422+03:00",
				"current_shipping_status": "DISPATCH PENDING",
				"current_status": "FAILED",
				"description": "â˜…COCAINEâ˜… 90% PURITY UNCUT AAA+ HQ NL/BE/DE SHIPPING :: 0.25 GRAM PREMIUM AAA+ HQ x 1",
				"number_of_messages": 1,
				"type": "ethereum",
				"updated_at": "2018-05-25T22:23:49.558654+03:00",
				"uuid": "0xf7f9b1e0cb45edc2947f7d9a662f278180461557",
				"vendor_username": "dutch-christiana",
				"updated_at_string": "25.05.2018 22:23",
				"created_at_string": "25.05.2018 10:21"
			}
		],
		"api_session": {
			"token": "8f29a2da0ed948c85b4506162473e4d6",
			"end_date": "2018-06-10T10:02:28.244376+03:00",
			"is_2fa_session": false,
			"is_2fa_completed": false
		}
	}

## Transaction Details

### GET http://tochka3evlj3sxdv.onion/api/payments/:tx_id?token=:token

Example

	> curl http://tochka3evlj3sxdv.onion/api/payments/:tx_id?token=:token

	{
		"captcha_id": "CQN5BvdYr2PL5Zf3Tipa",
		"thread": {
			"uuid": "idd",
			"is_read_by_reciever": false,
			"section": "",
			"section_id": 0,
			"parent_uuid": "",
			"title": "Transaction thread @txid",
			"text": "",
			"type": "transaction",
			"is_encrypted": false,
			"has_image": false,
			"created_at": "2018-05-25T10:21:19.966986+03:00",
			"updated_at": "2018-05-25T10:21:19.968172+03:00",
			"deleted_at": null,
			"sender": {
				"username": "test1230",
				"registration_date": "2018-04-15T18:02:46.0321+03:00",
				"last_login_date": "2018-06-13T00:13:11.373106+03:00",
				"bitmessage": "",
				"tox": "",
				"email": "",
				"pgp": "",
				"description": "",
				"long_description": "",
				"is_premium": false,
				"is_premium_plus": false,
				"is_possible_scammer": false,
				"vacation_mode": false,
				"is_seller": false,
				"is_trustedseller": false,
				"is_moderator": false,
				"is_admin": false,
				"is_staff": false
			},
			"text_html": "",
			"messages": [
				{
					"uuid": "idd",
					"is_read_by_reciever": false,
					"section": "",
					"section_id": 0,
					"parent_uuid": "",
					"title": "Transaction thread @txid",
					"text": "",
					"type": "transaction",
					"is_encrypted": false,
					"has_image": false,
					"created_at": "2018-05-25T10:21:19.966986+03:00",
					"updated_at": "2018-05-25T10:21:19.968172+03:00",
					"deleted_at": null,
					"sender": {
						"username": "test1230",
						"registration_date": "2018-04-15T18:02:46.0321+03:00",
						"last_login_date": "2018-06-13T00:13:11.373106+03:00",
						"bitmessage": "",
						"tox": "",
						"email": "",
						"pgp": "",
						"description": "",
						"long_description": "",
						"is_premium": false,
						"is_premium_plus": false,
						"is_possible_scammer": false,
						"vacation_mode": false,
						"is_seller": false,
						"is_trustedseller": false,
						"is_moderator": false,
						"is_admin": false,
						"is_staff": false
					},
					"text_html": ""
				},
				{
					"uuid": "idd-eefaa4e003da45a66668f115a7336cc1",
					"is_read_by_reciever": false,
					"section": "",
					"section_id": 0,
					"parent_uuid": "idd",
					"title": "",
					"text": "fghfgh",
					"type": "transaction",
					"is_encrypted": false,
					"has_image": false,
					"created_at": "2018-05-25T10:21:19.970829+03:00",
					"updated_at": "2018-05-25T10:21:19.974426+03:00",
					"deleted_at": null,
					"sender": {
						"username": "test1230",
						"registration_date": "2018-04-15T18:02:46.0321+03:00",
						"last_login_date": "2018-06-13T00:13:11.373106+03:00",
						"bitmessage": "",
						"tox": "",
						"email": "",
						"pgp": "",
						"description": "",
						"long_description": "",
						"is_premium": false,
						"is_premium_plus": false,
						"is_possible_scammer": false,
						"vacation_mode": false,
						"is_seller": false,
						"is_trustedseller": false,
						"is_moderator": false,
						"is_admin": false,
						"is_staff": false
					},
					"text_html": "\u003cp\u003efghfgh\u003c/p\u003e\n"
				}
			],
			"number_of_messages": 2,
			"number_of_unread_messages": 0,
			"is_read": false
		},
		"transaction": {
			"uuid": "txid",
			"transaction_status": [
				{
					"ID": 153508,
					"CreatedAt": "2018-05-25T10:21:19.982451+03:00",
					"UpdatedAt": "2018-05-25T22:23:49.568762+03:00",
					"DeletedAt": null,
					"time": "2018-05-25T10:21:19.81422+03:00",
					"amount": 0,
					"status": "PENDING",
					"comment": ""
				},
				{
					"ID": 153684,
					"CreatedAt": "2018-05-25T22:23:49.568177+03:00",
					"UpdatedAt": "2018-05-25T22:23:49.568177+03:00",
					"DeletedAt": null,
					"time": "2018-05-25T22:23:49.558654+03:00",
					"amount": 0,
					"status": "FAILED",
					"comment": "Escrow failed automatically"
				}
			],
			"shipping_status": [],
			"description": "desc",
			"type": "ethereum",
			"ethereum_transcation_uuid": "txid",
			"dispute_uuid": "",
			"ethereum_transcation": {
				"uuid": "txid",
				"amount": 0.04958546550835019
			},
			"vendor": {
				"username": "dutch-christiana",
				"registration_date": "2017-08-25T02:30:41.06101+03:00",
				"last_login_date": "2018-06-04T10:29:04.846495+03:00",
				"bitmessage": "",
				"tox": "",
				"email": "",
				"pgp": "fdd",
				"description": "ddd",
				"long_description": "dddd",
				"is_premium": true,
				"is_premium_plus": true,
				"is_possible_scammer": false,
				"vacation_mode": false,
				"is_seller": true,
				"is_trustedseller": true,
				"is_moderator": false,
				"is_admin": false,
				"is_staff": false
			},
			"buyer": {
				"username": "test1230",
				"registration_date": "2018-04-15T18:02:46.0321+03:00",
				"last_login_date": "2018-06-13T00:13:11.373106+03:00",
				"bitmessage": "",
				"tox": "",
				"email": "",
				"pgp": "",
				"description": "",
				"long_description": "",
				"is_premium": false,
				"is_premium_plus": false,
				"is_possible_scammer": false,
				"vacation_mode": false,
				"is_seller": false,
				"is_trustedseller": false,
				"is_moderator": false,
				"is_admin": false,
				"is_staff": false
			},
			"amount": "0.049585",
			"amount_to_pay": "0.049585",
			"current_amount_paid": "0",
			"created_at_string": "25.05.2018 10:21",
			"current_payment_status": "FAILED",
			"current_shipping_status": "DISPATCH PENDING",
			"fe_allowed": true,
			"is_failed": true,
			"number_of_messages": 2,
			"transaction_status_list": [
				{
					"time": "2 weeks ago",
					"amount": 0,
					"status": "PENDING",
					"comment": ""
				},
				{
					"time": "2 weeks ago",
					"amount": 0,
					"status": "FAILED",
					"comment": "Escrow failed automatically"
				}
			]
		},
		"api_session": {
			"token": "8f29a2da0ed948c85b4506162473e4d6",
			"end_date": "2018-06-10T10:02:28.244376+03:00",
			"is_2fa_session": false,
			"is_2fa_completed": false
		}
	}

## Transaction Chat

### POST http://tochka3evlj3sxdv.onion/api/payments/:tx_id?token=:token

Parameters:

* text
* captcha_id
* captcha

Example:

> curl --data "text=tx_chat&captcha_id=9mkx5dXi94QfSTmLbOTI&captcha=7005" http://tochka3evlj3sxdv.onion/api/payments/:tx_id?token=:token

## Wallet List & Balances

## GET http://tochka3evlj3sxdv.onion/api/wallet

Example

	> curl http://tochka3evlj3sxdv.onion/api/wallet

	{
		"btc_balance": {
			"balance": 0.00503501,
			"unconfirmed_balance": 0
		},
		"btc_wallet": {
			"public_key": "xxxx",
			"created_at": "2016-09-28T02:09:48.092677+03:00",
			"updated_at": "2018-07-10T03:17:12.420063+03:00"
		},
		"eth_balance": {
			"balance": 0
		},
		"eth_wallet": {
			"public_key": "xxxx",
			"updated_at": "2018-07-10T03:17:12.483727+03:00"
		},
		"bch_balance": {
			"balance": 0.00114601,
			"unconfirmed_balance": 0
		},
		"bch_wallet": {
			"public_key": "xxxx",
			"created_at": "2016-09-28T02:09:48.092677+03:00",
			"updated_at": "2018-07-10T03:17:12.45335+03:00"
		}
	}

## Bitcoin Withdraw 

## Bitcoin Cash Withdraw

## Ethereum Withdraw