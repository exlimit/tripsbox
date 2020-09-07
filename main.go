package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/martian/log"
	"github.com/sirupsen/logrus"
)

//BusLine -x
type BusLine struct {
	Name   string    `json:"name"`
	ID     string    `json:"id"`
	Vendor int       `json:"vendor"`
	Path   []Point   `json:"path"`
	Time   []float64 `json:"timestamps"`
}

//Point xy
type Point []float64

func main() {
	// AppendDatasets 添加datasets目录下未入库的数据集
	//遍历dir目录下所有.mbtiles
	counter := 0
	dir := "../data"
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Error(err)
	}

	var lines []BusLine
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		fileName := item.Name()
		ext := filepath.Ext(fileName)
		name := strings.TrimSuffix(strings.TrimPrefix(fileName, "20161212_"), ext)
		switch ext {
		case ".txt":
			pathfile := filepath.Join(dir, fileName)
			file, err := os.Open(pathfile)
			if err != nil {
				logrus.Warning(err)
				continue
			}
			defer file.Close()

			reader := csv.NewReader(file)
			reader.Comma = '\t'
			rows, err := reader.ReadAll()
			if err != nil {
				fmt.Print(err)
			}

			var bus BusLine
			var path []Point
			var times []float64
			for i, r := range rows {
				if bus.ID != r[2] {
					if i != 0 {
						bus.Name = name
						bus.Vendor = i % 2
						bus.Path = path
						bus.Time = times
						lines = append(lines, bus)
						//clean
						bus.ID = r[2]
						path = []Point{}
						times = []float64{}
					}
				}
				y, _ := strconv.ParseFloat(r[0], 64)
				x, _ := strconv.ParseFloat(r[1], 64)
				pt := Point{x, y}
				path = append(path, pt)
				bus.ID = r[2]
				tv, _ := strconv.ParseFloat(r[4][8:12], 64)
				times = append(times, tv)
			}
		}
	}

	jsondata, err := json.Marshal(lines) // convert to JSON
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(jsondata))

	jsonFile, err := os.Create("../data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsondata)
	log.Infof("proc %d datasets ~", counter)
}
