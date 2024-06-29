package util

import (
	"citadel/internal/middleware"
	"context"
)

func ShowPaymentMethodDialog(ctx context.Context) bool {
	return ctx.Value(middleware.CTX_KEY_SHOW_PAYMENT_METHOD_DIALOG).(bool)
}
