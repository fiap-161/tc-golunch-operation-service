package gateway

type AuthGateway interface {
	GenerateToken(userID string, userType string, additionalClaims map[string]any) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}
