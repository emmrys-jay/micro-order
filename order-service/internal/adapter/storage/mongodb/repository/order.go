package repository

import (
	"context"
	"time"

	"order-service/internal/adapter/config"
	"order-service/internal/adapter/storage/mongodb"
	"order-service/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/**
 * OrderRepository implements port.OrderRepository interface
 * and provides an access to the MongoDB database
 */
type OrderRepository struct {
	ordersCol *mongo.Collection
	itemsCol  *mongo.Collection
}

// NewOrderRepository creates a new order repository instance
func NewOrderRepository(db *mongodb.DB) *OrderRepository {
	return &OrderRepository{
		ordersCol: db.Client.Database(config.GetConfig().Database.Name).Collection("orders"),
		itemsCol:  db.Client.Database(config.GetConfig().Database.Name).Collection("orderItems"),
	}
}

func (or *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError) {
	session, err := or.ordersCol.Database().Client().StartSession()
	if err != nil {
		return nil, domain.NewInternalCError("error starting session: " + err.Error())
	}
	defer session.EndSession(ctx)

	orderItems := make([]domain.OrderItem, len(order.OrderItems))
	copy(orderItems, order.OrderItems)
	order.OrderItems = nil

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		order.ID = primitive.NewObjectID()
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()

		// Insert order first
		_, err := or.ordersCol.InsertOne(sc, order)
		if err != nil {
			session.AbortTransaction(sc)
			return domain.NewInternalCError("error inserting order: " + err.Error())
		}

		// Insert order items
		items := make([]interface{}, len(orderItems))
		for i := range orderItems {
			orderItems[i].OrderID = order.ID
			orderItems[i].ID = primitive.NewObjectID()
			orderItems[i].CreatedAt = time.Now()
			orderItems[i].UpdatedAt = time.Now()
			items[i] = orderItems[i]
		}

		_, err = or.itemsCol.InsertMany(sc, items)
		if err != nil {
			session.AbortTransaction(sc)
			return domain.NewInternalCError("error inserting order items: " + err.Error())
		}

		if err := session.CommitTransaction(sc); err != nil {
			return domain.NewInternalCError("error committing transaction: " + err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, domain.NewInternalCError("error creating order: " + err.Error())
	}

	order.OrderItems = orderItems
	return order, nil
}

func (or *OrderRepository) GetOrder(ctx context.Context, id primitive.ObjectID) (*domain.Order, domain.CError) {
	var order domain.Order

	err := or.ordersCol.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError("error finding order: " + err.Error())
	}

	cursor, err := or.itemsCol.Find(ctx, bson.M{"order_id": id})
	if err != nil {
		return nil, domain.NewInternalCError("error finding order items: " + err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item domain.OrderItem
		if err := cursor.Decode(&item); err != nil {
			return nil, domain.NewInternalCError("error decoding order item: " + err.Error())
		}
		order.OrderItems = append(order.OrderItems, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, domain.NewInternalCError("error iterating order items: " + err.Error())
	}

	return &order, nil
}

func (or *OrderRepository) ListOrders(ctx context.Context, userId primitive.ObjectID) ([]domain.Order, domain.CError) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userId}}},
		{primitive.E{Key: "$lookup", Value: bson.M{
			"from":         "orderItems",
			"localField":   "_id",
			"foreignField": "order_id",
			"as":           "order_items",
		}}},
	}

	cursor, err := or.ordersCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, domain.NewInternalCError("error aggregating orders: " + err.Error())
	}
	defer cursor.Close(ctx)

	var orders = []domain.Order{}
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, domain.NewInternalCError("error decoding order: " + err.Error())
		}
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		return nil, domain.NewInternalCError("error iterating orders: " + err.Error())
	}

	return orders, nil
}

func (or *OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) (*domain.Order, domain.CError) {
	order.UpdatedAt = time.Now()

	filter := bson.M{"_id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"status":     order.Status,
			"updated_at": order.UpdatedAt,
		},
	}

	_, err := or.ordersCol.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrDataNotFound
		}
		return nil, domain.NewInternalCError("error updating order: " + err.Error())
	}

	return order, nil
}

func (or *OrderRepository) UpdateOrderProducts(ctx context.Context, prod *domain.ProductUpdateFromQueue) (int64, domain.CError) {
	oId, err := primitive.ObjectIDFromHex(prod.Id)
	if err != nil {
		return -1, domain.NewBadRequestCError("Error parsing product id received: " + err.Error())
	}

	filter := bson.M{"product_id": oId}
	update := bson.M{
		"$set": bson.M{
			"product_name": prod.Name,
		},
	}

	res, err := or.itemsCol.UpdateMany(ctx, filter, update)
	if err != nil {
		return -1, domain.NewInternalCError("error updating order: " + err.Error())
	}

	return res.MatchedCount, nil
}
