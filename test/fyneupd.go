// fyneupd.go - Program to test updating Fyne GUI widgets.
// Mark Riordan  12-MAY-2024
package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func timeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func doTrivial(host string) (float64, error) {
	randomNumber := 100.0 * rand.Float64()
	return randomNumber, nil
}

func writeLog(msg string) {
	f, err := os.OpenFile("fyneupd.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	//rand.Seed(time.Now().UnixNano())
	a := app.New()
	w := a.NewWindow("Fyne Test")

	pingLabel := widget.NewLabel("")
	errorLabel := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	writeLog("fyneupd starting.")

	go func() {
		for {
			pingTime, err := doTrivial("google.com")
			strPingTime := fmt.Sprintf("%.1f", pingTime)
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

			if pingTime > 80.0 {
				errorLabel.Text = timeString() + " High number: " + strPingTime
				errorLabel.Refresh()
				fmt.Println(timeString() + " High High number: " + strPingTime)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	w.SetContent(container.NewVBox(pingLabel, errorLabel))
	w.ShowAndRun()
}
