package client

import ()

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
