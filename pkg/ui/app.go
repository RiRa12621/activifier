package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"activifier/pkg/jiggle"
)

func Run() {
	a := app.New()
	w := a.NewWindow("Activifier")

	// interval slider settings
	const (
		minSeconds = 5
		maxSeconds = 300
	)

	jiggler := jiggle.New(30 * time.Second)

	status := widget.NewLabel("Status: stopped")
	intervalLabel := widget.NewLabel("Interval: 30s")

	slider := widget.NewSlider(minSeconds, maxSeconds)
	slider.Step = 1
	slider.Value = 30

	var startStopBtn *widget.Button
	var desk desktop.App
	var haveDesk bool
	var trayStartStop *fyne.MenuItem

	secsFromSlider := func() int {
		return int(slider.Value + 0.5)
	}

	buildTrayMenu := func() {
		if !haveDesk {
			return
		}
		// Rebuild menu to ensure label updates show up reliably.
		m := fyne.NewMenu("Activifier",
			fyne.NewMenuItem("Show", func() { w.Show() }),
			trayStartStop,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				jiggler.Stop()
				a.Quit()
			}),
		)
		desk.SetSystemTrayMenu(m)
	}

	syncUI := func() {
		period := jiggler.Period()

		// Keep the label consistent with actual stored period.
		intervalLabel.SetText(fmt.Sprintf("Interval: %ds", int(period.Seconds()+0.5)))

		if jiggler.IsRunning() {
			status.SetText(fmt.Sprintf("Status: running (every %s)", period))
			startStopBtn.SetText("Stop")
			if trayStartStop != nil {
				trayStartStop.Label = "Stop"
			}
		} else {
			status.SetText(fmt.Sprintf("Status: stopped (set to %s)", period))
			startStopBtn.SetText("Start")
			if trayStartStop != nil {
				trayStartStop.Label = "Start"
			}
		}

		buildTrayMenu()
	}

	setIntervalFromSlider := func() {
		secs := secsFromSlider()
		intervalLabel.SetText(fmt.Sprintf("Interval: %ds", secs))

		jiggler.SetPeriod(time.Duration(secs) * time.Second)
		syncUI()
	}

	toggle := func() {
		period := time.Duration(secsFromSlider()) * time.Second
		if jiggler.IsRunning() {
			jiggler.Stop()
		} else {
			jiggler.Start(period)
		}
		syncUI()
	}
	
	slider.OnChanged = func(_ float64) {
		setIntervalFromSlider()
	}

	startStopBtn = widget.NewButton("Start", nil)
	startStopBtn.OnTapped = toggle

	w.SetContent(container.NewVBox(
		widget.NewLabel("Move mouse 1px up and 1px down at a configurable interval."),
		intervalLabel,
		slider,
		container.NewHBox(startStopBtn),
		status,
	))

	// Hide on close (tray style)
	w.SetCloseIntercept(func() { w.Hide() })

	// Tray menu
	if d, ok := a.(desktop.App); ok {
		desk = d
		haveDesk = true

		trayStartStop = fyne.NewMenuItem("Start", nil)
		trayStartStop.Action = toggle

		buildTrayMenu()
	}

	// init UI state from default slider value / jiggler period
	setIntervalFromSlider()
	w.ShowAndRun()
}
