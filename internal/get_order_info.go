package internal

import (
	"github.com/TrifonovDA/WB/internal/model"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type Request struct {
	OrderUid string `json:"order_uid"`
}
type Response struct {
	Order model.Order `json:"order"`
}

var tpl = template.Must(template.ParseFiles("cmd/index.html"))

func BasisHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func Get_order_info_handler(w http.ResponseWriter, r *http.Request) {
	var response Response
	//Парсинг URL
	url_str, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(200)
		_, bad_err := w.Write([]byte("Error by parsing url!"))
		log.Println("Error by parsing url: ", r.URL.String())
		if bad_err != nil {
			log.Printf("Error by parsing url %v and writing response err %v\n", err, bad_err)
		}
		return
	}
	//Получение параметров запроса
	params := url_str.Query()
	order_uid := params.Get("order_uid")

	//достаем значения ордера из кэша.
	order, err := model.Cache.Get_order(order_uid)
	if err != nil {
		w.WriteHeader(200)
		_, bad_err := w.Write([]byte("Error by getting order/order is not exist!"))
		if bad_err != nil {
			log.Printf("Error by getting order/order is not exist %v and writing response err %v\n", err, bad_err)
		}
		return
	}
	response.Order = order
	//отправляем html страницу пользователю
	err = tpl.Execute(w, &response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, bad_err := w.Write([]byte("Error by executing response!"))
		if bad_err != nil {
			log.Printf("Error by executing response %v and writing response err %v\n", err, bad_err)
		}
	}
	return
}
