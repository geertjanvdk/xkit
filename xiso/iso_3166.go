// Copyright (c) 2021, Geert JM Vanderkelen

package xiso

import (
	"regexp"
	"strings"
	"time"
)

// iso3166LastUpdate is used in tests to trigger an check to update the data.
// This is the best thing for now, as there is no good way to check this online.
var iso3166LastUpdate = time.Date(2021, 12, 5, 0, 0, 0, 0, time.UTC)

type Country struct {
	Name    string
	Alpha2  string
	Alpha3  string
	Numeric int
}

// ReCountryAlpha2 is a regular expression matching ISO 3166-1 Alpha-2
// country codes. This not checks it validity; merely the format.
var ReCountryAlpha2 = regexp.MustCompile("(?i)^[A-Z]{2}$")

// ReCountryAlpha3 is a regular expression matching ISO 3166-1 Alpha-3
// country codes. This not checks it validity; merely the format.
var ReCountryAlpha3 = regexp.MustCompile("(?i)^[A-Z]{3}$")

// IsEmpty returns whether c is empty. This can be used when looking up
// a country, and country is not available.
func (c Country) IsEmpty() bool {
	return c == (Country{})
}

type countries []Country

// CountryAlpha2 looks up country using its 2-letter code (Alpha-2).
// If not found, the empty country is returned, which can be checked using
// the method Country.IsEmpty. The look-up is case-insensitive.
// The caller is responsible in implementing some kind of caching mechanism.
func CountryAlpha2(code string) Country {
	if len(code) != 2 || !ReCountryAlpha2.MatchString(code) {
		return Country{}
	}
	code = strings.ToUpper(code)
	for _, c := range iso3166Countries {
		if c.Alpha2 == code {
			return c
		}
	}
	return Country{}
}

// CountryAlpha3 looks up country using its 3-letter code (Alpha-3).
// If not found, the empty country is returned, which can be checked using
// the method Country.IsEmpty. The look-up is case-insensitive.
// The caller is responsible in implementing some kind of caching mechanism.
func CountryAlpha3(code string) Country {
	if len(code) != 3 || !ReCountryAlpha3.MatchString(code) {
		return Country{}
	}
	code = strings.ToUpper(code)
	for _, c := range iso3166Countries {
		if c.Alpha3 == code {
			return c
		}
	}
	return Country{}
}

// iso3166Alpha2 holds ISO 3166 country information including the 2 and 3 letter code.
// Source: https://www.iso.org/iso-3166-country-codes.html
// 😉
// ([A-Z]{2})([A-Z]{3})(\d{3})\n -> $1\t$2\t$3\n
// (.*?)\t.*?\t(.*) -> $1\t$2
// ^(.*?)\t([A-Z]{2})\t([A-Z]{3})\t0{0,2}(\d{1,3})$ -> {\nName: "$1",\nAlpha2:"$2",\nAlpha3:"$3",\nNumeric:$4,\n},
var iso3166Countries = countries{
	{
		Name:    "Afghanistan",
		Alpha2:  "AF",
		Alpha3:  "AFG",
		Numeric: 4,
	},
	{
		Name:    "Albania",
		Alpha2:  "AL",
		Alpha3:  "ALB",
		Numeric: 8,
	},
	{
		Name:    "Algeria",
		Alpha2:  "DZ",
		Alpha3:  "DZA",
		Numeric: 12,
	},
	{
		Name:    "American Samoa",
		Alpha2:  "AS",
		Alpha3:  "ASM",
		Numeric: 16,
	},
	{
		Name:    "Andorra",
		Alpha2:  "AD",
		Alpha3:  "AND",
		Numeric: 20,
	},
	{
		Name:    "Angola",
		Alpha2:  "AO",
		Alpha3:  "AGO",
		Numeric: 24,
	},
	{
		Name:    "Anguilla",
		Alpha2:  "AI",
		Alpha3:  "AIA",
		Numeric: 660,
	},
	{
		Name:    "Antarctica",
		Alpha2:  "AQ",
		Alpha3:  "ATA",
		Numeric: 10,
	},
	{
		Name:    "Antigua and Barbuda",
		Alpha2:  "AG",
		Alpha3:  "ATG",
		Numeric: 28,
	},
	{
		Name:    "Argentina",
		Alpha2:  "AR",
		Alpha3:  "ARG",
		Numeric: 32,
	},
	{
		Name:    "Armenia",
		Alpha2:  "AM",
		Alpha3:  "ARM",
		Numeric: 51,
	},
	{
		Name:    "Aruba",
		Alpha2:  "AW",
		Alpha3:  "ABW",
		Numeric: 533,
	},
	{
		Name:    "Australia",
		Alpha2:  "AU",
		Alpha3:  "AUS",
		Numeric: 36,
	},
	{
		Name:    "Austria",
		Alpha2:  "AT",
		Alpha3:  "AUT",
		Numeric: 40,
	},
	{
		Name:    "Azerbaijan",
		Alpha2:  "AZ",
		Alpha3:  "AZE",
		Numeric: 31,
	},
	{
		Name:    "Bahamas (the)",
		Alpha2:  "BS",
		Alpha3:  "BHS",
		Numeric: 44,
	},
	{
		Name:    "Bahrain",
		Alpha2:  "BH",
		Alpha3:  "BHR",
		Numeric: 48,
	},
	{
		Name:    "Bangladesh",
		Alpha2:  "BD",
		Alpha3:  "BGD",
		Numeric: 50,
	},
	{
		Name:    "Barbados",
		Alpha2:  "BB",
		Alpha3:  "BRB",
		Numeric: 52,
	},
	{
		Name:    "Belarus",
		Alpha2:  "BY",
		Alpha3:  "BLR",
		Numeric: 112,
	},
	{
		Name:    "Belgium",
		Alpha2:  "BE",
		Alpha3:  "BEL",
		Numeric: 56,
	},
	{
		Name:    "Belize",
		Alpha2:  "BZ",
		Alpha3:  "BLZ",
		Numeric: 84,
	},
	{
		Name:    "Benin",
		Alpha2:  "BJ",
		Alpha3:  "BEN",
		Numeric: 204,
	},
	{
		Name:    "Bermuda",
		Alpha2:  "BM",
		Alpha3:  "BMU",
		Numeric: 60,
	},
	{
		Name:    "Bhutan",
		Alpha2:  "BT",
		Alpha3:  "BTN",
		Numeric: 64,
	},
	{
		Name:    "Bolivia (Plurinational State of)",
		Alpha2:  "BO",
		Alpha3:  "BOL",
		Numeric: 68,
	},
	{
		Name:    "Bonaire, Sint Eustatius and Saba",
		Alpha2:  "BQ",
		Alpha3:  "BES",
		Numeric: 535,
	},
	{
		Name:    "Bosnia and Herzegovina",
		Alpha2:  "BA",
		Alpha3:  "BIH",
		Numeric: 70,
	},
	{
		Name:    "Botswana",
		Alpha2:  "BW",
		Alpha3:  "BWA",
		Numeric: 72,
	},
	{
		Name:    "Bouvet Island",
		Alpha2:  "BV",
		Alpha3:  "BVT",
		Numeric: 74,
	},
	{
		Name:    "Brazil",
		Alpha2:  "BR",
		Alpha3:  "BRA",
		Numeric: 76,
	},
	{
		Name:    "British Indian Ocean Territory (the)",
		Alpha2:  "IO",
		Alpha3:  "IOT",
		Numeric: 86,
	},
	{
		Name:    "Brunei Darussalam",
		Alpha2:  "BN",
		Alpha3:  "BRN",
		Numeric: 96,
	},
	{
		Name:    "Bulgaria",
		Alpha2:  "BG",
		Alpha3:  "BGR",
		Numeric: 100,
	},
	{
		Name:    "Burkina Faso",
		Alpha2:  "BF",
		Alpha3:  "BFA",
		Numeric: 854,
	},
	{
		Name:    "Burundi",
		Alpha2:  "BI",
		Alpha3:  "BDI",
		Numeric: 108,
	},
	{
		Name:    "Cabo Verde",
		Alpha2:  "CV",
		Alpha3:  "CPV",
		Numeric: 132,
	},
	{
		Name:    "Cambodia",
		Alpha2:  "KH",
		Alpha3:  "KHM",
		Numeric: 116,
	},
	{
		Name:    "Cameroon",
		Alpha2:  "CM",
		Alpha3:  "CMR",
		Numeric: 120,
	},
	{
		Name:    "Canada",
		Alpha2:  "CA",
		Alpha3:  "CAN",
		Numeric: 124,
	},
	{
		Name:    "Cayman Islands (the)",
		Alpha2:  "KY",
		Alpha3:  "CYM",
		Numeric: 136,
	},
	{
		Name:    "Central African Republic (the)",
		Alpha2:  "CF",
		Alpha3:  "CAF",
		Numeric: 140,
	},
	{
		Name:    "Chad",
		Alpha2:  "TD",
		Alpha3:  "TCD",
		Numeric: 148,
	},
	{
		Name:    "Chile",
		Alpha2:  "CL",
		Alpha3:  "CHL",
		Numeric: 152,
	},
	{
		Name:    "China",
		Alpha2:  "CN",
		Alpha3:  "CHN",
		Numeric: 156,
	},
	{
		Name:    "Christmas Island",
		Alpha2:  "CX",
		Alpha3:  "CXR",
		Numeric: 162,
	},
	{
		Name:    "Cocos (Keeling) Islands (the)",
		Alpha2:  "CC",
		Alpha3:  "CCK",
		Numeric: 166,
	},
	{
		Name:    "Colombia",
		Alpha2:  "CO",
		Alpha3:  "COL",
		Numeric: 170,
	},
	{
		Name:    "Comoros (the)",
		Alpha2:  "KM",
		Alpha3:  "COM",
		Numeric: 174,
	},
	{
		Name:    "Congo (the Democratic Republic of the)",
		Alpha2:  "CD",
		Alpha3:  "COD",
		Numeric: 180,
	},
	{
		Name:    "Congo (the)",
		Alpha2:  "CG",
		Alpha3:  "COG",
		Numeric: 178,
	},
	{
		Name:    "Cook Islands (the)",
		Alpha2:  "CK",
		Alpha3:  "COK",
		Numeric: 184,
	},
	{
		Name:    "Costa Rica",
		Alpha2:  "CR",
		Alpha3:  "CRI",
		Numeric: 188,
	},
	{
		Name:    "Croatia",
		Alpha2:  "HR",
		Alpha3:  "HRV",
		Numeric: 191,
	},
	{
		Name:    "Cuba",
		Alpha2:  "CU",
		Alpha3:  "CUB",
		Numeric: 192,
	},
	{
		Name:    "Curaçao",
		Alpha2:  "CW",
		Alpha3:  "CUW",
		Numeric: 531,
	},
	{
		Name:    "Cyprus",
		Alpha2:  "CY",
		Alpha3:  "CYP",
		Numeric: 196,
	},
	{
		Name:    "Czechia",
		Alpha2:  "CZ",
		Alpha3:  "CZE",
		Numeric: 203,
	},
	{
		Name:    "Côte d'Ivoire",
		Alpha2:  "CI",
		Alpha3:  "CIV",
		Numeric: 384,
	},
	{
		Name:    "Denmark",
		Alpha2:  "DK",
		Alpha3:  "DNK",
		Numeric: 208,
	},
	{
		Name:    "Djibouti",
		Alpha2:  "DJ",
		Alpha3:  "DJI",
		Numeric: 262,
	},
	{
		Name:    "Dominica",
		Alpha2:  "DM",
		Alpha3:  "DMA",
		Numeric: 212,
	},
	{
		Name:    "Dominican Republic (the)",
		Alpha2:  "DO",
		Alpha3:  "DOM",
		Numeric: 214,
	},
	{
		Name:    "Ecuador",
		Alpha2:  "EC",
		Alpha3:  "ECU",
		Numeric: 218,
	},
	{
		Name:    "Egypt",
		Alpha2:  "EG",
		Alpha3:  "EGY",
		Numeric: 818,
	},
	{
		Name:    "El Salvador",
		Alpha2:  "SV",
		Alpha3:  "SLV",
		Numeric: 222,
	},
	{
		Name:    "Equatorial Guinea",
		Alpha2:  "GQ",
		Alpha3:  "GNQ",
		Numeric: 226,
	},
	{
		Name:    "Eritrea",
		Alpha2:  "ER",
		Alpha3:  "ERI",
		Numeric: 232,
	},
	{
		Name:    "Estonia",
		Alpha2:  "EE",
		Alpha3:  "EST",
		Numeric: 233,
	},
	{
		Name:    "Eswatini",
		Alpha2:  "SZ",
		Alpha3:  "SWZ",
		Numeric: 748,
	},
	{
		Name:    "Ethiopia",
		Alpha2:  "ET",
		Alpha3:  "ETH",
		Numeric: 231,
	},
	{
		Name:    "Falkland Islands (the) [Malvinas]",
		Alpha2:  "FK",
		Alpha3:  "FLK",
		Numeric: 238,
	},
	{
		Name:    "Faroe Islands (the)",
		Alpha2:  "FO",
		Alpha3:  "FRO",
		Numeric: 234,
	},
	{
		Name:    "Fiji",
		Alpha2:  "FJ",
		Alpha3:  "FJI",
		Numeric: 242,
	},
	{
		Name:    "Finland",
		Alpha2:  "FI",
		Alpha3:  "FIN",
		Numeric: 246,
	},
	{
		Name:    "France",
		Alpha2:  "FR",
		Alpha3:  "FRA",
		Numeric: 250,
	},
	{
		Name:    "French Guiana",
		Alpha2:  "GF",
		Alpha3:  "GUF",
		Numeric: 254,
	},
	{
		Name:    "French Polynesia",
		Alpha2:  "PF",
		Alpha3:  "PYF",
		Numeric: 258,
	},
	{
		Name:    "French Southern Territories (the)",
		Alpha2:  "TF",
		Alpha3:  "ATF",
		Numeric: 260,
	},
	{
		Name:    "Gabon",
		Alpha2:  "GA",
		Alpha3:  "GAB",
		Numeric: 266,
	},
	{
		Name:    "Gambia (the)",
		Alpha2:  "GM",
		Alpha3:  "GMB",
		Numeric: 270,
	},
	{
		Name:    "Georgia",
		Alpha2:  "GE",
		Alpha3:  "GEO",
		Numeric: 268,
	},
	{
		Name:    "Germany",
		Alpha2:  "DE",
		Alpha3:  "DEU",
		Numeric: 276,
	},
	{
		Name:    "Ghana",
		Alpha2:  "GH",
		Alpha3:  "GHA",
		Numeric: 288,
	},
	{
		Name:    "Gibraltar",
		Alpha2:  "GI",
		Alpha3:  "GIB",
		Numeric: 292,
	},
	{
		Name:    "Greece",
		Alpha2:  "GR",
		Alpha3:  "GRC",
		Numeric: 300,
	},
	{
		Name:    "Greenland",
		Alpha2:  "GL",
		Alpha3:  "GRL",
		Numeric: 304,
	},
	{
		Name:    "Grenada",
		Alpha2:  "GD",
		Alpha3:  "GRD",
		Numeric: 308,
	},
	{
		Name:    "Guadeloupe",
		Alpha2:  "GP",
		Alpha3:  "GLP",
		Numeric: 312,
	},
	{
		Name:    "Guam",
		Alpha2:  "GU",
		Alpha3:  "GUM",
		Numeric: 316,
	},
	{
		Name:    "Guatemala",
		Alpha2:  "GT",
		Alpha3:  "GTM",
		Numeric: 320,
	},
	{
		Name:    "Guernsey",
		Alpha2:  "GG",
		Alpha3:  "GGY",
		Numeric: 831,
	},
	{
		Name:    "Guinea",
		Alpha2:  "GN",
		Alpha3:  "GIN",
		Numeric: 324,
	},
	{
		Name:    "Guinea-Bissau",
		Alpha2:  "GW",
		Alpha3:  "GNB",
		Numeric: 624,
	},
	{
		Name:    "Guyana",
		Alpha2:  "GY",
		Alpha3:  "GUY",
		Numeric: 328,
	},
	{
		Name:    "Haiti",
		Alpha2:  "HT",
		Alpha3:  "HTI",
		Numeric: 332,
	},
	{
		Name:    "Heard Island and McDonald Islands",
		Alpha2:  "HM",
		Alpha3:  "HMD",
		Numeric: 334,
	},
	{
		Name:    "Holy See (the)",
		Alpha2:  "VA",
		Alpha3:  "VAT",
		Numeric: 336,
	},
	{
		Name:    "Honduras",
		Alpha2:  "HN",
		Alpha3:  "HND",
		Numeric: 340,
	},
	{
		Name:    "Hong Kong",
		Alpha2:  "HK",
		Alpha3:  "HKG",
		Numeric: 344,
	},
	{
		Name:    "Hungary",
		Alpha2:  "HU",
		Alpha3:  "HUN",
		Numeric: 348,
	},
	{
		Name:    "Iceland",
		Alpha2:  "IS",
		Alpha3:  "ISL",
		Numeric: 352,
	},
	{
		Name:    "India",
		Alpha2:  "IN",
		Alpha3:  "IND",
		Numeric: 356,
	},
	{
		Name:    "Indonesia",
		Alpha2:  "ID",
		Alpha3:  "IDN",
		Numeric: 360,
	},
	{
		Name:    "Iran (Islamic Republic of)",
		Alpha2:  "IR",
		Alpha3:  "IRN",
		Numeric: 364,
	},
	{
		Name:    "Iraq",
		Alpha2:  "IQ",
		Alpha3:  "IRQ",
		Numeric: 368,
	},
	{
		Name:    "Ireland",
		Alpha2:  "IE",
		Alpha3:  "IRL",
		Numeric: 372,
	},
	{
		Name:    "Isle of Man",
		Alpha2:  "IM",
		Alpha3:  "IMN",
		Numeric: 833,
	},
	{
		Name:    "Israel",
		Alpha2:  "IL",
		Alpha3:  "ISR",
		Numeric: 376,
	},
	{
		Name:    "Italy",
		Alpha2:  "IT",
		Alpha3:  "ITA",
		Numeric: 380,
	},
	{
		Name:    "Jamaica",
		Alpha2:  "JM",
		Alpha3:  "JAM",
		Numeric: 388,
	},
	{
		Name:    "Japan",
		Alpha2:  "JP",
		Alpha3:  "JPN",
		Numeric: 392,
	},
	{
		Name:    "Jersey",
		Alpha2:  "JE",
		Alpha3:  "JEY",
		Numeric: 832,
	},
	{
		Name:    "Jordan",
		Alpha2:  "JO",
		Alpha3:  "JOR",
		Numeric: 400,
	},
	{
		Name:    "Kazakhstan",
		Alpha2:  "KZ",
		Alpha3:  "KAZ",
		Numeric: 398,
	},
	{
		Name:    "Kenya",
		Alpha2:  "KE",
		Alpha3:  "KEN",
		Numeric: 404,
	},
	{
		Name:    "Kiribati",
		Alpha2:  "KI",
		Alpha3:  "KIR",
		Numeric: 296,
	},
	{
		Name:    "Korea (the Democratic People's Republic of)",
		Alpha2:  "KP",
		Alpha3:  "PRK",
		Numeric: 408,
	},
	{
		Name:    "Korea (the Republic of)",
		Alpha2:  "KR",
		Alpha3:  "KOR",
		Numeric: 410,
	},
	{
		Name:    "Kuwait",
		Alpha2:  "KW",
		Alpha3:  "KWT",
		Numeric: 414,
	},
	{
		Name:    "Kyrgyzstan",
		Alpha2:  "KG",
		Alpha3:  "KGZ",
		Numeric: 417,
	},
	{
		Name:    "Lao People's Democratic Republic (the)",
		Alpha2:  "LA",
		Alpha3:  "LAO",
		Numeric: 418,
	},
	{
		Name:    "Latvia",
		Alpha2:  "LV",
		Alpha3:  "LVA",
		Numeric: 428,
	},
	{
		Name:    "Lebanon",
		Alpha2:  "LB",
		Alpha3:  "LBN",
		Numeric: 422,
	},
	{
		Name:    "Lesotho",
		Alpha2:  "LS",
		Alpha3:  "LSO",
		Numeric: 426,
	},
	{
		Name:    "Liberia",
		Alpha2:  "LR",
		Alpha3:  "LBR",
		Numeric: 430,
	},
	{
		Name:    "Libya",
		Alpha2:  "LY",
		Alpha3:  "LBY",
		Numeric: 434,
	},
	{
		Name:    "Liechtenstein",
		Alpha2:  "LI",
		Alpha3:  "LIE",
		Numeric: 438,
	},
	{
		Name:    "Lithuania",
		Alpha2:  "LT",
		Alpha3:  "LTU",
		Numeric: 440,
	},
	{
		Name:    "Luxembourg",
		Alpha2:  "LU",
		Alpha3:  "LUX",
		Numeric: 442,
	},
	{
		Name:    "Macao",
		Alpha2:  "MO",
		Alpha3:  "MAC",
		Numeric: 446,
	},
	{
		Name:    "Madagascar",
		Alpha2:  "MG",
		Alpha3:  "MDG",
		Numeric: 450,
	},
	{
		Name:    "Malawi",
		Alpha2:  "MW",
		Alpha3:  "MWI",
		Numeric: 454,
	},
	{
		Name:    "Malaysia",
		Alpha2:  "MY",
		Alpha3:  "MYS",
		Numeric: 458,
	},
	{
		Name:    "Maldives",
		Alpha2:  "MV",
		Alpha3:  "MDV",
		Numeric: 462,
	},
	{
		Name:    "Mali",
		Alpha2:  "ML",
		Alpha3:  "MLI",
		Numeric: 466,
	},
	{
		Name:    "Malta",
		Alpha2:  "MT",
		Alpha3:  "MLT",
		Numeric: 470,
	},
	{
		Name:    "Marshall Islands (the)",
		Alpha2:  "MH",
		Alpha3:  "MHL",
		Numeric: 584,
	},
	{
		Name:    "Martinique",
		Alpha2:  "MQ",
		Alpha3:  "MTQ",
		Numeric: 474,
	},
	{
		Name:    "Mauritania",
		Alpha2:  "MR",
		Alpha3:  "MRT",
		Numeric: 478,
	},
	{
		Name:    "Mauritius",
		Alpha2:  "MU",
		Alpha3:  "MUS",
		Numeric: 480,
	},
	{
		Name:    "Mayotte",
		Alpha2:  "YT",
		Alpha3:  "MYT",
		Numeric: 175,
	},
	{
		Name:    "Mexico",
		Alpha2:  "MX",
		Alpha3:  "MEX",
		Numeric: 484,
	},
	{
		Name:    "Micronesia (Federated States of)",
		Alpha2:  "FM",
		Alpha3:  "FSM",
		Numeric: 583,
	},
	{
		Name:    "Moldova (the Republic of)",
		Alpha2:  "MD",
		Alpha3:  "MDA",
		Numeric: 498,
	},
	{
		Name:    "Monaco",
		Alpha2:  "MC",
		Alpha3:  "MCO",
		Numeric: 492,
	},
	{
		Name:    "Mongolia",
		Alpha2:  "MN",
		Alpha3:  "MNG",
		Numeric: 496,
	},
	{
		Name:    "Montenegro",
		Alpha2:  "ME",
		Alpha3:  "MNE",
		Numeric: 499,
	},
	{
		Name:    "Montserrat",
		Alpha2:  "MS",
		Alpha3:  "MSR",
		Numeric: 500,
	},
	{
		Name:    "Morocco",
		Alpha2:  "MA",
		Alpha3:  "MAR",
		Numeric: 504,
	},
	{
		Name:    "Mozambique",
		Alpha2:  "MZ",
		Alpha3:  "MOZ",
		Numeric: 508,
	},
	{
		Name:    "Myanmar",
		Alpha2:  "MM",
		Alpha3:  "MMR",
		Numeric: 104,
	},
	{
		Name:    "Namibia",
		Alpha2:  "NA",
		Alpha3:  "NAM",
		Numeric: 516,
	},
	{
		Name:    "Nauru",
		Alpha2:  "NR",
		Alpha3:  "NRU",
		Numeric: 520,
	},
	{
		Name:    "Nepal",
		Alpha2:  "NP",
		Alpha3:  "NPL",
		Numeric: 524,
	},
	{
		Name:    "Netherlands (the)",
		Alpha2:  "NL",
		Alpha3:  "NLD",
		Numeric: 528,
	},
	{
		Name:    "New Caledonia",
		Alpha2:  "NC",
		Alpha3:  "NCL",
		Numeric: 540,
	},
	{
		Name:    "New Zealand",
		Alpha2:  "NZ",
		Alpha3:  "NZL",
		Numeric: 554,
	},
	{
		Name:    "Nicaragua",
		Alpha2:  "NI",
		Alpha3:  "NIC",
		Numeric: 558,
	},
	{
		Name:    "Niger (the)",
		Alpha2:  "NE",
		Alpha3:  "NER",
		Numeric: 562,
	},
	{
		Name:    "Nigeria",
		Alpha2:  "NG",
		Alpha3:  "NGA",
		Numeric: 566,
	},
	{
		Name:    "Niue",
		Alpha2:  "NU",
		Alpha3:  "NIU",
		Numeric: 570,
	},
	{
		Name:    "Norfolk Island",
		Alpha2:  "NF",
		Alpha3:  "NFK",
		Numeric: 574,
	},
	{
		Name:    "North Macedonia",
		Alpha2:  "MK",
		Alpha3:  "MKD",
		Numeric: 807,
	},
	{
		Name:    "Northern Mariana Islands (the)",
		Alpha2:  "MP",
		Alpha3:  "MNP",
		Numeric: 580,
	},
	{
		Name:    "Norway",
		Alpha2:  "NO",
		Alpha3:  "NOR",
		Numeric: 578,
	},
	{
		Name:    "Oman",
		Alpha2:  "OM",
		Alpha3:  "OMN",
		Numeric: 512,
	},
	{
		Name:    "Pakistan",
		Alpha2:  "PK",
		Alpha3:  "PAK",
		Numeric: 586,
	},
	{
		Name:    "Palau",
		Alpha2:  "PW",
		Alpha3:  "PLW",
		Numeric: 585,
	},
	{
		Name:    "Palestine, State of",
		Alpha2:  "PS",
		Alpha3:  "PSE",
		Numeric: 275,
	},
	{
		Name:    "Panama",
		Alpha2:  "PA",
		Alpha3:  "PAN",
		Numeric: 591,
	},
	{
		Name:    "Papua New Guinea",
		Alpha2:  "PG",
		Alpha3:  "PNG",
		Numeric: 598,
	},
	{
		Name:    "Paraguay",
		Alpha2:  "PY",
		Alpha3:  "PRY",
		Numeric: 600,
	},
	{
		Name:    "Peru",
		Alpha2:  "PE",
		Alpha3:  "PER",
		Numeric: 604,
	},
	{
		Name:    "Philippines (the)",
		Alpha2:  "PH",
		Alpha3:  "PHL",
		Numeric: 608,
	},
	{
		Name:    "Pitcairn",
		Alpha2:  "PN",
		Alpha3:  "PCN",
		Numeric: 612,
	},
	{
		Name:    "Poland",
		Alpha2:  "PL",
		Alpha3:  "POL",
		Numeric: 616,
	},
	{
		Name:    "Portugal",
		Alpha2:  "PT",
		Alpha3:  "PRT",
		Numeric: 620,
	},
	{
		Name:    "Puerto Rico",
		Alpha2:  "PR",
		Alpha3:  "PRI",
		Numeric: 630,
	},
	{
		Name:    "Qatar",
		Alpha2:  "QA",
		Alpha3:  "QAT",
		Numeric: 634,
	},
	{
		Name:    "Romania",
		Alpha2:  "RO",
		Alpha3:  "ROU",
		Numeric: 642,
	},
	{
		Name:    "Russian Federation (the)",
		Alpha2:  "RU",
		Alpha3:  "RUS",
		Numeric: 643,
	},
	{
		Name:    "Rwanda",
		Alpha2:  "RW",
		Alpha3:  "RWA",
		Numeric: 646,
	},
	{
		Name:    "Réunion",
		Alpha2:  "RE",
		Alpha3:  "REU",
		Numeric: 638,
	},
	{
		Name:    "Saint Barthélemy",
		Alpha2:  "BL",
		Alpha3:  "BLM",
		Numeric: 652,
	},
	{
		Name:    "Saint Helena, Ascension and Tristan da Cunha",
		Alpha2:  "SH",
		Alpha3:  "SHN",
		Numeric: 654,
	},
	{
		Name:    "Saint Kitts and Nevis",
		Alpha2:  "KN",
		Alpha3:  "KNA",
		Numeric: 659,
	},
	{
		Name:    "Saint Lucia",
		Alpha2:  "LC",
		Alpha3:  "LCA",
		Numeric: 662,
	},
	{
		Name:    "Saint Martin (French part)",
		Alpha2:  "MF",
		Alpha3:  "MAF",
		Numeric: 663,
	},
	{
		Name:    "Saint Pierre and Miquelon",
		Alpha2:  "PM",
		Alpha3:  "SPM",
		Numeric: 666,
	},
	{
		Name:    "Saint Vincent and the Grenadines",
		Alpha2:  "VC",
		Alpha3:  "VCT",
		Numeric: 670,
	},
	{
		Name:    "Samoa",
		Alpha2:  "WS",
		Alpha3:  "WSM",
		Numeric: 882,
	},
	{
		Name:    "San Marino",
		Alpha2:  "SM",
		Alpha3:  "SMR",
		Numeric: 674,
	},
	{
		Name:    "Sao Tome and Principe",
		Alpha2:  "ST",
		Alpha3:  "STP",
		Numeric: 678,
	},
	{
		Name:    "Saudi Arabia",
		Alpha2:  "SA",
		Alpha3:  "SAU",
		Numeric: 682,
	},
	{
		Name:    "Senegal",
		Alpha2:  "SN",
		Alpha3:  "SEN",
		Numeric: 686,
	},
	{
		Name:    "Serbia",
		Alpha2:  "RS",
		Alpha3:  "SRB",
		Numeric: 688,
	},
	{
		Name:    "Seychelles",
		Alpha2:  "SC",
		Alpha3:  "SYC",
		Numeric: 690,
	},
	{
		Name:    "Sierra Leone",
		Alpha2:  "SL",
		Alpha3:  "SLE",
		Numeric: 694,
	},
	{
		Name:    "Singapore",
		Alpha2:  "SG",
		Alpha3:  "SGP",
		Numeric: 702,
	},
	{
		Name:    "Sint Maarten (Dutch part)",
		Alpha2:  "SX",
		Alpha3:  "SXM",
		Numeric: 534,
	},
	{
		Name:    "Slovakia",
		Alpha2:  "SK",
		Alpha3:  "SVK",
		Numeric: 703,
	},
	{
		Name:    "Slovenia",
		Alpha2:  "SI",
		Alpha3:  "SVN",
		Numeric: 705,
	},
	{
		Name:    "Solomon Islands",
		Alpha2:  "SB",
		Alpha3:  "SLB",
		Numeric: 90,
	},
	{
		Name:    "Somalia",
		Alpha2:  "SO",
		Alpha3:  "SOM",
		Numeric: 706,
	},
	{
		Name:    "South Africa",
		Alpha2:  "ZA",
		Alpha3:  "ZAF",
		Numeric: 710,
	},
	{
		Name:    "South Georgia and the South Sandwich Islands",
		Alpha2:  "GS",
		Alpha3:  "SGS",
		Numeric: 239,
	},
	{
		Name:    "South Sudan",
		Alpha2:  "SS",
		Alpha3:  "SSD",
		Numeric: 728,
	},
	{
		Name:    "Spain",
		Alpha2:  "ES",
		Alpha3:  "ESP",
		Numeric: 724,
	},
	{
		Name:    "Sri Lanka",
		Alpha2:  "LK",
		Alpha3:  "LKA",
		Numeric: 144,
	},
	{
		Name:    "Sudan (the)",
		Alpha2:  "SD",
		Alpha3:  "SDN",
		Numeric: 729,
	},
	{
		Name:    "Suriname",
		Alpha2:  "SR",
		Alpha3:  "SUR",
		Numeric: 740,
	},
	{
		Name:    "Svalbard and Jan Mayen",
		Alpha2:  "SJ",
		Alpha3:  "SJM",
		Numeric: 744,
	},
	{
		Name:    "Sweden",
		Alpha2:  "SE",
		Alpha3:  "SWE",
		Numeric: 752,
	},
	{
		Name:    "Switzerland",
		Alpha2:  "CH",
		Alpha3:  "CHE",
		Numeric: 756,
	},
	{
		Name:    "Syrian Arab Republic (the)",
		Alpha2:  "SY",
		Alpha3:  "SYR",
		Numeric: 760,
	},
	{
		Name:    "Taiwan (Province of China)",
		Alpha2:  "TW",
		Alpha3:  "TWN",
		Numeric: 158,
	},
	{
		Name:    "Tajikistan",
		Alpha2:  "TJ",
		Alpha3:  "TJK",
		Numeric: 762,
	},
	{
		Name:    "Tanzania, the United Republic of",
		Alpha2:  "TZ",
		Alpha3:  "TZA",
		Numeric: 834,
	},
	{
		Name:    "Thailand",
		Alpha2:  "TH",
		Alpha3:  "THA",
		Numeric: 764,
	},
	{
		Name:    "Timor-Leste",
		Alpha2:  "TL",
		Alpha3:  "TLS",
		Numeric: 626,
	},
	{
		Name:    "Togo",
		Alpha2:  "TG",
		Alpha3:  "TGO",
		Numeric: 768,
	},
	{
		Name:    "Tokelau",
		Alpha2:  "TK",
		Alpha3:  "TKL",
		Numeric: 772,
	},
	{
		Name:    "Tonga",
		Alpha2:  "TO",
		Alpha3:  "TON",
		Numeric: 776,
	},
	{
		Name:    "Trinidad and Tobago",
		Alpha2:  "TT",
		Alpha3:  "TTO",
		Numeric: 780,
	},
	{
		Name:    "Tunisia",
		Alpha2:  "TN",
		Alpha3:  "TUN",
		Numeric: 788,
	},
	{
		Name:    "Turkey",
		Alpha2:  "TR",
		Alpha3:  "TUR",
		Numeric: 792,
	},
	{
		Name:    "Turkmenistan",
		Alpha2:  "TM",
		Alpha3:  "TKM",
		Numeric: 795,
	},
	{
		Name:    "Turks and Caicos Islands (the)",
		Alpha2:  "TC",
		Alpha3:  "TCA",
		Numeric: 796,
	},
	{
		Name:    "Tuvalu",
		Alpha2:  "TV",
		Alpha3:  "TUV",
		Numeric: 798,
	},
	{
		Name:    "Uganda",
		Alpha2:  "UG",
		Alpha3:  "UGA",
		Numeric: 800,
	},
	{
		Name:    "Ukraine",
		Alpha2:  "UA",
		Alpha3:  "UKR",
		Numeric: 804,
	},
	{
		Name:    "United Arab Emirates (the)",
		Alpha2:  "AE",
		Alpha3:  "ARE",
		Numeric: 784,
	},
	{
		Name:    "United Kingdom of Great Britain and Northern Ireland (the)",
		Alpha2:  "GB",
		Alpha3:  "GBR",
		Numeric: 826,
	},
	{
		Name:    "United States Minor Outlying Islands (the)",
		Alpha2:  "UM",
		Alpha3:  "UMI",
		Numeric: 581,
	},
	{
		Name:    "United States of America (the)",
		Alpha2:  "US",
		Alpha3:  "USA",
		Numeric: 840,
	},
	{
		Name:    "Uruguay",
		Alpha2:  "UY",
		Alpha3:  "URY",
		Numeric: 858,
	},
	{
		Name:    "Uzbekistan",
		Alpha2:  "UZ",
		Alpha3:  "UZB",
		Numeric: 860,
	},
	{
		Name:    "Vanuatu",
		Alpha2:  "VU",
		Alpha3:  "VUT",
		Numeric: 548,
	},
	{
		Name:    "Venezuela (Bolivarian Republic of)",
		Alpha2:  "VE",
		Alpha3:  "VEN",
		Numeric: 862,
	},
	{
		Name:    "Viet Nam",
		Alpha2:  "VN",
		Alpha3:  "VNM",
		Numeric: 704,
	},
	{
		Name:    "Virgin Islands (British)",
		Alpha2:  "VG",
		Alpha3:  "VGB",
		Numeric: 92,
	},
	{
		Name:    "Virgin Islands (U.S.)",
		Alpha2:  "VI",
		Alpha3:  "VIR",
		Numeric: 850,
	},
	{
		Name:    "Wallis and Futuna",
		Alpha2:  "WF",
		Alpha3:  "WLF",
		Numeric: 876,
	},
	{
		Name:    "Western Sahara",
		Alpha2:  "EH",
		Alpha3:  "ESH",
		Numeric: 732,
	},
	{
		Name:    "Yemen",
		Alpha2:  "YE",
		Alpha3:  "YEM",
		Numeric: 887,
	},
	{
		Name:    "Zambia",
		Alpha2:  "ZM",
		Alpha3:  "ZMB",
		Numeric: 894,
	},
	{
		Name:    "Zimbabwe",
		Alpha2:  "ZW",
		Alpha3:  "ZWE",
		Numeric: 716,
	},
	{
		Name:    "Åland Islands",
		Alpha2:  "AX",
		Alpha3:  "ALA",
		Numeric: 248,
	},
}
