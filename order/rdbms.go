package order

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RdbmsAccess interface {
	Insert(or OrderRequest) error
}

type rdbmsAccess struct {
	pool *pgxpool.Pool
}

func (o *rdbmsAccess) Insert(or OrderRequest) error{
	_, err := o.pool.Exec(context.Background(),
		"INSERT INTO orders(id, customerId, productId, orderDesc) VALUES($1, $2, $3, $4)",
		or.Id,or.CustomerId,or.ProductId,or.OrderDesc)

	if err != nil  {
		switch v := (err).(type) {
		case *pgconn.PgError:
			if v.Code == "23505" {
				return nil
			}
			return v
		default:
			fmt.Print("Error adding order to DB", err)
			return err
		}
	}
	return nil
}

func NewRdbmsAccess(pool *pgxpool.Pool) RdbmsAccess {
	return &rdbmsAccess{pool}
}





