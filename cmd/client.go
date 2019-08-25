package main

import (
	"context"
	"fmt"
	"github.com/ios116/regservice/config"
	"github.com/ios116/regservice/session"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.NewConfig()
	logger, _ := cfg.CreateLogger()
	sugar := logger.Sugar()
	defer logger.Sync()
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		sugar.Fatal("Can't connect to grpc= ", address)
	}

	sessManager := session.NewAuthCheckerClient(conn)
	ctx := context.Background()
	// создаем сессию
	sessId, err := sessManager.Create(ctx,
		&session.Session{
			Login:     "rvasily",
			Useragent: "chrome",
		})
	sugar.Info("create sessId=", sessId, err)

	// проеряем сессию
	sess, err := sessManager.Check(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})
	sugar.Info("check session=", sess, err)

	// удаляем сессию
	_, err = sessManager.Delete(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess, err = sessManager.Check(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})
	sugar.Info("check session repeat", sess, err)

	defer conn.Close()
}
