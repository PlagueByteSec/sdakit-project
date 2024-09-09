package analysis

func (check *SubdomainCheck) Misconfigurations() {
	check.hostHeaders()         // Host header injections
	check.cookieInjectionPath() // session hijacking, xss
}

func (check *SubdomainCheck) Purpose() {
	check.mailServer()   // General, Exchange
	check.api()          // Content types, API versions, rate limit
	check.login()        // Scan response body for login indicators
	check.basicWebpage() // HTML response (no other results)
	check.cms()          // Top 20 CMS
}
