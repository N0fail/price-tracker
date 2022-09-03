package config

const (
	ProductCreateTopic     = "product_create_topic"
	ProductDeleteTopic     = "product_delete_topic"
	PriceTimeStampAddTopic = "price_time_stamp_add_topic"
	ProductListTopic       = "product_list_topic"
	PriceHistoryTopic      = "price_history_topic"
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

type ProductListRequest struct {
	PageNumber     uint32 `json:"page_number"`
	ResultsPerPage uint32 `json:"results_per_page"`
	OrderBy        int32  `json:"order_by"`
}

type PriceHistoryRequest struct {
	Code string `json:"code"`
}

type Response struct {
	Ok bool `json:"ok"`
}
