package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	gohook "github.com/robotn/gohook"
	"gopkg.in/yaml.v3"
)

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
    configPtr := flag.Bool("config", false, "Configure the mouse positions")
    flag.Parse()

    if *configPtr {
        err := configureMousePositions()
        check(err)
        return
    }

    positions := readConfigYaml()
    mainPosition := positions[0]
    secPositions := positions[1:]
    breakPoint := 0
    for si := 0; ; si = (si+1) % len(secPositions) {
        robotgo.MoveMouse(mainPosition.x, mainPosition.y)
        fmt.Println("Cliking 15 times at main position")
        for i := 0; i < 15; i++ {
            robotgo.Click("left")
            time.Sleep(25 * time.Millisecond)
        }

        fmt.Printf("Cliking at position %d\n", si)
        robotgo.MoveClick(secPositions[si].x, secPositions[si].y, "left")
        time.Sleep(25 * time.Millisecond)

        breakPoint++
        if breakPoint == 50 {
            break
        }
    }
}

func configureMousePositions() error {
    positions, ok := getMousePositions()
    if !ok {
        return errors.New("Error during mouse positions registration.")
    }
    writeConfigYaml(positions)
    return nil
}

func getMousePositions() ([]mousePosition, bool) {
    fmt.Println("Click on the positions to register, starting by the main position. Right click to cancel.")
    positions := make([]mousePosition, 0)

    eventHook := gohook.Start()
    for e := range eventHook {
        if e.Kind == gohook.MouseDown {
            switch e.Button {
            case gohook.MouseMap["left"]:
                position := mousePosition{int(e.X), int(e.Y)}
                positions = append(positions, position)
                fmt.Printf("Position x:%d y:%d registered\n", e.X, e.Y)
            // for some reason, my mouse's right button returns code 3, which
            // is the code of the center button
            case gohook.MouseMap["center"]:
                fmt.Println("End of positions registration")
                return positions, true
            }
        }
    }
    return positions, false
}

func writeConfigYaml(positions []mousePosition) {
    positionsMap := positionsToMap(positions)

    data, err := yaml.Marshal(&positionsMap)
    check(err)

    err = ioutil.WriteFile("config.yaml", data, 0664)
    check(err)
}

func readConfigYaml() []mousePosition {
    yfile, err := ioutil.ReadFile("config.yaml")
    check(err)

    data := make([]map[string]int, 0)

    err = yaml.Unmarshal(yfile, &data)
    check(err)

    return mapToPositions(data)
}

func positionsToMap(positions []mousePosition) []map[string]int{
    positionsMapList := make([]map[string]int, len(positions))
    for i, position := range positions {
        positionsMapList[i] = make(map[string]int)
        positionsMapList[i]["x"] = position.x
        positionsMapList[i]["y"] = position.y
    }
    return positionsMapList
}

func mapToPositions(positionsMapList []map[string]int) []mousePosition {
    positions := make([]mousePosition, len(positionsMapList))
    for i, position := range positionsMapList {
        positions[i] = mousePosition{position["x"], position["y"]}
    }
    return positions
}
