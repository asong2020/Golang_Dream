package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var (
	result map[string]*requestBody
	analysis map[string]*requestBody
)

type requestBody struct {
	count int32
	query string
	time float64
}

type cellValue struct {
	sheet string
	cell string
	value string
}

func main()  {
	result = make(map[string]*requestBody,100)
	analysis = make(map[string]*requestBody,100)
	file := openFile()
	logDeal(file)
	analysisBody()
	exportExcel()
}

func openFile() *os.File {
	file,err := os.Open("./request.log")
	if err != nil{
		log.Println("open log err: ",err)
	}
	return file
}

func logDeal(file *os.File)  {
	// 按行读取
	br := bufio.NewReader(file)
	for{
		line,_,err := br.ReadLine()
		// file read complete
		if err == io.EOF{
			log.Println("file read complete")
			return
		}
		//json deal
		var data interface{}
		err = json.Unmarshal(line,&data)
		if err != nil{
			fmt.Errorf("json marshal error")
		}
		deal(data)
	}
}

func deal(data interface{})  {
	var request string
	var query string
	var time float64
	value,ok := data.(map[string]interface{})
	if ok{
		for k,v := range value{
			if k == "httpRequest"{
				switch v1 := v.(type) {
				case map[string]interface{}:
					for k1,v11 := range v1{
						if k1 == "request"{
							switch val := v11.(type) {
							case string:
								request = val
								//fmt.Println(request)
							}
						}
					}
				}
			}
			if k == "params"{
				switch v1 := v.(type) {
				case map[string]interface{}:
					for k1,v11 := range v1{
						if k1 == "query"{
							switch val := v11.(type) {
							case string:
								query = val
								//fmt.Println(query)
							}
						}
					}
				}
			}
			if k == "timings"{
				switch v1 := v.(type) {
				case map[string]interface{}:
					for k1,v11 := range v1{
						if k1 == "evalTotalTime"{
							switch val := v11.(type) {
							case float64:
								time = val
							//	fmt.Println(time)
							}
						}
					}
				}
			}
		}
		b := &requestBody{
			query: query,
			time: time,
		}
		if _,o := result[request];o{
			b.count = result[request].count + 1
			b.time = b.time + result[request].time
			result[request] = b
		}else {
			b.count = 1
			result[request] = b
		}
	}
}
//analysis data
func analysisBody()  {
	for k,v := range result{
		req := &requestBody{}
		req.time = v.time / float64(v.count)
		req.count = v.count
		req.query = v.query
		analysis[k] = req
	}
}

//export excel
func exportExcel()  {
	file := excelize.NewFile()
	//insert title
	cellValues := make([]*cellValue,0)
	cellValues = append(cellValues,&cellValue{
		sheet: "sheet1",
		cell: "A1",
		value: "request",
	},&cellValue{
		sheet: "sheet1",
		cell: "B1",
		value: "count",
	},&cellValue{
		sheet: "sheet1",
		cell: "C1",
		value: "query",
	},&cellValue{
		sheet: "sheet1",
		cell: "D1",
		value: "avgTime",
	})
	index := file.NewSheet("Sheet1")
	// 设置工作簿的默认工作表
	file.SetActiveSheet(index)
	for _, cellValue := range cellValues {
		file.SetCellValue(cellValue.sheet, cellValue.cell, cellValue.value)
	}
	//insert data
	cnt := 1
	for k,v := range analysis{
		cnt = cnt + 1
		for k1,v1 := range cellValues{
			switch k1 {
			case 0:
				v1.cell = fmt.Sprintf("A%d",cnt)
				v1.value = k
			case 1:
				v1.cell = fmt.Sprintf("B%d",cnt)
				v1.value = fmt.Sprintf("%d",v.count)
			case 2:
				v1.cell = fmt.Sprintf("C%d",cnt)
				v1.value = v.query
			case 3:
				v1.cell = fmt.Sprintf("D%d",cnt)
				v1.value = strconv.FormatFloat(v.time,'f',-1,64)
			}
		}
		for _,vc := range cellValues{
			file.SetCellValue(vc.sheet,vc.cell,vc.value)
		}
	}

	//generate file
	err := file.SaveAs("./log.xlsx")
	if err != nil{
		fmt.Errorf("generate excel error")
	}
}