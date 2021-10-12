package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	gohook "github.com/robotn/gohook"
	"gopkg.in/yaml.v3"
)

const SCREEN = 1
const MIN_X = 3470
const MIN_Y = 114
const MAX_X = 3480
const MAX_Y = 124
var joinCh chan struct{}

type mousePosition struct {
    x int
    y int
}

func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    minPoint := image.Point{MIN_X, MIN_Y}
    maxPoint := image.Point{MAX_X, MAX_Y}
    bounds := image.Rectangle{minPoint, maxPoint}

    go func () {
        eventHook := gohook.Start()
        for e := range eventHook {
            if e.Kind == gohook.KeyDown {
                fmt.Println("Stoping auto clicker.")
                os.Exit(0)
            }
        }
    }()

    for {
        prestige := prestigeLocationChanged(bounds)

        ascentionPositions, ok := readConfigYaml("ascending-config.yaml")
        if !ok {
            raiseConfigError()
        }


        upgradePositions, ok := readConfigYaml("upgrade-config.yaml")
        if !ok {
            raiseConfigError()
        }

        if prestige {
            execAscention(ascentionPositions)
        } else {
            upgrade(upgradePositions)
        }

        time.Sleep(time.Second)
    }
}

func execAscention(positions []mousePosition) {
    fmt.Println("Ascending")

    // Clink in "legacy"
    robotgo.MoveClick(positions[0].x, positions[0].y, "left")
    time.Sleep(100 * time.Millisecond)

    // Click "yes" to confirm
    robotgo.MoveClick(positions[1].x, positions[1].y, "left")
    // Wait the animation
    time.Sleep(7 * time.Second)

    // Click in "reincarnate"
    robotgo.MoveClick(positions[2].x, positions[2].y, "left")
    time.Sleep(100 * time.Millisecond)

    // Click "yes" to confirm
    robotgo.MoveClick(positions[3].x, positions[3].y, "left")
    time.Sleep(1 * time.Second)

    // Click to close the popup notifications
    robotgo.MoveClick(positions[4].x, positions[4].y, "left")
    time.Sleep(1 * time.Second)
}

func upgrade(positions []mousePosition) {
    fmt.Println("Upgrading")

    setBuy100Position := positions[0]
    upgradesPositions := positions[1:]

    robotgo.MoveClick(setBuy100Position.x, setBuy100Position.y, "left")

    i := 0
    for ; i < len(upgradesPositions); i++ {
        time.Sleep(25 * time.Millisecond)
        robotgo.MoveClick(upgradesPositions[i].x, upgradesPositions[i].y, "left")
    }

    // This is for wait more and try buy upgrades more times
    i--
    for j := 0; j < 3; j++ {
        time.Sleep(50 * time.Millisecond)
        robotgo.MoveClick(upgradesPositions[i].x, upgradesPositions[i].y, "left")
    }

}

func readConfigYaml(filename string) ([]mousePosition, bool) {
    yfile, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, false
    }

    data := make([]map[string]int, 0)

    err = yaml.Unmarshal(yfile, &data)
    if err != nil {
        return nil, false
    }

    if len(data) == 0 {
        raiseConfigError()
    }

    return mapToPositions(data), true
}

func raiseConfigError() {
    fmt.Println("Error with config file.")
    os.Exit(0)
}

func mapToPositions(positionsMapList []map[string]int) []mousePosition {
    positions := make([]mousePosition, len(positionsMapList))
    for i, position := range positionsMapList {
        for k := range position {
            if k != "x" && k != "y" {
                raiseConfigError()
            }
        }
        if _, ok := position["x"]; !ok {
            raiseConfigError()
        }
        if _, ok := position["y"]; !ok {
            raiseConfigError()
        }
        positions[i] = mousePosition{position["x"], position["y"]}
    }
    return positions
}

func prestigeLocationChanged(bounds image.Rectangle) bool {
    file, err := os.Open("/home/iago/Documents/repos/golang-autoclicker/cookie-auto-ascend/no_prestige.png")
    check(err)
    defer file.Close()

    fileImg, err := png.Decode(file)
    check(err)
    noPrestige := fileImg.(*image.RGBA)

    currentImage, err := screenshot.CaptureRect(bounds)
    check(err)

    for i := 0; i < len(currentImage.Pix); i++ {
        if currentImage.Pix[i] != noPrestige.Pix[i] {
            return true
        }
    }

    return false
}

//func configureMousePositions() error {
    //positions, ok := getMousePositions()
    //if !ok {
        //return errors.New("Error during mouse positions registration.")
    //}
    //writeConfigYaml(positions)
    //return nil
//}

//func getMousePositions() ([]mousePosition, bool) {
    //fmt.Println("Click on the positions to register, starting by the main position. Right click to cancel.")
    //positions := make([]mousePosition, 0)

    //eventHook := gohook.Start()
    //for e := range eventHook {
        //if e.Kind == gohook.MouseDown {
            //switch e.Button {
            //case gohook.MouseMap["left"]:
                //position := mousePosition{int(e.X), int(e.Y)}
                //positions = append(positions, position)
                //fmt.Printf("Position x:%d y:%d registered\n", e.X, e.Y)
            //// for some reason, my mouse's right button returns code 3, which
            //// is the code of the center button
            //case gohook.MouseMap["center"]:
                //fmt.Println("End of positions registration")
                //return positions, true
            //}
        //}
    //}
    //return positions, false
//}

//func writeConfigYaml(positions []mousePosition) {
    //positionsMap := positionsToMap(positions)

    //data, err := yaml.Marshal(&positionsMap)
    //check(err)

    //err = ioutil.WriteFile("aaa.yaml", data, 0664)
    //check(err)
//}

//func positionsToMap(positions []mousePosition) []map[string]int{
    //positionsMapList := make([]map[string]int, len(positions))
    //for i, position := range positions {
        //positionsMapList[i] = make(map[string]int)
        //positionsMapList[i]["x"] = position.x
        //positionsMapList[i]["y"] = position.y
    //}
    //return positionsMapList
//}
