package shared

import (
	"errors"
	"fmt"
	"strings"
)

type BuildPathParam struct {
	Key   string
	Value string
}

// BuildPath This function is used to build a URL path by replacing placeholders in the path template with actual values from the params map.
// Example: Path template: "/api/v1/users/{userId}/orders/{orderId}"
// Params: map[string]string{"userId": "123", "orderId": "456"}
// Result: "/api/v1/users/123/orders/456"
func BuildPath(pathTemplate string, params []BuildPathParam) (string, error) {
	for _, v := range params {
		placeholder := fmt.Sprintf("{%v}", v.Key)
		if !strings.Contains(pathTemplate, placeholder) {
			return "", errors.New("placeholder not found in path template")
		}
		pathTemplate = strings.ReplaceAll(pathTemplate, placeholder, v.Value)
	}
	return pathTemplate, nil
}
