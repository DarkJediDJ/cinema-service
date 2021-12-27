package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/darkjedidj/cinema-service/internal"
	db "github.com/darkjedidj/cinema-service/internal/queries"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

//Authentification ...
func Authentification(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var user config.User
	json.NewDecoder(request.Body).Decode(&user)
	dbUser := db.SelectUser(user)
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}

	jwtToken, err := GenerateJWT(dbUser.AddHalls, dbUser.AddMovies, dbUser.AddSessions, dbUser.UserID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	saveErr := CreateAuth(user.UserID, jwtToken)
	if saveErr != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + saveErr.Error() + `"}`))
	}
	response.Write([]byte(`{"token":"` + jwtToken.AccessToken + `"}`))
	response.Write([]byte(`{"refresh token":"` + jwtToken.RefreshToken + `"}`))
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func GenerateJWT(addHalls bool, addMovies bool, addSessions bool, userID int) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 2).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()
	var err error
	e := godotenv.Load(".env")
	if e != nil {
		log.Fatal("Error loading .env file")
	}
	atClaims := jwt.MapClaims{}
	atClaims["userID"] = userID
	atClaims["addHalls"] = addHalls
	atClaims["addMovies"] = addMovies
	atClaims["addSessions"] = addSessions
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

var Client *redis.Client

func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	Client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
func CreateAuth(userid int, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := Client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := Client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 { 
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

type AccessDetails struct {
	AccessUuid  string
	AddHalls    bool
	AddMovies   bool
	AddSessions bool
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		addHalls, ok := claims["addHalls"].(bool)
		if !ok {
			return nil, err
		}
		addMovies, ok := claims["addMovies"].(bool)
		if !ok {
			return nil, err
		}
		addSessions, ok := claims["addSession"].(bool)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid:  accessUuid,
			AddHalls:    addHalls,
			AddMovies:   addMovies,
			AddSessions: addSessions,
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *AccessDetails) (string, error) {
	uuid, err := Client.Get(authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return uuid, nil
}
