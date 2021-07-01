module huaweichain

go 1.15

require (
	git.huawei.com/poissonsearch/wienerchain/contract/sdk v0.0.0
	git.huawei.com/poissonsearch/wienerchain/proto v0.0.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/zerolog v1.22.0 // indirect
)

replace (
	git.huawei.com/poissonsearch/wienerchain/contract/sdk => ./contract/sdk
	git.huawei.com/poissonsearch/wienerchain/proto => ./proto
)
