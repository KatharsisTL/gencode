package main

import (
	"os"
	"strings"
	"log"
	"bufio"
	"io/ioutil"
)

type InputParams struct {
	In string
	Out string
	Sett string
}

func main() {
	params := readInputParams(os.Args[1:])
	s := readInFile(params)
	m := readSettFile(params)
	s = processTemplate(s, m)
	saveOutFile(params, s)
}

//Считываем входные параметры
func readInputParams(args []string) InputParams {
	s := InputParams{}
	for _, elem := range args {
		if strings.Contains(elem, "-in=") {
			s.In = getValue(elem)
		} else if strings.Contains(elem, "-out") {
			s.Out = getValue(elem)
		} else if strings.Contains(elem, "-sett") {
			s.Sett = getValue(elem)
		}
	}
	return s
}

//Из строки name=value возвращает name
func getName(s string) (res string) {
	i := strings.Index(s, "=")
	if i > 0 {
		res = s[:i]
	}
	return res
}

//Из строки name=value возвращает value
func getValue(s string) (res string) {
	i := strings.Index(s, "=")
	if i > 0 {
		res = s[i+1:]
	}
	return res
}

//Считываем файл .sett
func readSettFile(params InputParams) (m map[string]string) {
	m = make(map[string]string)
	file, err := os.Open(params.Sett)
	if err != nil {
		log.Println(err.Error())
	} else {
		f := bufio.NewReader(file)
		b := true
		for b {
			line, err := f.ReadString('\n')
			if err != nil {
				b = false
			}
			m[getName(line)] = getValue(line[:len(line)-1])
		}
	}
	return m
}

//Считываем файл шаблона .tmpl
func readInFile(params InputParams) (s string) {
	b, err := ioutil.ReadFile(params.In)
	if err != nil {
		log.Println(err.Error())
	} else {
		s = string(b)
	}
	return s
}

//Заменяем {{Name}} в файле .tmpl на value из строки Name=value файла .sett
func processTemplate(template string, m map[string]string) (s string) {
	s = template
	for key, value := range m {
		s = strings.Replace(s, "{{"+key+"}}", value, -1)
	}
	return s
}

//Сохраняем обработанные данные в -out
func saveOutFile(params InputParams, s string) {
	if _, err := os.Stat(params.Out); os.IsNotExist(err) {
		f, err := os.Create(params.Out)
		if err != nil {
			log.Println(err.Error())
		} else {
			defer f.Close()
			f.WriteString(s)
		}
	} else {
		log.Println("Файл "+params.Out+" уже существует")
	}
}