package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// 对称密钥（HS256）
	secretKey = []byte("your-secret-key-change-in-production")

	// RSA密钥对（RS256）
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	// Token过期时间
	tokenExpiration = 24 * time.Hour

	// 签名方法：HS256 或 RS256
	signingMethod = "HS256"
)

// Claims JWT声明
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// SetSecretKey 设置JWT对称密钥（HS256）
func SetSecretKey(key string) {
	secretKey = []byte(key)
	signingMethod = "HS256"
}

// SetRSAPrivateKey 设置RSA私钥（RS256）
func SetRSAPrivateKey(privateKeyPath string) error {
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return errors.New("failed to decode PEM block")
	}

	// 支持PKCS1和PKCS8格式
	var key *rsa.PrivateKey
	if block.Type == "RSA PRIVATE KEY" {
		// PKCS1格式
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse PKCS1 private key: %w", err)
		}
	} else if block.Type == "PRIVATE KEY" {
		// PKCS8格式
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse PKCS8 private key: %w", err)
		}
		var ok bool
		key, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return errors.New("not an RSA private key")
		}
	} else {
		return fmt.Errorf("unsupported key type: %s", block.Type)
	}

	privateKey = key
	signingMethod = "RS256"
	return nil
}

// SetRSAPublicKey 设置RSA公钥（RS256）
func SetRSAPublicKey(publicKeyPath string) error {
	keyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return errors.New("failed to decode PEM block")
	}

	var pub interface{}
	if block.Type == "PUBLIC KEY" {
		// PKIX格式
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse PKIX public key: %w", err)
		}
	} else if block.Type == "RSA PUBLIC KEY" {
		// PKCS1格式
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse PKCS1 public key: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported key type: %s", block.Type)
	}

	var ok bool
	publicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}

	return nil
}

// SetTokenExpiration 设置Token过期时间
func SetTokenExpiration(duration time.Duration) {
	tokenExpiration = duration
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint64, username string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "flowgo",
		},
	}

	var token *jwt.Token
	if signingMethod == "RS256" {
		if privateKey == nil {
			return "", errors.New("RSA private key not set")
		}
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		return token.SignedString(privateKey)
	} else {
		// 默认使用HS256
		token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(secretKey)
	}
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	var token *jwt.Token
	var err error

	if signingMethod == "RS256" {
		if publicKey == nil {
			return nil, errors.New("RSA public key not set")
		}
		token, err = jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
	} else {
		// 默认使用HS256
		token, err = jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})
	}

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
