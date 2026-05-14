package payment

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

// VerifySucceededPaymentIntent checks Stripe for a succeeded PI matching the order total.
// Card numbers never touch this server — only Stripe IDs and server-side secret key.
func VerifySucceededPaymentIntent(secretKey, paymentIntentID string, expectedTotal float64, currency string) error {
	if strings.TrimSpace(secretKey) == "" || strings.TrimSpace(paymentIntentID) == "" {
		return errors.New("missing stripe credentials or payment intent id")
	}
	stripe.Key = secretKey
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return fmt.Errorf("stripe: %w", err)
	}
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return fmt.Errorf("payment intent status %s", pi.Status)
	}
	want := int64(math.Round(expectedTotal * 100))
	if pi.Amount != want {
		return fmt.Errorf("amount mismatch: intent %d vs order %d cents", pi.Amount, want)
	}
	cur := strings.ToLower(string(pi.Currency))
	if cur != strings.ToLower(currency) {
		return fmt.Errorf("currency mismatch: %s vs %s", cur, currency)
	}
	return nil
}
