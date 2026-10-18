package csrf

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	apiDelivery "github.com/Timur1414/Smart-Catch-Up/internal/app/api_gateway/delivery/http"
	"github.com/Timur1414/Smart-Catch-Up/pkg/logger"
	"go.uber.org/zap"
)

type Service interface {
	ValidateCsrf(ctx context.Context, csrfToken string, accessToken string) (bool, error)
	SetCsrfCookie(ctx context.Context, w http.ResponseWriter, r *http.Request)
	GetCsrfFromCookie(ctx context.Context, r *http.Request) string
	GetCsrfFromHeader(r *http.Request) string
	GetAccessToken(ctx context.Context, r *http.Request) (string, error)
	ValidateSecFetchSite(r *http.Request) bool
}

const (
	HeaderName           = "X-CSRF-Token"
	CookieName           = "csrf_token"
	CookieExpirationTime = time.Hour
)

type Csrf struct {
	secret          []byte
	randNonceLength int
	mu              sync.RWMutex
}

func NewCsrf(secretKey string) (*Csrf, error) {
	if secretKey == "" {
		return &Csrf{}, ErrCsrfSecret
	}
	return &Csrf{
		secret:          []byte(secretKey),
		randNonceLength: 32,
	}, nil
}

func (obj *Csrf) ValidateCsrf(ctx context.Context, csrfToken string, accessToken string) (bool, error) {
	log := logger.GetLoggerWithRequestId(ctx)
	parts := strings.Split(csrfToken, ".")
	if len(parts) != 2 {
		log.Warn("Get invalid CSRF token")
		return false, ErrInvalidCsrf
	}
	hmacValue := parts[0]
	randomValue := parts[1]
	message := fmt.Sprintf("%d!%s!%d!%s", len(accessToken), accessToken, len(randomValue), randomValue)
	hmacObj := hmac.New(sha256.New, obj.secret)
	hmacObj.Write([]byte(message))
	expectedHmac := hex.EncodeToString(hmacObj.Sum(nil))
	if !hmac.Equal([]byte(hmacValue), []byte(expectedHmac)) {
		log.Warn("Get wrong CSRF token")
		return false, ErrInvalidCsrfSignature
	}
	return true, nil
}

func (obj *Csrf) SetCsrfCookie(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	token, err := obj.GetAccessToken(ctx, r)
	if err != nil {
		return
	}
	signedCsrfToken, err := obj.generateSignedCsrf(ctx, token)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    signedCsrfToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(CookieExpirationTime),
	})
	rawCsrfToken, err := obj.getCsrfValueFromToken(ctx, signedCsrfToken)
	if err != nil {
		return
	}
	w.Header().Set(HeaderName, rawCsrfToken)
}

func (obj *Csrf) GetCsrfFromCookie(ctx context.Context, r *http.Request) string {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		log := logger.GetLoggerWithRequestId(ctx)
		log.Warn("Failed to get CSRF cookie", zap.Error(err))
		return ""
	}
	return cookie.Value
}

func (obj *Csrf) GetCsrfFromHeader(r *http.Request) string {
	return r.Header.Get(HeaderName)
}

func (obj *Csrf) GetAccessToken(ctx context.Context, r *http.Request) (string, error) {
	cookie, err := r.Cookie(apiDelivery.AccessTokenName)
	if err != nil {
		log := logger.GetLoggerWithRequestId(ctx)
		log.Warn("Failed to get access token", zap.Error(err))
		return "", err
	}
	return cookie.Value, nil
}

func (obj *Csrf) ValidateSecFetchSite(r *http.Request) bool {
	log := logger.GetLoggerWithRequestId(r.Context())
	secFetchSite := r.Header.Get("Sec-Fetch-Site")
	if secFetchSite == "" {
		log.Info("[CSRF] Sec-Fetch-Site header missing, relying on CSRF tokens")
		return true
	}
	if r.Method == http.MethodGet || r.Method == http.MethodOptions || r.Method == http.MethodHead || r.Method == http.MethodTrace || r.Method == http.MethodConnect {
		return true
	}
	switch secFetchSite {
	case "same-origin":
		log.Info("[CSRF] Sec-Fetch-Site: same-origin, request allowed")
		return true
	case "same-site":
		log.Info("[CSRF] Sec-Fetch-Site: same-site, request allowed")
		return true
	case "cross-site":
		log.Warn("[CSRF] Sec-Fetch-Site: cross-site, request blocked")
		return false
	case "none":
		log.Info("[CSRF] Sec-Fetch-Site: none, request allowed")
		return true
	default:
		log.Warn("[CSRF] Sec-Fetch-Site: unknown, request blocked",
			zap.String("Sec-Fetch-Site", secFetchSite))
		return false
	}
}

func (obj *Csrf) generateRandValue(ctx context.Context) (string, error) {
	data := make([]byte, obj.getRandNonceLength())
	_, err := rand.Read(data)
	if err != nil {
		log := logger.GetLoggerWithRequestId(ctx)
		log.Error("Error generating random value", zap.Error(err))
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func (obj *Csrf) getCsrfValueFromToken(ctx context.Context, csrf string) (string, error) {
	parts := strings.Split(csrf, ".")
	if len(parts) != 2 {
		log := logger.GetLoggerWithRequestId(ctx)
		log.Warn("Get invalid CSRF value")
		return "", ErrInvalidCsrf
	}
	return parts[1], nil
}

func (obj *Csrf) generateSignedCsrf(ctx context.Context, token string) (string, error) {
	randomValue, err := obj.generateRandValue(ctx)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("%d!%s!%d!%s", len(token), token, len(randomValue), randomValue)
	hmacObj := hmac.New(sha256.New, obj.getSecret())
	hmacObj.Write([]byte(message))
	hmacValue := hex.EncodeToString(hmacObj.Sum(nil))
	csrfToken := hmacValue + "." + randomValue
	return csrfToken, nil
}

func (obj *Csrf) getSecret() []byte {
	obj.mu.RLock()
	defer obj.mu.RUnlock()
	return obj.secret
}

func (obj *Csrf) getRandNonceLength() int {
	obj.mu.RLock()
	defer obj.mu.RUnlock()
	return obj.randNonceLength
}
