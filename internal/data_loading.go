package internal

import (
	"context"
	"fmt"
	"github.com/TrifonovDA/WB/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

const update_cache = "select oi.order_uid, oi.track_number, oi.entry, ui.name, ui.phone, ui.zip, ui.city, ui.address, ui.region, ui.email, " +
	"pi.transaction, pi.request_id, pi.currency, pi.provider, pi.amount, pi.payment_dt, pi.bank, pi.delivery_cost, pi.goods_total, pi.custom_fee, " +
	"ii.chrt_id, ii.track_number, ii.price, ii.rid, ii.name, ii.sale, ii.size, ii.total_price, ii.nm_id, ii.brand, ii.status, " +
	"oi.locale, oi.internal_signature, oi.customer_id, oi.delivery_service, oi.shardkey, oi.sm_id, oi.date_created, oi.oof_shard " +
	"from order_info oi inner join payment_info pi on pi.transaction = order_uid inner join user_info ui on oi.user_uid = ui.user_uid  left join items_info ii on ii.chrt_id = ANY(oi.items) order by order_uid;"

func LoadData(pool *pgxpool.Pool, mu *sync.Mutex, ctx context.Context, ch_error chan error, cache *model.Simple_cache) {
	mu.Lock()
	defer mu.Unlock()
	rows, err := pool.Query(ctx, update_cache)
	if err != nil {
		log.Printf("Error by getting rows: %v", err)
		ch_error <- err
	}
	defer rows.Close()

	var pre_value string
	for rows.Next() {
		var order model.Order
		var item model.Item
		err := rows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
			&order.Delivery.Region, &order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
			&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee, &item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
			&order.Locale, &order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
		if err != nil {
			ch_error <- err
		}
		if order.OrderUid == pre_value {
			cache_order, err := cache.Get_order(order.OrderUid)
			if err != nil {
				ch_error <- fmt.Errorf("get order from cache error. Strange bad case")
			}
			new_items := append(cache_order.Items, item)
			new_order := cache_order
			new_order.Items = new_items
			cache.Insert(new_order)
			continue
		}

		order.Items = append(order.Items, item)
		cache.Insert(order)
		pre_value = order.OrderUid
	}
}
