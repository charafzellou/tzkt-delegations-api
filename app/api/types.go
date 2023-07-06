package main

import "time"

type Request struct {
	Year int `json:"year,omitempty"`
}

type Response struct {
	Data []DelegationApi `json:"data,omitempty"`
}

// Define a DelegationsApi struct
//
//	"timestamp": "2022-05-05T06:29:14Z",
//	"amount": "125896",
//	"delegator": "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
//	"block": "2338084"
type DelegationApi struct {
	Timestamp time.Time `json:"timestamp"`
	Amount    float64   `json:"amount"`
	Delegator string    `json:"delegator"`
	Block     int       `json:"block"`
}

// Define a Delegation struct
type Delegation struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Level     int    `json:"level"`
	Timestamp string `json:"timestamp"`
	Block     string `json:"block"`
	Hash      string `json:"hash"`
	Counter   int    `json:"counter"`
	Initiator struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"initiator"`
	Sender struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"sender"`
	SenderCodeHash int     `json:"senderCodeHash"`
	Nonce          int     `json:"nonce"`
	GasLimit       int     `json:"gasLimit"`
	GasUsed        int     `json:"gasUsed"`
	StorageLimit   int     `json:"storageLimit"`
	BakerFee       int     `json:"bakerFee"`
	Amount         float64 `json:"amount"`
	PrevDelegate   struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"prevDelegate"`
	NewDelegate struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"newDelegate"`
	Status string `json:"status"`
	Errors []struct {
		Type string `json:"type"`
	} `json:"errors"`
	Quote struct {
		Btc int `json:"btc"`
		Eur int `json:"eur"`
		Usd int `json:"usd"`
		Cny int `json:"cny"`
		Jpy int `json:"jpy"`
		Krw int `json:"krw"`
		Eth int `json:"eth"`
		Gbp int `json:"gbp"`
	} `json:"quote"`
}
