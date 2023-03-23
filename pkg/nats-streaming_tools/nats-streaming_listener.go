package nats_streaming_tools

import (
	"context"
	"encoding/json"
	"github.com/TrifonovDA/WB/internal"
	model "github.com/TrifonovDA/WB/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/stan.go"
	"sync"
	"time"
)

func HandleNatsMessage(conn stan.Conn, ch_err chan error, db_conn *pgxpool.Pool, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	var datasaver_mu sync.Mutex

	for {
		var order model.Order
		sc, err := conn.Subscribe("test-cluster", func(m *stan.Msg) {
			err := json.Unmarshal(m.Data, &order)
			if err != nil {
				ch_err <- err
				return
			}
			//добавляем значение в кэш
			go model.Cache.Insert(order)
			//распараллелить, отправить ошибки в канал!
			go internal.SaveData(order, db_conn, context.Background(), &datasaver_mu, ch_err)

		}, stan.DurableName("my-durable"))
		if err != nil {
			ch_err <- err // запись в канал
		}
		time.Sleep(1 * time.Second)
		sc.Close()
	}

}
