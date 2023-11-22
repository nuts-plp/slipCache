package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myapp := app.New()
	www := myapp.NewWindow("hello world")
	www.SetContent(widget.NewLabel("jksajda"))
	www.ShowAndRun()

}
