package config

const (
	ProductCreateTopic     = "product_create_topic"
	ProductDeleteTopic     = "product_delete_topic"
	PriceTimeStampAddTopic = "price_time_stamp_add_topic"
)

var (
	Brokers = []string{"localhost:19091", "localhost:29091", "localhost:39091"}
)

type ProductCreateRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProductDeleteRequest struct {
	Code string `json:"code"`
}

type PriceTimeStampAddRequest struct {
	Code  string  `json:"code"`
	Ts    int64   `json:"ts"`
	Price float64 `json:"price"`
}

type Response struct {
	Ok bool `json:"ok"`
}
