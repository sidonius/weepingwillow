package main

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AccountDoc struct {
	ID       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	RoleID   int    `json:"roleid" bson:"roleid"`
	Stat     int    `json:"stat" bson:"stat"`
}

func (x AccountDoc) toResponsedAccount() *ResponsedAccount {
	return &ResponsedAccount{ID: x.ID, Name: x.Name, RoleID: x.RoleID}
}

type ResponsedAccount struct {
	ID     string `json:"id" bson:"id"`
	Name   string `json:"name" bson:"name"`
	RoleID int    `json:"roleid" bson:"roleid"`
}

type LoginCredential struct {
	ID       string `json:"id" form:"id" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

const (
	jwtIdentityKey = "id"
	jwtNameKey     = "name"
	jwtRoleIDKey   = "roleid"
	jwtRealm       = "weepingwillow"

	cookieName   = "jwt"
	cookieDomain = "ems.com.tw"
	// cookieDomain = "localhost"
)

var jwtSecret = []byte("junoatweepingwillow")

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"eid":  claims[jwtIdentityKey],
		"name": claims[jwtNameKey],
		"text": "Hello World.",
		"time": time.Now().Format(time.RFC3339),
	})
}

func payloadHandler(data interface{}) jwt.MapClaims {
	if v, ok := data.(*ResponsedAccount); ok {
		return jwt.MapClaims{
			jwtIdentityKey: v.ID,
			jwtNameKey:     v.Name,
			jwtRoleIDKey:   v.RoleID,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	if v, ok := claims[jwtIdentityKey].(string); ok {
		return &ResponsedAccount{
			ID: v,
		}
	}
	return nil
}

func loginResponseHandler(c *gin.Context, code int, token string, t time.Time) {
	cookie, err := c.Cookie(cookieName)
	if err != nil {
		elog.Warn("failed to get cookie: ", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": t.Format(time.RFC3339),
		"cookie": cookie,
	})
}

func authenticatorHandler(c *gin.Context) (interface{}, error) {
	// var login LoginCredential
	// // jsonData, err := ioutil.ReadAll(c.Request.Body)
	// // if err != nil {
	// // 	elog.Error("failed to read request body: ", err)
	// // }
	// // elog.Info(string(jsonData))
	// if err := c.ShouldBindJSON(&login); err != nil {
	// 	// if there is an EOF error, comment out ioutil.ReadAll(c.Request.Body)
	// 	elog.Error("failed to bind: ", err)
	// 	return "", jwt.ErrMissingLoginValues
	// }

	// var account AccountDoc
	// err := mdb.Collection("account").FindOne(
	// 	context.TODO(),
	// 	bson.D{
	// 		{Key: "$and",
	// 			Value: bson.A{
	// 				bson.D{{Key: "eid", Value: login.EID}},
	// 				bson.D{{Key: "stat", Value: bson.D{{Key: "$ne", Value: -1}}}},
	// 			},
	// 		},
	// 	},
	// ).Decode(&account)
	// if err != nil {
	// 	elog.Errorf("query for %s: %v", login.EID, err)
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	// 	return nil, jwt.ErrFailedAuthentication
	// }

	// good := false
	// if account.Stat == 1 {
	// 	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(login.Password)); err == nil {
	// 		good = true
	// 	} else {
	// 		elog.Error("invalid password: ", err)
	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	// 	}
	// } else if account.Stat == 0 {
	// 	if account.Password == login.Password && account.Password != "" {
	// 		if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost); err == nil {
	// 			result, err := mdb.Collection("account").UpdateOne(
	// 				context.TODO(),
	// 				bson.D{{Key: "eid", Value: account.EID}},
	// 				bson.D{
	// 					{Key: "$set", Value: bson.D{{Key: "password", Value: string(hashedPassword)}}},
	// 					{Key: "$set", Value: bson.D{{Key: "stat", Value: 1}}},
	// 				},
	// 			)
	// 			if err == nil && result.ModifiedCount != 0 {
	// 				good = true
	// 			} else {
	// 				elog.Error("update account data: ", err)
	// 				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	// 			}
	// 		} else {
	// 			elog.Error("generate hash: ", err)
	// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	// 		}
	// 	} else {
	// 		elog.Errorf("invalid password(%s) for eid(%s): ", login.Password, login.EID)
	// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	// 	}
	// }

	// if good {
	// 	return account.toResponsedAccount(), nil
	// } else {
	return nil, jwt.ErrFailedAuthentication
	// }
}

func NewJwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           jwtRealm,
		Key:             jwtSecret,
		Timeout:         time.Minute * time.Duration(cnf.WebService.JwtTimeout),
		MaxRefresh:      time.Minute * time.Duration(cnf.WebService.JwtRefreshTime),
		IdentityKey:     jwtIdentityKey,
		PayloadFunc:     payloadHandler,
		IdentityHandler: identityHandler,
		Authenticator:   authenticatorHandler,

		// User can define own LoginResponse func.
		LoginResponse: loginResponseHandler,

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		// Optionally return the token as a cookie
		SendCookie: true,

		// CookieName allow cookie name change for development
		CookieName: cookieName,

		// Allow cookie domain change for development
		CookieDomain: cookieDomain,
	})
}
