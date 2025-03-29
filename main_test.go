package persiandate_go_test

import (
	"testing"

	persiandate "github.com/NothingMotion/PersianDate-GO"
)

func TestPersianDate(t *testing.T) {
	pd := persiandate.NewPersianDate("YYYY/MM/DD")

	t.Log(pd.JalaliFullNow())
}
