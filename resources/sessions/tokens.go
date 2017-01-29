package sessions

import (
    "github.com/dgrijalva/jwt-go"
    "gopkg.in/mgo.v2/bson"
    "time"
)

type ITokenClaims interface {
    GetJTI() string
    GetExpireAt() time.Time
}

type IToken interface {
    IsValid() bool
}

type IRevokedToken interface {
}

func NewTokenClaims(userID string) *TokenClaims {
    return &TokenClaims{UserId: userID}
}

type TokenClaims struct {
    UserId   string
    ExpireAt time.Time
    IssuedAt time.Time
    JTI      string
}

func (claim *TokenClaims) GetJTI() string {
    return claim.JTI
}

func (claim *TokenClaims) GetExpireAt() time.Time {
    return claim.ExpireAt
}

type RevokedToken struct {
    ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
    ExpiryDate  time.Time     `json:"exp" bson:"exp"`
    RevokedDate time.Time     `json:"revokedat" bson:"revokedat"`
}

func NewToken(token *jwt.Token) *Token {
    return &Token{token}
}

type Token struct {
    *jwt.Token
}

func (token *Token) IsValid() bool {
    return token.Valid
}
