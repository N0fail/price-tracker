package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	internalConfig "gitlab.ozon.dev/N0fail/price-tracker/internal/config"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/kafka/config"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/kafka/counter"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/api"
	"gitlab.ozon.dev/N0fail/price-tracker/internal/pkg/core/product/models"
	pb "gitlab.ozon.dev/N0fail/price-tracker/pkg/api"
	"log"
	_ "net/http/pprof"
	"time"
)

type IncomeHandler struct {
	productApi      api.Interface
	requestsCounter *counter.Counter
	errorsCounter   *counter.Counter
	redisClient     *redis.Client
}

func (i *IncomeHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *IncomeHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

type ProductCreateHandler struct {
	IncomeHandler
}

func (h *ProductCreateHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.requestsCounter.Inc()
		var request config.ProductCreateRequest
		err := json.Unmarshal(msg.Value, &request)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("ProductCreateHandler ConsumeClaim Unmarshal error: %v", err.Error())
			continue
		}

		err = h.productApi.ProductCreate(context.Background(), models.Product{
			Code: request.Code,
			Name: request.Name,
		})
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("ProductCreateHandler ProductCreate error: %v", err.Error())
			return err
		}

		logrus.Infof("ProductCreateHandler success ProductCreate code: %v, name: %v", request.Code, request.Name)
	}
	return nil
}

type ProductDeleteHandler struct {
	IncomeHandler
}

func (h *ProductDeleteHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.requestsCounter.Inc()
		var request config.ProductDeleteRequest
		err := json.Unmarshal(msg.Value, &request)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("ProductDeleteHandler ConsumeClaim Unmarshal error: %v", err.Error())
			continue
		}

		err = h.productApi.ProductDelete(context.Background(), request.Code)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("ProductDeleteHandler ProductDelete error: %v", err.Error())
			return err
		}

		logrus.Infof("ProductDeleteHandler success ProductDelete code: %v", request.Code)
	}
	return nil
}

type PriceTimeStampAddHandler struct {
	IncomeHandler
}

func (h *PriceTimeStampAddHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.requestsCounter.Inc()
		var request config.PriceTimeStampAddRequest
		err := json.Unmarshal(msg.Value, &request)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("PriceTimeStampAddHandler ConsumeClaim Unmarshal error: %v", err.Error())
			continue
		}

		err = h.productApi.PriceTimeStampAdd(context.Background(), request.Code, models.PriceTimeStamp{
			Price: request.Price,
			Date:  time.Unix(request.Ts, 0).UTC(),
		})
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("PriceTimeStampAddHandler PriceTimeStampAdd error: %v", err.Error())
			return err
		}

		logrus.Infof("PriceTimeStampAddHandler success PriceTimeStampAdd code: %v, price: %v, date: %v", request.Code, request.Price, time.Unix(request.Ts, 0).Format("2 Jan 2006 15:04"))
	}
	return nil
}

type ProductListHandler struct {
	IncomeHandler
}

func (h *ProductListHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.requestsCounter.Inc()
		var request config.ProductListRequest
		err := json.Unmarshal(msg.Value, &request)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("PriceTimeStampAddHandler ConsumeClaim Unmarshal error: %v", err.Error())
			continue
		}

		var orderBy string
		if request.OrderBy == int32(pb.ProductListRequest_name) {
			orderBy = "name"
		} else {
			orderBy = "code"
		}

		productSnapShots, err := h.productApi.ProductList(context.Background(), request.PageNumber, request.ResultsPerPage, orderBy)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("PriceTimeStampAddHandler PriceTimeStampAdd error: %v", err.Error())
			return err
		}

		result := make([]*pb.ProductSnapShot, 0, len(productSnapShots))
		for _, productSnapShot := range productSnapShots {
			var priceTimeStamp *pb.PriceTimeStamp
			if !productSnapShot.LastPrice.IsEmpty() {
				priceTimeStamp = &pb.PriceTimeStamp{
					Price: productSnapShot.LastPrice.Price,
					Ts:    productSnapShot.LastPrice.Date.Unix(),
				}
			}

			result = append(result, &pb.ProductSnapShot{
				Code:           productSnapShot.Code,
				Name:           productSnapShot.Name,
				PriceTimeStamp: priceTimeStamp,
			})
		}

		message, err := json.Marshal(pb.ProductListResponse{
			ProductSnapShots: result,
		})

		if err != nil {
			logrus.Errorf("PriceHistoryHandler json.Marshal error: %v", err.Error())
			return err
		}
		h.redisClient.Publish(context.Background(), config.ProductListTopic, message)

		logrus.Infof("ProductListHandler success ProductList PageNumber: %v, ResultsPerPage: %v, orderBy: %v", request.PageNumber, request.ResultsPerPage, orderBy)
	}
	return nil
}

type PriceHistoryHandler struct {
	IncomeHandler
}

func (h *PriceHistoryHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.requestsCounter.Inc()
		var request config.PriceHistoryRequest
		err := json.Unmarshal(msg.Value, &request)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("PriceTimeStampAddHandler ConsumeClaim Unmarshal error: %v", err.Error())
			continue
		}

		priceHistory, err := h.productApi.PriceHistory(context.Background(), request.Code)
		if err != nil {
			h.errorsCounter.Inc()
			logrus.Errorf("PriceHistoryHandler PriceHistory error: %v", err.Error())
			return err
		}

		result := make([]*pb.PriceTimeStamp, 0, len(priceHistory))
		for _, priceTimeStamp := range priceHistory {
			result = append(result, &pb.PriceTimeStamp{
				Price: priceTimeStamp.Price,
				Ts:    priceTimeStamp.Date.Unix(),
			})
		}

		message, err := json.Marshal(pb.PriceHistoryResponse{
			PriceHistory: result,
		})

		if err != nil {
			logrus.Errorf("PriceHistoryHandler json.Marshal error: %v", err.Error())
			return err
		}
		h.redisClient.Publish(context.Background(), config.PriceHistoryTopic, message)

		logrus.Infof("PriceHistoryHandler success PriceHistory code: %v", request.Code)
	}
	return nil
}

func startHandler(ctx context.Context, topic string, handler sarama.ConsumerGroupHandler) {
	cfg := sarama.NewConfig()

	income, err := sarama.NewConsumerGroup(config.Brokers, topic, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		if err := income.Consume(ctx, []string{topic}, handler); err != nil {
			logrus.Infof("on consume: <%v>", err)
			time.Sleep(time.Second * 5)
		}
	}
}

func Run(productApi api.Interface) {
	ctx := context.Background()
	requestsCounter := counter.New("requestsCounter")
	errorsCounter := counter.New("errorsCounter")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost" + internalConfig.RedisPort,
		Password: "",
		DB:       0,
	})

	go startHandler(ctx, config.ProductCreateTopic, &ProductCreateHandler{
		IncomeHandler{
			productApi:      productApi,
			requestsCounter: requestsCounter,
			errorsCounter:   errorsCounter,
		},
	})

	go startHandler(ctx, config.ProductDeleteTopic, &ProductDeleteHandler{
		IncomeHandler{
			productApi:      productApi,
			requestsCounter: requestsCounter,
			errorsCounter:   errorsCounter,
		},
	})

	go startHandler(ctx, config.PriceTimeStampAddTopic, &PriceTimeStampAddHandler{
		IncomeHandler{
			productApi:      productApi,
			requestsCounter: requestsCounter,
			errorsCounter:   errorsCounter,
		},
	})

	go startHandler(ctx, config.ProductListTopic, &ProductListHandler{
		IncomeHandler{
			productApi:      productApi,
			requestsCounter: requestsCounter,
			errorsCounter:   errorsCounter,
			redisClient:     redisClient,
		},
	})

	startHandler(ctx, config.PriceHistoryTopic, &PriceHistoryHandler{
		IncomeHandler{
			productApi:      productApi,
			requestsCounter: requestsCounter,
			errorsCounter:   errorsCounter,
			redisClient:     redisClient,
		},
	})
}
