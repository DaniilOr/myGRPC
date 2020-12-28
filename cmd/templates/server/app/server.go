package app

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	serverPb "github.com/DaniilOr/myGRPC/pkg/server"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Server struct {
	db  *pgxpool.Pool
	ctx context.Context
}
func NewServer(db *pgxpool.Pool, ctx context.Context) *Server {
	return &Server{db: db, ctx: ctx}
}

func (s*Server) GetById(ctx context.Context, request *serverPb.GetRequest) (*serverPb.AutoPay, error){
	var response serverPb.AutoPay
	var created time.Time
	var updated time.Time
		row := s.db.QueryRow(ctx, `
		SELECT * FROM auto_payments WHERE payment_id=$1
	`, request.PaymentId)
		err := row.Scan(&response.PaymentId, &response.Name, &response.Number, &created, &updated)
		if err != nil{
			log.Println(err)
			return nil, err
		}
		response.TimeCreated = timestamppb.New(created)
		response.TimeUpdated = timestamppb.New(updated)
		return &response, nil
}

func (s*Server) GetAll(ctx context.Context, request *emptypb.Empty) (*serverPb.AllResponse, error){
	var response serverPb.AutoPay
	var responses []*serverPb.AutoPay
	var created time.Time
	var updated time.Time
	rows, err := s.db.Query(s.ctx, `
	SELECT * FROM auto_payments LIMIT 50 `,
	)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		rows.Scan(
			&response.PaymentId,
			&response.Name,
			&response.Number,
			&created,
			&updated,
		)
		response.TimeCreated = timestamppb.New(created)
		response.TimeUpdated = timestamppb.New(updated)
		responses = append(responses, &response)
	}
	if rows.Err() != nil{
		log.Println(rows.Err())
		return nil, rows.Err()
	}
	return &serverPb.AllResponse{Items: responses}, nil
}

func (s*Server) Create(ctx context.Context, request *serverPb.CreateRequest) (*emptypb.Empty, error){
	_, err := s.db.Exec(ctx, `
	INSERT INTO auto_payments VALUES (DEFAULT, $1, $2)
	`, request.Name, request.Number)
	if err != nil{
		log.Println(err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
func (s*Server) DeleteById(ctx context.Context, request *serverPb.DeleteRequest) (*serverPb.Response, error){
	res, err := s.db.Exec(ctx, `
	DELETE FROM auto_payments WHERE payment_id=$1
	`, request.PaymentId)
	if err != nil{
		log.Println(err)
		return &serverPb.Response{Affected: 0}, err
	}
	return &serverPb.Response{Affected: res.RowsAffected()}, nil
}

func (s*Server) UpdateById(ctx context.Context,  request *serverPb.UpdateRequest) (*serverPb.Response, error){
	var affected int64
	if request.Name != ""{
		result, err := s.db.Exec(ctx, `UPDATE auto_payments SET name=$1, timeUpdated=$3 WHERE payment_id=$2`,
									request.Name, request.PaymentId, time.Now().Format(time.RFC3339))
		if err != nil{
			log.Println(err)
			return nil, err
		}
		affected = result.RowsAffected()
	}
	if request.Number != ""{
		result, err := s.db.Exec(ctx, `UPDATE auto_payments SET number=$1, timeUpdated=$3 WHERE payment_id=$2`,
									request.Number, request.PaymentId, time.Now().Format(time.RFC3339))
		if err != nil{
			log.Println(err)
			return nil, err
		}
		affected = result.RowsAffected()
	}
	// в affected 2 раза будет записано одно и то же число, при срабатывании любого if-а
	return &serverPb.Response{Affected: affected}, nil
}

