package main

import (
	"context"
	serverPb "github.com/DaniilOr/myGRPC/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"os"
	"time"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()

	client := serverPb.NewPayServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	response, err := client.GetById(ctx, &serverPb.GetRequest{PaymentId: 1})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}

	log.Printf("Считал: %v\n", response)
	_, err = client.Create(ctx, &serverPb.CreateRequest{Number: "123456789", Name: "Anton"})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}
	response, err = client.GetById(ctx, &serverPb.GetRequest{PaymentId: 2})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}
	log.Printf("Добавил нового пользователя: %v\n", response)
	multiplePayments, err := client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}
	log.Printf("Все пользователи: %v\n", multiplePayments)
	_, err = client.DeleteById(ctx, &serverPb.DeleteRequest{PaymentId: 1})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}
	multiplePayments, err = client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}
	log.Printf("После удаления первого пользователя остались: %v\n", multiplePayments)
	_, err = client.UpdateById(ctx, &serverPb.UpdateRequest{PaymentId: 2, Name: "Dima"})
	multiplePayments, err = client.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Print(st.Code())
			log.Print(st.Message())
		}
		return err
	}
	log.Printf("После изменения имени пользователя: %v\n", multiplePayments)
	return nil
}
