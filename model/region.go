package model

import "gorm.io/gorm"

type Country struct {
	gorm.Model
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Latitude       string     `json:"latitude"`
	Longitude      string     `json:"longitude"`
	ISO3           string     `json:"iso3"`
	ISO2           string     `json:"iso2"`
	NumericCode    string     `json:"numeric_code"`
	PhoneCode      string     `json:"phone_code"`
	Capital        string     `json:"capital"`
	Currency       string     `json:"currency"`
	CurrencyName   string     `json:"currency_name"`
	CurrencySymbol string     `json:"currency_symbol"`
	TLD            string     `json:"tld"`
	Native         string     `json:"native"`
	Region         string     `json:"region"`
	Subregion      string     `json:"subregion"`
	Timezones      []Timezone `json:"timezones" gorm:"many2many:country_timezones;"`
	Emoji          string     `json:"emoji"`
	EmojiUnicode   string     `json:"emojiU"`
	States         []State    `json:"states" gorm:"foreignKey:CountryID"`
}

type Timezone struct {
	gorm.Model
	ZoneName      string `json:"zoneName"`
	GMTOffset     int    `json:"gmtOffset"`
	GMTOffsetName string `json:"gmtOffsetName"`
	Abbreviation  string `json:"abbreviation"`
	TimezoneName  string `json:"tzName"`
}

type State struct {
	gorm.Model
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	CountryID uint
	StateCode string `json:"state_code"`
	Cities    []City `json:"cities" gorm:"foreignKey:StateID"`
}

type City struct {
	gorm.Model
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	StateID   uint
}
