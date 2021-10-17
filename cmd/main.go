package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/jackc/pgx/v4/pgxpool"
	"gokit/order"
	"net/http"
	"os"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)

	var svc order.OrderService

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("eu-west-1")},
		Profile: "default",
	})

	s := sqs.New(sess)

	fmt.Println("started 1")
	p, err := pgxpool.Connect(context.Background(), "postgresql://postgres:ozgotozgot1@database-1.cptj1r7jikob.eu-west-1.rds.amazonaws.com:5432")
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err); err != nil {
			return
		}
		os.Exit(1)
	}
	defer p.Close()

	rdbmsAccess := order.NewRdbmsAccess(p)
	sqsAccess := order.NewSqsAccess(s,"https://sqs.eu-west-1.amazonaws.com/236584826472/StandardQueue",int64(10))
	orderService := order.NewOrderService(rdbmsAccess,sqsAccess)
	orderService = order.LoggingMiddleware(logger)(svc)

	addOrderHandler := httptransport.NewServer(
		order.MakeAddOrderEndpoint(orderService),
		order.DecodeOrderRequest,
		order.EncodeResponse,
	)

	http.Handle("/order", addOrderHandler)
	http.ListenAndServe("localhost:8080", nil)
}

