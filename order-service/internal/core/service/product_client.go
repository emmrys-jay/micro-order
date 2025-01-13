package service

import (
	"context"
	"fmt"
	"order-service/internal/adapter/config"
	"order-service/internal/adapter/logger"
	"order-service/internal/core/domain"
	"order-service/internal/core/service/product"
	"order-service/internal/core/util"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cacheDuraction time.Duration

func newProductClient(conf *config.DiscoveryConfiguration) (*grpc.ClientConn, product.ProductClient, error) {
	conn, err := grpc.NewClient(conf.ProductUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a new connection: %w", err)
	}

	productClient := product.NewProductClient(conn)
	return conn, productClient, nil
}

func (os *OrderService) GetProduct(ctx context.Context, productId primitive.ObjectID) (*product.ProductResponse, domain.CError) {
	log := logger.FromCtx(ctx)

	log.Info("Checking if product exists in cache", zap.String("product_id", productId.Hex()))
	cacheKey := util.GenerateCacheKey("product", productId.Hex())
	cachedProd, err := os.cache.Get(ctx, cacheKey)
	if err == nil {
		log.Info("Product with id found in cache", zap.String("product_id", productId.Hex()))

		var product product.ProductResponse
		err := util.Deserialize(cachedProd, &product)
		if err != nil {
			log.Error("Error deserializing found product in cache", zap.Error(err))
			return nil, domain.ErrInternal
		}

		return &product, nil
	}
	log.Info("Product not found in cache")

	grpcConn, grpcClient, err := newProductClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating product client", zap.Error(err))
		return nil, domain.ErrInternal
	}
	defer grpcConn.Close()

	log.Info("Created new grpc product client")
	log.Info("Making request to fetch product", zap.String("product_id", productId.Hex()))

	product, err := grpcClient.Get(context.Background(), &product.ProductRequest{ProductId: productId.Hex()})
	if err != nil {
		log.Error("Error getting product", zap.Error(err))
		return nil, domain.ErrInternal
	}

	log.Info("Successfully fetched user")

	log.Info("Saving returned product to cache")
	serialProd, err := util.Serialize(product)
	if err != nil {
		log.Error("Error serializing product", zap.Error(err))
		return nil, domain.ErrInternal
	}

	err = os.cache.Set(ctx, cacheKey, serialProd, os.cacheTtl)
	if err != nil {
		log.Error("Error saving returned product to cache", zap.Error(err))
		return nil, domain.ErrInternal
	}

	return product, nil
}

func (os *OrderService) GetProductsByIDs(ctx context.Context, validProductIDs []string) (products []*product.ProductResponse, cerr domain.CError) {
	log := logger.FromCtx(ctx)

	cacheKeys := make([]string, len(validProductIDs))
	for i, v := range validProductIDs {
		cacheKeys[i] = util.GenerateCacheKey("product", v)
	}
	log.Info("Checking if products exist in cache")

	results, err := os.cache.MGet(ctx, cacheKeys)
	if err != nil {
		log.Error("Error fetching products from cache", zap.Error(err))
		return nil, domain.ErrInternal
	}

	fetchedProducts, nilIndexes, err := os.deseralizeProducts(results)
	if err != nil {
		log.Error("Error deserializing cache products", zap.Error(err))
		return nil, domain.ErrInternal
	}

	// If all requested products are in the cache
	if len(nilIndexes) == 0 {
		log.Info("Successfully fetched all products from cache")
		return fetchedProducts, nil
	}

	log.Info(fmt.Sprintf("Fetched %v out of %v products from cache", len(fetchedProducts)-len(nilIndexes), len(fetchedProducts)))

	// Get only the indexes that were nil when returned from cache
	idsToFetch := make([]string, len(nilIndexes))
	for i, v := range nilIndexes {
		idsToFetch[i] = validProductIDs[v]
	}

	pConn, pClient, err := newProductClient(&config.GetConfig().Discovery)
	if err != nil {
		log.Error("Error creating product client", zap.Error(err))
		return nil, domain.ErrInternal
	}
	defer pConn.Close()

	log.Info("Created new grpc product client")
	log.Info("Making request to fetch many products", zap.Any("product_ids", idsToFetch))

	prods, err := pClient.GetMany(context.Background(), &product.ProductsRequest{ProductIds: idsToFetch})
	if err != nil {
		log.Error("Error fetching products", zap.Error(err))
		return nil, domain.ErrInternal
	}

	log.Info("Successfully fetched products")

	grpcFetchedProducts := make(map[string]*product.ProductResponse)
	for _, v := range prods.Products {
		grpcFetchedProducts[v.Id] = v
	}

	// Run through the nil indexes from cache and get their ids from validProductIDs
	// Then get the returned products and store them in the empty indexes in fetchedProducts
	for _, v := range nilIndexes {
		fetchedProducts[v] = grpcFetchedProducts[validProductIDs[v]]
	}

	// Store the retrieved products in cache
	log.Info("Storing the retrieved products to cache")
	err = os.saveProductsToCache(grpcFetchedProducts)
	if err != nil {
		log.Error("Retrived grpc products could not be saved to cache", zap.Error(err))
	} else {
		log.Info("Successfully saved retrieved grpc products to cache")
	}

	return fetchedProducts, nil
}

func (os *OrderService) deseralizeProducts(products [][]byte) (validProducts []*product.ProductResponse, nilIndex []int, err error) {
	for i, v := range products {
		var p product.ProductResponse

		if v == nil {
			nilIndex = append(nilIndex, i)
			validProducts = append(validProducts, nil)
			continue
		}

		err := util.Deserialize(v, &p)
		if err != nil {
			return nil, nil, err
		}
		validProducts = append(validProducts, &p)
	}

	return
}

func (os *OrderService) saveProductsToCache(products map[string]*product.ProductResponse) error {
	pMap := make(map[string][]byte)
	for k, v := range products {
		cacheKey := util.GenerateCacheKey("product", k)
		cacheValue, err := util.Serialize(v)
		if err != nil {
			return err
		}

		pMap[cacheKey] = cacheValue
	}

	return os.cache.MSet(context.Background(), pMap, os.cacheTtl)
}
