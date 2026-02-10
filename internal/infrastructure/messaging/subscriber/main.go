package subscriber

import "github.com/go-redis/redis/v8"

var rdb *redis.Client

func Run() {
	go subCreateUser()
	go subCreateLogsConfirmacao()
	go subCreateMercadoProdutos()
	go SubUpdateLevelUser()
}
