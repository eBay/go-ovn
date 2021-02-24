module github.com/ebay/go-ovn

go 1.12

require (
	github.com/cenk/hub v1.0.1 // indirect
	github.com/ebay/libovsdb v0.0.0-20190718202342-e49b8c4e1142
	github.com/google/uuid v1.1.1
	github.com/stretchr/testify v1.4.0
)

replace github.com/ebay/libovsdb => github.com/amorenoz/libovsdb v0.0.0-20210331101800-0be550fb92be
