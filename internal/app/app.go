package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/stan.go"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	_ "orderProcessor/docs"
	httphandler "orderProcessor/internal/delivery/http"
	natsstr "orderProcessor/internal/delivery/stan"
	"orderProcessor/internal/repository"
	"orderProcessor/internal/repository/postgres"
	redi "orderProcessor/internal/repository/redis"
	"time"
)

func Run() {
	sc, err := stan.Connect("test-cluster", "orderProcessor", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		slog.Info("error: ", err)
		panic(err)
	}
	defer sc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, "postgres://postgres:password@localhost:5434/postgres")
	if err != nil {
		slog.Info("error: ", err)
		panic(err)
	}
	defer conn.Close()

	postgresDB := postgres.NewPostgres(conn)

	err = postgres.CreateTables(ctx, conn)
	if err != nil {
		slog.Info("error: ", err)
		panic(err)
	}

	opt, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	redisCache := redi.NewRedis(client)

	repo := repository.NewRepository(postgresDB, redisCache)
	err = repo.CacheRecovery(ctx)
	if err != nil {
		slog.Info("error: ", err)
		panic(err)
	}

	handler := httphandler.NewHandler(repo)

	srv := &http.Server{
		Addr: "localhost:8089",
	}

	http.HandleFunc("/order/", handler.GetOrder)
	http.HandleFunc("/order", handler.GetAllOrders)

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	go func() {
		slog.Info("server started")
		if err := srv.ListenAndServe(); err != nil {
			slog.Info("error: ", err)
			panic(err)
		}
		slog.Info("server stopped ")
	}()

	subscriber := natsstr.NewSubscriber(repo)
	err = subscriber.Subscribe(sc)
	if err != nil {
		slog.Info("error: ", err)
		panic(err)
	}

	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	//
	//// Ждем сигнала остановки
	//<-stop
	//
	//// Останавливаем HTTP-сервер
	//if err := srv.Shutdown(context.Background()); err != nil {
	//	slog.Info("error: ", err)
	//	panic(err)
	//}
	//
	//// Закрываем соединение с базой данных
	//if err := conn.Close(ctx); err != nil {
	//	slog.Info("error: ", err)
	//	panic(err)
	//}

	select {}
}
