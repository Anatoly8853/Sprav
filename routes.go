// routes.go

package main

import "fmt"

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	//router.Use(setUserStatus())

	// Handle the index route
	//router.GET("/", showIndexPage)
	router.GET("/", AuthUser, handlerIndex)

	// Group user related routes together
	userRoutes := router.Group("/user")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		// Страница авторизации
		//userRoutes.GET("/login", ensureNotLoggedIn(), showLoginPage)
		userRoutes.GET("/login", AuthUser, handlerAuthorization)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		//userRoutes.POST("/login", ensureNotLoggedIn(), performLogin)
		userRoutes.POST("/login", AuthUser, handlerUserAuthorization)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		//userRoutes.GET("/logout", ensureLoggedIn(), logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		//userRoutes.GET("/register", ensureNotLoggedIn(), showRegistrationPage)
		userRoutes.GET("/register", handlerRegistration)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		//userRoutes.POST("/register", ensureNotLoggedIn(), register)
		userRoutes.POST("/register", handlerUserRegistration)

		userRoutes.GET("/deletecookie", AuthUser, DeleteCookie)

	}

	// Group article related routes together
	//articleRoutes := router.Group("/article")
	articleRoutes := router.Group("/spravochnik")
	{
		// Handle GET requests at /article/view/some_article_id
		//articleRoutes.GET("/view/:article_id", getArticle)
		articleRoutes.GET("/update/:id", AuthUser, EditPage)

		// Handle the GET requests at /article/create
		// Show the article creation page
		// Ensure that the user is logged in by using the middleware
		//articleRoutes.GET("/create", ensureLoggedIn(), showArticleCreationPage)
		articleRoutes.POST("/update/:id", AuthUser, handlerUpdate)

		// Handle POST requests at /article/create
		// Ensure that the user is logged in by using the middleware
		//articleRoutes.POST("/create", ensureLoggedIn(), createArticle)
		articleRoutes.POST("/delete/:id", AuthUser, handlerDelete)

		articleRoutes.POST("/new", AuthUser, postNewContactSpravochnik)

		articleRoutes.GET("/new", AuthUser, getNewContactSpravochnik)

		articleRoutes.POST("/poisc", AuthUser, postPoiscContactSpravochnik)

		// Отображение HTML страницы с формой для загрузки файла
		articleRoutes.GET("/upload", AuthUser, getNewUpload)

		// Обработка загрузки файла
		articleRoutes.POST("/upload", AuthUser, postNewUpload)

		articleRoutes.GET("/export", AuthUser, getExportSpravochnik)

		articleRoutes.POST("/export", AuthUser, postExportSpravochnik)

	}
	fmt.Println("Сервер запустился теперь переходим в браузер и открываем адрес http://127.0.0.1:8080")
}
