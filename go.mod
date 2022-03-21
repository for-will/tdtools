module robot

go 1.18

require (
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

require (
	github.com/ahmetb/go-linq/v3 v3.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/mod v0.5.1 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/myPuffer/gotosql v0.1.0 => ../../../vmshares/local/gotosql
