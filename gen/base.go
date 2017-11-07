package gen

import (
	"os"
	"log"
	"bufio"
	"bytes"
	"strings"
	"io/ioutil"
)

//Из строки name=value возвращает name
func GetName(s string) (res string) {
	i := strings.Index(s, "=")
	if i > 0 {
		res = s[:i]
	}
	return res
}

//Из строки name=value возвращает value
func GetValue(s string) (res string) {
	i := strings.Index(s, "=")
	if i > 0 {
		res = s[i+1:]
	}
	return res
}

//Считываем файл .sett
func ReadSettFile(settFilePath string) (m map[string]string) {
	m = make(map[string]string)
	file, err := os.Open(settFilePath)
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
			bts := bytes.TrimSuffix([]byte(line), []byte("\n"))
			line = string(bts)
			m[GetName(line)] = GetValue(line)
		}
	}
	return m
}

//Считываем файл шаблона .tmpl
func ReadTemplateFile(templateFilePath string) (s string) {
	b, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		log.Println(err.Error())
	} else {
		s = string(b)
	}
	return s
}

//Заменяем {{Name}} в файле .tmpl на value из строки Name=value файла .sett
func ProcessTemplate(template string, m map[string]string) (s string) {
	s = template
	for key, value := range m {
		s = strings.Replace(s, "{{"+key+"}}", value, -1)
	}
	return s
}

//Сохраняем обработанные данные в -out
func SaveOutFile(outPath string, s string) {
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		f, err := os.Create(outPath)
		if err != nil {
			log.Println(err.Error())
		} else {
			defer f.Close()
			f.WriteString(s)
		}
	} else {
		log.Println("Файл "+outPath+" уже существует")
	}
}