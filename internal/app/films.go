package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"os/signal"
	"syscall"
	"vk_quests/config"
	v1 "vk_quests/internal/delivery/http/v1"
	"vk_quests/internal/delivery/http/v1/handlers"
	qr "vk_quests/internal/repository/quest"
	ur "vk_quests/internal/repository/user"
	qu "vk_quests/internal/usecase/quest"
	uu "vk_quests/internal/usecase/user"
	"vk_quests/pkg/server"

	_ "github.com/lib/pq"
)

func Run(cfg *config.Config) {
	// Logger
	l, logFile := prepareLogger(cfg.LoggerInfo)

	defer func() {
		if logFile != nil {
			_ = logFile.Close()
		}
		_ = l.Sync()
	}()

	// Postgres
	pg, err := sqlx.Open("postgres", cfg.Postgres.URL)
	if err != nil {
		l.Fatal("[App] Init - postgres.New: %s", err)
	}
	defer pg.Close()

	if err := pg.Ping(); err != nil {
		l.Fatal("[App] Init - can't check connection to sql with error %s", err)
	}
	l.Info("[App] Init - success check connection to postgresql")

	// Repository
	questRepository := qr.NewPostgresQuest(pg)
	userRepository := ur.NewPostgresUser(pg)

	// Use-cases
	questUsecase := qu.NewQuestUsecase(questRepository)
	userUsecase := uu.NewUserUsecase(userRepository, questRepository)

	// Handlers
	questHandlers := handlers.NewQuestHandlers(questUsecase)
	userHandlers := handlers.NewUserHandlers(userUsecase)

	// routes
	router, err := v1.NewRouter("/api", l, prepareRoutes(userHandlers, questHandlers))
	if err != nil {
		l.Fatal("[App] Init - init handler error: %s", err)
	}

	httpServer := server.New(router, server.Port(cfg.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("[App] Start - server started")

	select {
	case s := <-interrupt:
		l.Info("[App] Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("[App] Run - httpServer.Notify: %s", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("[App] Stop - httpServer.Shutdown: %s", err))
	}

	l.Info("[App] Stop - server stopped")
}
