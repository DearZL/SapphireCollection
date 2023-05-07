package global

import (
	"P/handler"
	"P/model"
	"P/repository"
	"P/service"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB               *gorm.DB
	Redis            *redis.Client
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
	dsn := "root:123456@tcp(127.0.0.1:3306)/p?parseTime=true&charset=utf8&parseTime=true&loc=Local"
	DB, dbErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
		&model.Commodity{},
		&model.Order{},
		&model.Block{},
	)
	if err != nil {
		panic(err.Error())
		return
	}
}

func initRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		panic(err.Error())
		return
	}
}

func initHandler() {
	UserHandler = handler.UserHandler{
		Redis: Redis,
		UserSrvI: &service.UserService{
			UserRepo: &repository.UserRepository{
				DB: DB,
			},
			SessionRepo: &repository.SessionRepository{
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
	PayHandler = handler.PayHandler{}
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
		},
	}
	CommodityHandler = handler.CommodityHandler{
		CommoditySrvI: &service.CommodityService{
			CommodityRepo: &repository.CommodityRepository{
				DB: DB,
			},
		},
	}

}

func init() {
	initDataBase()
	initRedis()
	initHandler()
}
