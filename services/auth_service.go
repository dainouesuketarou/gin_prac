package services

import (
	"first_gin_app/models"
	"first_gin_app/repositories"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(email string, password string) error
	Login(email string, password string) (*string, error)
	GetUserFromToken(tokenString string) (*models.User, error)
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) Signup(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	return s.repository.CreateUser(user)
}

func (s *AuthService) Login(email string, password string) (*string, error) {
	foundUser, err := s.repository.FindUser(email)
	if err != nil {
		return nil, err
	}

	// 第一引数(データベースに保存されているパスワード)と第二引数(ユーザーが入力したパスワード)の比較
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	token, err := CreateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func CreateToken(userId uint, email string) (*string, error) {
	// jwt.SigningMethodHS256は*jwt.SigningMethodHMAC型である。
	// このことはでコードするときに必要な知識である。(GetUserFromToken)
	// var jwt.SigningMethodHS256 *jwt.SigningMethodHMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId,
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	})

	// SignedStringメソッドでトークンを著名し、SECRET_KEYを使ってトークンを暗号化し、
	// 第三者がトークンの内容を改ざんできなくなる。
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

// 参考
// var token *jwt.Token
// &jwt.Token{
//     Raw: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEyMywiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNzAwMDAwMDAwfQ.tV2Cn2Ne8ihGftPUN-Kcnp3rNaPvgMlVJ45x_jVsNmU",
//     Method: &jwt.SigningMethodHMAC{},
//     Header: map[string]interface{}{
//         "alg": "HS256",
//         "typ": "JWT",
//     },
//     Claims: jwt.MapClaims{
//         "sub":   123,
//         "email": "test@example.com",
//         "exp":   1700000000,
//     },
//     Signature: "tV2Cn2Ne8ihGftPUN-Kcnp3rNaPvgMlVJ45x_jVsNmU",
//     Valid: true,
// }

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	// jwt.ParseのkeyFuncは解析されたトークンを受け取り、署名を検証するための暗号化キーを返す必要があります。
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// トークンのハッシュアルゴリズムがjwt.SigningMethodHMAC型かどうかをokに返します。(型アサーション)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	var user *models.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 有効期限が切れているとき
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, jwt.ErrTokenExpired
		}
		user, err = s.repository.FindUser(claims["email"].(string))
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
