package operation

import "strings"

func buildOperation(layer, service, method string) string {
	return strings.Join(
		[]string{layer, service, method},
		".",
	)
}

func ServicesOperation(service, method string) string {
	const layer = "Services"
	return buildOperation(layer, service, method)
}
