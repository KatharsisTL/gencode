package gen

import (
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
)

type Field struct {
	Name string
	Type string
	AdditionPrefix string
	AdditionPostfix string
	//Блок rules
	Required bool
	//Участвует в формировании текста сообщения
	UserDescr string
	//Если тип 'int'
	//Если мин или макс не равны нулю, то создается правило на макс и мин ограничение
	MinIntVal int
	MaxIntVal int
	//Если необходима кастомная сортировка
	//Имя поля в Js
	JsFldName string
	//Имя поля в бд с указанием имени внешней таблицы (при необходимости)
	SortFldDbName string

	//Формирование контроллера
	//Флаг, участвует ли поле в проверке создания объекта
	CheckCreate bool
	//Флаг, участвует ли поле в проверке редактирования объекта
	CheckUpdate bool
}

type Model struct {
	Name string
	UserDescr string
	TableName string
	Fields []Field
	ControllerSettings ControllerSettings
}

type ControllerSettings struct {
	//go. bool-выражение
	GetRight string
	//sql
	SelectOrder string
	//go. bool-выражение
	UpdateRight string

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
	//Заполняем блок Entity---------------------------------------------------------------
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
	//-----------------------------------------------------------------------------------

	//Заполняем блок Rules---------------------------------------------------------------
	rulesStr := "rules: {\n"
	//Для каждого поля
	for _, fld := range m.Fields {
		if fld.Required {
			rulesStr += "\t\t\t\t"+strings.Title(fld.Name)+": [\n\t\t\t\t\t{required: true, message: \"Введите "+fld.UserDescr+"\"},\n"
			if fld.Type == "int" {
				rulesStr += "\t\t\t\t\t{type: 'number', message: 'Введите число'},\n"
				if fld.MinIntVal != 0 || fld.MaxIntVal != 0 {
					rulesStr += "\t\t\t\t\t{type: 'number', min:"+strconv.Itoa(fld.MinIntVal)+", max:"+strconv.Itoa(fld.MaxIntVal)+", message: 'Введите число в диапазоне от "+strconv.Itoa(fld.MinIntVal)+" до "+strconv.Itoa(fld.MaxIntVal)+"'},\n"
				}
			}
			rulesStr += "\t\t\t\t],\n"
		}
	}
	rulesStr += "\t\t\t},\n"
	tmpl = strings.Replace(tmpl, "{{Rules}}", rulesStr, -1)
	//-----------------------------------------------------------------------------------

	//Заполняем блок EntityTitle---------------------------------------------------------
	tmpl = strings.Replace(tmpl, "{{EntityTitle}}", m.GetNameTitle(), -1)
	//-----------------------------------------------------------------------------------

	//Заполняем блок Sort----------------------------------------------------------------
	sortStr := ""
	for _, fld := range m.Fields {
		if fld.SortFldDbName != "" {
			if fld.JsFldName != "" {
				sortStr += "\t\t\t\tcase \""+strings.Title(fld.JsFldName)+"\":\n"
			} else {
				sortStr += "\t\t\t\tcase \""+strings.Title(fld.Name)+"\":\n"
			}
			sortStr += "\t\t\t\t\tthis.sort.prop = \""+fld.SortFldDbName+"\";\n\t\t\t\t\tbreak;\n"
		}
	}
	tmpl = strings.Replace(tmpl, "{{Sort}}", sortStr, -1)
	//-----------------------------------------------------------------------------------

	//Создаем функцию change-------------------------------------------------------------
	changeStr := "\t\tpublic change(id: number = 0) {\n"
	changeStr += "\t\t\t//Если добавление записи\n\t\t\tif(id == 0) {\n"
	/*for _, fld := range m.Fields {
		if fld.Type == "string" {
			changeStr += "\t\t\t\tthis.editDialog."+strings.ToLower(m.Name)+"."+strings.Title(fld.Name)+" = \"\";\n"
		} else if fld.Type == "int" {
			changeStr += "\t\t\t\tthis.editDialog."+strings.ToLower(m.Name)+"."+strings.Title(fld.Name)+" = 0;\n"
		}
	}*/
	changeStr += "\t\t\t\tthis.editDialog.title = \""+strings.Title(m.UserDescr)+". Добавление\";\n"
	changeStr += "\t\t\t} else {\n"
	changeStr += "\t\t\t\t//Если редактирование записи\n"
	changeStr += "\t\t\t\tFetch.Json(this.getOneUrl(id)).then((json)=>{\n"
	changeStr += "\t\t\t\t\tthis.editDialog."+strings.ToLower(m.Name)+" = json;\n"
	changeStr += "\t\t\t\t\tthis.editDialog.title = `"+strings.Title(m.UserDescr)+" #${json.Id}. Редактирование`;\n"
	changeStr += "\t\t\t\t});\n"

	changeStr += "\t\t\t}\n\t\t\tthis.editDialog.visible = true;\n\t\t}\n"
	tmpl = strings.Replace(tmpl, "{{Change}}", changeStr, -1)
	//-----------------------------------------------------------------------------------

	//Создаем функцию changeClose--------------------------------------------------------
	changeCloseStr := "\t\tpublic changeClose() {\n"
	changeCloseStr += "\t\t\tthis.editDialog.title = \"\";\n"
	for _, fld := range m.Fields {
		if fld.Type == "string" {
			changeStr += "\t\t\t\tthis.editDialog."+strings.ToLower(m.Name)+"."+strings.Title(fld.Name)+" = \"\";\n"
		} else if fld.Type == "int" {
			changeStr += "\t\t\t\tthis.editDialog."+strings.ToLower(m.Name)+"."+strings.Title(fld.Name)+" = 0;\n"
		}
	}
	changeCloseStr += "\t\t\tthis.editDialog.visible = false;\n"
	changeCloseStr += "\t\t}\n"
	tmpl = strings.Replace(tmpl, "{{ChangeClose}}", changeCloseStr, -1)
	//-----------------------------------------------------------------------------------

	tmpl = strings.Replace(tmpl, "{{EntityName}}", strings.ToLower(m.Name), -1)

	delDialogFillStr := "this.delDialog.title = \""+strings.Title(m.UserDescr)+". Удаление\";\n"
	delDialogFillStr += "\t\t\tthis.delDialog.message = `Вы действительно хотите удалить объект \""+strings.ToLower(m.UserDescr)+" #${json.Id}\"`;"
	tmpl = strings.Replace(tmpl, "{{DelDialogFill}}", delDialogFillStr, -1)

	//Сохраняем сформированный файл------------------------------------------------------
	SaveOutFile(outFilePath, tmpl)
}

//Генерация файла объекта
func (m *Model) GenerateObjFile(templateFilePath string, outFilePath string) {
	tmpl := ReadTemplateFile(templateFilePath)

	tmpl = strings.Replace(tmpl, "{{EntityTitle}}", m.GetNameTitle(), -1)
	tmpl = strings.Replace(tmpl, "{{EntityName}}", m.GetNameLower(), -1)

	//Сохраняем сформированный файл------------------------------------------------------
	SaveOutFile(outFilePath, tmpl)
}

//Генерация файла контроллера
func (m *Model) GenerateControllerFile(templateFilePath string, outFilePath string) {
	tmpl := ReadTemplateFile(templateFilePath)

	tmpl = strings.Replace(tmpl, "{{EntityName}}", m.GetNameTitle(), -1)
	tmpl = strings.Replace(tmpl, "{{GetRight}}", m.ControllerSettings.GetRight, -1)
	tmpl = strings.Replace(tmpl, "{{SelectOrder}}", m.ControllerSettings.SelectOrder, -1)
	tmpl = strings.Replace(tmpl, "{{UpdateRight}}", m.ControllerSettings.UpdateRight, -1)

	//Формирование проверки создания объекта-------------------------------------------------------------------------
	createWhereString := ""
	createWhereParams := ""
	for _, fld := range m.Fields {
		if fld.CheckCreate {
			pre := ""
			pre1 := ""
			if createWhereString != "" {
				pre = " and "
				pre1 = ", "
			}
			if fld.SortFldDbName != "" {
				createWhereString += pre+strings.ToLower(fld.SortFldDbName)+" = ?"
			} else {
				createWhereString += pre+strings.ToLower(fld.Name)+" = ?"
			}
			createWhereParams += pre1+"obj."+strings.Title(fld.Name)
		}
	}
	tmpl = strings.Replace(tmpl, "{{CreateWhereString}}", createWhereString, -1)
	tmpl = strings.Replace(tmpl, "{{CreateWhereParams}}", createWhereParams, -1)
	//---------------------------------------------------------------------------------------------------------------

	//Формирование проверки редактирования объекта-------------------------------------------------------------------------
	updateWhereString := ""
	updateWhereParams := ""
	for _, fld := range m.Fields {
		if fld.CheckUpdate {
			pre := ""
			pre1 := ""
			if updateWhereString != "" {
				pre = " and "
				pre1 = ", "
			}
			if fld.SortFldDbName != "" {
				updateWhereString += pre+strings.ToLower(fld.SortFldDbName)+" = ?"
			} else {
				updateWhereString += pre+strings.ToLower(fld.Name)+" = ?"
			}
			updateWhereParams += pre1+"obj."+strings.Title(fld.Name)
		}
	}
	if updateWhereString != "" {
		updateWhereString += " and id <> ?"
	} else {
		updateWhereString += " id <> ?"
	}

	if updateWhereParams != "" {
		updateWhereParams += ", obj.Id"
	} else {
		updateWhereParams += "obj.Id"
	}
	tmpl = strings.Replace(tmpl, "{{UpdateWhereString}}", updateWhereString, -1)
	tmpl = strings.Replace(tmpl, "{{UpdateWhereParams}}", updateWhereParams, -1)
	//---------------------------------------------------------------------------------------------------------------

	//Сохраняем сформированный файл------------------------------------------------------
	SaveOutFile(outFilePath, tmpl)
}