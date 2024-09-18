package delivery

import (
	"fmt"
	"polen/config"
	"polen/delivery/controller/api"
	"polen/docs"

	"polen/delivery/middleware"
	"polen/manager"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	ucManager manager.UseCaseManager
	engine    *gin.Engine
	host      string
	log       *logrus.Logger
}

func (s *Server) Run() {
	s.initMiddlewares()
	s.initControllers()
	s.swagDocs()
	err := s.engine.Run(s.host)
	if err != nil {
		panic(err)
	}
}

func (s *Server) initMiddlewares() {
	s.engine.Use(middleware.LogRequestMiddleware(s.log))
}

func (s *Server) initControllers() {
	rg := s.engine.Group("/api/v1")
	api.NewAuthController(s.ucManager.UserUseCase(), s.ucManager.AuthUseCase(), rg).Route()
	api.NewBiodataController(s.ucManager.BiodataUserUseCase(), rg).Route()
	api.NewTopUpController(s.ucManager.TopUpUsecase(), s.ucManager.BiodataUserUseCase(), rg).Route()
	api.NewDepositeInterestController(s.ucManager.DepositerInterestUseCase(), rg).Route()
	api.NewLoanInterestController(s.ucManager.LoanInterestUseCase(), rg).Route()
	api.NewSaldoController(s.ucManager.SaldoUsecase(), rg).Route()
	api.NewDepositeController(s.ucManager.DepositeUsecase(), rg).Route()
	api.NewAppHandlingCostController(s.ucManager.AppHandlingCostUseCase(), rg).Route()
	api.NewLoanController(s.ucManager.LoanUsecase(), rg).Route()
	api.NewLatePaymentFeeController(s.ucManager.LatePaymentFee(), rg).Route()
}

func (s *Server) swagDocs() {
	docs.SwaggerInfo.Title = "Polen p2p Landing App"
	docs.SwaggerInfo.Version = "v1"
	docs.SwaggerInfo.BasePath = "/api/v1"
	s.engine.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
	}
	infraManager, err := manager.NewInfraManager(cfg)
	if err != nil {
		fmt.Println(err)
	}
	repoManager := manager.NewRepoManager(infraManager)
	useCaseManager := manager.NewUsecaseManager(repoManager)

	host := fmt.Sprintf("%s:%s", cfg.ApiHost, cfg.ApiPort)
	log := logrus.New()

	engine := gin.Default()

	// Routing untuk endpoint Swagger UI
	// engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// // Jadwal pembaruan database setiap 24 jam
	// updateInterval := 24 * time.Hour

	// // Mulai penjadwalan pembaruan database
	// ticker := time.NewTicker(updateInterval)

	// // Loop tak terbatas untuk menjalankan pembaruan database
	// go func() {
	// 	err := useCaseManager.DepositeUsecase().Update()
	// 	fmt.Println("hollaa")
	// 	if err != nil {
	// 		log.Error("Gagal melakukan pembaruan database: ", err)
	// 	}
	// 	fmt.Println("Melakukan pembaruan database...")
	// 	for range ticker.C {
	// 		// Panggil fungsi untuk melakukan pembaruan database di sini
	// 		err := useCaseManager.DepositeUsecase().Update()
	// 		if err != nil {
	// 			log.Error("Gagal melakukan pembaruan database: ", err)
	// 		}
	// 		fmt.Println("Melakukan pembaruan database...")
	// 	}
	// }()

	return &Server{
		ucManager: useCaseManager,
		engine:    engine,
		host:      host,
		log:       log,
	}
}
