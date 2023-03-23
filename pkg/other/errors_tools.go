package other

import (
	"context"
	"errors"
	"log"
	"time"
)

var UpdatedNoRowsErr error = errors.New("no rows in result set")

func SelectionErrors(exit chan error, nofatal_err chan error, ctx context.Context) {
	for {
		select {
		case some_error, ok := <-exit:
			if ok == false {
				time.Sleep(time.Second * 1)
				continue
			}
			log.Fatalln("Something get fatal wrong!: ", some_error.Error())
		case some_error, ok := <-nofatal_err:
			if ok == false {
				time.Sleep(time.Second * 1)
				continue
			}
			log.Println("Something get nofatal wrong!: ", some_error.Error())
		case <-ctx.Done():
			log.Fatalln("Background context was done!")
		default:
			time.Sleep(time.Second * 1)
		}
	}
}
