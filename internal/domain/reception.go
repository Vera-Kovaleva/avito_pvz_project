package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var _ ReceptionsInterface = (*ReceptionService)(nil)

var (
	errReception                         = errors.New("reception service error")
	ErrAvitoServiceReceptionInvalidPVZID = errors.Join(
		errReception,
		errors.New("invalid pvz id"),
	)
	errAvitoServiceCreateReception = errors.Join(
		errReception,
		errors.New("create failed error"),
	)
	ErrAvitoServiceCreateReception = errors.Join(
		errAvitoServiceCreateReception,
		errors.New("create failed"),
	)
	ErrAvitoServiceCreateReceptionFindActive = errors.Join(
		errAvitoServiceCreateReception,
		errors.New("find active failed"),
	)
	errAvitoServiceCloseReception           = errors.Join(errReception, errors.New("close failed"))
	ErrAvitoServiceCloseReceptionFindActive = errors.Join(
		errAvitoServiceCloseReception,
		errors.New("find active failed"),
	)
	ErrAvitoServiceCloseReception = errors.Join(
		errAvitoServiceCloseReception,
		errors.New("close failed"),
	)

	errProduct                         = errors.New("products service error")
	ErrAvitoServiceProductInvalidPVZID = errors.Join(errProduct, errors.New("invalid pvz id"))
	errAvitoServiceCreateProduct       = errors.Join(
		errProduct,
		errors.New("create product failed"),
	)
	ErrAvitoServiceCreateProduct = errors.Join(
		errAvitoServiceCreateProduct,
		errors.New("create product failed"),
	)
	ErrAvitoServiceCreateProductFindActive = errors.Join(
		errAvitoServiceCreateProduct,
		errors.New("find active failed"),
	)
	errAvitoServiceDeleteProduct = errors.Join(
		errProduct,
		errors.New("delete product failed"),
	)
	ErrAvitoServiceDeleteProduct = errors.Join(
		errAvitoServiceDeleteProduct,
		errors.New("delete product failed"),
	)
	ErrAvitoServiceDeleteProductFindActive = errors.Join(
		errAvitoServiceDeleteProduct,
		errors.New("find active failed"),
	)
)

type ReceptionService struct {
	provider      ConnectionProvider
	receptionRepo ReceptionsRepository
	productRepo   ProductsRepository
}

func NewReceptionService(
	provider ConnectionProvider,
	receptionRepo ReceptionsRepository,
	productRepo ProductsRepository,
) *ReceptionService {
	return &ReceptionService{
		provider:      provider,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}
}

func validPVZID(id PVZID) error {
	if uuid.UUID(id) == uuid.Nil || uuid.Validate(id.String()) != nil {
		return errors.New("uuid is not valid")
	}
	return nil
}

func (s *ReceptionService) Create(
	ctx context.Context,
	authUser AuthenticatedUser,
	pvzID PVZID,
) (Reception, error) {
	//* Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» может инициировать приём товара.
	var reception Reception

	errValidID := validPVZID(pvzID)
	if errValidID != nil {
		return reception, errors.Join(errValidID, ErrAvitoServiceReceptionInvalidPVZID)
	}

	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, c Connection) error {
		reception = Reception{
			ID:    ReceptionID(uuid.New()),
			PVZID: pvzID,
		}
		return s.receptionRepo.Create(ctx, c, reception)
	})
	if err != nil {
		return reception, errors.Join(ErrAvitoServiceCreateReception, err)
	}
	err = s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		reception, err = s.receptionRepo.FindActive(ctx, c, pvzID)
		return err
	})
	if err != nil {
		return reception, errors.Join(ErrAvitoServiceCreateReceptionFindActive, err)
	}
	return reception, nil
}

func (s *ReceptionService) Close(
	ctx context.Context,
	authUser AuthenticatedUser,
	pvzID PVZID,
) (Reception, error) {
	var reception Reception

	errValidID := validPVZID(pvzID)
	if errValidID != nil {
		return reception, errors.Join(errValidID, ErrAvitoServiceReceptionInvalidPVZID)
	}

	var errFindActive error
	err := s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		reception, errFindActive = s.receptionRepo.FindActive(ctx, c, pvzID)
		return errFindActive
	})
	if err != nil {
		return reception, errors.Join(ErrAvitoServiceCloseReceptionFindActive, err)
	}
	err = s.provider.ExecuteTx(ctx, func(ctx context.Context, c Connection) error {
		return s.receptionRepo.Close(ctx, c, reception.ID)
	})
	if err != nil {
		return reception, errors.Join(ErrAvitoServiceCloseReception, err)
	}
	return reception, nil
}

func (s *ReceptionService) CreateProduct(
	ctx context.Context,
	authUser AuthenticatedUser,
	pvzID PVZID,
	productType ProductType,
) (Product, error) {
	/* * Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» может добавлять товары после его осмотра.
	 * Если же нет новой незакрытой приёмки товаров, то в таком случае должна возвращаться ошибка, и товар не должен добавляться в систему.*/
	var product Product

	errValidID := validPVZID(pvzID)
	if errValidID != nil {
		return product, errors.Join(errValidID, ErrAvitoServiceProductInvalidPVZID)
	}

	var errFindActive error
	var reception Reception

	err := s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		reception, errFindActive = s.receptionRepo.FindActive(ctx, c, pvzID)
		return errFindActive
	})
	if err != nil {
		return product, errors.Join(ErrAvitoServiceCreateProductFindActive, err)
	}
	err = s.provider.ExecuteTx(ctx, func(ctx context.Context, c Connection) error {
		product = Product{
			ID:          ProductID(uuid.New()),
			ReceptionID: reception.ID,
			Type:        productType,
			CreatedAt:   time.Now(),
		}

		return s.productRepo.Create(ctx, c, product)
	})
	if err != nil {
		return product, errors.Join(ErrAvitoServiceCreateProduct, err)
	}
	return product, nil
}

func (s *ReceptionService) DeleteLastProduct(
	ctx context.Context,
	authUser AuthenticatedUser,
	pvzID PVZID,
) error {
	/* * Только авторизованный пользователь системы с ролью «сотрудник ПВЗ» может удалять товары, которые были добавленыв рамках текущей приёмки на ПВЗ.*/
	errValidID := validPVZID(pvzID)
	if errValidID != nil {
		return errors.Join(errValidID, ErrAvitoServiceProductInvalidPVZID)
	}

	var errFindActive error
	var reception Reception

	err := s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		reception, errFindActive = s.receptionRepo.FindActive(ctx, c, pvzID)
		return errFindActive
	})
	if err != nil {
		return errors.Join(ErrAvitoServiceDeleteProductFindActive, err)
	}

	err = s.provider.ExecuteTx(ctx, func(ctx context.Context, c Connection) error {
		return s.productRepo.DeleteLast(ctx, c, reception.ID)
	})
	if err != nil {
		return errors.Join(ErrAvitoServiceDeleteProduct, err)
	}
	return nil
}
