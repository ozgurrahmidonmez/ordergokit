package order

type ServiceMiddleware func(service OrderService) OrderService

type OrderService interface {
	Add(order OrderRequest) error
}

type orderService struct {
	rdbmsAccess RdbmsAccess
	sqsAccess SqsAccess
}

func (s *orderService) Add(order OrderRequest) error {
	if err := s.rdbmsAccess.Insert(order) ; err != nil {
		return err
	}
	if err := s.sqsAccess.Enqueue(order) ; err != nil {
		return err
	}
	return nil
}

func NewOrderService(rdbmsAccess RdbmsAccess,sqsAccess SqsAccess) OrderService {
	return &orderService{rdbmsAccess,sqsAccess}
}


