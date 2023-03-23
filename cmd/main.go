package main

import (
	"context"
	"github.com/TrifonovDA/WB/internal"
	"github.com/TrifonovDA/WB/internal/model"
	"github.com/TrifonovDA/WB/pkg/database_tools"
	nats_streaming_tools "github.com/TrifonovDA/WB/pkg/nats-streaming_tools"
	"github.com/TrifonovDA/WB/pkg/other"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"sync"
)

const get_order_info = "/get_order_info"
const basis = "/"

// создание мультиплексора для http сервера
var mux = http.NewServeMux()

// подключение стилей к http серверу
var fs = http.FileServer(http.Dir("cmd/assets"))

func main() {
	var mu_loaddata sync.Mutex
	var mu sync.Mutex
	//создание кэша
	//var cache model.Cache
	//cache = cache.NewCache()
	//создание каналов для обработки двух типов ошибок
	exit_ch := make(chan error)
	errors_ch := make(chan error, 4)
	ctx := context.Background()

	//запуск селектора ошибок из каналов
	go other.SelectionErrors(exit_ch, errors_ch, ctx)

	//Создание пула коннектов к бд
	dbConnection := database_tools.NewConnection(ctx)
	//загрузка данных из бд в кэш
	internal.LoadData(dbConnection, &mu_loaddata, ctx, errors_ch, &model.Cache)
	//коннект к nats_streaming серверу
	var nats_streaming_conn, err_nats_streaming_conn = stan.Connect("test-cluster", "sub_2")
	if err_nats_streaming_conn != nil {
		log.Fatalln("err conn", err_nats_streaming_conn)
	}

	//обработка сообщений из nats_streaming сервера
	go nats_streaming_tools.HandleNatsMessage(nats_streaming_conn, errors_ch, dbConnection, &mu)

	//поддержка метода получения информации по ордеру и базового пути "/", а также поддержка стилей
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc(get_order_info, internal.Get_order_info_handler)
	mux.HandleFunc(basis, internal.BasicHandler)

	log.Printf("start listening http")
	err := http.ListenAndServe(":60601", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
