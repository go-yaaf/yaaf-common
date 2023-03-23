// Token Utils tests

package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/utils"
	"testing"
	"time"
)

func TestBounds(t *testing.T) {

	format := "YYYY-MMM-DD HH:mm:ss.sss"
	now := entity.Now()
	tu := utils.TimeUtils(now.Add(time.Hour))
	fmt.Println("Now    ", tu.Format(format))

	from := tu.LowerBound(time.Minute).Format(format)
	to := tu.UpperBound(time.Minute).Format(format)
	fmt.Println("Minute ", from, to)

	from = tu.LowerBound(time.Hour).Format(format)
	to = tu.UpperBound(time.Hour).Format(format)
	fmt.Println("Hour   ", from, to)

	from = tu.LowerBound(time.Hour * 24).Format(format)
	to = tu.UpperBound(time.Hour * 24).Format(format)
	fmt.Println("Day    ", from, to)

}

func TestSeries(t *testing.T) {

	format := "YYYY-MMM-DD HH:mm:ss.sss"
	// 60 minutes series
	//printSeries(time.Hour, time.Minute, format)
	//printFrames(time.Hour, time.Minute, format)

	// 24 hors series
	//printSeries(time.Hour*24, time.Hour, format)
	printFrames(time.Hour*24, time.Hour, format)

	// 7 days series
	//printSeries(time.Hour*24*7, time.Hour*24, format)
	//printFrames(time.Hour*24*7, time.Hour*24, format)

}

func printSeries(period time.Duration, interval time.Duration, format string) {

	from := utils.TimeUtils(entity.Now()).LowerBound(period).Get()
	to := utils.TimeUtils(entity.Now()).UpperBound(period).Get()
	minutes := utils.TimeUtils(from).Series(to, interval)
	for i, m := range minutes {
		fmt.Println(i, m, utils.TimeUtils(m).Format(format))
	}

	fmt.Println("-------------------")

	from = utils.TimeUtils(entity.Now()).UpperBound(period).Get()
	to = utils.TimeUtils(entity.Now()).LowerBound(period).Get()
	minutes = utils.TimeUtils(from).Series(to, interval)
	for i, m := range minutes {
		fmt.Println(i, m, utils.TimeUtils(m).Format(format))
	}
}

func printFrames(period time.Duration, interval time.Duration, format string) {

	from := utils.TimeUtils(entity.Now()).LowerBound(period).Get()
	to := utils.TimeUtils(entity.Now()).UpperBound(period).Get()
	frames := utils.TimeUtils(from).TimeFrames(to, interval)
	for i, f := range frames {
		fmt.Println(i, f.String(format))
	}

	fmt.Println("-------------------")

	from = utils.TimeUtils(entity.Now()).UpperBound(period).Get()
	to = utils.TimeUtils(entity.Now()).LowerBound(period).Get()
	frames = utils.TimeUtils(from).TimeFrames(to, interval)
	for i, f := range frames {
		fmt.Println(i, f.String(format))
	}
}
