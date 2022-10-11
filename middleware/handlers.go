package middleware

import (
	"bytes"
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"io"
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"time"
	"truecaller/models" // models package where Response and Rupifi_FE scheme are defined

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// created connection with postgres db
func createConnection() *sql.DB {
	// loaded .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Opened the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	// checked the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

/* ------------------------- handler functions ------------------------- */

func callback(user models.Response) int64 {

	// created the postgres db connection
	db := createConnection()

	// closed the db connection
	defer db.Close()

	// created the insert sql query
	sqlStatement := `INSERT INTO truecaller (requestID, accessToken, endpoint) VALUES ($1, $2, $3) RETURNING userid`
	// returning userid will return the id of the selected user

	// the inserted id will be stored in this id
	var id int64

	// execute the sql statement
	err := db.QueryRow(sqlStatement, user.RequestID, user.AccessToken, user.Endpoint).Scan(&id)
	// Scan function will save the insert id in the id

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Passed the selected user to Callback URL successfully with id: %v\n", id)
	return id
}

func Callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// created an empty user of type models.Response
	var user models.Response

	// decoded the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// called callback() function and passed the user
	insertID := callback(user)

	// formatted a response object
	res := response{
		ID:      insertID,
		Message: "User passed to callback URL successfully!",
	}

	// sent the response
	json.NewEncoder(w).Encode(res)

	requestid := user.RequestID
	accesstoken := user.AccessToken
	url := user.Endpoint
	fmt.Println("")
	fmt.Println("requestID:", requestid)
	fmt.Println("accessToken:", accesstoken)
	requestBody, _ := MakeApiCall("GET", url, nil, accesstoken)
	str_requestBody := string(requestBody)
	fmt.Println(str_requestBody)
}

func details_FE() (models.Rupifi_FE, error) {
	// created the postgres db connection
	db := createConnection()

	// closed the db connection
	defer db.Close()

	// created the select sql query
	// select Accesstoken from Rupifi's Database
	sqlStatement := `SELECT * FROM rupifi`

	// executed the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// closed the statement
	defer rows.Close()

	var user_FE models.Rupifi_FE

	// ****************************  Verification Logic  ****************************

	user_FE.AccessToken = "eyJhbGciOiJSUzI1NiJ9.eyJ1c2VySWQiOjE5OTgxLCJ1c2VyUm9sZSI6IkJPUlJPV0VSIiwidXNlclBob25lIjoiOTgxNTg5NzgwMCIsImNsaWVudElkIjoiTlVSVFVSRV9GQVJNIiwiY2xpZW50SW50ZXJuYWxJZCI6IjRiMDA1ZWYyLTkzM2UtNGIyYS05YWUzLTE5ZThjYWFiMmE1NSIsImxlbmRlcklkIjoiV0VTVEVSTl9DQVBJVEFMIiwibGVuZGVySW50ZXJuYWxJZCI6Ijc2NjNjNjNkLTVkNjktNDdjNi04OGZiLTI1ZGYyZGM1M2JiNCIsInN1YiI6IlJ1cGlmaSBVc2VyIiwiaWF0IjoxNjYzNTg3NjgzLCJleHAiOjE2NjM2NzQwODMsImp0aSI6ImYyZTlmYWMyLTVlOTItNDVlNS1iOTVjLTlmNmZlZWZjNDc0YiJ9.m5HBEuAko0Cn-0ZrNVviGw8SlHm5rTX7yMUCUOVnNhQ8Mm7-BufTD9FQMlcwPRqikoUl2p4J2H-zbq0U5IZtlr151VCrUXXVMi_boWk5kmXV1v57bsf7DdTfQlI9woVBba4vo543br_lXfjLlIRlJPuZxpMwYlBfnQJ2IL5lUg5oAroHT8KNzIgP8jKSf9wdoMzQXmc2W1fxCMX8PcTnByP2ZeCES24TdfaL46wtRoZbcm3lokiLPVA2FsqUBfyUDgavsmgbtKY0m8cb6jYYdv7e4jMytqPNBdE1I9xhvflhPMe0ghIBMTNork4GSrmX0XyaFoyuKF7p-_ATqCk3uw"

	// append the user in the users slice
	// users = append(users, user)

	// return empty user on error
	return user_FE, err
}

func Details(w http.ResponseWriter, r *http.Request) {
	allowedHeaders := "Accept, Content-Type, x-custom-header, origin, authorization"

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Content-Type", "application/json")

	// get all the users in the db
	users, err := details_FE()

	if err != nil {
		log.Fatalf("Unable to get the accesstoken from the Database. %v", err)
	}

	// sent all the users as response
	json.NewEncoder(w).Encode(users)
}

func MakeApiCall(method string, path string, payload []byte, token string) ([]byte, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	bearer := "Bearer " + token

	if method == "GET" {
		req.Header = http.Header{
			"Authorization": {bearer},
		}
	}

	fmt.Println("endpoint:", path)
	fmt.Println(string(payload))

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return bodyBytes, nil
}
