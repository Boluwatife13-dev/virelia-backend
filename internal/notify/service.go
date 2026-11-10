package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const brevoURL = "https://api.brevo.com/v3/smtp/email"

type Brevo struct {
	APIKey string
	From   string
}

func (b Brevo) Send(
	ctx context.Context,
	to, subject, html string,
) error {
	if b.APIKey == "" {
		log.Printf("EMAIL to %s: %s\n%s", to, subject, html)
		return nil
	}

	body, err := json.Marshal(map[string]any{
		"sender":      map[string]string{"name": "Virelia", "email": b.From},
		"to":          []map[string]string{{"email": to}},
		"subject":     subject,
		"htmlContent": html,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx, "POST", brevoURL, bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("api-key", b.APIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("brevo returned %d", res.StatusCode)
	}

	return nil
}

func (b Brevo) SendAll(
	ctx context.Context,
	to []string,
	subject, html string,
) {
	for _, email := range to {
		err := b.Send(ctx, email, subject, html)
		if err != nil {
			log.Printf("email to %s failed: %v", email, err)
		}
	}
}
