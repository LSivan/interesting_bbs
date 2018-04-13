package test

import (
	"testing"
	"time"
	"fmt"
)

func TestTimeFeature(t *testing.T) {
	now := time.Now()
	lastYear := now.AddDate(0,0,-2)
	fmt.Println(now.Truncate(24 * time.Hour).Sub(lastYear.Truncate(24 * time.Hour)).Hours() / 24)
}