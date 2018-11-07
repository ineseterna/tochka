package marketplace

import (
	"fmt"
)

/*
	Models
*/

type VendorStats struct {
	NumberOfReleasedTransactions int
	NumberOfDisputes             int
	NumberOfWonDisputes          int
	NumberOfLostDisputes         int
	NumberOfWarnings             int
	NumberOfReviews              int
	NumberOfPositiveReviews      int
	NumberOfNegativeReviews      int
	NumberOfNeutralReviews       int
}

/*
	Cache
*/

func CacheGetVendorStats(userUuid string) VendorStats {

	queryStats := func() VendorStats {
		return VendorStats{
			NumberOfReleasedTransactions: CountNumberOfReleasedTransactionsForVendor(userUuid),
			NumberOfDisputes:             CountDisputesForUserUuid(userUuid, ""),
			NumberOfWonDisputes:          CountDisputesForUserUuid(userUuid, "RESOLVED TO VENDOR"),
			NumberOfLostDisputes:         CountDisputesForUserUuid(userUuid, "RESOLVED TO BUYER"),
			NumberOfWarnings:             CountWarningsForUser(userUuid),
			NumberOfReviews:              CountRatingReviewsForVendor(userUuid),
			NumberOfPositiveReviews:      CountPositiveRatingReviewsForVendor(userUuid),
			NumberOfNegativeReviews:      CountNegativeRatingReviewsForVendor(userUuid),
			NumberOfNeutralReviews:       CountNeutralRatingReviewsForVendor(userUuid),
		}
	}

	key := fmt.Sprintf("vendor-stats-%s", userUuid)
	cStats, _ := cache15m.Get(key)
	if cStats == nil {
		stats := queryStats()
		cache15m.Set(key, stats)
		return stats
	}

	return cStats.(VendorStats)
}
