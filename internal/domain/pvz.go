package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var _ PVZsInterface = (*PVZService)(nil)

var (
	errPVZ                   = errors.New("pvz service error")
	ErrAvitoServiceCreatePVZ = errors.Join(
		errPVZ,
		errors.New("create pvz failed"),
	)
	ErrAvitoServiceFindAllPVZ = errors.Join(
		errPVZ,
		errors.New("find all last pvz failed"),
	)
	errAvitoServiceFindPVZReceptionProducts = errors.Join(
		errPVZ,
		errors.New("search failed"),
	)
	ErrAvitoServiceFindPVZReceptionProductsSearchProducts = errors.Join(
		errAvitoServiceFindPVZReceptionProducts,
		errors.New("search products failed"),
	)
	ErrAvitoServiceFindPVZReceptionProductsSearchPVZs = errors.Join(
		errAvitoServiceFindPVZReceptionProducts,
		errors.New("search pvzs failed"),
	)
	ErrAvitoServiceFindPVZReceptionProductsSearchReceptions = errors.Join(
		errAvitoServiceFindPVZReceptionProducts,
		errors.New("search reseptions failed"),
	)
)

type PVZService struct {
	provider      ConnectionProvider
	pvzRepo       PVZsRepository
	productRepo   ProductsRepository
	receptionRepo ReceptionsRepository
}

func NewPVZService(
	provider ConnectionProvider,
	pvzRepo PVZsRepository,
	productRepo ProductsRepository,
	receptionRepo ReceptionsRepository,
) *PVZService {
	return &PVZService{
		provider:      provider,
		pvzRepo:       pvzRepo,
		productRepo:   productRepo,
		receptionRepo: receptionRepo,
	}
}

func (s *PVZService) Create(
	ctx context.Context,
	authUser AuthenticatedUser,
	pvzCity PVZCity,
) (PVZ, error) {
	//* Только пользователь с ролью «модератор» может завести ПВЗ в системе.
	var pvz PVZ
	err := s.provider.ExecuteTx(ctx, func(ctx context.Context, c Connection) error {
		pvz = PVZ{
			ID:           PVZID(uuid.New()),
			City:         pvzCity,
			RegisteredAt: time.Now(),
		}

		return s.pvzRepo.Create(ctx, c, pvz)
	})
	if err != nil {
		return pvz, errors.Join(ErrAvitoServiceCreatePVZ, err)
	}
	return pvz, nil
}

func (s *PVZService) FindAll(ctx context.Context) ([]PVZ, error) {
	var pvzsAll []PVZ
	err := s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		var findAllError error
		pvzsAll, findAllError = s.pvzRepo.FindAll(ctx, c)
		return findAllError
	})
	if err != nil {
		return nil, errors.Join(ErrAvitoServiceFindAllPVZ, err)
	}
	return pvzsAll, nil
}

func (s *PVZService) FindPVZReceptionProducts(
	ctx context.Context,
	authUser AuthenticatedUser,
	from *time.Time,
	to *time.Time,
	page *int,
	limit *int,
) ([]PVZReceptionsProducts, error) {
	var result []PVZReceptionsProducts

	var products []Product
	var searchError error
	err := s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		products, searchError = s.productRepo.Search(ctx, c, from, to, page, limit)
		return searchError
	})
	if err != nil {
		return result, errors.Join(err, ErrAvitoServiceFindPVZReceptionProductsSearchProducts)
	}

	var receptionIDs []PVZID
	for _, product := range products {
		receptionIDs = append(receptionIDs, product.ReceptionID)
	}
	var receptions []Reception
	var findReceptionError error
	err = s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		receptions, findReceptionError = s.receptionRepo.FindByIDs(ctx, c, receptionIDs)
		return findReceptionError
	})
	if err != nil {
		return result, errors.Join(err, ErrAvitoServiceFindPVZReceptionProductsSearchReceptions)
	}

	var pvzIDs []PVZID
	for _, reception := range receptions {
		pvzIDs = append(pvzIDs, reception.PVZID)
	}

	var pvzs []PVZ
	var findPVZError error
	err = s.provider.Execute(ctx, func(ctx context.Context, c Connection) error {
		pvzs, findPVZError = s.pvzRepo.FindByIDs(ctx, c, pvzIDs)
		return findPVZError
	})
	if err != nil {
		return result, errors.Join(err, ErrAvitoServiceFindPVZReceptionProductsSearchPVZs)
	}

	return Builder(products, receptions, pvzs), nil
}

func Builder(products []Product, receptions []Reception, pvzs []PVZ) []PVZReceptionsProducts {
	productsToReceptionsByID := make(map[ReceptionID][]Product)
	for _, product := range products {
		productsToReceptionsByID[product.ReceptionID] = append(
			productsToReceptionsByID[product.ReceptionID],
			product,
		)
	}

	receptionsToPVZsByID := make(map[PVZID][]ReceptionsProducts)
	for _, reception := range receptions {
		receptionsToPVZsByID[reception.PVZID] = append(
			receptionsToPVZsByID[reception.PVZID],
			ReceptionsProducts{
				Reception: reception,
				Products:  productsToReceptionsByID[reception.ID],
			},
		)
	}

	var result []PVZReceptionsProducts
	for _, pvz := range pvzs {
		result = append(result, PVZReceptionsProducts{
			PVZ:        pvz,
			Receptions: receptionsToPVZsByID[pvz.ID],
		})
	}

	return result
}
