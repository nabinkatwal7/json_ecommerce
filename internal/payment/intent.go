package payment

import (
	"fmt"
	"math"
	"strings"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

// CreatePaymentIntent creates a PaymentIntent for the checkout total. The client uses the returned client secret
// with Stripe.js; card data is collected and tokenized by Stripe (SAQ A).
func CreatePaymentIntent(secretKey string, amount float64, currency, orderID string) (clientSecret, paymentIntentID string, err error) {
	if strings.TrimSpace(secretKey) == "" {
		return "", "", fmt.Errorf("missing stripe secret")
	}
	stripe.Key = secretKey
	amountCents := int64(math.Round(amount * 100))
	if amountCents <= 0 {
		return "", "", fmt.Errorf("invalid amount")
	}
	cur := strings.ToLower(strings.TrimSpace(currency))
	if cur == "" {
		cur = "usd"
	}
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(cur),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Metadata: map[string]string{
			"orderId": orderID,
		},
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		return "", "", err
	}
	return pi.ClientSecret, pi.ID, nil
}
