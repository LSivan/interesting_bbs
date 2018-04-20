package test

import (
	"testing"
	"time"
	"fmt"
	"math"
)

func TestTimeFeature(t *testing.T) {
	now := time.Now()
	lastYear := now.AddDate(0,0,0)
	fmt.Println(math.Log10(now.Truncate(24 * time.Hour).Sub(lastYear.Truncate(24 * time.Hour)).Hours()+11))
}