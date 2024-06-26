package graph

import "GoExpertPostGrad-Orders-Challenge/internal/usecase"

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	OrderCreationUseCase usecase.OrderCreationUseCase
	OrderListingUseCase  usecase.OrderListingUseCase
}
