// @title MCA Bank Gateway API
// @version 1.0
// @description API Gateway aggregating all MCA Bank services (Auth, User, etc.).
// @host localhost:5004
// @BasePath /
// @schemes http
package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/cors"
	pbPaymentClient "github.com/rasteiro11/MCABankGateway/gen/proto/go/payment"
	pbAuthClient "github.com/rasteiro11/MCABankGateway/gen/proto/go/user"
	authRestClient "github.com/rasteiro11/MCABankGateway/pkg/rest/auth"
	customerRestClient "github.com/rasteiro11/MCABankGateway/pkg/rest/customer"
	"github.com/rasteiro11/MCABankGateway/pkg/transport/http/middleware"
	authService "github.com/rasteiro11/MCABankGateway/src/auth/service"
	customerHttp "github.com/rasteiro11/MCABankGateway/src/customer/delivery/http"

	_ "github.com/rasteiro11/MCABankGateway/docs"
	authHttp "github.com/rasteiro11/MCABankGateway/src/auth/delivery/http"
	balanceService "github.com/rasteiro11/MCABankGateway/src/balance/service"
	customerService "github.com/rasteiro11/MCABankGateway/src/customer/service"
	paymentHttp "github.com/rasteiro11/MCABankGateway/src/payment/delivery/http"
	paymentService "github.com/rasteiro11/MCABankGateway/src/payment/service"
	"github.com/rasteiro11/PogCore/pkg/config"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/server"
	"github.com/rasteiro11/PogCore/pkg/telemetry/opentelemetry"
	"github.com/rasteiro11/PogCore/pkg/telemetry/tracer"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	provider := opentelemetry.NewProvider(ctx)
	tracer.SetGlobal(provider)

	ctx, span := tracer.Instance().Start(ctx, "test")
	span.RecordError(errors.New("Simple test"))
	span.End()

	credentials := insecure.NewCredentials()
	authConn, err := grpc.Dial(config.Instance().RequiredString("AUTH_GRPC_SERVICE"),
		grpc.WithTransportCredentials(credentials))
	if err != nil {
		logger.Of(ctx).Fatalf(
			"[main] grpc.Dial returned error: err=%+v", err)
	}

	authClient := pbAuthClient.NewAuthServiceClient(authConn)

	paymentConn, err := grpc.Dial(config.Instance().RequiredString("PAYMENT_GRPC_SERVICE"),
		grpc.WithTransportCredentials(credentials))
	if err != nil {
		logger.Of(ctx).Fatalf(
			"[main] grpc.Dial returned error: err=%+v", err)
	}

	balanceClient := pbPaymentClient.NewBalanceServiceClient(paymentConn)
	paymentClient := pbPaymentClient.NewPaymentServiceClient(paymentConn)

	app := server.NewServer()
	app.AddHandler("/swagger/*", "", http.MethodGet, fiberSwagger.WrapHandler)
	app.Use("/clientes", middleware.ValidateUserMiddleware(authClient))
	app.Use("/*", cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	customerClient := customerRestClient.New(config.Instance().RequiredString("CUSTOMER_SERVICE_URL"))
	authRestClient := authRestClient.New(config.Instance().RequiredString("AUTH_SERVICE_URL"))

	balanceSvc := balanceService.NewBalanceService(balanceClient)
	paymentSvc := paymentService.NewPaymentService(paymentClient)
	customerSvc := customerService.NewCustomerService(customerClient, balanceSvc)
	authSvc := authService.NewAuthService(authRestClient)

	customerHttp.NewHandler(app, customerHttp.WithCustomerService(customerSvc))
	authHttp.NewHandler(app, authHttp.WithAuthService(authSvc))
	paymentHttp.NewHandler(app, paymentHttp.WithPaymentService(paymentSvc))

	app.PrintRouter()

	port := config.Instance().RequiredString("SERVER_PORT")
	if err := app.Start(port); err != nil {
		logger.Of(ctx).Fatalf("[main] server.Start() returned error: %+v\n", err)
	}
}
