package internal

import (
	"context"
	"fmt"
	"github.com/TrifonovDA/WB/internal/model"
	"github.com/TrifonovDA/WB/pkg/other"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

const query_in_order = "insert into public.order_info(order_uid, track_number, entry, user_uid, items, locale, internal_signature, customer_id,delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES " +
	"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);"
const query_in_user_info = "insert into public.user_info(user_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
const query_get_user_info = "select user_uid from public.user_info where name = $1 and phone = $2 and zip = $3 and city = $4 and address = $5 and region = $6 and email = $7;"
const query_in_payment_info = "insert into public.payment_info(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);"
const query_in_items_info = "insert into public.items_info(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
const query_get_item = "select chrt_id from public.items_info where chrt_id = $1;"

func SaveData(order model.Order, db_conn *pgxpool.Pool, ctx context.Context, mu *sync.Mutex, ch_err chan error) {
	mu.Lock()
	defer mu.Unlock()
	var items []int

	for _, elem := range order.Items {
		var chrt_id pgtype.Int4
		items = append(items, elem.ChrtId)
		//добавлена проверка по айди, если есть - ничего не делаем, если нет - добавить строку. Чтобы не перегружать базу данными
		check_item_row := db_conn.QueryRow(ctx, query_get_item, elem.ChrtId)
		check_item_err := check_item_row.Scan(&chrt_id)
		//Если нет такого итема в бд, добавляем строчку
		if check_item_err != nil {
			if check_item_err.Error() == other.UpdatedNoRowsErr.Error() {
				row_item, err_item := db_conn.Query(ctx, query_in_items_info, elem.ChrtId, elem.TrackNumber, elem.Price, elem.Rid, elem.Name, elem.Sale, elem.Size, elem.TotalPrice, elem.NmId, elem.Brand, elem.Status)
				defer row_item.Close()
				if err_item != nil {
					if pgErr, ok := err_item.(*pgconn.PgError); ok { //обработка ошибок бд
						newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Code: %s, SQLState: %%", pgErr.Message, pgErr.Detail, pgErr.Code, pgErr.SQLState())
						ch_err <- fmt.Errorf("Item inserting. Error in database: %v ", newErr)
						return
					} else {
						log.Printf("Item inserting. Unknown error from db by inserting order! Error: %v", err_item)
						ch_err <- err_item
						return
					}
				}
			} else {
				ch_err <- fmt.Errorf("Get item query error.Unknown error in database_tools connection! Error: %v", check_item_err)
				return
			}
		}
	}

	//инсертим инфу по юзеру. Если такая уже есть, не инсертим.
	var user_uid pgtype.UUID
	new_user_uid, err_uid := uuid.NewUUID()
	if err_uid != nil {
		ch_err <- fmt.Errorf("Generating user uuid error: %v", err_uid)
		return
	}
	var new_user_flg bool
	check_user_row := db_conn.QueryRow(ctx, query_get_user_info, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	check_user_err := check_user_row.Scan(&user_uid)

	if check_user_err != nil {
		if check_user_err.Error() == other.UpdatedNoRowsErr.Error() {
			new_user_flg = true
			row_user, err_user := db_conn.Query(ctx, query_in_user_info, new_user_uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
			defer row_user.Close()
			if err_user != nil {
				if pgErr, ok := err_user.(*pgconn.PgError); ok { //обработка ошибок бд
					newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Code: %s, SQLState: %%", pgErr.Message, pgErr.Detail, pgErr.Code, pgErr.SQLState())
					ch_err <- fmt.Errorf("User inserting. Error in database: %v ", newErr)
					return
				} else {
					ch_err <- fmt.Errorf("User inserting. Unknown error from db by inserting order! Error: %v", err_user)
					return
				}
			}
		} else {
			ch_err <- fmt.Errorf("Get user query error.Unknown error in database_tools connection! Error: %v", check_user_err)
			return
		}
	}

	//инсертим ордеры
	var user_uid_for_order uuid.UUID
	switch new_user_flg {
	case true:
		user_uid_for_order = new_user_uid
	default:
		user_uid_for_order = user_uid.Bytes
	}
	//if New_user_uid
	row_order, err_order := db_conn.Query(ctx, query_in_order, order.OrderUid, order.TrackNumber, order.Entry, user_uid_for_order, items, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	defer row_order.Close()
	if err_order != nil {
		if pgErr, ok := err_order.(*pgconn.PgError); ok { //обработка ошибок бд
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Code: %s, SQLState: %%", pgErr.Message, pgErr.Detail, pgErr.Code, pgErr.SQLState())
			ch_err <- fmt.Errorf("Order inserting. Error in database: %v ", newErr)
			return
		} else {
			ch_err <- fmt.Errorf("Order inserting. Unknown error from db by inserting order! Error: %v", err_order)
			return
		}
	}

	//инсертим платежную информацию
	row_payment, err_payment := db_conn.Query(ctx, query_in_payment_info, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee)
	defer row_payment.Close()
	if err_payment != nil {
		if pgErr, ok := err_payment.(*pgconn.PgError); ok { //обработка ошибок бд
			newErr := fmt.Errorf("Payment inserting. SQL Error: %s, Detail: %s, Code: %s, SQLState: %%", pgErr.Message, pgErr.Detail, pgErr.Code, pgErr.SQLState())
			ch_err <- fmt.Errorf("Error in database: %v ", newErr)
			return
		} else {
			ch_err <- fmt.Errorf("Payment inserting. Unknown error from db by inserting order! Error: %v", err_payment)
			return
		}
	}
	return
}
