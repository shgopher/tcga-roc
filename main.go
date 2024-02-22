/*
 * @Author: shgopher shgopher@gmail.com
 * @Date: 2024-02-21 20:47:01
 * @LastEditors: shgopher shgopher@gmail.com
 * @LastEditTime: 2024-02-22 21:03:14
 * @FilePath: /r1/main.go
 * @Description:
 *
 * Copyright (c) 2024 by shgopher, All Rights Reserved.
 */
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var (
	timePath = "./b.tsv"
	genPath  = "./gen"
	jsonPath = "./a.json"
	genName ="PTK2"
)

func main() {
	cs := NewCase(jsonPath)
	for _, c := range cs {
		c.GetGens(genName)
		c.GetTime(timePath)
	}
	Build(Only(cs))
}

type Case struct {
	CaseID              string // 病人的id
	FileName            string // 病人id对应的cart 文件名
	GenFpkmUqUnstranded string // 某基因的表达值
	Time                string // 存活时间
	Alive               bool   // 是否活着
	pass                bool   // 判断数据是否可用
}

// 传入要处理的metadata文件地址，输出一个 CaseID 的结构体类型
func NewCase(f string) []*Case {
	data, err := os.ReadFile(f)
	if err != nil {
		fmt.Println("开头读取json错误", err)
	}
	cas := parseJSON(data)
	fmt.Println("在原始数据的meta中一共读取了", len(cas), "个数据")
	return cas
}

// 读取 caseid 中的 id，然后找到 filename中对应基因的表达值
func (c *Case) GetGens(gensName string) {
	// 通过fileName，找到相对应csv文件中的基因表达值
	pathmeta := retrieveData(genPath)
	var realPath string
	if value, ok := ifFile(c.FileName, pathmeta); ok {
		realPath = value
	} else {
		fmt.Println("数据是什么？", value)
		panic("未找到文件，错误！！！")
	}
	file, _ := os.Open(realPath)
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("在读取gen时的 tsv 文件时错误，错误信息是：", err)
	}
	var isFind int
	for _, record := range records {
		if len(record) > 8 && record[1] == gensName {
			c.GenFpkmUqUnstranded = record[8]
			isFind++
		}
	}
	if isFind != 1 {
		panic("未找到基因，错误！！！")
	}
}

// 读取case id 中的id，找到对应的存活时间，以及是否存活
func (c *Case) GetTime(filePath string) {
	// vital_status 是第15行（均从零开始算）
	// days_to_death 是第9行
	// days_to_follow 是第49行

	file, _ := os.Open(filePath)
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	// 读取第一行并跳过
	if _, err := reader.Read(); err != nil {
		fmt.Println("第一行跳过错误", err)
	}
	records, _ := reader.ReadAll()

	var isFind int
	for _, record := range records {
		if c.CaseID == record[0] {
			isFind++

			if record[15] == "Alive" {
				c.Alive = true
			} else if record[15] == "Dead" {
				c.Alive = false
			} else {
				break
			}

			if c.Alive {
				d := deathNumber(record[49])

				if d != "0" {
					num, _ := strconv.ParseFloat(record[49], 64)
					//fmt.Println("测试num1：",num,num/365)
					c.Time = fmt.Sprintf("%f", num/365)
				}
			} else {
				d := deathNumber(record[9])

				if d != "0" {
					num, _ := strconv.ParseFloat(record[9], 64)
					//fmt.Println("测试num2：",num,num/365)
					c.Time = fmt.Sprintf("%f", num/365)
				}
			}
			if isFind == 1 {
				break
			}
		}

	}
	if isFind == 0 {
		fmt.Println(isFind)
		panic("未找到用户名称，错误！！！")
	}
}

// 输出为vsc格式，按照 id 存货时间，基因表达值
func Build(cas []*Case) {

	// 创建CSV文件
	csvFile, err := os.Create("data.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	// 创建CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// 写入表头
	header := []string{"caseid", "futime", "fustat", "gene"}
	err = writer.Write(header)
	if err != nil {
		panic(err)
	}
	// 写入数据
	i := 0
	for _, cas := range cas {
		if cas.Time == "" {
			continue
		}
		a := "0"
		if cas.Alive {
			a = "0"
		} else {
			a = "1"
		}
		row := []string{
			cas.CaseID,
			cas.Time,
			a,
			cas.GenFpkmUqUnstranded,
		}
		err = writer.Write(row)
		if err != nil {
			panic(err)
		}
		i++
	}
	fmt.Println("可用数据", i)
}

// 解析函数
func parseJSON(data []byte) []*Case {
	var cas []*Case
	// 定义Response切片
	var resp Response

	// 解析JSON到Response

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(data, &resp); err != nil {
		panic(err)
	}
	// 遍历切片解析每个Item
	for _, item := range resp {
		// 访问文件名
		var c = new(Case)
		c.FileName = item.FileName
		c.CaseID = item.AssociatedEntities[0].CaseID
		cas = append(cas, c)
		// 遍历AssociatedEntities获取case id
	}
	return cas
}

// 定义结构体匹配JSON结构
type Entity struct {
	CaseID string `json:"case_id"`
}

type Item struct {
	FileName           string   `json:"file_name"`
	AssociatedEntities []Entity `json:"associated_entities"`
}

type Response []Item

func retrieveData(path string) []string {
	var values []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// if the file is noe regular, it mean the file is done,you should return
		if !info.Mode().IsRegular() {
			return nil
		}
		values = append(values, path)
		return nil
	})
	if err != nil {
		fmt.Println("读取路径错误", err)
	}

	return values
}

func ifFile(fileName string, path []string) (string, bool) {
	for _, v := range path {
		if filepath.Base(v) == fileName {
			return v, true
		}
	}
	return "", false
}

func deathNumber(number string) string {

	_, err := strconv.ParseFloat(number, 64)
	if err != nil {
		fmt.Println("转化是什么错误：", err, number)
		return "0"
	}

	return "1"
}

func Only(cases []*Case) []*Case {
	ma := make(map[string]int)
	r := make([]*Case, 0)
	for _, v := range cases {
		if ma[v.CaseID] < 1 {
			r = append(r, &Case{
				GenFpkmUqUnstranded: v.GenFpkmUqUnstranded,
				Time:                v.Time,
				Alive:               v.Alive,
				pass:                v.pass,
				CaseID:              v.CaseID,
				FileName:            v.FileName,
			})
			ma[v.CaseID]++
		}
	}
	return r
}
