package auth

import _ "embed"

const (
	RuleAny            = "rule_any"
	RuleAdminOnly      = "rule_admin_only"
	RuleUserOnly       = "rule_user_only"
	RuleAdminOrSubject = "rule_admin_or_subject"
)

//go:embed rego/authorization.rego
var regoAuthorizationScript string
