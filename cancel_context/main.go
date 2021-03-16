package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	printHelloUsers(ctx)
}

func printHelloUsers(ctx context.Context) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	type response struct {
		result string
		err error
	}
	ch := make(chan response) // goroutinからデータを受信するためのchanel
	defer close(ch)

	go func() {
		<- time.After(2*time.Second) // 太郎さんは5秒後に呼ばれる
		helloUser, err := getHelloUser(childCtx, 1)
		ch <- response{
			result: helloUser,
			err:    err,
		}
		return
	}()

	go func() {
		helloUser, err := getHelloUser(childCtx, 2) // 花子さんはすぐ呼ばれる
		ch <- response{
			result: helloUser,
			err:    err,
		}
		return
	}()

	var helloUsers []string
	for i := 0; i < 2; i++ {
		r := <- ch
		if r.err != nil {
			fmt.Println(r.err)
			cancel() // エラーを発見した時点でキャンセルする
		} else {
			helloUsers = append(helloUsers, r.result)
		}
	}
	fmt.Println(helloUsers)
}

func getHelloUser(ctx context.Context, id int) (string, error) {
	user, err := getUser(ctx, id)
	if err != nil {
		return "", fmt.Errorf("error id=%d detail: %v", id, err)
	}
	return "こんにちは" + user, nil
}

func getUser(ctx context.Context, id int) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	switch id {
	case 1:
		return "太郎", nil
	case 2:
		return "花子", nil
	default:
		return "", errors.New("not found")
	}
}
