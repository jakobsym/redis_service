package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/jakobsym/redis_service/model"
)

var (
	ErrNotExist = errors.New("order does not exist")
)

type RedisRepo struct {
	Client *redis.Client
}

type FindAllPage struct {
	Size   uint
	Offset uint
	//maybe cursor?
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

// because Redis is key/value store we need to generate keys for our data
// which well perform grouping for us
func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

// idiomatic to have context as 1st parameter
func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	// Transaction
	// groups commands making them Atomic
	// I.E: All will work, or none at all
	// we can avoid partial states
	// as we may insert into db succesfully, but fail to add to the order set

	// Marshall returns byte array
	client := r.Client.WithContext(ctx)
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}
	key := orderIDKey(order.OrderID)

	// start a new transaction in a pipeline
	txn := client.TxPipeline()

	//.Set() will overwrite data
	// .SetNX() does not
	res := txn.SetNX(key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard() // discard any potential changes
		return fmt.Errorf("failed to set: %w", err)
	}
	// store Ids in a set
	if err := txn.SAdd("orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	// commit our commands using Exec(); enables data gurantee
	// avoids partial states
	if _, err := txn.Exec(); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	var order model.Order

	key := orderIDKey(id)
	client := r.Client.WithContext(ctx)
	val, err := client.Get(key).Result()

	// return custom error, if value at given id is Nil in DB
	//return empty Orders, as cheaper on memory
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	// store order into a model.Order
	err = json.Unmarshal([]byte(val), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)
	client := r.Client.WithContext(ctx)
	txn := client.TxPipeline()

	err := txn.Del(key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("delete oder: %w", err)
	}

	if err := txn.SRem("orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from orders set: %w", err)
	}

	// commit our commands using Exec(); enables data gurantee
	// avoids partial states
	if _, err := txn.Exec(); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

// only update existing records
func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("unable to decode json order: %w", err)
	}
	client := r.Client.WithContext(ctx)
	key := orderIDKey(order.OrderID)
	// SetXX() == only set if that value already exists
	err = client.SetXX(key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

// bad practice to retrieve EVERY record from redis
// this can be expensive on memory thus...
// we perform 'pagination'
// which breaks results down into seperate pages
// client can request more data everytime without having to get all of it out of the db

// FindResult as return defines orders, and next cursor so caller knows where to continue paging
func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	client := r.Client.WithContext(ctx)
	// * := everything in set
	res := client.SScan("orders", uint64(page.Offset), "*", int64(page.Size))
	keys, cursor, err := res.Result()

	// records returned in random order
	// side effect of using a set
	// redis has a sorted set, if you need data sorted
	if err != nil {
		return FindResult{}, fmt.Errorf("fialed to get order ids: %w", err)
	}

	// check key size before doing extra work
	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	// pull values from set of keys
	xs, err := client.MGet(keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get values from set: %w", err)
	}
	// make slice w/ same len of resulting slice above `xs`
	orders := make([]model.Order, len(xs))

	// cast e/a element into a string
	// unmarshall into order struct
	// store into order slice
	for i, x := range xs {
		x := x.(string)
		var order model.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}
		orders[i] = order
	}

	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
