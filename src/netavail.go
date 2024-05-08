// netavail.go - Simple network availability monitor.
// Pings a host periodically and reports the ping time in a window;
// also logs the result to netavail.log.
// Mark Riordan  08-MAY-2024
package main

import (
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-ping/ping"
)

func timeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func pingHost(host string) (string, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return "", err
	}

	pinger.Count = 1
	err = pinger.Run()
	if err != nil {
		return "", err
	}

	stats := pinger.Statistics()
	pingTime := fmt.Sprintf("%v", stats.AvgRtt)

	return pingTime, nil
}

func writeLog(pingTime string) {
	f, err := os.OpenFile("netavail.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	logLine := fmt.Sprintf("%s %s\n", timeString(), pingTime)

	if _, err := f.WriteString(logLine); err != nil {
		fmt.Println(err)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Ping Monitor")

	pingLabel := widget.NewLabel("")

	go func() {
		for {
			pingTime, err := pingHost("google.com")
			if err != nil {
				fmt.Println(timeString() + " Error: " + err.Error())
				writeLog("!! Error: " + err.Error())
				pingLabel.SetText(timeString() + " Error: " + err.Error())
				time.Sleep(1 * time.Second)
				continue
			}

			pingLabel.SetText(timeString() + " " + pingTime)
			writeLog(pingTime)

			time.Sleep(5 * time.Second)
		}
	}()

	w.SetContent(container.NewVBox(pingLabel))
	w.ShowAndRun()
}
