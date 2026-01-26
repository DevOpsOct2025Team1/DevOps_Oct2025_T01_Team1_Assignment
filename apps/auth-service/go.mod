module github.com/provsalt/DOP_P01_Team1/auth-service

go 1.25.5

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/provsalt/DOP_P01_Team1/common v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.47.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/provsalt/DOP_P01_Team1_Assignment/user-service => ../user-service

replace github.com/provsalt/DOP_P01_Team1/common => ../common
