package model

import (
	"fmt"
	"sync"
	"time"
)

// var Mu_cache sync.Mutex
// var Cache_struct = Cache{Mu_cache, make(map[string]Order)}
var cache Simple_cache
var Cache = cache.NewCache()

type Simple_cache struct {
	mu_cache  sync.Mutex
	Cache_map map[string]Order
}

func (cache *Simple_cache) NewCache() Simple_cache {
	return Simple_cache{
		Cache_map: make(map[string]Order),
	}
}

func (cache *Simple_cache) Insert(order Order) {
	cache.mu_cache.Lock()
	defer cache.mu_cache.Unlock()

	cache.Cache_map[order.OrderUid] = order
}
func (cache *Simple_cache) Get_order(order_uid string) (Order, error) {
	cache.mu_cache.Lock()
	defer cache.mu_cache.Unlock()

	order, ok := cache.Cache_map[order_uid]
	if !ok {
		return order, fmt.Errorf("Order is not exist.")
	}
	return order, nil
}

type Item struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
type Order struct {
	OrderUid    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	Delivery    struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction"`
		RequestId    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    int    `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	} `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}
