package store

import (
	"context"
	"fmt"
	"time"

	"rate-limiter-service/config"
	pb "rate-limiter-service/proto/ratelimiter"

	"github.com/RussellLuo/slidingwindow"
	"github.com/hashicorp/consul/api"
	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/proto"
)

type RateLimiterStore struct {
	cli *api.Client
}

func New() (*RateLimiterStore, error) {
	cfg := config.GetConfig()
	db := cfg.DB
	dbport := cfg.DBPort

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &RateLimiterStore{
		cli: client,
	}, nil
}

func (rs *RateLimiterStore) Get(ctx context.Context, id string) (*pb.RateLimiter, error) {
	kv := rs.cli.KV()

	key := fmt.Sprintf("rateLimiters/%s", id)

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	rateLimiter := &pb.RateLimiter{}
	err = proto.Unmarshal(pair.Value, rateLimiter)
	if err != nil {
		return nil, err
	}

	return rateLimiter, nil
}

func (rs *RateLimiterStore) GetAll(ctx context.Context) (*pb.ListOfRateLimiters, error) {
	kv := rs.cli.KV()
	data, _, err := kv.List("rateLimiters", nil)
	if err != nil {
		return nil, err
	}

	rateLimiters := []*pb.RateLimiter{}
	for _, pair := range data {
		rateLimiter := &pb.RateLimiter{}
		err = proto.Unmarshal(pair.Value, rateLimiter)
		if err != nil {
			return nil, err
		}
		rateLimiters = append(rateLimiters, rateLimiter)
	}

	return &pb.ListOfRateLimiters{
		Limiters: rateLimiters,
	}, nil
}

func (rs *RateLimiterStore) Create(ctx context.Context, limiter *pb.CreateRateLimiterRequest) (*pb.RateLimiter, error) {
	kv := rs.cli.KV()

	id := ""
	if limiter.RateLimiter.Name == "sistem" {
		id = fmt.Sprintf("%s", limiter.RateLimiter.Name)
	} else {
		id = fmt.Sprintf("%s-%s", limiter.RateLimiter.Name, limiter.RateLimiter.UserName)
	}
	limiter.RateLimiter.Id = id

	key := fmt.Sprintf("rateLimiters/%s", id)

	data, err := proto.Marshal(limiter.RateLimiter)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return limiter.RateLimiter, nil

}

func (rs *RateLimiterStore) Update(ctx context.Context, limiter *pb.UpdateRateLimiterRequest) (*pb.RateLimiter, error) {
	kv := rs.cli.KV()

	id := limiter.RateLimiter.Id

	key := fmt.Sprintf("rateLimiters/%s", id)

	data, err := proto.Marshal(limiter.RateLimiter)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return limiter.RateLimiter, nil
}

func (rs *RateLimiterStore) Delete(ctx context.Context, id string) (*pb.DeleteRateLimiterResponse, error) {
	kv := rs.cli.KV()

	key := fmt.Sprintf("rateLimiters/%s", id)

	_, err := kv.Delete(key, nil)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteRateLimiterResponse{
		Deleted: true,
	}, nil
}

func (rs *RateLimiterStore) IsRequestAllowed(ctx context.Context, id string) (*pb.AllowResponse, error) {
	rateLimiter, err := rs.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	switch rateLimiter.Type {
	case "tokenBucket":
		limiter := rate.NewLimiter(rate.Limit(rateLimiter.ReqPerSec), int(rateLimiter.Burst))
		allowed := limiter.Allow()
		return &pb.AllowResponse{
			Allowed: allowed,
		}, nil
	case "leakyBucket":
		limiter := ratelimit.New(int(rateLimiter.ReqPerSec))
		limiter.Take()
		allowed := true
		return &pb.AllowResponse{
			Allowed: allowed,
		}, nil
	case "slidingWindow":
		limiter, _ := slidingwindow.NewLimiter(time.Second, rateLimiter.ReqPerSec, func() (slidingwindow.Window, slidingwindow.StopFunc) {
			// NewLocalWindow returns an empty stop function, so it's
			// unnecessary to call it later.
			return slidingwindow.NewLocalWindow()
		})
		allowed := limiter.Allow()
		return &pb.AllowResponse{
			Allowed: allowed,
		}, nil

	default:
		return nil, fmt.Errorf("algorithm is not implemented")
	}

}
