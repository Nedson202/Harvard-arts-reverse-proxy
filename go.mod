module github.com/nedson202/harvard-arts-reverse-proxy

go 1.12

require (
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/subosito/gotenv v1.2.0
)

replace github.com/nedson202/harvard-arts-reverse-proxy => ./../harvard-arts-reverse-proxy
