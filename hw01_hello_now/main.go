package main

import (
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	const layout = "2006-01-02 15:04:05 -0700 MST" //1945-05-09 10:03:00 +0000 UTC
	ntpTime, err := ntp.Time("ru.pool.ntp.org")
	if err != nil {
		panic(err)
	}

	fmt.Printf("current time: %s\nexact time: %s\n", time.Now().Format(layout), ntpTime.Format(layout))
}
