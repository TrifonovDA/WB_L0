package main

import (
	"encoding/json"
	"fmt"
	"github.com/TrifonovDA/WB/internal/model"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"sync"
	"sync/atomic"
	"time"
)

var conn, _ = stan.Connect("test-cluster", "sub_1")

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	defer conn.Close()
	var sum int32

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			err := publishing(conn, &mu, i)
			atomic.AddInt32(&sum, int32(i))
			if err != nil {
				fmt.Println("publishing err:", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(sum)
}

func publishing(conn stan.Conn, mutex *sync.Mutex, i int) error {
	DateCreated, err := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")
	if err != nil {
		fmt.Println("DateCreated time parsing error")
	}
	OrderUid, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("UUID generating error")
	}
	var Request = model.Order{
		OrderUid:    OrderUid.String(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: struct {
			Name    string `json:"name"`
			Phone   string `json:"phone"`
			Zip     string `json:"zip"`
			City    string `json:"city"`
			Address string `json:"address"`
			Region  string `json:"region"`
			Email   string `json:"email"`
		}{
			Name:    "Test Testovv",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: struct {
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
		}{
			Transaction:  OrderUid.String(),
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []model.Item{{
			ChrtId:      9934930,
			TrackNumber: "WBILMTESTTRACK" + fmt.Sprintf("%v", i),
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmId:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
			{ChrtId: 9934931,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmId:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202},
		},

		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       DateCreated,
		OofShard:          "1",
	}
	marshalling_req, err := json.Marshal(Request)
	if err != nil {
		fmt.Println("marshalling error")
	}
	//mutex.Lock()
	//defer mutex.Unlock()
	//подумать над PublishAsync
	err = conn.Publish("test-cluster", marshalling_req) //, ackHandler)
	//fmt.Println(txt)
	if err != nil {
		return err
	}
	//fmt.Println("blya", i)
	return nil
}
