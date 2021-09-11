package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	gohook "github.com/robotn/gohook"
	"gopkg.in/yaml.v3"
)

type mousePosition struct {
    x int16
    y int16
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
    fmt.Println(positions)
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
                position := mousePosition{e.X, e.Y}
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

    data := make([]map[string]int16, 0)

    err = yaml.Unmarshal(yfile, &data)
    check(err)

    return mapToPositions(data)
}

func positionsToMap(positions []mousePosition) []map[string]int16{
    positionsMapList := make([]map[string]int16, len(positions))
    for i, position := range positions {
        positionsMapList[i] = make(map[string]int16)
        positionsMapList[i]["x"] = position.x
        positionsMapList[i]["y"] = position.y
    }
    return positionsMapList
}

func mapToPositions(positionsMapList []map[string]int16) []mousePosition {
    positions := make([]mousePosition, len(positionsMapList))
    for i, position := range positionsMapList {
        positions[i] = mousePosition{position["x"], position["y"]}
    }
    return positions
}
