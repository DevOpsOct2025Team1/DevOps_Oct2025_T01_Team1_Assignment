module github.com/provsalt/DOP_P01_Team1/auth-service

go 1.25.5

require (
	github.com/golang-jwt/jwt/v5 v5.2.3
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/provsalt/DOP_P01_Team1/common v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.46.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/fatih/color v1.16.0 // indirect
	golang.org/x/exp v0.0.0-20250808145144-a408d31f581a // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/provsalt/DOP_P01_Team1_Assignment/user-service => ../user-service

replace github.com/provsalt/DOP_P01_Team1/common => ../common
