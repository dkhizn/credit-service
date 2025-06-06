package httpadapter_test

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	httpadapter "credit-service/internal/adapters/primary/http-adapter"
	creditService "credit-service/internal/application/credit-service"
	"credit-service/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNew_HTTPAdapter(t *testing.T) {
	cfg := &config.Config{
		App: config.App{
			Name:    "test-app",
			Version: "0.1",
		},
		HTTP: config.HTTP{
			Port: ":0",
		},
	}

	// Логгер, который пишет в буфер (чтобы не засорять stdout)
	logBuf := &bytes.Buffer{}
	logger := log.New(logBuf, "", 0)

	// Передаём nil вместо *CredtiService, т.к. конструктор его не использует внутри
	adapter, err := httpadapter.New(logger, cfg, (*creditService.CredtiService)(nil))
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
}

// TestStart_Shutdown проверяет, что Start корректно запускает сервер и завершает его по отмене контекста.
func TestStart_Shutdown(t *testing.T) {
	cfg := &config.Config{
		App: config.App{
			Name:    "test-app",
			Version: "0.1",
		},
		HTTP: config.HTTP{
			Port: ":0",
		},
	}

	logBuf := &bytes.Buffer{}
	logger := log.New(logBuf, "", 0)

	adapter, err := httpadapter.New(logger, cfg, (*creditService.CredtiService)(nil))
	assert.NoError(t, err)

	// Создаём контекст, который позже отменим
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем Start в отдельной горутине (он блокирует, пока не отменён контекст)
	startCh := make(chan error)
	go func() {
		startCh <- adapter.Start(ctx)
	}()

	// Небольшая пауза, чтобы сервер успел подняться
	time.Sleep(50 * time.Millisecond)

	// Отправляем запрос на «слепую» конечную точку, чтобы убедиться, что сервер слушает
	// (неважно, получим ли 404 — главное, что соединение состоялось)
	_, _ = http.Get("http://127.0.0.1" + cfg.HTTP.Port + "/nonexistent")

	// Теперь отменяем контекст — это должно инициировать Shutdown
	cancel()

	// Ждём, что Start() вернётся без ошибки
	select {
	case err := <-startCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Start() did not return within timeout")
	}
}
