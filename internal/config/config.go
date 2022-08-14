package config

import "time"

const (
	CommandDelimeter  = "_"                // разделитель для аргументов комманд
	MinNameLength     = 4                  // минимальное число(байт) символов в имени продукта
	DateFormat        = "2 Jan 2006 15:04" // формат времени для хранения дат
	GetUpdatesTimeout = 60                 // таймаут на запрос получения обновлений чата
	GrpcPort          = ":8081"            // порт, на котором работает grpc сервер
	RESTPort          = ":8080"            // порт, на котором работает REST и сваггер

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
)
