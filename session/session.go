package session

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	hankoJwt "github.com/teamhanko/hanko/crypto/jwt"
	"time"
)

type Generator struct {
	jwtGenerator  *hankoJwt.Generator
	sessionLength time.Duration
}

func NewGenerator(privateKey *jwk.Key) (*Generator, error) {
	g, err := hankoJwt.NewGenerator(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create session generator: %w", err)
	}
	return &Generator{
		jwtGenerator:  g,
		sessionLength: time.Minute * 60, // TODO: should come from config
	}, nil
}

func (g *Generator) Generate(userId uuid.UUID) (string, error) {
	issuedAt := time.Now()
	expiration := issuedAt.Add(g.sessionLength)

	token := jwt.New()
	_ = token.Set(jwt.SubjectKey, userId.String())
	_ = token.Set(jwt.IssuedAtKey, issuedAt)
	_ = token.Set(jwt.ExpirationKey, expiration)
	//_ = token.Set(jwt.AudienceKey, []string{"http://localhost"})

	signed, err := g.jwtGenerator.Sign(token)
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func (g *Generator) Verify(token string) (jwt.Token, error) {
	parsedToken, err := g.jwtGenerator.Verify([]byte(token))
	if err != nil {
		return nil, fmt.Errorf("failed to verify session token: %w", err)
	}

	return parsedToken, nil
}
