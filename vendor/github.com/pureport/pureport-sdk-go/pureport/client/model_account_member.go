/*
 * Pureport Control Plane
 *
 * Pureport API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

type AccountMember struct {
	Account *Link  `json:"account"`
	Href    string `json:"href,omitempty"`
	Roles   []Link `json:"roles"`
	User    *Link  `json:"user"`
}
