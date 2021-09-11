package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func raiseConfigError() {
    fmt.Println("Error with config file. Run the program with \"-config\" flag again to fix it.")
    os.Exit(0)
}

func main() {
    configPtr := flag.Bool("config", false, "Configure the mouse positions")
    flag.Parse()

    if *configPtr {
        err := configureMousePositions()
        check(err)
        return
    }

    go func () {
        eventHook := gohook.Start()
        for e := range eventHook {
            if e.Kind == gohook.KeyDown {
                fmt.Println("Stoping auto clicker.")
                os.Exit(0)
            }
        }
    }()

    positions, ok := readConfigYaml()
    if !ok {
        raiseConfigError()
    }

    fmt.Println("Auto clicker started. Press any key to stop it.")
    fmt.Println("*clicking*")
    mainPosition := positions[0]
    secPositions := positions[1:]
    for si := 0; ; si = (si+1) % len(secPositions) {
        robotgo.MoveMouse(mainPosition.x, mainPosition.y)
        for i := 0; i < 15; i++ {
            robotgo.Click("left")
            time.Sleep(25 * time.Millisecond)
        }

        robotgo.MoveClick(secPositions[si].x, secPositions[si].y, "left")
        time.Sleep(25 * time.Millisecond)
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

func readConfigYaml() ([]mousePosition, bool) {
    yfile, err := ioutil.ReadFile("config.yaml")
    if err != nil {
        return nil, false
    }

    data := make([]map[string]int, 0)

    err = yaml.Unmarshal(yfile, &data)
    if err != nil {
        return nil, false
    }

    return mapToPositions(data), true
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

func positionsToMap(positions []mousePosition) []map[string]int{
    positionsMapList := make([]map[string]int, len(positions))
    for i, position := range positions {
        positionsMapList[i] = make(map[string]int)
        positionsMapList[i]["x"] = position.x
        positionsMapList[i]["y"] = position.y
    }
    return positionsMapList
}
