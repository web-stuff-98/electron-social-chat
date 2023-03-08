package helpers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/web-stuff-98/electron-social-chat/pkg/db"
	"github.com/web-stuff-98/electron-social-chat/pkg/db/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
	Same stuff as Go-Social-Media, except the session is kept on redis
*/

func createCookie(token string, expiry time.Time) http.Cookie {
	return http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Expires:  expiry,
		MaxAge:   120,
		Secure:   os.Getenv("PRODUCTION") == "true",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}
}

func GetClearedCookie() http.Cookie {
	return http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
		Secure:   os.Getenv("PRODUCTION") == "true",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}
}

func GenerateCookieAndSession(ctx context.Context, uid primitive.ObjectID, collections db.Collections, redisClient *redis.Client) (http.Cookie, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in generate cookie helper function")
		}
	}()
	sid := uuid.New()
	sessionDuration := time.Minute * 2
	expiry := primitive.NewDateTimeFromTime(time.Now().Add(sessionDuration))
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    sid.String(),
		ExpiresAt: expiry.Time().Unix(),
	})
	token, err := claims.SignedString([]byte(os.Getenv("SECRET")))
	cookie := createCookie(token, expiry.Time())
	cmd := redisClient.Set(ctx, sid.String(), uid.Hex(), sessionDuration)
	if cmd.Err() != nil {
		return cookie, err
	}
	return cookie, err
}

func GetUserFromRequest(r *http.Request, ctx context.Context, collections db.Collections, redisClient *redis.Client) (*models.User, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic in get user from request helper function")
		}
	}()
	originalCookie, err := r.Cookie("refresh_token")
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(originalCookie.Value, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	sessionID := token.Claims.(*jwt.StandardClaims).Issuer
	if sessionID == "" {
		return nil, fmt.Errorf("Empty cookie value")
	}
	val, err := redisClient.Get(ctx, sessionID).Result()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving session")
	}
	uid, err := primitive.ObjectIDFromHex(val)
	if err != nil {
		return nil, fmt.Errorf("Invalid ID in session")
	}
	var user models.User
	if err := collections.UserCollection.FindOne(context.TODO(), bson.M{"_id": uid}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Your account could not be found")
		} else {
			return nil, fmt.Errorf("Internal error")
		}
	}
	return &user, nil
}

func DownloadURL(inputURL string) io.ReadCloser {
	_, err := url.Parse(inputURL)
	if err != nil {
		log.Fatal("Failed to parse image url")
	}
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}
	resp, err := client.Get(inputURL)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body
}
func DownloadRandomImage(pfp bool) io.ReadCloser {
	if !pfp {
		return DownloadURL("https://picsum.photos/1100/500")
	} else {
		return DownloadURL("https://100k-faces.glitch.me/random-image")
	}
}

func RemoveDuplicates(strings []string) []string {
	uniqueStrings := make(map[string]bool)
	var unique []string
	for _, str := range strings {
		if !uniqueStrings[str] {
			uniqueStrings[str] = true
			unique = append(unique, str)
		}
	}
	return unique
}
