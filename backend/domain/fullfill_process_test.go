package domain_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/tochamateusz/machine_auction/domain"
	"github.com/tochamateusz/machine_auction/domain/auction"
)

func TestNewFullfil(t *testing.T) {

	t.Run("Should not be finish", func(t *testing.T) {
		mockAcution := auction.NewAuction("id", " image", " name", " year", " price", "endDate")
		process := domain.NewFullfillPorcess(*mockAcution)

		assert.False(t, process.IsDone())

	})
}
