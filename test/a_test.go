package test

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestSetFBProperties(t *testing.T) {
	jsonStr := `
	{
		"$url": "https://example.com",
		"total": 99.9,
		"key_word": "keyword123",
		"product_title": "Test Product",
		"title": "Product Title",
		"quantity": 2,
		"product_id": "prod_001",
		"order_id": "order_001",
		"client_id": 1750905219951,
		"price": 19.95,
		"platform": "web",
		"user_agent": "Mozilla/5.0",
		"currency": "USD",
		"AD_variant_ids": ["123", "456"],
		"content_ids": ["abc", "def"],
		"template_name": "theme-x",
		"theme_version": "1.0",
		"client": "shoplaza",
		"AD_fbp": "fb.1.123456",
		"AD_fbc": "fb.1.654321",
		"event_time": 1750905219951,
		"AD_event_id": "event_123",
		"AD_em": "email@example.com",
		"AD_ph": "1234567890",
		"AD_ln": "Doe",
		"AD_fn": "John",
		"AD_ct": "City",
		"AD_st": "State",
		"AD_zp": "12345",
		"AD_cc": "US",
		"$referrer": "https%3A%2F%2Freferrer.com"
	}`

	js, err := simplejson.NewJson([]byte(jsonStr))
	assert.NoError(t, err)

	props, err := SetFBProperties(js)
	assert.NoError(t, err)

	assert.Equal(t, "1750905219951", props.EventTime)
	assert.Equal(t, "event_123", props.EventID)
	assert.Equal(t, int64(1750905219951), props.ClientID2)

	fmt.Println(cast.ToString(props.EventTime2))
}

type FBProperties struct {
	URL_         string      `json:"_url,omitempty"`
	Total        float64     `json:"total,omitempty"`
	KeyWord      string      `json:"key_word,omitempty"`
	ProductTitle string      `json:"product_title,omitempty"`
	Quantity     int         `json:"quantity,omitempty"`
	ProductID    string      `json:"product_id,omitempty"`
	OrderID      string      `json:"order_id,omitempty"`
	Price        float64     `json:"price,omitempty"`
	EventID      string      `json:"event_id,omitempty"`
	Currency     string      `json:"currency,omitempty"`
	Platform     string      `json:"platform,omitempty"`
	VariantIds   interface{} `json:"variant_ids,omitempty"`
	ContentIds   interface{} `json:"content_ids,omitempty"`
	Title        string      `json:"title"`
	ThemeName    string      `json:"theme_name"`
	ThemeVersion string      `json:"theme_version"`
	Client       string      `json:"client"`
	ClientID     string      `json:"-"`
	ClientID2    int64       `json:"-"`
	Referrer     string      `json:"-"`
	EventTime    string
	EventTime2   int64
}

func SetFBProperties(p *simplejson.Json) (*FBProperties, error) {
	properties := new(FBProperties)
	var err error
	properties.URL_, err = p.Get("$url").String()
	properties.KeyWord, err = p.Get("key_word").String()
	properties.ProductTitle, err = p.Get("product_title").String()
	properties.Title, err = p.Get("title").String()
	properties.ProductID, err = p.Get("product_id").String()
	properties.OrderID, err = p.Get("order_id").String()
	properties.Platform, err = p.Get("platform").String() // platform
	properties.Currency, err = p.Get("currency").String()
	properties.VariantIds = p.Get("AD_variant_ids").Interface()
	properties.ContentIds = p.Get("content_ids").Interface()

	properties.ThemeName, err = p.Get("template_name").String()
	properties.ThemeVersion, err = p.Get("theme_version").String()
	properties.Client, err = p.Get("client").String()

	// pagetView 没有AD_event_id
	emptyJson := simplejson.Json{}
	eventId := p.Get("AD_event_id")
	if emptyJson == *eventId {
		eventId = p.Get("event_id")
	}

	properties.EventTime, err = p.Get("event_time").String()
	properties.EventID, err = eventId.String()

	properties.EventTime2 = p.Get("event_time").MustInt64()

	if referrer, err := p.Get("$referrer").String(); err == nil {
		properties.Referrer, err = url.QueryUnescape(referrer)
	}
	clientID := p.Get("client_id")
	properties.ClientID2 = clientID.MustInt64()

	return properties, err
}
