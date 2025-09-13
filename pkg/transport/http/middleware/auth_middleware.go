package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	pbAuthClient "github.com/rasteiro11/MCABankGateway/gen/proto/go/user"
	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
	"github.com/rasteiro11/PogCore/pkg/transport/rest"
)

var ErrNotAuthorized = errors.New("not authorized")

func ValidateUserMiddleware(authClient pbAuthClient.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwtToken := c.GetReqHeaders()
		tok, ok := jwtToken["Authorization"]
		if !ok {
			return rest.NewStatusUnauthorized(c, ErrNotAuthorized)
		}

		tok = strings.ReplaceAll(tok, "Bearer ", "")
		res, err := authClient.VerifySession(c.Context(), &pbAuthClient.VerifySessionRequest{
			Token: tok,
		})
		if err != nil {
			return rest.NewStatusUnauthorized(c, err)
		}

		claims := &domain.Claims{
			UserID: uint(res.GetUserId()),
		}
		c.Context().SetUserValue("user", claims)

		return c.Next()
	}
}
