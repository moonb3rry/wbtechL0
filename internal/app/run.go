package app

import (
	"WBTech0/config"
	"WBTech0/internal/controller/http"
	"WBTech0/internal/repository/cache"
	"WBTech0/internal/repository/nats_queue"
	"WBTech0/internal/repository/pgrepo"
	"WBTech0/internal/service"
	"WBTech0/pkg/httpserver"
	"WBTech0/pkg/nats"
	"WBTech0/pkg/postgres"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {

	pg, err := postgres.New(postgres.GetConnString(&cfg.Db), postgres.MaxPoolSize(cfg.Db.MaxPoolSize))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pg.Close()

	err = pg.Pool.Ping(context.Background())
	if err != nil {
		fmt.Print(err)
	}

	sc, err := nats.New("test-cluster", "client-1", &cfg.Nats.URL)
	if err != nil {
		log.Fatalf("Unable to connect to NATS: %v", err)
	}

	orderRepo := pgrepo.NewOrderRepo(pg)
	cacheRepo := cache.NewOrderCache(orderRepo)
	natsRepo := nats_queue.NewNatsRepository(sc.Conn)
	natsService := service.NewNatsService(natsRepo, cacheRepo, orderRepo)

	_, err = natsService.StartListening("subject", "queueGroup")
	if err != nil {
		log.Fatal("Nats initialization problem: %v", err)
	}

	orderService := service.NewOrderService(orderRepo, cacheRepo)
	httpController := http.NewOrderController(orderService)
	if err != nil {
		log.Fatal("Controller initialization problem: %v", err)
	}

	httpServer := httpserver.New(httpController,
		httpserver.Port(cfg.HttpServer.Addr),
		httpserver.ReadTimeout(cfg.HttpServer.ReadTimeout),
		httpserver.WriteTimeout(cfg.HttpServer.WriteTimeout),
		httpserver.ShutdownTimeout(cfg.HttpServer.ShutdownTimeout),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Printf("Running app:%v version:%v", cfg.App.Name, cfg.App.Version)

	select {
	case s := <-interrupt:
		log.Printf("App running signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Fatal(fmt.Errorf("App HTTP server notify: %v", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("App HTTP server shutdown: %v", err))
	}
	fmt.Printf("server down")
}
