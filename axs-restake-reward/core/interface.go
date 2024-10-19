package core

import "math/big"

type coreManager interface {
	GetBalance() (*big.Float, error)
	GetStakingAmount() (*big.Float, error)
	GetTotalStaked() (*big.Float, error)
	GetUserRewardInfo() (UserRewardResult, error)
	CreatePeriodicalTelegramMessage()
}
