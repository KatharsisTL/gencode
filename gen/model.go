package gen

import (
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Field struct {
	Name string
	Type string
	AdditionPrefix string
	AdditionPostfix string
	Required bool
	UserDescr string
	MinIntVal int
	MaxIntVal int
}

type Model struct {
	Name string
	TableName string
	Fields []Field
}

//Сохраняет файл на жесткий диск
func (m *Model) Save(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			log.Println(err.Error())
		} else {
			defer f.Close()
			bytes, _ := json.Marshal(m)
			f.Write(bytes)
		}
	} else {
		log.Println("Файл "+filePath+" уже существует")
	}
}

//Считывает файл с жесткого диска
func (m *Model) Load(filePath string) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(bytes, m)
	}
}

//Генерация файла модели
func (m *Model) GenerateModelFile(filePath string) {
	s := "package models\n\ntype "+m.GetNameTitle()+" struct {\n"
	for _, fld := range m.Fields {
		s += "\t"+fld.Name+" "+fld.Type
		if fld.AdditionPrefix != "" {
			s += " `"+fld.AdditionPrefix
			if fld.AdditionPostfix != "" {
				s+=":\""+fld.AdditionPostfix+"\""
			}
			s+= "`"
		}
		s += "\n"
	}
	s += "}\n"

	if m.TableName != "" {
		s += "func (o *"+m.GetNameTitle()+") TableName() string {\n"
		s += "\treturn \""+m.TableName+"\"\n}\n"
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			log.Println(err.Error())
		} else {
			defer f.Close()
			f.WriteString(s)
		}
	} else {
		log.Println("Файл "+filePath+" уже существует")
	}
}

//Переделывает test в Test
func (m *Model) GetNameTitle() string {
	return strings.Title(strings.ToLower(m.Name))
}

func (m *Model) GetNameLower() string {
	return strings.ToLower(m.Name)
}

//Генерация файла таблицы
func (m *Model) GenerateTableFile(templateFilePath string, outFilePath string) {
	//Заполняем блок Entity
	entityStr := m.GetNameLower()+": {\n"
	for _, fld := range m.Fields {
		switch(fld.Type) {
		case "int":
			entityStr += "\t\t\t\t"+fld.Name+": 0,\n"
		case "string":
			entityStr += "\t\t\t\t"+fld.Name+": \"\",\n"
		}
	}
	entityStr += "\t\t\t},"
	tmpl := ReadTemplateFile(templateFilePath)
	tmpl = strings.Replace(tmpl, "{{Entity}}", entityStr, -1)
	//Заполняем блок Rules
	rulesStr := "rules: {\n"
	rulesStr += "\t\t\t},\n"
	tmpl = strings.Replace(tmpl, "{{Rules}}", rulesStr, -1)
	//Заполняем блок EntityTitle
	tmpl = strings.Replace(tmpl, "{{EntityTitle}}", m.GetNameTitle(), -1)

	SaveOutFile(outFilePath, tmpl)
}