// @Author : liguoyu
// @Date: 2019/10/29 15:42
package protocol

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"errors"

	log "github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
)

type UserClaim struct {
	jwt.StandardClaims
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
	//Verify   bool   `json:"verify"`
}

func (c *UserClaim) Valid() (err error) {
	if err = c.StandardClaims.Valid(); err != nil {
		return err
	}
	if err = c.ValidSub(); err != nil {
		return err
	}
	return
}

func (c *UserClaim) ValidSub() error {
	if c.Subject == "" {
		return errors.New("missing subject")
	}
	return nil
}

// MakeJWTTokenToContext 返回JWTTokenToContext 直接解析jwt token payload的kv到 context
func MakeJWTTokenToContext() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		token, ok := ctx.Value(kitjwt.JWTTokenContextKey).(string)
		if !ok {
			return ctx
		}
		log.Debug("action ", "解析jwt ", "token ", token)
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			return ctx
		}

		bs, err := DecodeSegment(parts[1])
		if err != nil {
			log.Error("action ", "解析jwt. DecodeSegment 失败 ", "err ", err)
			return ctx
		}
		log.Debug("action ", "解析jwt ", "bs ", string(bs))

		var data map[string]interface{}
		if err := json.Unmarshal(bs, &data); err != nil {
			log.Error("action ", "解析jwt. Unmarshal 失败 ", "err ", err)
			return ctx
		}
		for k, v := range data {
			ctx = context.WithValue(ctx, k, v)
		}

		return ctx
	}
}

func MakeJWTClaimsToContext(ctx context.Context, r *http.Request) context.Context {
	ctx = kitjwt.HTTPToContext()(ctx, r)
	tokenString, ok := ctx.Value(kitjwt.JWTTokenContextKey).(string)
	if !ok {
		return ctx
	}
	token, _ := jwt.ParseWithClaims(tokenString, &UserClaim{}, nil)
	if err := token.Claims.Valid(); err != nil {
		return ctx
	}
	ctx = context.WithValue(ctx, kitjwt.JWTClaimsContextKey, token.Claims)
	return ctx
}

func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
