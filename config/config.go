package config

var (
	DOMAIN string = "localhost"
	FRONTEND_URL string = "http://localhost:3000"
	BACKEND_URL string = "localhost:8080"
	
	// Database fields (if changed update in Makefile as well)
	CONNECTION_STRING string = "postgresql://root:password@localhost:5432/wellnus?sslmode=disable"
)
