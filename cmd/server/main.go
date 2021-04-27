package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"os"

	"github.com/llinder/svtauth/internal/routes"

	"github.com/akamensky/argparse"
	"github.com/gin-gonic/gin"
)

type Grant struct {
	CallerIdentityUrl string `form:"caller_identity_url" binding:"required"`
}

type CallerIdentity struct {
	XMLName xml.Name `xml:"GetCallerIdentityResult"`
	ARN     string   `xml:"Arn"`
	UserId  string   `xml:"UserId"`
	Account int      `xml:"Account"`
}

type CallerIdentityResponse struct {
	XMLName        xml.Name       `xml:"GetCallerIdentityResponse"`
	CallerIdentity CallerIdentity `xml:"GetCallerIdentityResult"`
}

func main() {
	parser := argparse.NewParser("svtclient", "Service Token Client")
	address := parser.String("a", "address", &argparse.Options{Default: "127.0.0.1:8080", Help: "Server listen address"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(parser.Usage(err))
	}

	client := &http.Client{}

	router := gin.Default()

	router.POST("/grant", routes.PostGrant(client))
	router.GET("/health", routes.GetHealth())

	server := &http.Server{
		Addr:    *address,
		Handler: router,
	}

	server.ListenAndServe()
}
