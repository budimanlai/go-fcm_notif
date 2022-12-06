package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	services "github.com/budimanlai/go-cli-service"
	"github.com/eqto/dbm"
	"google.golang.org/api/option"
)

func StartService(mctx *services.Service) {
	jsonFile := mctx.Config.GetString("fcm.json_config")

	ctx := context.Background()
	opts := []option.ClientOption{option.WithCredentialsFile(jsonFile)}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		mctx.Log("new firebase app: %s", err)
		return
	}

	fcmClient, err1 := app.Messaging(ctx)
	if err != nil {
		mctx.Log("messaging: %s", err1)
		return
	}

	for {
		result, e := mctx.Db.Select("SELECT * FROM fcm_messages WHERE status = 'pending' LIMIT 500")
		if e != nil {
			mctx.Log(e.Error())
		}

		if len(result) > 0 {
			var messages = []*messaging.Message{}

			for _, item := range result {
				messages = append(messages, &messaging.Message{
					Notification: &messaging.Notification{
						Title: item.String("title"),
						Body:  item.String("body"),
					},
					Data:  convertData(item.StringOr("data", "")),
					Token: item.String("token"),
				})
			}

			_, err := fcmClient.SendAll(ctx, messages)
			if err != nil {
				mctx.Log(err)
				updateError(mctx, result, err.Error())
			} else {
				mctx.Log("Send all:", len(result))
				updateDone(mctx, result)
			}

			messages = nil
			time.Sleep(1 * time.Second)
		} else {
			mctx.Log("Sleep...")
			time.Sleep(2 * time.Second)
		}

		if mctx.IsStopped {
			mctx.Log("Exit from loop StartService")
			break
		}
	}
}

func updateDone(ctx *services.Service, item []dbm.Resultset) {
	for _, item := range item {
		_, e1 := ctx.Db.Exec("UPDATE fcm_messages SET status = 'done', sended_at = now(), response_log = 'OK' WHERE id = ?",
			item.Int("id"))
		if e1 != nil {
			ctx.Log(e1.Error())
		}
	}
}

func updateError(ctx *services.Service, item []dbm.Resultset, err_msg string) {
	for _, item := range item {
		_, e1 := ctx.Db.Exec("UPDATE fcm_messages SET status = 'error', sended_at = now(), response_log = ? WHERE id = ?",
			err_msg, item.Int("id"))
		if e1 != nil {
			ctx.Log(e1.Error())
		}
	}
}

func convertData(data string) map[string]string {
	var params = make(map[string]string)

	if len(data) == 0 {
		return params
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)

	for key, element := range result {
		switch t := element.(type) {
		case string:
			params[key] = t
			break
		case int:
			params[key] = strconv.Itoa(t)
			break
		case bool:
			params[key] = strconv.FormatBool(t)
			break
		case float64:
			params[key] = fmt.Sprintf("%v", t)
			break
		}
	}

	return params
}

func StopService(mctx *services.Service) {
	mctx.Log("Stop Service")
	mctx.IsStopped = true
}
