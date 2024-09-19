package shared

import "github.com/PlagueByteSec/sentinel-project/v2/pkg"

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

func poolRemoveDuplicates(pool []string) []string {
	temp := make(map[string]bool)
	revisedPool := make([]string, 0, len(pool))
	for idx := 0; idx < len(pool); idx++ {
		value := pool[idx]
		if !temp[value] {
			revisedPool = append(revisedPool, value)
			temp[value] = true
		}
	}
	return revisedPool
}

func PoolsCleanupCore(pools *PoolBase) {
	// General
	pools.PoolIPv4Addresses = poolRemoveDuplicates(pools.PoolIPv4Addresses)
	pools.PoolIPv6Addresses = poolRemoveDuplicates(pools.PoolIPv6Addresses)
	pools.PoolSubdomains = poolRemoveDuplicates(pools.PoolSubdomains)
	pools.PoolHttpSuccessSubdomains = poolRemoveDuplicates(pools.PoolHttpSuccessSubdomains)
}

func PoolsCleanupSummary(pools *PoolBase) {
	// Purpose/Service detection
	pools.PoolMailSubdomains = poolRemoveDuplicates(pools.PoolMailSubdomains)
	pools.PoolApiSubdomains = poolRemoveDuplicates(pools.PoolApiSubdomains)
	pools.PoolLoginSubdomains = poolRemoveDuplicates(pools.PoolLoginSubdomains)
	pools.PoolCmsSubdomains = poolRemoveDuplicates(pools.PoolCmsSubdomains)
	// Security testing
	pools.PoolCorsSubdomains = poolRemoveDuplicates(pools.PoolCorsSubdomains)
	pools.PoolCookieInjection = poolRemoveDuplicates(pools.PoolCookieInjection)
	pools.PoolHeaderInjection = poolRemoveDuplicates(pools.PoolHeaderInjection)
	pools.PoolRequestSmuggling = poolRemoveDuplicates(pools.PoolRequestSmuggling)
}

func PoolAppendValue(subdomain string, pool *[]string) {
	if !pkg.IsInSlice(subdomain, *pool) {
		*pool = append(*pool, subdomain)
	}
}
