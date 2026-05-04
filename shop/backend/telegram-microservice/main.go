package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AndreyLebedev1998/admin-grpc"
	auth "github.com/AndreyLebedev1998/auth-grpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var bot *tgbotapi.BotAPI
	var ctx = context.Background()

	connAuth, err := grpc.Dial("authorization-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				MaxDelay:   5 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
			},

			MinConnectTimeout: 5 * time.Second,
		}),
	)

	if err != nil {
		log.Fatal("gRPC dial error:", err)
	}

	defer connAuth.Close()

	connAdmin, err := grpc.Dial("admin-microservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				MaxDelay:   5 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
			},

			MinConnectTimeout: 5 * time.Second,
		}),
	)

	if err != nil {
		log.Fatal("gRPC dial error:", err)
	}

	defer connAdmin.Close()

	clientAuth := auth.NewAuthServiceClient(connAuth)
	clientAdmin := admin.NewAdminServiceClient(connAdmin)

	// Получение телегарм бота
	token, err := clientAdmin.GetTelegramToken(ctx, &admin.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	if token != nil {
		bot, err = tgbotapi.NewBotAPI(token.Token)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		log.Fatal("Token not found")
		return
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	var tokenTgUser string

	go func() {
		for update := range updates {
			if update.Message != nil {
				fmt.Println(update.Message.Chat.ID)
				args := strings.Split(update.Message.Text, " ")
				fmt.Println(update.Message.Text)
				if len(args) > 1 {
					tokenTgUser = args[1]
					var tgChatId = update.Message.Chat.ID
					resp, err := clientAuth.UpdateChatIdTgUser(ctx, &auth.ChatIdTgUser{
						ChatId: update.Message.Chat.ID,
						Token:  tokenTgUser,
					})

					if err != nil {
						msg := tgbotapi.NewMessage(tgChatId, err.Error())
						_, err = bot.Send(msg)
						if err != nil {
							log.Println(err)
						}
					}

					if resp != nil {
						msg := tgbotapi.NewMessage(tgChatId, "Вы успешно привязали свой аакатунт к боту")

						_, err = bot.Send(msg)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}()

	mux := http.NewServeMux()
	http.ListenAndServe(":8080", mux)
}
