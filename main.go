package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Sprav struct {
	ID           int    `json:"ID"`
	City         string `json:"City"`
	Organization string `json:"Organization"`
	Dolgnost     string `json:"Dolgnost"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"lastName"`
	MiddleName   string `json:"MiddleName"`
	Contacts     string `json:"Contacts"`
	Email        string `json:"Email"`
}
type Poiscsprav struct {
	Poisc string `json:"poisc"`
}

var router *gin.Engine
var db *sql.DB

//var autorUsername string

const connectionString = "test.db"

func main() {

	var e error
	db, e = sql.Open("sqlite3", connectionString)
	if e != nil {
		fmt.Println(e)
		return
	}
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.Static("/assets/", "front/")
	router.LoadHTMLGlob("html/templates/*")
	// Initialize the routes
	initializeRoutes()
	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// функция выводит все данные из базы
func getAllArticles() []Sprav {

	row, e := db.Query(`SELECT "ID", TRIM("City"), TRIM("Organization"), TRIM("Dolgnost"), TRIM("FirstName"), TRIM("LastName"), TRIM("MiddleName"), TRIM("Contacts"), TRIM("Email")
FROM "Spravochnic" ORDER BY "Organization" ASC, "FirstName" ASC, "LastName" ASC, "MiddleName" ASC`)

	if e != nil {
		fmt.Println("Ошибка при чтении таблицы getAllArticles", e)
		//return e
	}

	defer row.Close()

	products := []Sprav{}

	for row.Next() {
		p := Sprav{}
		err := row.Scan(&p.ID, &p.City, &p.Organization, &p.Dolgnost, &p.FirstName, &p.LastName, &p.MiddleName, &p.Contacts, &p.Email)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	return products
}

// функция выводит все данные из базы
func (p Poiscsprav) postPoiscAllSpravochnik() []Sprav {

	row, e := db.Query(`SELECT ID, City, Organization, Dolgnost, FirstName, LastName, MiddleName,
       Contacts, Email FROM Spravochnic WHERE City = $1 OR Organization = $2 OR Dolgnost = $3 OR FirstName = $4 OR LastName = $5 OR MiddleName = $6 OR Contacts =$7 OR Email = $8`,
		p.Poisc, p.Poisc, p.Poisc, p.Poisc, p.Poisc, p.Poisc, p.Poisc, p.Poisc)

	if e != nil {
		fmt.Println("Ошибка при чтении таблицы postPoiscAllSpravochnik", e)
		return nil
	}

	defer row.Close()

	products := []Sprav{}

	for row.Next() {
		s := Sprav{}
		err := row.Scan(&s.ID, &s.City, &s.Organization, &s.Dolgnost, &s.FirstName, &s.LastName, &s.MiddleName, &s.Contacts, &s.Email)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, s)
	}

	return products
}

// Обнавление значений в базе справочник
func (s Sprav) postUpdateId() []Sprav {
	//res, err := strconv.Atoi(s.ID)
	//if err != nil {
	//	panic(err)
	//}
	row := db.QueryRow("SELECT ID, City, Organization, Dolgnost, FirstName, LastName, MiddleName, Contacts, Email FROM Spravochnic WHERE ID =$1", s.ID)

	products := []Sprav{}

	p := Sprav{}
	err := row.Scan(&p.ID, &p.City, &p.Organization, &p.Dolgnost, &p.FirstName, &p.LastName, &p.MiddleName, &p.Contacts, &p.Email)
	if err != nil {
		fmt.Println(err)
	}

	products = append(products, p)

	return products
}

// Главная страница с выводои всех значений из базы
func handlerIndex(c *gin.Context) {

	products := getAllArticles()

	render(c, gin.H{
		"Title":    "Телефонный справочник",
		"Classind": "current",
		"payload":  products}, "index.html")
}

// Регистрация нового юзера get запрос
func handlerRegistration(c *gin.Context) {
	render(c, gin.H{
		"Title":    "Страница регистрации",
		"Classreg": "current",
	}, "registration.html")
}

// Авторизация юзеров get запрос
func handlerAuthorization(c *gin.Context) {
	render(c, gin.H{
		"Title": "Страница авторизации",
	}, "authorization.html")
}

// Регистрация нового юзера post запрос c проверкой заполнения всех полей
func handlerUserRegistration(c *gin.Context) {
	var user User

	e := c.Bind(&user)

	if e != nil {
		c.HTML(200, "registration.html", gin.H{
			"ErrorTitle":   "Oшибка ",
			"Classreg":     "current",
			"ErrorMessage": "Произошла ошибка попробуйте её "})
		return
	}

	if user.Login == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" {
		c.HTML(200, "registration.html", gin.H{
			"ErrorTitle":   "Oшибка ",
			"Classreg":     "current",
			"ErrorMessage": "Все поля должны быть заполнены"})
		return
	}

	if len(user.Login) < 3 || len(user.Password) < 3 {
		c.HTML(200, "registration.html", gin.H{
			"ErrorTitle":   "Oшибка ",
			"Classreg":     "current",
			"ErrorMessage": "Минимальная длина логина и пароля 3 знака"})
		return
	}

	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])

	e = user.Create()
	if e != nil {
		c.HTML(200, "registration.html", gin.H{
			"ErrorTitle":   "Oшибка ",
			"Classreg":     "current",
			"ErrorMessage": "Введите другие данные"})
		return
	}

	//c.Set("is_logged_reg", true)

	c.HTML(200, "registration.html", gin.H{
		"ErrorTitle":   " ",
		"Classreg":     "current",
		"ErrorMessage": "Вы успешно зарегистрированы ",
		"ErrorM":       user.FirstName,
		"ErrorN":       user.LastName,
		"RegTitle":     "Ok",
	})

}

// Авторизация юзеров post запрос c проверкой
func handlerUserAuthorization(c *gin.Context) {
	var user User
	user.Login = c.PostForm("Login")
	user.Password = c.PostForm("Password")
	//fmt.Println("Выводим user = ", user)

	if user.Login == "" || user.Password == "" {
		c.HTML(200, "authorization.html", gin.H{
			"ErrorTitle":   "Oшибка входа",
			"Classavt":     "current",
			"ErrorMessage": "Требуется ввести и логин и пароль"})
		return
	}

	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])

	res, e := user.Select()
	if e != nil {
		c.HTML(200, "authorization.html", gin.H{
			"ErrorTitle":   "Oшибка входа",
			"Classavt":     "current",
			"ErrorMessage": "Предоставлены неверные учетные данные"})
		return
	}

	if res == "" {
		return
	}
	user.LastName = res

	// Получаем параметр запроса "username"
	username := c.PostForm("Login")
	if username != "" {
		SetCookie(c, username)
	}
	AuthUser(c)

	//c.Set("is_logged_in", true)

	render(c, gin.H{
		"Title":        "Страница авторизации",
		"ErrorTitle":   "Сообщение",
		"ErrorMessage": "Вы успешно авторизированы " + user.LastName,
		"AvtorizTitle": "ec",
		"Classavt":     "current",
	}, "authorization.html")

}

func EditPage(c *gin.Context) {
	//c.Set("is_logged_in", true)
	id := c.Param("id")
	res, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	row := db.QueryRow("SELECT ID, City, Organization, Dolgnost, FirstName, LastName, MiddleName, Contacts, Email FROM Spravochnic WHERE ID =$1", res)

	products := []Sprav{}

	p := Sprav{}
	err = row.Scan(&p.ID, &p.City, &p.Organization, &p.Dolgnost, &p.FirstName, &p.LastName, &p.MiddleName, &p.Contacts, &p.Email)
	if err != nil {
		render(c, gin.H{
			"Title":   "Редактирование справочника",
			"Message": "Нет данных в базе по этому адресу",
		}, "update.html")
		return
	}
	products = append(products, p)

	render(c, gin.H{
		"Title":   "Редактирование справочника",
		"payload": products}, "update.html")

}

func handlerUpdate(c *gin.Context) {
	var sprav Sprav

	e := c.Bind(&sprav)

	if e != nil {
		products := sprav.postUpdateId()
		render(c, gin.H{
			"Title":   "Редактирование справочника",
			"Message": "Произошла ошибка попробуйте её ",
			"payload": products}, "update.html")
		return
	}
	if sprav.LastName == "" || sprav.FirstName == "" || sprav.Contacts == "" {
		products := sprav.postUpdateId()
		render(c, gin.H{
			"Title":   "Редактирование справочника",
			"Message": "Поля Фамилия, Имя, Телефон не должны быть пустыми ",
			"payload": products}, "update.html")
		return
	}

	er := sprav.Update()
	if er != nil {
		products := sprav.postUpdateId()
		render(c, gin.H{
			"Title":   "Редактирование справочника",
			"Message": "Произошла ошибка попробуйте её ",
			"payload": products}, "update.html")
		return
	}

	products := sprav.postUpdateId()

	render(c, gin.H{
		"Title":   "Редактирование справочника",
		"Message": "Данные успешно сохранены",
		"payload": products}, "update.html")
}

func handlerDelete(c *gin.Context) {
	var sprav Sprav
	/*
		    id := c.Param("id")
			res, err := strconv.Atoi(id)
			if err != nil {
				panic(err)
			}
	*/
	e := c.Bind(&sprav)
	products := sprav.postUpdateId()

	if e != nil {
		render(c, gin.H{
			"Title":   "Редактирование справочника",
			"Message": "Произошла ошибка попробуйте её ",
			"payload": products}, "update.html")
		return
	}

	er := sprav.Delete()
	if er != nil {
		render(c, gin.H{
			"Title":   "Редактирование справочника",
			"Message": "Произошла ошибка попробуйте её "},
			"update.html")
		return
	}

	render(c, gin.H{
		"Title":   "Редактирование справочника",
		"Message": "Данные успешно удалены",
	}, "update.html")
}

func postNewContactSpravochnik(c *gin.Context) {

	var sprav Sprav

	e := c.Bind(&sprav)

	if e != nil {
		render(c, gin.H{
			"Title":    "Добавить новый контакт",
			"Classnew": "current",
			"Message":  "Произошла ошибка попробуйте её раз ",
		}, "newcontact.html")
		return
	}
	if sprav.LastName == "" || sprav.FirstName == "" || sprav.Contacts == "" {

		render(c, gin.H{
			"Title":    "Добавить новый контакт",
			"Classnew": "current",
			"Message":  "Поля Фамилия, Имя, Телефон не должны быть пустыми ",
		}, "newcontact.html")
		return
	}

	res, er := sprav.New()
	if er != nil {
		c.HTML(500, "newcontact.html", gin.H{
			"Title":    "Редактирование справочника",
			"Classnew": "current",
			"Message":  "Произошла ошибка попробуйте ещё"})
		return
	}

	sprav.ID = res

	products := sprav.postUpdateId()

	render(c, gin.H{
		"Title":    "Редактирование справочника",
		"Classnew": "current",
		"Message":  "Данные успешно сохранены",
		"payload":  products}, "newcontact.html")

}

func getNewContactSpravochnik(c *gin.Context) {
	render(c, gin.H{
		"Title":    "Добавить новый контакт",
		"Classnew": "current",
		"Message":  "Поля Фамилия, Имя, Телефон не должны быть пустыми ",
	}, "newcontact.html")
}

// Функция post поиск контактов на главной странице в справочнике
func postPoiscContactSpravochnik(c *gin.Context) {
	var poisc Poiscsprav
	poisc.Poisc = c.PostForm("poisc")

	if poisc.Poisc == "" {
		products := getAllArticles()

		render(c, gin.H{
			"Title":    "Телефонный справочник",
			"Classind": "current",
			"Message":  "Поле не должно быть пустым",
			"payload":  products}, "index.html")
		return
	}

	products := poisc.postPoiscAllSpravochnik()

	render(c, gin.H{
		"Title":    "Телефонный справочник",
		"Classind": "current",
		"Message":  "Вот что нашлось",
		"payload":  products}, "index.html")

}

// Функция get открывает страницу для импорта справочника
func getNewUpload(c *gin.Context) {
	render(c, gin.H{
		"Title":    "Импорт справочника",
		"Classupl": "current",
		"Message":  "Для импорта, название столбцов в таком порядке №, Город, Орг-ция, Дол-ть, Фамилия, Имя, Отчество, Телефон, Email строки Фамилия, Имя, Телефон, должны быть обязательно заполнены",
	}, "upload.html")
}

// Функция post импортирует справочник
func postNewUpload(c *gin.Context) {
	var sprav Sprav
	// Открываем файл XLS
	file, err := c.FormFile("file")
	if err != nil {
		render(c, gin.H{
			"Classupl": "current",
			"Message":  "Ошибка : Не удалось загрузить файл",
		}, "upload.html")
		return
	}

	// Сохраняем файл на сервере
	err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
	if err != nil {
		render(c, gin.H{
			"Classupl": "current",
			"Message":  "Ошибка : Не удалось сохранить файл",
		}, "upload.html")
		return
	}

	// Открываем файл XLS
	xlsx, err := excelize.OpenFile("./uploads/" + file.Filename)
	if err != nil {
		fmt.Println("Failed to open XLSX file:", err)
		return
	}

	var rowss [][]string
	// Получаем список всех листов в файле
	sheetList := xlsx.GetSheetList()

	// Перебираем все листы и выводим их имена
	for _, sheetName := range sheetList {
		//fmt.Printf("Sheet %d: %s\n", i+1, sheetName)

		// Получаем первый лист из файла XLS
		//sheetName := xlsx.GetSheetName(1)
		rows, er := xlsx.GetRows(sheetName)
		if er != nil {
			fmt.Println("Незнаю что читать ", er)
			return
		}
		rowss = rows
	}
	// Перебираем строки листа
	for _, row := range rowss {

		// Проверяем, что строка не пуста
		if len(row) == 0 {
			//fmt.Println("Skipping empty row:", row)
			continue
		}

		if row[0] == "№" && row[1] == "Город" && row[2] == "Орг-ция" && row[3] == "Дол-ть" && row[4] == "Фамилия" && row[5] == "Имя" && row[6] == "Отчество" && row[7] == "Телефон" && row[8] == "Email" {
			//fmt.Println("Это заголовки таблицы пропускаем = ", row)
			continue
		}

		// Проверяем, что в строке достаточно элементов перед тем, как обращаться к обязательным полям
		if len(row) >= 9 && (row[4] == "" || row[5] == "" || row[7] == "") {
			//fmt.Println("Skipping row with missing required data:", row)
			continue
		}

		// Извлекаем данные из строки
		sprav.City = row[1]
		sprav.Organization = row[2] // Второй столбец
		sprav.Dolgnost = row[3]     // Третий столбец
		sprav.FirstName = row[4]    // Четвертый столбец
		sprav.LastName = row[5]     // Пятый столбец
		sprav.MiddleName = row[6]   // Шестой столбец
		sprav.Contacts = row[7]     // Седьмой столбец
		if len(row) >= 9 {
			sprav.Email = row[8] // Восьмой столбец
		} else {
			sprav.Email = " "
		}

		_, er := sprav.New()
		if er != nil {
			render(c, gin.H{
				"Classupl": "current",
				"Message":  "Ошибка : Произошла ошибка импорта в базу данных",
			}, "upload.html")
			return
		}
	}
	if sprav.FirstName == "" && sprav.LastName == "" {
		render(c, gin.H{
			"Classupl": "current",
			"Message":  "Сообщение : Нет данных для импорта в базу данных",
		}, "upload.html")

	} else {
		render(c, gin.H{
			"Classupl": "current",
			"Message":  "Сообщение : Успешно импортировали в базу данных",
		}, "upload.html")
	}
}

// Функция get экспорт справочника
func getExportSpravochnik(c *gin.Context) {

	render(c, gin.H{
		"Title":    "Экспорт справочника",
		"Classexp": "current",
		"Message":  "Экспортируем справочник из базы",
	}, "export.html")
}

// Функция post экспорт справочника
func postExportSpravochnik(c *gin.Context) {
	// Запрос данных из базы
	rows, err := db.Query(`SELECT TRIM("City"), TRIM("Organization"), TRIM("Dolgnost"), TRIM("FirstName"), TRIM("LastName"), TRIM("MiddleName"), TRIM("Contacts"), TRIM("Email")
FROM "Spravochnic" ORDER BY "Organization" ASC, "FirstName" ASC, "LastName" ASC, "MiddleName" ASC`)
	if err != nil { // Обработка ошибки при выполнении запроса
		log.Fatal(err)
	}
	defer rows.Close() // Закрытие соединения с базой данных после использования

	// Создание нового файла Excel
	file := excelize.NewFile()
	sheetName := "Sheet1"
	_, er := file.NewSheet(sheetName) // Создание нового листа в файле Excel
	if er != nil {                    // Обработка ошибки при создании нового листа
		log.Fatal(er)
	}

	// Названия столбцов для Excel
	colNames := []string{"Город", "Орг-ция", "Дол-ть", "Фамилия", "Имя", "Отчество", "Телефон", "Email"}

	// Создание стиля для границ
	borderStyle, err := file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center", // Выравнивание по центру
			Vertical:   "distributed",
		},
	})
	if err != nil { // Обработка ошибки при создании стиля
		log.Fatal(err)
	}

	// Создание стиля для границ и жирного шрифта
	borderStylecolNames, er := file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
		},
		Font: &excelize.Font{
			Size: 14,
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center", // Выравнивание по центру
			Vertical:   "center",
		},
	})
	if er != nil { // Обработка ошибки при создании стиля
		log.Fatal(err)
	}

	// Запись заголовков столбцов в файл Excel и применение стиля к ним
	for i, colName := range colNames {
		cell := columnIndexToAlpha(i+2) + "1" // Начиная с второго столбца, чтобы учесть столбец с порядковыми номерами
		file.SetCellValue(sheetName, cell, colName)
		file.SetCellStyle(sheetName, cell, cell, borderStylecolNames)
	}
	file.SetColWidth(sheetName, "B", "B", 12)
	file.SetColWidth(sheetName, "C", "C", 15)
	file.SetColWidth(sheetName, "D", "D", 27)
	file.SetColWidth(sheetName, "E", "E", 20)
	file.SetColWidth(sheetName, "F", "F", 12)
	file.SetColWidth(sheetName, "G", "G", 25)
	file.SetColWidth(sheetName, "H", "H", 15)
	file.SetColWidth(sheetName, "I", "I", 30)

	// Добавление столбца с порядковыми номерами и применение к нему стиля
	file.SetCellValue(sheetName, "A1", "№")
	file.SetCellStyle(sheetName, "A1", "A1", borderStylecolNames)
	file.SetColWidth(sheetName, "A", "A", 8)

	// Создание среза для сканирования значений из базы данных
	numColumns := len(colNames)
	data := make([]interface{}, numColumns)
	for i := range data {
		var value interface{}
		data[i] = &value
	}

	// Запись данных из базы данных в файл Excel
	rowIndex := 2
	num := 1 // Порядковый номер для каждой строки
	for rows.Next() {
		err := rows.Scan(data...) // Сканирование значений из строки результата запроса в срез
		if err != nil {           // Обработка ошибки при сканировании
			log.Fatal(err)
		}

		// Запись порядкового номера в первый столбец
		file.SetCellValue(sheetName, "A"+strconv.Itoa(rowIndex), num)
		file.SetCellStyle(sheetName, "A"+strconv.Itoa(rowIndex), "A"+strconv.Itoa(rowIndex), borderStyle)
		num++

		// Запись данных из среза в соответствующие ячейки Excel и применение стиля границ к ним
		for i, valuePtr := range data {
			value := *valuePtr.(*interface{})                             // Получение значения из указателя интерфейса
			cell := columnIndexToAlpha(i+2) + fmt.Sprintf("%d", rowIndex) // Начиная с третьего столбца
			file.SetCellValue(sheetName, cell, value)
			file.SetCellStyle(sheetName, cell, cell, borderStyle) // Применение стиля границ к ячейке
		}
		rowIndex++
	}

	// Сохранение файла Excel
	err = file.SaveAs("exported_data.xlsx")
	if err != nil { // Обработка ошибки при сохранении файла
		log.Fatal(err)
	}

	// Отправка файла пользователю для загрузки
	c.Header("Content-Disposition", "attachment; filename=exported_data.xlsx")
	c.Header("Content-Type", "application/octet-stream")
	c.File("exported_data.xlsx")
}

// Функция для установки персонализированной куки
func SetCookie(c *gin.Context, username string) {
	// Создаем новую куку с временем жизни один час
	expiration := time.Now().Add(24 * time.Hour)

	// URL-кодируем значение куки
	encodedValue := url.QueryEscape(username)

	cookie := http.Cookie{
		Name:    "token",
		Value:   encodedValue,
		Expires: expiration,
		Path:    "/",         // Указываем путь куки
		Domain:  "127.0.0.1", // Указываем домен куки
	}

	// Устанавливаем куку в ответ
	c.Writer.Header().Add("Set-Cookie", cookie.String())
}

// Функция для удаления куки
func DeleteCookie(c *gin.Context) {
	var user User
	cookie := http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",         // Указываем путь куки
		Domain:  "127.0.0.1", // Указываем домен куки
	}
	// Устанавливаем пустую куку с истекшим сроком действия
	c.Writer.Header().Add("Set-Cookie", cookie.String())

	render(c, gin.H{
		"Title":        "Выход",
		"ErrorTitle":   "Сообщение",
		"ErrorMessage": "Вы успешно вышли " + user.LastName,
	}, "DeleteCookie.html")
}

// Функция для проверки наличия куки и вывода приветствия
func AuthUser(c *gin.Context) {
	var user User
	// Получаем куку из запроса
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.Set("is_logged_in", false)
		return
	}

	c.Set("is_logged_in", true)

	if user.LastName == "" {
		user.Login = cookie.Value
		res, e := user.Cookie()
		if e != nil {
			return
		}

		if res == "" {
			return
		}
		user.LastName = res
	}
}

func render(c *gin.Context, data gin.H, templateName string) {
	AuthUser(c)
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

// Пример использования функции axis для преобразования числового индекса в строковый формат для столбца
func columnIndexToAlpha(index int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for index > 0 {
		index-- // Индексы начинаются с 1, а не с 0
		result = string(letters[index%26]) + result
		index /= 26
	}
	return result
}
