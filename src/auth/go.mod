module eago-auth

go 1.14

replace eago-common => ../common

require (
	eago-common v0.0.0-00010101000000-000000000000
	github.com/Unknwon/goconfig v0.0.0-20200908083735-df7de6a44db8
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.4.2
	github.com/itsjamie/gin-cors v0.0.0-20160420130702-97b4a9da7933
	github.com/jda/go-crowd v0.0.0-20180225080536-9c6f17811dc6
	github.com/micro/go-micro/v2 v2.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	google.golang.org/protobuf v1.23.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.11
)
