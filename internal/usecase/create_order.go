package usecase

import (
	"GoExpertPostGrad-Orders-Challenge/internal/entity"
	"GoExpertPostGrad-Orders-Challenge/pkg/events"
)

type OrderInputDTO struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
	EventDispatcher   events.EventDispatcherInterface
}

func NewCreateOrderUseCase(orderRepository entity.OrderRepositoryInterface, orderCreatedEvent events.EventInterface, eventDispatcher events.EventDispatcherInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository:   orderRepository,
		OrderCreatedEvent: orderCreatedEvent,
		EventDispatcher:   eventDispatcher,
	}
}

func (uc *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	// Create order
	order, err := toEntityOrder(input)
	if err != nil {
		return OrderOutputDTO{}, err
	}
	err = order.CalculateFinalPrice()
	if err != nil {
		return OrderOutputDTO{}, err
	}

	// Save order
	if err := uc.OrderRepository.Save(order); err != nil {
		return OrderOutputDTO{}, err
	}

	// Create and Dispatch Order Created event
	dto, err := uc.dispatchOrderCreatedEvent(order)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	return dto, nil
}

func toEntityOrder(input OrderInputDTO) (*entity.Order, error) {
	order, err := entity.NewOrder(input.Name, input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *CreateOrderUseCase) dispatchOrderCreatedEvent(order *entity.Order) (OrderOutputDTO, error) {
	dto := &OrderOutputDTO{
		ID:         order.ID,
		Name:       order.Name,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
	uc.OrderCreatedEvent.SetPayload(dto)
	err := uc.EventDispatcher.Dispatch(uc.OrderCreatedEvent)
	if err != nil {
		return OrderOutputDTO{}, err
	}
	return *dto, nil
}
