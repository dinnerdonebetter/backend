package aes

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/cryptography"
)

// aesImpl is the standard EncryptorDecryptor implementation.
type aesImpl struct {
	tracer tracing.Tracer
	logger logging.Logger
	key    [32]byte
}

func NewEncryptorDecryptor(tracerProvider tracing.TracerProvider, logger logging.Logger, key []byte) (cryptography.EncryptorDecryptor, error) {
	if len(key) != 32 {
		return nil, cryptography.ErrIncorrectKeyLength
	}

	return &aesImpl{
		logger: logging.EnsureLogger(logger).WithName("encryptor"),
		tracer: tracing.NewTracer(tracerProvider.Tracer("encryptor")),
	}, nil
}