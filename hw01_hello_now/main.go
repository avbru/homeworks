package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	const layout = "2006-01-02 15:04:05 -0700 MST"
	ntpTime, err := ntp.Time("ru.pool.ntp.org")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("current time: %s\nexact time: %s\n", time.Now().Format(layout), ntpTime.Format(layout))
}
