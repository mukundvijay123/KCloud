package metadatarouter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTMiddleWare struct {
	secretKey      interface{}
	signingMethod  jwt.SigningMethod
	MiddlewareFunc func(jwt.MapClaims, context.Context) (context.Context, bool, error) // return enriched ctx
	logger         *log.Logger
}

// SetSigningMethod sets the signing method (e.g., jwt.SigningMethodHS256)
func (j *JWTMiddleWare) SetSigningMethod(method jwt.SigningMethod) {
	j.signingMethod = method
}

// SetSecretKey sets the secret key used for signing and verifying tokens
func (j *JWTMiddleWare) SetSecretKey(key interface{}) {
	j.secretKey = key
}

// SetMiddlewareFunc sets the custom verification and ctx enrichment function
func (j *JWTMiddleWare) SetMiddlewareFunc(f func(jwt.MapClaims, context.Context) (context.Context, bool, error)) {
	j.MiddlewareFunc = f
}

// SetLogger sets a custom logger
func (j *JWTMiddleWare) SetLogger(l *log.Logger) {
	j.logger = l
}

// GenerateToken creates a JWT for a given user_id
func (j *JWTMiddleWare) GenerateToken(userID string) (string, error) {
	if j.signingMethod == nil {
		j.signingMethod = jwt.SigningMethodHS256
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString(j.secretKey)
}

// JWTMiddleware validates the JWT and passes the request to the next handler
func (j *JWTMiddleWare) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			if j.logger != nil {
				j.logger.Printf("[WARN] Missing Authorization header from %s", r.RemoteAddr)
			}
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method != j.signingMethod {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			if j.logger != nil {
				j.logger.Printf("[ERROR] Invalid token from %s: %v", r.RemoteAddr, err)
			}
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			if j.logger != nil {
				j.logger.Printf("[ERROR] Invalid token claims from %s", r.RemoteAddr)
			}
			return
		}

		// start with existing request ctx
		ctx := r.Context()
		if j.MiddlewareFunc != nil {
			var valid bool
			ctx, valid, err = j.MiddlewareFunc(claims, ctx)
			if err != nil {
				http.Error(w, "Error verifying token", http.StatusInternalServerError)
				if j.logger != nil {
					j.logger.Printf("[ERROR] MiddlewareFunc failed: %v", err)
				}
				return
			}
			if !valid {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				if j.logger != nil {
					j.logger.Printf("[WARN] Token claims rejected for %s", r.RemoteAddr)
				}
				return
			}
		}

		r = r.WithContext(ctx)

		if j.logger != nil {
			j.logger.Printf("[INFO] Authorized request from %s (user_id=%v) in %v",
				r.RemoteAddr, claims["user_id"], time.Since(start))
		}

		next.ServeHTTP(w, r)
	})
}
