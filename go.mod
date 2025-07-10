module github.com/pavandhadge/goFileShare

go 1.23.5

replace github.com/pavandhadge/goFileShare/server => ./server

replace github.com/pavandhadge/goFileShare/handlers => ./handlers

replace github.com/pavandhadge/goFileShare/auth => ./auth

replace github.com/pavandhadge/goFileShare/utils => ./utils

require github.com/cloudinary/cloudinary-go/v2 v2.9.1

require (
	github.com/creasty/defaults v1.7.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gorilla/schema v1.4.1 // indirect
)
