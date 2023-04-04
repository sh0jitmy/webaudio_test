package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func update(s socketio.Conn,server *socketio.Server,o map[string]string) {
	log.Println("update ch", o)
	s.Emit("/","ch",o)
	server.BroadcastToNamespace("/","ch",o)
} 

func main() {
	//subscriber list
	subscList := map[string]int{}  //SID,status
	//pubscriber list
	pubList := map[string] string{}  //SID,ChName
	
	//ch list(To client)
	wschList := map[string]string{} //ChName,SID)  

	router := gin.New()
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "pub:ch", func(s socketio.Conn, msg string) {
		log.Println("pub socketid", s.ID())
		log.Println("pub ch", msg)
		wschList[msg] = s.ID()
		update(s,server,wschList)
		s.Join(msg)
	})
	
	server.OnEvent("/", "sub:connect", func(s socketio.Conn, msg string) {
		log.Println("sub socketid", s.ID())
		log.Println("sub connect", msg)
		subscList[s.ID()] = 1
		update(s,server,wschList)
		s.Join(msg)
	})
	
	server.OnEvent("/", "sub:join", func(s socketio.Conn, msg string) {
		log.Println("sub socketid", s.ID())
		log.Println("sub join", msg)
		s.Join(msg)
	})
	server.OnEvent("/", "sub:leave", func(s socketio.Conn, msg string) {
		log.Println("sub socketid", s.ID())
		log.Println("sub leave", msg)
		s.Leave(msg)
	})
	
	server.OnEvent("/", "audio", func(s socketio.Conn, msg string) {
		log.Println("audio", msg)
		server.BroadcastToRoom("/",wschList[s.ID()],"audio",msg)
	});

	server.OnEvent("/", "disconnect", func(s socketio.Conn, msg string) {
		if _,ok := subscList[s.ID()]; ok {
			subscList[s.ID()] = 0	
		}  
		if val,ok := pubList[s.ID()]; ok {
			wschList[val] = ""				
		}  
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
		if _,ok := subscList[s.ID()]; ok {
			subscList[s.ID()] = 0	
		}  
		if val,ok := pubList[s.ID()]; ok {
			wschList[val] = ""				
		}  
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	router.Use(GinMiddleware("http://localhost:58080"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("../asset"))

	if err := router.Run(":9999"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
