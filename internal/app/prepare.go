package app

import (
	"io"
	"log"
	"net/http"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"vk_quests/config"
	_ "vk_quests/docs"
	v1 "vk_quests/internal/delivery/http/v1"
	"vk_quests/internal/delivery/http/v1/handlers"
	"vk_quests/internal/pkg/prepare"
	"vk_quests/pkg/logger"
)

func prepareLogger(cfg config.LoggerInfo) (*logger.Logger, *os.File) {
	var logOut io.Writer
	var logFile *os.File
	var err error

	if cfg.Directory != "" {
		logFile, err = prepare.OpenLogDir(cfg.Directory)
		if err != nil {
			log.Fatalf("[App] Init - create logger error: %s", err)
		}

		logOut = logFile
	} else {
		logOut = os.Stderr
		logFile = nil
	}

	l := logger.New(
		logger.Params{
			AppName:                  cfg.AppName,
			LogDir:                   cfg.Directory,
			Level:                    cfg.Level,
			UseStdAndFile:            cfg.UseStdAndFile,
			AddLowPriorityLevelToCmd: cfg.AllowShowLowLevel,
		},
		logOut,
	)

	return l, logFile
}

func prepareRoutes(userHandlers *handlers.UserHandlers, questHandlers *handlers.QuestHandlers) v1.Routes {
	return v1.Routes{
		//"Index"
		v1.Route{
			Method:      http.MethodGet,
			Pattern:     "/swagger/*any",
			HandlerFunc: ginSwagger.WrapHandler(swaggerFiles.Handler),
		},

		// "CreateUser"
		v1.Route{
			Method:      http.MethodPost,
			Pattern:     "/user",
			HandlerFunc: userHandlers.CreateUser,
		},

		// "DeleteUser"
		v1.Route{
			Method:      http.MethodDelete,
			Pattern:     "/user/:" + handlers.UserIdField,
			HandlerFunc: userHandlers.DeleteUser,
		},

		// "UpdateUser"
		v1.Route{
			Method:      http.MethodPut,
			Pattern:     "/user/:" + handlers.UserIdField,
			HandlerFunc: userHandlers.UpdateUser,
		},

		// "GetUserHistory"
		v1.Route{
			Method:      http.MethodGet,
			Pattern:     "/user/:" + handlers.UserIdField + "/history",
			HandlerFunc: userHandlers.GetUserHistory,
		},

		// "GetUsers"
		v1.Route{
			Method:      http.MethodGet,
			Pattern:     "/user/list",
			HandlerFunc: userHandlers.GetUsers,
		},

		// "CompleteQuest"
		v1.Route{
			Method:      http.MethodPost,
			Pattern:     "/user/complete",
			HandlerFunc: userHandlers.CompleteQuest,
		},

		// "CreateQuest"
		v1.Route{
			Method:      http.MethodPost,
			Pattern:     "/quest",
			HandlerFunc: questHandlers.CreateQuest,
		},

		// "DeleteQuest"
		v1.Route{
			Method:      http.MethodDelete,
			Pattern:     "/quest/:" + handlers.QuestIdField,
			HandlerFunc: questHandlers.DeleteQuest,
		},

		// "UpdateQuest"
		v1.Route{
			Method:      http.MethodPut,
			Pattern:     "/quest/:" + handlers.QuestIdField,
			HandlerFunc: questHandlers.UpdateQuest,
		},

		// "GetQuest"
		v1.Route{
			Method:      http.MethodGet,
			Pattern:     "/quest/:" + handlers.QuestIdField,
			HandlerFunc: questHandlers.GetQuest,
		},

		// "GetQuests",
		v1.Route{
			Method:      http.MethodGet,
			Pattern:     "/quest/list",
			HandlerFunc: questHandlers.GetQuests,
		},
	}
}
