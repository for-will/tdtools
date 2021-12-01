module robot

go 1.16

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/json-iterator/go v1.1.11
	github.com/myPuffer/gotosql v0.1.0
	go.uber.org/zap v1.17.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/mysql v1.0.6
	gorm.io/gorm v1.21.9
)

replace (
	github.com/myPuffer/gotosql v0.1.0 => ../../../vmshares/local/gotosql
)