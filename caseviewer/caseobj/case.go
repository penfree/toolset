package caseviewer

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

//wrap a json object
type EmrJsonDoc struct {
	Data map[string]interface{}
}

//get attribute from json
func (this *EmrJsonDoc) Get(key, defValue string) (value interface{}) {
	if v, ok := this.Data[key]; ok {
		value = v
	} else {
		value = defValue
	}
	return value
}

func (this *EmrJsonDoc) HasKey(key string) (hasKey bool) {
	if _, ok := this.Data[key]; ok {
		hasKey = true
	} else {
		hasKey = false
	}
	return hasKey
}

//基本信息
type EmrBasic struct {
	EmrJsonDoc
}

//诊断
type EmrDiagnosis struct {
	EmrJsonDoc
}

//医嘱
type EmrOrder struct {
	EmrJsonDoc
}

//化验
type EmrLab struct {
	EmrJsonDoc
}

//检查
type EmrExam struct {
	EmrJsonDoc
}

//生命体征
type EmrVital struct {
	EmrJsonDoc
}

//emr
type EmrDoc struct {
	EmrJsonDoc
}

//一份完整病历
type EmrCase struct {
	ID        string
	Info      EmrBasic
	Diagnosis []EmrDiagnosis
	Order     []EmrOrder
	Lab       []EmrLab
	Exam      []EmrExam
	Vital     []EmrVital
	Doc       []EmrDoc
}

func (this *EmrCase) load(jsonStr, cf, col string) {
	jsonStr = jsonStr[:len(jsonStr)-3]
	var data [](map[string]interface{})
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		data = nil
		fmt.Println(err)
	}
	if cf == "info" && col == "visit" && len(data) != 0 {
		this.Info = EmrBasic{EmrJsonDoc{data[0]}}
	} else if cf == "info" && col == "diagnosis" {
		for _, v := range data {
			this.Diagnosis = append(this.Diagnosis, EmrDiagnosis{EmrJsonDoc{v}})
		}
	} else if cf == "order" && col == "data" {
		for _, v := range data {
			this.Order = append(this.Order, EmrOrder{EmrJsonDoc{v}})
		}
	} else if cf == "lab" && col == "report" {
		for _, v := range data {
			this.Lab = append(this.Lab, EmrLab{EmrJsonDoc{v}})
		}
	} else if cf == "exam" && col == "data" {
		for _, v := range data {
			this.Exam = append(this.Exam, EmrExam{EmrJsonDoc{v}})
		}
	} else if cf == "exam" && col == "vital" {
		for _, v := range data {
			this.Vital = append(this.Vital, EmrVital{EmrJsonDoc{v}})
		}
	} else if cf == "emr" && col == "data" {
		for _, v := range data {
			this.Doc = append(this.Doc, EmrDoc{EmrJsonDoc{v}})
		}
	}
}

//一个病人的病历
type PatientCase struct {
	Case map[string]*EmrCase
}

func (this PatientCase) B64Encode(s string) (encodedString string) {
	encodedString = base64.StdEncoding.EncodeToString([]byte(s))
	return encodedString
}

//load json string to map
func (this *PatientCase) Load(filePath string) {
	this.Case = make(map[string]*EmrCase)
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		words := strings.Split(string(line), string('\t'))
		if len(words) == 4 {
			visitID := words[0]
			if _, ok := this.Case[visitID]; ok {
				this.Case[visitID].load(words[3], words[1], words[2])
			} else {
				this.Case[visitID] = &EmrCase{}
				this.Case[visitID].load(words[3], words[1], words[2])
			}
		}
		if err == io.EOF {
			break
		}
	}
}
