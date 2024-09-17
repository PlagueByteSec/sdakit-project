package analysis

func (check *SubdomainCheck) Misconfigurations() {
	check.hostHeaders()     // Host header injections
	check.cookieInjection() // session hijacking, xss
	check.CORS()
}

func (check *SubdomainCheck) PurposeHTTP() {
	check.api()   // Content types, API versions, rate limit
	check.login() // Scan response body for login indicators
	check.cms()   // Top 20 CMS
}

func (check *SubdomainCheck) PurposeNonHTTP() {
	check.MailServer()
}
