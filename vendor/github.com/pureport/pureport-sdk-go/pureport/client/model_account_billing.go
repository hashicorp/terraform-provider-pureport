/*
 * Pureport Control Plane
 *
 * Pureport API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

type AccountBilling struct {
	Account        *Link            `json:"account,omitempty"`
	Address        *PhysicalAddress `json:"address"`
	Email          string           `json:"email"`
	Href           string           `json:"href,omitempty"`
	Name           string           `json:"name"`
	StripeExpiry   string           `json:"stripeExpiry"`
	StripeLastFour string           `json:"stripeLastFour"`
	StripeToken    string           `json:"stripeToken"`
}
