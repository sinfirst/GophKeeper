package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var sugar zap.SugaredLogger

// ResponseData содержит данные о gRPC ответе
type GrpcResponseData struct {
	Status     string
	StatusCode int
	Duration   time.Duration
	Method     string
}

// loggingUnaryServerInterceptor возвращает UnaryServerInterceptor для логирования
func LoggingUnaryInterceptor(logger zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Вызываем обработчик
		res, err := handler(ctx, req)

		// Получаем статус ответа
		st, _ := status.FromError(err)
		statusCode := st.Code()
		statusMsg := st.Message()

		// Формируем данные для логирования
		responseData := &GrpcResponseData{
			Status:     statusMsg,
			StatusCode: int(statusCode),
			Duration:   time.Since(start),
			Method:     info.FullMethod,
		}

		// Логируем информацию
		logger.Infoln(
			"\n",
			"-----GRPC REQUEST-----\n",
			"Method:", responseData.Method, "\n",
			"Status:", responseData.Status, "\n",
			"Duration:", responseData.Duration, "\n",
			"Status code:", responseData.StatusCode, "\n",
		)

		return res, err
	}
}

// NewLogger конструктор для структуры
func NewLogger() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	sugar = *logger.Sugar()

	return sugar
}
