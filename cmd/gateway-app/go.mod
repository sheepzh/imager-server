module gateway-app

go 1.15

require pkg v0.0.0

replace pkg => ../../pkg

require (
	gateway v0.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
)

replace gateway => ../../gateway
