package main

import (
	"github.com/ZhomartZhan/common_lib_hw31"
	users "github.com/ZhomartZhan/users_hw31"
	"github.com/djumanoff/amqp"
	"log"
)

func main() {
	usersMongoStore, err := users.NewUsersStore(common_lib_hw31.MongoConfig{
		Host:           "localhost",
		Port:           "27017",
		Database:       "ivi",
		CollectionName: "users",
	})
	if err != nil {
		log.Fatal(err)
	}
	usersAmqpEndpoints := users.NewUsersAmqpEndpoints(usersMongoStore)
	rabbitConfig := amqp.Config{
		Host:     "localhost",
		Port:     5672,
		LogLevel: 5,
	}
	serverConfig := amqp.ServerConfig{
		ResponseX: "response",
		RequestX:  "request",
	}

	sess := amqp.NewSession(rabbitConfig)
	err = sess.Connect()
	if err != nil {
		panic(err)
		return
	}
	srv, err := sess.Server(serverConfig)
	if err != nil {
		panic(err)
		return
	}
	srv.Endpoint("users.create", usersAmqpEndpoints.CreateUserAmqpEndpoint())
	srv.Endpoint("users.getById", usersAmqpEndpoints.GetUserByIdAmqpEndpoint())
	srv.Endpoint("users.getByUsernameAndPassword", usersAmqpEndpoints.GetByUsernameAndPasswordAmqpEndpoint())
	err = srv.Start()
	if err != nil {
		panic(err)
		return
	}
}
