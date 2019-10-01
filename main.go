package main

import (
	"log"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"
)

func main() {
	core.QCoreApplication_SetApplicationName("qt_test")
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	gui.NewQGuiApplication(len(os.Args), os.Args)
	view := quick.NewQQuickView(nil)
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)

	p := NewPingTimeChart(nil)

	view.RootContext().SetContextProperty("PingChartBridge", p.bridge)
	log.Println("set context property")
	view.SetSource(core.NewQUrl3("qrc:/qml/main.qml", 0))
	log.Println("set source")

	view.Show()
	log.Println("show")
	gui.QGuiApplication_Exec()
}
