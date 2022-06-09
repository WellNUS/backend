package config

var (
	DOMAIN string = "localhost"
	FRONTEND_URL string = "http://localhost:3000"
	BACKEND_URL string = "localhost:8080"
	
	// Database fields (if changed update in Makefile as well)
	CONNECTION_STRING string = "postgresql://root:password@0.0.0.0:49730/wellnus?sslmode=disable"
)
