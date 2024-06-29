package restserver

type SeverOption func(*Server)

func WithPort(port int) SeverOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithMode(mode string) SeverOption {
	return func(s *Server) {
		s.mode = mode
	}
}

func WithHealthz(healthz bool) SeverOption {
	return func(s *Server) {
		s.healthz = healthz
	}
}

func WithEnableProfiling(enableProfiling bool) SeverOption {
	return func(s *Server) {
		s.enableProfiling = enableProfiling
	}
}

func WithMiddlewares(middlewares []string) SeverOption {
	return func(s *Server) {
		s.middlewares = middlewares
	}
}

func WithJwt(jwt *JwtInfo) SeverOption {
	return func(s *Server) {
		s.jwt = jwt
	}
}

func WithTransName(transName string) SeverOption {
	return func(s *Server) {
		s.transName = transName
	}
}
