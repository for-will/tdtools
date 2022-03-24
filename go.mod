module market

go 1.16

require (
	github.com/ahmetb/go-linq/v3 v3.2.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/json-iterator/go v1.1.11
	github.com/myPuffer/gotosql v0.1.0
	go.uber.org/zap v1.17.0
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	golang.org/x/tools v0.1.8
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/mysql v1.0.6
	gorm.io/gorm v1.21.9
)

replace github.com/myPuffer/gotosql v0.1.0 => ../../../vmshares/local/gotosql
