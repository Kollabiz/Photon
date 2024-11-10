package Utils

import (
	"github.com/Kollabiz/GoColor"
)

func Log(str string) {
	GoColor.Print256("[LOG]: ", 0, 7)
	GoColor.Println256(str, 7, 0)
}

func LogError(str string) {
	GoColor.Print256("[ERR]: ", 0, 9)
	GoColor.Println256(str, 9, 0)
}

func LogWarning(str string) {
	GoColor.Print256("[WRN]: ", 0, 11)
	GoColor.Println256(str, 11, 0)
}

func LogSuccess(str string) {
	GoColor.Print256("[SUC]: ", 0, 10)
	GoColor.Println256(str, 10, 0)
}
