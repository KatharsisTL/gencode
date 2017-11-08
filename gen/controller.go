package gen

/*import (
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
	"strings"
)

type TemplateSettingsString struct {
	Name string
	Val string
}

type ControllerSettings struct {
	TemplatePath string
	TemplateSettingsPath string
	OutPath string
	templateSettings []TemplateSettingsString
}

//Сохраняет файл на жесткий диск
func (c *ControllerSettings) Save(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			log.Println(err.Error())
		} else {
			defer f.Close()
			bytes, _ := json.Marshal(c)
			f.Write(bytes)
		}
	} else {
		log.Println("Файл "+filePath+" уже существует")
	}

	if len(c.templateSettings) > 0 {
		if c.TemplateSettingsPath != "" {
			if _, err := os.Stat(c.TemplateSettingsPath); os.IsNotExist(err) {
				f, err := os.Create(c.TemplateSettingsPath)
				if err != nil {
					log.Println(err.Error())
				} else {
					defer f.Close()
					for num, ss := range c.templateSettings {
						s := ""
						if num == 0 {
							s += ss.Name+"="
						} else {
							s += "\n"+ss.Name+"="
						}
						f.WriteString(s)
					}
				}
			} else {
				log.Println("Файл "+c.TemplateSettingsPath+" уже существует")
			}
		}
	}
}

//Считывает файл с жесткого диска
func (c *ControllerSettings) Load(filePath string) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(bytes, c)
		if c.TemplateSettingsPath != "" {
			c.ReadSettingsStringsFromTemplate()
		}
	}
}

//Считываем имена переменных из шаблона
func (c *ControllerSettings) ReadSettingsStringsFromTemplate() {
	b, err := ioutil.ReadFile(c.TemplatePath)
	if err != nil {
		log.Println(err.Error())
	} else {
		s := string(b)
		//log.Println(s)
		b := true

		for b {
			i := strings.Index(s, "{{")
			if i > -1 {
				s = s[i+2:]
				j := strings.Index(s, "}}")
				if j > -1 {
					//Проверяем текущие строки
					str := s[0:j]
					//log.Println(str)
					b1 := true
					for _, st := range c.templateSettings {
						if st.Name == str {
							b1 = false
						}
					}
					if b1 {
						ss := TemplateSettingsString{Name: str}
						c.templateSettings = append(c.templateSettings, ss)
					}

					//Обрезаем строку
					s = s[j+2:]
				}
				//b = false
			} else {
				b = false
			}
		}
	}
}

func (c *ControllerSettings) GenerateControllerFile() {
	m := ReadSettFile(c.TemplateSettingsPath)
	tmpl := ReadTemplateFile(c.TemplatePath)
	s1 := ProcessTemplate(tmpl, m)
	SaveOutFile(c.OutPath, s1)
}*/