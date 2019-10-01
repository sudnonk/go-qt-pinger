package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/quick"
)

func main() {
	core.QCoreApplication_SetApplicationName("go-qt-pinger")
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	gui.NewQGuiApplication(len(os.Args), os.Args)
	view := quick.NewQQuickView(nil)
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)

	p := NewPingTimeChart(nil)

	view.RootContext().SetContextProperty("PingChartBridge", p.bridge)
	view.SetSource(core.NewQUrl3("qrc:/qml/main.qml", 0))

	view.Show()
	gui.QGuiApplication_Exec()
}
