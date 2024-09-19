package shared

import "github.com/PlagueByteSec/sentinel-project/v2/pkg"

type PoolBase struct {
	PoolIPv4Addresses         []string
	PoolIPv6Addresses         []string
	PoolSubdomains            []string
	PoolMailSubdomains        []string
	PoolApiSubdomains         []string
	PoolLoginSubdomains       []string
	PoolCorsSubdomains        []string
	PoolHttpSuccessSubdomains []string
	PoolCmsSubdomains         []string
	PoolCookieInjection       []string
	PoolHeaderInjection       []string
	PoolRequestSmuggling      []string
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

func PoolsInit(pools *PoolBase) {
	pools.PoolIPv4Addresses = make([]string, 0)
	pools.PoolIPv6Addresses = make([]string, 0)
	pools.PoolSubdomains = make([]string, 0)
	pools.PoolMailSubdomains = make([]string, 0)
	pools.PoolApiSubdomains = make([]string, 0)
	pools.PoolLoginSubdomains = make([]string, 0)
	pools.PoolCorsSubdomains = make([]string, 0)
	pools.PoolHttpSuccessSubdomains = make([]string, 0)
	pools.PoolCmsSubdomains = make([]string, 0)
}

func PoolsCleanupCore(pools *PoolBase) {
	pools.PoolIPv4Addresses = poolRemoveDuplicates(pools.PoolIPv4Addresses)
	pools.PoolIPv6Addresses = poolRemoveDuplicates(pools.PoolIPv6Addresses)
	pools.PoolSubdomains = poolRemoveDuplicates(pools.PoolSubdomains)
}

func PoolsCleanupSummary(pools *PoolBase) {
	pools.PoolMailSubdomains = poolRemoveDuplicates(pools.PoolMailSubdomains)
	pools.PoolApiSubdomains = poolRemoveDuplicates(pools.PoolApiSubdomains)
	pools.PoolLoginSubdomains = poolRemoveDuplicates(pools.PoolLoginSubdomains)
	pools.PoolCorsSubdomains = poolRemoveDuplicates(pools.PoolCorsSubdomains)
	pools.PoolHttpSuccessSubdomains = poolRemoveDuplicates(pools.PoolHttpSuccessSubdomains)
	pools.PoolCmsSubdomains = poolRemoveDuplicates(pools.PoolCmsSubdomains)
	pools.PoolCookieInjection = poolRemoveDuplicates(pools.PoolCookieInjection)
	pools.PoolHeaderInjection = poolRemoveDuplicates(pools.PoolHeaderInjection)
}

func PoolAppendValue(subdomain string, pool *[]string) {
	if !pkg.IsInSlice(subdomain, *pool) {
		*pool = append(*pool, subdomain)
	}
}
