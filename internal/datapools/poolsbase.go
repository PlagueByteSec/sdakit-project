package datapools

type PoolBase struct {
	// General
	PoolIPv4Addresses         []string
	PoolIPv6Addresses         []string
	PoolSubdomains            []string
	PoolHttpSuccessSubdomains []string
	// Purpose/Service detection
	PoolMailSubdomains  []string
	PoolApiSubdomains   []string
	PoolLoginSubdomains []string
	PoolCmsSubdomains   []string
	// Security testing
	PoolCorsSubdomains   []string
	PoolCookieInjection  []string
	PoolHeaderInjection  []string
	PoolRequestSmuggling []string
}

func PoolsInit(pools *PoolBase) {
	// General
	pools.PoolIPv4Addresses = make([]string, 0)
	pools.PoolIPv6Addresses = make([]string, 0)
	pools.PoolSubdomains = make([]string, 0)
	pools.PoolHttpSuccessSubdomains = make([]string, 0)
	// Purpose/Service detection
	pools.PoolMailSubdomains = make([]string, 0)
	pools.PoolApiSubdomains = make([]string, 0)
	pools.PoolLoginSubdomains = make([]string, 0)
	pools.PoolCmsSubdomains = make([]string, 0)
	// Security testing
	pools.PoolCorsSubdomains = make([]string, 0)
	pools.PoolHeaderInjection = make([]string, 0)
	pools.PoolCookieInjection = make([]string, 0)
	pools.PoolRequestSmuggling = make([]string, 0)
}
