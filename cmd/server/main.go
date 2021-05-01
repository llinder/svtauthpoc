package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/llinder/svtauthpoc/internal/roles"
	"github.com/llinder/svtauthpoc/internal/routes"

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
	rolesFile := parser.String("r", "roles", &argparse.Options{Default: "st-only-roles.yaml", Help: "Roles config file"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(parser.Usage(err))
	}

	repo, err := roles.GetRepo(*rolesFile) // GetRepo(rolesFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to load roles config %s", *rolesFile), err)
	}

	fmt.Println(repo.GetTarget("generic", "arn:aws:iam::344541377886:role/STPreprodPowerUser"))

	// rolesConf, err := yaml.Decoder
	// if err != nil {
	// 	log.Fatal(fmt.Sprintf("unable to load roles config %s", *rolesFile), err)
	// }

	// fmt.Println(rolesConf.String())

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
