package grant

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/llinder/svtauth/internal/model"
)

type CallerIdentity struct {
	XMLName xml.Name `xml:"GetCallerIdentityResult"`
	ARN     string   `xml:"Arn"`
	UserId  string   `xml:"UserId"`
	Account int      `xml:"Account"`
}

type callerIdentityResponse struct {
	XMLName        xml.Name       `xml:"GetCallerIdentityResponse"`
	CallerIdentity CallerIdentity `xml:"GetCallerIdentityResult"`
}

type errorDetail struct {
	XMLName xml.Name `xml:"Error"`
	Type    string   `xml:"Type"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

type errorResponse struct {
	XMLName   xml.Name    `xml:"ErrorResponse"`
	Error     errorDetail `xml:"Error"`
	RequestId string      `xml:"RequestId"`
}

func getCallerIdentity(client *http.Client, grant *model.GrantRequest) (ident CallerIdentity, err error) {
	fmt.Printf("identity url %s\n", grant.CallerIdentityUrl)

	resp, err := client.Get(grant.CallerIdentityUrl)
	if err != nil {
		return ident, &GrantError{0, "Get caller identity call failed", err}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ident, &GrantError{resp.StatusCode, "Failed to read response", err}
	}

	if resp.StatusCode == 200 {
		result := callerIdentityResponse{}
		unmarshalError := xml.Unmarshal(body, &result)
		if unmarshalError != nil {
			return ident, &GrantError{resp.StatusCode, "Failed to decode caller identity respose", err}
		}
		fmt.Println(result.CallerIdentity.ARN)

		// TODO grant token with
		return result.CallerIdentity, err

	} else {
		result := errorResponse{}
		unmarshalError := xml.Unmarshal(body, &result)
		if unmarshalError != nil {
			return ident, &GrantError{
				resp.StatusCode,
				"Get caller identity failed with unmarshalling error",
				unmarshalError,
			}
		} else {
			return ident, &GrantError{
				resp.StatusCode,
				fmt.Sprintf("Get caller identity failed with %s and message %s", result.Error.Code, result.Error.Message),
				errors.New("error response"),
			}
		}
	}
}
