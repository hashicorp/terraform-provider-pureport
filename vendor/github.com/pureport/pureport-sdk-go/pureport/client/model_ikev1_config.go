/*
 * Pureport Control Plane
 *
 * Pureport API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

type Ikev1Config struct {
	Esp *Ikev1EspConfig `json:"esp"`
	Ike *Ikev1IkeConfig `json:"ike"`
}
