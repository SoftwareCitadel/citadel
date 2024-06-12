package util

import "context"

func ShowPaymentMethodDialog(ctx context.Context) bool {
	return ctx.Value("SHOW_PAYMENT_METHOD_DIALOG").(bool)
}
