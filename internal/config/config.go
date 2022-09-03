package config

import "time"

const (
	CommandDelimeter   = "_"                // разделитель для аргументов комманд
	MinNameLength      = 4                  // минимальное число(байт) символов в имени продукта
	InvalidCodeSymbols = "\n"               // символы, которые запрещено использовать в коде продукта
	CacheValidTime     = 60                 // время в течении которого кэш считается актуальным
	DateFormat         = "2 Jan 2006 15:04" // формат времени для хранения дат
	GetUpdatesTimeout  = 60                 // таймаут на запрос получения обновлений чата
	RedisPort          = ":6379"            // порт, на котором работает редис
	GrpcPort           = ":8081"            // порт, на котором работает grpc сервер
	RESTPort           = ":8100"            // порт, на котором работает REST и сваггер
	RequestTimeout     = time.Second / 2    // время, которое выделяется на выполнение запроса, по истечении этого времени запрос прерывается

	DefaultResultsPerPage = 20     // кол-во результатов на на страницу по умолчанию (ф-я List)
	DefaultOrderBy        = "code" //поле сортировки по умолчанию (ф-я List)

	DbHost     = "localhost"
	DbPort     = 5432
	DbUser     = "user"
	DbPassword = "password"
	DbName     = "price-tracker"

	DbMaxConnIdleTime = time.Minute
	DbMaxConnLifetime = time.Hour
	DbMinConns        = 2
	DbMaxConns        = 4

	MemcachedAddress = "localhost:11211"
)
