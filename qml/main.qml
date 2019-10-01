import QtQuick 2.0
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.3
import QtQuick.Controls.Material 2.0
import QtCharts 2.0

Rectangle{
    id: window
    width: 640
    height: 480
    visible: true

    Item {
        Text {
            id: element
            x: 101
            y: 28
            text: qsTr("Pingを打ちます")
            font.pixelSize: 41
        }

        TextField {
            id: inputObject
            x: 149
            y: 94
            width: 277
            height: 40
            text: qsTr("IPアドレスを入力してください")
            renderType: Text.QtRendering
        }

        Text {
            id: lossRate
            x: 125
            y: 307
            width: 407
            height: 53
            Connections {
                target: PingChartBridge
                onUpdateLossRate: lossRate.text = data
            }
            font.pixelSize: 41
        }

        Button {
            id: startPing
            x: 114
            y: 150
            text: qsTr("スタート")
            font.wordSpacing: 0.1
            onClicked: PingChartBridge.startPing(inputObject.text)
        }

        Button {
            id: stopPing
            x: 301
            y: 150
            text: qsTr("ストップ")
            font.wordSpacing: 0.1
            onClicked: PingChartBridge.stopPing()
        }

        Text {
            id: rtt
            x: 125
            y: 227
            width: 407
            height: 53
            font.pixelSize: 41
            Connections {
                target: PingChartBridge
                onAddPoint: rtt.text = x + ": "+ y + " ms"
            }
        }

    }
}






