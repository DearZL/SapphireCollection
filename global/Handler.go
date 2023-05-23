package global

import (
	"P/conf"
	"P/handler"
	"P/model"
	"P/repository"
	"P/service"
	"encoding/gob"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB               *gorm.DB
	Redis            *redis.Client
	Redis1           *redis.Client
	UserHandler      handler.UserHandler
	SessionHandler   handler.SessionHandler
	CommonHandler    handler.CommonHandler
	PayHandler       handler.PayHandler
	BlockHandler     handler.BlockHandler
	OrderHandler     handler.OrderHandler
	CommodityHandler handler.CommodityHandler
)

func initDataBase() {
	var dbErr error
	var user = viper.GetString("database.username") + ":" + viper.GetString("database.password")
	var port = viper.GetString("database.port")
	var host = viper.GetString("database.host")
	dsn := user + "@tcp(" + host + ":" + port + ")/p?parseTime=true&charset=utf8&parseTime=true&loc=Local"
	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	},
	)
	if dbErr != nil {
		panic(dbErr.Error())
		return
	}
	err := DB.AutoMigrate(
		&model.File{},
		&model.User{},
		&model.Order{},
		&model.Block{},
		&model.Commodity{},
	)
	if err != nil {
		panic(err.Error())
		return
	}
}

func initRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host" + ":" + viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       0,
	})

	Redis1 = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host" + ":" + viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       3,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		panic(err.Error())
		return
	}
}

func initHandler() {
	UserHandler = handler.UserHandler{
		UserSrvI: &service.UserService{
			Redis: Redis,
			UserRepo: &repository.UserRepository{
				DB: DB,
			},
			SessionRepo: &repository.SessionRepository{
				DB: DB,
			},
			CommodityRepo: &repository.CommodityRepository{
				DB: DB,
			},
		},
	}

	SessionHandler = handler.SessionHandler{
		SessionSrvI: &service.SessionService{
			SessionRepo: &repository.SessionRepository{
				DB: DB,
			},
		},
	}

	CommonHandler = handler.CommonHandler{
		CommonSrvI: &service.CommonService{
			CommonRepo: &repository.CommonRepository{
				DB: DB,
			},
		},
	}

	PayHandler = handler.PayHandler{
		PaySrvI: &service.PayService{},
		OrderSrvI: &service.OrderService{
			OrderRepo: &repository.OrderRepository{
				DB: DB,
			},
		},
	}

	BlockHandler = handler.BlockHandler{
		BlockSrvI: &service.BlockService{
			BlockRepo: &repository.BlockRepository{
				DB: DB,
			},
		},
	}

	OrderHandler = handler.OrderHandler{
		OrderSrvI: &service.OrderService{
			OrderRepo: &repository.OrderRepository{
				DB: DB,
			},
			CommodityRepo: &repository.CommodityRepository{
				DB: DB,
			},
			PaySrv: &service.PayService{},
		},
	}

	CommodityHandler = handler.CommodityHandler{
		CommoditySrvI: &service.CommodityService{
			CommodityRepo: &repository.CommodityRepository{
				DB: DB,
			},
		},
		BlockSrvI: &service.BlockService{
			BlockRepo: &repository.BlockRepository{
				DB: DB,
			},
		},
	}

}

func init() {
	conf.LogConf()
	conf.ViperConf()
	gob.Register(&model.User{})
	initDataBase()
	initRedis()
	initHandler()
}
