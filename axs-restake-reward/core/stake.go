package core

import (
	"math/big"
)

type UserRewardResult struct {
	DebitedRewards   *big.Int
	CreditedRewards  *big.Int
	LastClaimedBlock *big.Int
}

const StakingManagerContract = "0x8bd81a19420bad681b7bfc20e703ebd8e253782d"
