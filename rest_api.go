package main

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RestServiceFunc(ctx context.Context) {
	defer wg.Done()

	if cnf.WebService.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCH"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin", "Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	router.Use(cors.New(corsConfig))

	authMiddleware, err := NewJwtMiddleware()
	if err != nil {
		elog.Fatal("JWT Error: ", err.Error())
	}

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		elog.Info("NoRoute claims: ", claims)
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "weeping willow.")
	})

	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/refresh_token", authMiddleware.RefreshHandler)

	router.GET("/hello", GetHelloHandler)

	httpd := &http.Server{
		Addr:    ":" + cnf.WebService.Port,
		Handler: router,
	}

	go func() {
		if err := httpd.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			elog.Fatalf("launch simple server: %v", err)
		}
	}()

	<-ctx.Done()
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := httpd.Shutdown(c); err != nil {
		elog.Error("shutdown Web Server: ", err)
	}
	cancel()
	elog.Info("<<")
}

//

func GetNodeListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// func GetNodeListHandler(c *gin.Context) {
// 	// cursor, err := mdb.Collection(Coll_Node).Find(context.TODO(),
// 	// 	bson.D{{Key: "stat", Value: bson.D{{Key: "$eq", Value: 0}}}})
// 	// if err != nil {
// 	// 	elog.Errorf("find tenders: %v", err)
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
// 	// 	return
// 	// }

// 	// var nodes []Node
// 	// if err = cursor.All(context.TODO(), &nodes); err != nil {
// 	// 	elog.Error("cursor.All tenders: ", err)
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
// 	// 	return
// 	// }

// 	set := make(map[int]bool)

// 	for i := 0; i < 5; i++ {
// 		set[rand.Intn(len(nodelist))] = true
// 	}

// 	var nodes []Node
// 	for k := range set {
// 		nodes = append(nodes, nodelist[k])
// 	}

// 	c.JSON(http.StatusOK, TiedNode{Data: nodes})
// }
