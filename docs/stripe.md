# Stripe setup

Install/Update `stripe` cli

    brew upgrade stripe

Always use latest stripe cli version, when starting the client to forward webhook petitions to localhost, you'll get
the `STRIPE_WEBHOOK_SECRET` that you will need to tests.

Check if you see the webhook listed in `https://dashboard.stripe.com/test/webhooks` or similar page, this one is when using
test-mode, but perhaps in sandbox is diferent.

## Create the custumer ID

Very important to create an `stripe_customer_id` before. So use the customer's name and email to do so, and store the
`stripe_customer_id` in the database.

It is a good practice to store the `map key` name in constants

```go
const (
	MetadataKeyAppUserID    = "app_user_id"
	MetadataKeyAppAccountID = "app_account_id"
)
```

This is an example using a PocketBase database

```go
...
	stripeCustomerID := accountRecord.GetString("sub_stripe_customer_id")
	if stripeCustomerID == "" {
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(user.GetString("email")),
			Name:  stripe.String(user.GetString("name")),
			// Add any other relevant metadata
			Metadata: map[string]string{
				MetadataKeyAppUserID:    user.Id,
				MetadataKeyAppAccountID: accountID,
			},
		}
		newCustomer, err := customer.New(customerParams) // Use the customer package here
		if err != nil {
			fmt.Printf("Error creating Stripe customer: %v\n", err)
			return c.JSON(http.StatusInternalServerError, Map{"error": "Failed to create customer record"})
		}
		stripeCustomerID = newCustomer.ID
		accountRecord.Set("sub_stripe_customer_id", stripeCustomerID)
		// Save the updated user record with the new Stripe Customer ID
		if err := app.PB.Dao().SaveRecord(accountRecord).Error; err != nil {
			fmt.Printf("Error saving Stripe Customer ID to account %s: %v\n", accountRecord.Id, err)
			// Return an error to be safe.
			return c.JSON(http.StatusInternalServerError, Map{"error": "Failed to update account record"})
		}
		fmt.Printf("Created and saved Stripe Customer ID %s for account %s\n", stripeCustomerID, accountRecord.Id)
	}
...
```

##

After a successful Stripe subscription checkout, the following typically occurs:

- The customer is redirected to the success URL specified in your Checkout Session.
- Stripe sends a webhook event to your server to confirm the successful subscription.
- Your server should listen for this webhook event (usually checkout.session.completed) and update your database accordingly.
- The frontend can then poll your backend API to check for the updated subscription status.

In our provided code:

- The CreateCheckoutSession function in `/server/cmd/web/handlers/subscriptions.go`
  sets up the Checkout Session with success and cancel URLs:

```go
successURL := domain + "/subscription/success?session_id={CHECKOUT_SESSION_ID}"
cancelURL := domain + "/subscription/cancel"
```

The `StripeWebhook` function in `/server/cmd/web/handlers/stripe.go` handles incoming webhook events from Stripe, including subscription-related events.

The `SubscriptionService` in `/angular/src/app/services/subscription.service.ts` provides methods to check the subscription status:

```typescript
loadSubscriptionStatus(): Observable<SubscriptionStatus> {
  const headers = this.getAuthHeaders();
  return this.http.get<SubscriptionStatus>(`${this.apiUrl}/subscriptions/current`, { headers })
    .pipe(
      tap(status => this.subscriptionStatusSubject.next(status)),
      catchError(err => {
        this.subscriptionStatusSubject.next(null);
        return this.handleError(err);
      })
    );
}
```

To fully implement this flow:

- Ensure your webhook handler (StripeWebhook) updates your database when it receives a checkout.session.completed event.
- Implement the success page to call loadSubscriptionStatus() periodically until the subscription is confirmed.
- Update your UI based on the subscription status returned by the backend.

## Create Checkout Session

Notes on Metadata. To set metadata on a Subscription created from a Checkout
Session, you need to use the subscription_data.metadata parameter when creating
the Checkout Session. Here's how you can do it:

```go
...
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(stripeCustomerID), // Use the potentially newly created customer ID
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(reqBody.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		// Add any other relevant metadata
		Metadata: map[string]string{
			MetadataKeyAppUserID:    user.Id,
			MetadataKeyAppAccountID: accountID,
		},
		// When creating a Checkout Session for a subscription, include the
		// subscription_data.metadata parameter with the metadata you want to set on
		// the resulting Subscription.
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				MetadataKeyAppUserID:    user.Id,
				MetadataKeyAppAccountID: accountID,
				// Add any other metadata you want on the Subscription
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		// Optionally allow promotion codes
		AllowPromotionCodes: stripe.Bool(true),
	}
...
```
