package repository

import (
	"database/sql"
	"libretaxi/objects"
	"log"
)

type Repository struct {
	db *sql.DB
}

func (repo *Repository) FindUser(userId int64) *objects.User {
	user := &objects.User{}

	rows, err := repo.db.Query(`select "userId", "menuId", "username", "firstName", "lastName", "lon", "lat", "languageCode", "reportCnt" from users where "userId" = $1 limit 1`, userId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cnt := 0
	for rows.Next() {
		cnt++
		rows.Scan(&user.UserId, &user.MenuId, &user.Username, &user.FirstName, &user.LastName, &user.Lon, &user.Lat, &user.LanguageCode, &user.ReportCnt)
	}

	if cnt == 0 {
		return nil
	}

	return user
}

func (repo *Repository) SaveUser(user *objects.User) {
	// Upsert syntax: https://stackoverflow.com/questions/1109061/insert-on-duplicate-update-in-postgresql
	// Geo populate syntax: https://gis.stackexchange.com/questions/145007/creating-geometry-from-lat-lon-in-table-using-postgis/145009
	result, err := repo.db.Query(`INSERT INTO users ("userId", "menuId", "username", "firstName", "lastName", "lon", "lat", "geog", "languageCode", "reportCnt")
		VALUES ($1, $2, $3, $4, $5, $6, $7, ST_SetSRID(ST_MakePoint($7, $6), 4326), $8, $9)
		ON CONFLICT ("userId") DO UPDATE
		  SET "menuId" = $2,
		      "username" = $3,
		      "firstName" = $4,
		      "lastName" = $5,
		      "lon" = $6,
		      "lat" = $7,
		      "languageCode" = $8,
		      "reportCnt" = $9,
		      "geog" = ST_SetSRID(ST_MakePoint($6, $7), 4326)
		  `, user.UserId, user.MenuId, user.Username, user.FirstName, user.LastName, user.Lon, user.Lat, user.LanguageCode, user.ReportCnt)
	defer result.Close()

	if err != nil {
		log.Println(err)
	} else {
		log.Println("User saved")
	}
}

func (repo *Repository) SaveNewPost(post *objects.Post) {
	result, err := repo.db.Query(`INSERT INTO posts ("userId", "text", "lon", "lat", "geog", "reportCnt") VALUES ($1, $2, $3, $4, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)`,
		post.UserId, post.Text, post.Lat, post.Lon, post.ReportCnt)
	defer result.Close()

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Post saved")
	}
}

func (repo *Repository) UserIdsAround(lon float64, lat float64) (userIds []int64) {
	// select "userId", ST_Distance(c.x, "geog") AS distance from users, (SELECT ST_MakePoint(-122.415561, 37.633141)::geography) as c(x) where ST_DWithin(c.x, "geog", 25000)
	result, err := repo.db.Query(`select "userId" from users, (SELECT ST_MakePoint($1, $2)::geography) as c(x) where ST_DWithin(c.x, "geog", 25000)`,
		lon, lat)
	defer result.Close()

	if err != nil {
		log.Println(err)
		return nil
	}

	for result.Next() {
		var userId int64
		err := result.Scan(&userId)
		if err != nil {
			log.Println("Error getting userId")
		} else {
			userIds = append(userIds, userId)
		}
	}
	log.Printf("Found %d users around\n", len(userIds))
	return userIds
}

func NewRepository(db *sql.DB) *Repository {
	repo := &Repository{db: db}
	return repo
}