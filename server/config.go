package server

var (
	TCP_PORT			= "9090"
	HTTP_PORT			= "8080"
	MAX_CLIENTS			= 100
	MAX_REQUEST_COUNT	= 100
	CLIENT_ID			= 1
	REQUEST_PER_SECOND	= 30

	EXTERNAL_APIS		= []string{ "http://127.0.0.1:8081/delay", "http://127.0.0.1:8081/fail"}
)