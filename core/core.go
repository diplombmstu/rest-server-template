package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"
	"github.com/diplombmstu/rest-server-template/middlewares/mongodb"
	"github.com/diplombmstu/rest-server-template/middlewares/renderer"
	"github.com/diplombmstu/rest-server-template/server"
	"github.com/diplombmstu/rest-server-template/middlewares/context"
	"github.com/diplombmstu/rest-server-template/resources/sessions"
	"crypto/x509"
	"crypto/rsa"
	"encoding/pem"
	"github.com/golang/glog"
	"github.com/diplombmstu/rest-server-template/resources/events"
	"github.com/diplombmstu/rest-server-template/core/settings"
	"github.com/diplombmstu/rest-server-template/sse-service/sse_server"
	"github.com/diplombmstu/rest-server-template/resources/users"
)

func getPrivateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateSigningKeyData, err := ioutil.ReadFile(Settings.GetRsaPrivateKey())
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error()))
	}

	privPemBlock, _ := pem.Decode(privateSigningKeyData)
	privateSigningKey, err := x509.ParsePKCS1PrivateKey(privPemBlock.Bytes)

	publicSigningKeyData, err := ioutil.ReadFile(Settings.GetRsaPublicKey())
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error()))
	}

	pubPemBlock, _ := pem.Decode(publicSigningKeyData)
	publicSigningKey, err := x509.ParsePKIXPublicKey(pubPemBlock.Bytes)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error()))
	}

	return privateSigningKey, publicSigningKey.(*rsa.PublicKey), nil
}

func initLogging() {
	//log.SetFlags()
}

func initSettings(cfgFile string) {
	var err error
	Settings, err = settings.NewSettings(cfgFile)

	if err != nil {
		message := fmt.Sprintf("Failed to load settings %v", err.Error())
		glog.Errorln(message)
		panic(message)
	}
}

func Start(pars Parameters) {
	initLogging()
	initSettings(pars.ConfigFile)

	eventBroker := sse_server.NewBroker()

	glog.Infoln("Loading key pair...")
	privateSigningKey, publicSigningKey, err := getPrivateKeyPair()
	if err != nil {
		panic(err.Error())
	}

	ctx := context.New()

	glog.Infoln("Setting up DB session...")
	db := mongodb.New(&mongodb.Options{
		ServerName:   Settings.GetDataBaseServerName(),
		DatabaseName: Settings.GetDataBaseName(),
	})
	_ = db.NewSession()

	glog.Infoln("Setting up Renderer (unrolled_render)...")
	concreteRenderer := renderer.New(&renderer.Options{
		IndentJSON: true,
	}, renderer.JSON)

	glog.Infoln("Setting up router...")
	ac := server.NewAccessController(ctx, concreteRenderer)
	router := server.NewRouter(ctx, ac)

	glog.Infoln("Setting up resources...")

	glog.Infoln("...users resource...")
	var userRepoFactory sessions.IUserRepositoryFactory
	usersResource := users.NewResource(ctx, &users.Options{
		Database: db,
		Renderer: concreteRenderer,
	})

	userRepoFactory = usersResource.UserRepositoryFactory
	router.AddResources(usersResource)

	glog.Infoln("...session resource...")
	sessionsResource := sessions.NewResource(ctx, &sessions.Options{
		PrivateSigningKey:     privateSigningKey,
		PublicSigningKey:      publicSigningKey,
		Database:              db,
		Renderer:              concreteRenderer,
		UserRepositoryFactory: userRepoFactory,
	})

	glog.Infoln("...events resource...")
	eventsResource := events.NewResource(ctx, &events.Options{
		Broker:  eventBroker,
		Renderer:concreteRenderer,
	})

	glog.Infoln("Initializing the server...")
	s := server.NewServer(&server.Config{
		Context: ctx,
	})

	router.AddResources(sessionsResource, eventsResource)

	s.UseMiddleware(sessionsResource.NewAuthenticator())
	s.UseRouter(router)

	glog.Infoln("Running the server...")
	s.Run(fmt.Sprintf(":%d", Settings.GetPortToServe()), server.Options{
		Timeout: Settings.GetTimeout() * time.Second,
	})
}
