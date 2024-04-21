package main

// User ...
type User struct {
	ID        int    `json:"ID,-"`
	Login     string `json:"username"`
	Password  string `json:"password,-"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Role      string `json:"Role,-"`
}

// Create создание нового пользователя в базе
func (u User) Create() error {
	row := db.QueryRow(`INSERT INTO "User" ("Login", "Password", "FirstName", "LastName", "Role") VALUES ($1, $2, $3, $4, 'manager') RETURNING "ID"`,
		u.Login, u.Password, u.FirstName, u.LastName)
	e := row.Scan(&u.ID)
	if e != nil {
		return e
	}

	return nil
}

func (u User) Select() (string, error) {
	row := db.QueryRow(`SELECT "Role", "FirstName", "LastName" FROM "User" WHERE "Login"=$1 AND "Password"=$2`,
		u.Login, u.Password)
	e := row.Scan(&u.Role, &u.FirstName, &u.LastName)
	if e != nil {
		return u.FirstName, e
	}

	return u.FirstName, nil
}

func (u User) Cookie() (string, error) {
	row := db.QueryRow(
		`SELECT "Role", "FirstName", "LastName" FROM "User" WHERE "Login"=$1 `, u.Login, u.Password)
	e := row.Scan(&u.Role, &u.FirstName, &u.LastName)
	if e != nil {
		return u.FirstName, e
	}

	return u.FirstName, nil
}

func (s Sprav) Update() error {

	_, er := db.Exec(`UPDATE "Spravochnic" SET "City" = $1, "Organization" = $2, Dolgnost = $3, FirstName = $4,
                         LastName = $5, MiddleName = $6, Contacts = $7, Email = $8 WHERE "ID"=$9 `, s.City, s.Organization,
		s.Dolgnost, s.FirstName, s.LastName, s.MiddleName, s.Contacts, s.Email, s.ID)
	if er != nil {
		return er
	}

	return nil
}

func (s Sprav) Delete() error {

	_, er := db.Exec(`DELETE FROM "Spravochnic" WHERE  "ID" = $1 `, s.ID)
	if er != nil {

		return er
	}

	return nil
}

func (s Sprav) New() (int, error) {

	row := db.QueryRow(`INSERT INTO "Spravochnic" ("City", "Organization", "Dolgnost" , "FirstName",
                         "LastName", "MiddleName", "Contacts", "Email") VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING ID`, s.City, s.Organization,
		s.Dolgnost, s.FirstName, s.LastName, s.MiddleName, s.Contacts, s.Email)
	e := row.Scan(&s.ID)
	if e != nil {
		return s.ID, e
	}

	return s.ID, nil
}
