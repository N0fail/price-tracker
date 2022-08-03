package config

const CommandDelimeter = "_"          // разделитель для аргументов комманд
const MinNameLength = 4               // минимальное число(байт) символов в имени продукта
const DateFormat = "2 Jan 2006 15:04" // формат времени для хранения дат
const GetUpdatesTimeout = 60          // таймаут на запрос получения обновлений чата
const GrpcPort = ":8081"              // порт, на котором работает grpc сервер
const RESTPort = ":8080"              // порт, на котором работает REST и сваггер
