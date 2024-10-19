package core

import (
	"math/big"
	"time"
)

type CoreManager interface {
	GetBalance(string) (*big.Float, error)
	GetStakingAmount(string) (*big.Float, error)
	GetTotalStaked() (*big.Float, error)
	GetUserRewardInfo() (UserRewardResult, error)
	AutoCompoundRewards() (string, error)
	GetLastClaimedTime() time.Time
}

type UserRewardResult struct {
	DebitedRewards   *big.Int
	CreditedRewards  *big.Int
	LastClaimedBlock *big.Int
}
