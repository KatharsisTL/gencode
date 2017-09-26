# gencode

Консольное приложение для генерации текстовых файлов по шаблону, написанное на [go](https://golang.org/).

## Установка

```console
/*go get -u github.com/KatharsisTL/gencode*/
git clone https://github.com/KatharsisTL/gencode.git
cd ./gencode
go build ./
sudo ln -s ./gencode /usr/local/bin/gencode
```

## Описание

Для генерации файлов необходим файл шаблона .tmpl, файл настроек .sett и имя выходного файла, в который будет записан результат.

### Параметры приложения

Значения параметров начинаются со знака "минус" ```-```, имя параметра отделяется от значения знаком "равно" ```=```.

Все параметры являются обязательными.

| Параметр | Описание |
| --- | --- |
| -in      | Пусть к входному файлу шаблона ```.tmpl```. Можно указывать как абсолютный путь ```/test.tmpl```, так и относительный ```./test.tmpl``` |
| -sett    | Путь к файлу настроек |
| -out     | Путь к результирующему файлу |

### Файл шаблона

Файл содержит текст, который попадет в выходной файл без изменений и экранированные слова, например ``` {{EntityName}} ```, которым задается соответствие замены в файле настроек .sett

### Файл настроек

Файл содержит строки, задающие соответстие экранированных слов в шаблоне .tmpl и слов, которые необходимо вставить в шаблон. Имя и значение разделяются знаком "равно" ```=``` без пробелов.
```
EntityName=TestEntity
EntityDescription=Test description
```

### Пример

**Файл шаблона template.tmpl**

```go

type {{EntityName}}Entity struct {
    Name string
    {{EntityParam}} string
}

func (e *{{EntityName}}Entity) GetName() string {
    return e.Name
}

func (e *{{EntityName}}Entity) Get{{EntityParam}}() string {
    return e.{{EntityParam}}
}
```

**Файл настроек settings.sett**

```
EntityName=Test
EntityParam=Description
```

**Использование gencode**

```console
gencode -in=./template.tmpl -sett=./settings.sett -out=./TestEntity.go
```

**Результат TestEntity.go**

```go
type TestEntity struct {
    Name string
    Description string
}

func (e *TestEntity) GetName() string {
    return e.Name
}

func (e *TestEntity) GetDescription() string {
    return e.Description
}
```

