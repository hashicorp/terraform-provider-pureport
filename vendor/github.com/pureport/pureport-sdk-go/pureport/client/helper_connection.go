package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
)

const logReqMsg = `Pureport API Request Details:
---[ REQUEST ]---------------------------------------
%s
-----------------------------------------------------`

const logRespMsg = `Pureport API Response Details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------`

// ValidateConnection validates the content is a valid connection type
// and returns the content if it is valid, otherwise error
func ValidateConnection(content interface{}) (interface{}, error) {

	switch content.(type) {
	case AwsDirectConnectConnection,
		AzureExpressRouteConnection,
		GoogleCloudInterconnectConnection,
		DummyConnection,
		SiteIpSecVpnConnection:
		return content, nil

	default:
		return nil, reportError("body should be valid Connection")
	}
}

// DecodeConnectionData takes the HTTP Request content body and decodes it for the specific
// Connection Type. Since Golang doesn't have inheritance, the decode won't un-marshal parameters
// unless the specific connection type is defined.
func DecodeConnectionData(cli *APIClient, b []byte, contentType string) (interface{}, error) {

	var decodedConnection interface{}

	// Decode as a Base Connection first to get the type
	var base Connection
	err := cli.decode(&base, b, contentType)
	if err != nil {
		return decodedConnection, err
	}
	// Check the Connection type and decode as sub type
	switch base.Type_ {
	case "AWS_DIRECT_CONNECT":

		var c = AwsDirectConnectConnection{}
		err = cli.decode(&c, b, contentType)
		decodedConnection = c

	case "AZURE_EXPRESS_ROUTE":
		var c = AzureExpressRouteConnection{}
		err = cli.decode(&c, b, contentType)
		decodedConnection = c

	case "GOOGLE_CLOUD_INTERCONNECT":
		var c = GoogleCloudInterconnectConnection{}
		err = cli.decode(&c, b, contentType)
		decodedConnection = c

	case "SITE_IPSEC_VPN":
		var c = SiteIpSecVpnConnection{}
		err = cli.decode(&c, b, contentType)
		decodedConnection = c

	default:
		var c = DummyConnection{}
		err = cli.decode(&c, b, contentType)
		decodedConnection = c
	}

	return decodedConnection, err
}

// LogRequest should pretty print the HTTP Request being sent to the REST API.
// req - The HTTP Request
func logRequest(req *http.Request) {
	reqData, err := httputil.DumpRequest(req, true)
	if err == nil {
		log.Debugf(logReqMsg, prettyPrintJsonLines(reqData))
	} else {
		log.Errorf("Pureport API Request error: %#v", err)
	}
}

// LogResponse should pretty print the HTTP Response return from the REST API.
// resp - The HTTP Response
func logResponse(resp *http.Response) {
	respData, err := httputil.DumpResponse(resp, true)
	if err == nil {
		log.Debugf(logRespMsg, prettyPrintJsonLines(respData))
	} else {
		log.Errorf("Pureport API Response error: %#v", err)
	}
}

// prettyPrintJsonLines iterates through a []byte line-by-line,
// transforming any lines that are complete json into pretty-printed json.
func prettyPrintJsonLines(b []byte) string {
	parts := strings.Split(string(b), "\n")
	for i, p := range parts {
		if b := []byte(p); json.Valid(b) {
			var out bytes.Buffer
			json.Indent(&out, b, "", " ")
			parts[i] = out.String()
		}
	}
	return strings.Join(parts, "\n")
}
