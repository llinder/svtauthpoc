package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/akamensky/argparse"
)

type STSPresignedAPI interface {
	PresignGetCallerIdentity(
		ctx context.Context,
		params *sts.GetCallerIdentityInput,
		optFns ...func(*sts.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type STSStandardAPI interface {
	AssumeRole(
		ctx context.Context,
		params *sts.AssumeRoleInput,
		optFns ...func(*sts.Options)) (*sts.AssumeRoleOutput, error)
}

func GetPresignedCallerIdentityURL(c context.Context, api STSPresignedAPI, input *sts.GetCallerIdentityInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetCallerIdentity(c, input)
}

type Token struct {
	Token string
}

func main() {
	parser := argparse.NewParser("svtclient", "Service Token Client")

	server := parser.String("s", "server", &argparse.Options{Default: "http://localhost:8080"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(parser.Usage(err))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := sts.NewFromConfig(cfg)
	psClient := sts.NewPresignClient(client)

	identResp, err := GetPresignedCallerIdentityURL(context.TODO(), psClient, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}

	hc := http.Client{}
	form := url.Values{}
	form.Set("caller_identity_url", identResp.URL)
	req, err := http.NewRequest("POST", *server+"/grant", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := hc.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(string(body))
	} else {
		defer resp.Body.Close()
		token := new(Token)
		err := json.NewDecoder(resp.Body).Decode(token)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(token.Token)
	}

}
