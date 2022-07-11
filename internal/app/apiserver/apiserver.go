package apiserver

import "github.com/sirupsen/logrus"

type APIServer struct {
	config *Config
	logger *logrus.Logger
}

func NewAPIServer(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
	}
}

func (a *APIServer) Start() error {
	if err := a.configLogger(); err != nil {
		return err
	}

	a.logger.Info("Starting api server")

	return nil
}

func (a *APIServer) configLogger() error {
	level, err := logrus.ParseLevel(a.config.LogLevel)
	if err != nil {
		return err
	}

	a.logger.SetLevel(level)
	return nil
}
