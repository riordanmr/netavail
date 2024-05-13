// netavail.go - Simple network availability monitor.
// Pings a host periodically and reports the ping time in a window;
// also logs the result to netavail.log.
// Mark Riordan  08-MAY-2024
package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-ping/ping"
)

func timeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func pingHost(host string) (float64, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return -1.0, err
	}

	pinger.Count = 1
	err = pinger.Run()
	if err != nil {
		return -1.0, err
	}

	stats := pinger.Statistics()
	pingTime := 1000.0 * stats.AvgRtt.Seconds()

	return pingTime, nil
}

func writeLog(msg string) {
	f, err := os.OpenFile("netavail.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	logLine := fmt.Sprintf("%s %v\n", timeString(), msg)

	if _, err := f.WriteString(logLine); err != nil {
		fmt.Println(err)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Ping Monitor")

	pingLabel := widget.NewLabel("")
	errorLabel := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	hostname, _ := os.Hostname()
	writeLog("netavail starting on " + hostname)

	go func() {
		for {
			pingTime, err := pingHost("google.com")
			strPingTime := fmt.Sprintf("%.2f ms", pingTime)
			if err != nil {
				errorLabel.Text = timeString() + " Error: " + err.Error()
				// For some reason, Text object doesn't refresh automatically,
				// even though the Label object does. So we have to do it manually.
				errorLabel.Refresh()
				fmt.Println(timeString() + " Error: " + err.Error())
				writeLog("!! Error: " + err.Error())
				pingLabel.SetText(timeString() + " Error: " + err.Error())
				time.Sleep(2 * time.Second)
				continue
			}

			pingLabel.SetText(timeString() + " " + strPingTime)
			writeLog(strPingTime)

			if pingTime > 400.0 {
				errorLabel.Text = timeString() + " High ping time: " + strPingTime
				errorLabel.Refresh()
				fmt.Println(timeString() + " High ping time: " + strPingTime)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	w.SetContent(container.NewVBox(pingLabel, errorLabel))
	w.ShowAndRun()
}
