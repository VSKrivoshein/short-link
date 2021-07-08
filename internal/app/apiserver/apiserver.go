package apiserver

type APIServer struct {
	config *Config
}

func NewAPIServer(config *Config) *APIServer {
	return &APIServer{
		config: config,
	}
}

func (s *APIServer) Start() error {
	return nil
}
