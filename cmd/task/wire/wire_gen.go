// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"sakura/internal/repository"
	"sakura/internal/server"
	"sakura/internal/task"
	"sakura/pkg/app"
	"sakura/pkg/log"
	"sakura/pkg/sid"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, logger *log.Logger) (*app.App, func(), error) {
	db := repository.NewDB(viperViper, logger)
	repositoryRepository := repository.NewRepository(logger, db)
	transaction := repository.NewTransaction(repositoryRepository)
	sidSid := sid.NewSid()
	taskTask := task.NewTask(transaction, logger, sidSid)
	userRepository := repository.NewUserRepository(repositoryRepository)
	userTask := task.NewUserTask(taskTask, userRepository)
	taskServer := server.NewTaskServer(logger, userTask)
	appApp := newApp(taskServer)
	return appApp, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewTransaction, repository.NewUserRepository)

var taskSet = wire.NewSet(task.NewTask, task.NewUserTask)

var serverSet = wire.NewSet(server.NewTaskServer)

// build App
func newApp(task2 *server.TaskServer,
) *app.App {
	return app.NewApp(app.WithServer(task2), app.WithName("demo-task"))
}
