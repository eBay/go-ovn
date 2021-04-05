module github.com/ebay/go-ovn

go 1.12

require (
	github.com/cenkalti/hub v1.0.1 // indirect
	github.com/cenkalti/rpc2 v0.0.0-20210220005819-4a29bc83afe1 // indirect
	github.com/ebay/libovsdb v0.2.1-0.20210331070800-9dd672970aef
	github.com/google/uuid v1.1.1
	github.com/stretchr/testify v1.4.0
)

replace github.com/ebay/libovsdb => ../libovsdb
