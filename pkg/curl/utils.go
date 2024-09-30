package curl

import (
	"net/url"
)

func GetInternalUrl(microserviceUrl, endpoint string) string {

	baseURL, err := url.Parse("https://" + microserviceUrl)
	if err != nil {
		return ""
	}
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}

	finalURL := baseURL.ResolveReference(endpointURL)
	return finalURL.String()
}
