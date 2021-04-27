package grant

import (
	"fmt"
	"net/http"

	"github.com/llinder/svtauth/internal/model"
)

type GrantError struct {
	StatusCode    int
	Message       string
	OriginalError error
}

func (grantErr *GrantError) Error() string {
	return fmt.Sprintf("Grant failed with Error: %v", grantErr.OriginalError)
}

func DoGrant(client *http.Client, grant *model.GrantRequest) (string, error) {

	caller, err := getCallerIdentity(client, grant)
	if err != nil {
		return "", err
	} else {
		fmt.Println(caller.ARN)

		// TODO create token
		return "foo", err
	}

}
