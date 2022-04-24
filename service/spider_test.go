package service

import (
	"github.com/paleblueyk/logger"
	"testing"
)

func TestGetPrizeInformation(t *testing.T) {
	logger.Info(GetNextNum())
}

func TestGetNewPrize(t *testing.T) {
	GetNewPrize()
}