package core

import "math/big"

type CoreManager interface {
	GetBalance() (*big.Float, error)
	GetStakingAmount() (*big.Float, error)
	GetTotalStaked() (*big.Float, error)
	GetUserRewardInfo() (UserRewardResult, error)
	AutoCompoundRewards() (string, error)
}

type UserRewardResult struct {
	DebitedRewards   *big.Int
	CreditedRewards  *big.Int
	LastClaimedBlock *big.Int
}
