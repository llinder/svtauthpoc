package model

type GrantRequest struct {
	CallerIdentityUrl string `form:"caller_identity_url" binding:"required"`
}

type ErrorResponse struct {
	RequestId string
	Error     ErrorDetail
}

type ErrorDetail struct {
	Code    string
	Message string
	Details []ErrorDetail
}
