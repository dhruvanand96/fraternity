package main

import (
	"app/server/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
)

type UserData struct {
	UserArray []User `json:"user_array"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"Name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Payload struct {
	AccessToken  string `json:"access-token"`
	RefreshToken string `json:"refresh-token"`
	UserObj      User   `json:"user"`
	Error        string `json:"error"`
}

type AllUsersPayload struct {
	UserArray []User `json:"user_array"`
	Success   string `json:"success"`
	Error     string `json:"error"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

var client *redis.Client

func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/dummy-data", sendData).Methods("GET")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/show-users", ShowUsers).Methods("GET")
	r.HandleFunc("/create-user", CreateUser).Methods("POST")
	r.HandleFunc("/delete", DeleteUser).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("GET")

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowedOrigins: []string{"http://localhost:4200"},
	})

	handler := c.Handler(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + os.Getenv("PORT"),
	}

	log.Fatal(srv.ListenAndServe())
}

func sendData(w http.ResponseWriter, r *http.Request) {
	var data = struct {
		Title string `json:"title"`
	}{
		Title: "Golang + Angular Starter Kit",
	}

	jsonBytes, err := utils.StructToJSON(data)
	if err != nil {
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	return
}

func CreateAuth(userid uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func CreateToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func login(w http.ResponseWriter, r *http.Request) {

	var u User
	_ = json.NewDecoder(r.Body).Decode(&u)
	p := Payload{}

	fmt.Println("login params email %s and password %s ", u.Email, u.Password)

	fileExists()
	err, userData := readDataFromFile()

	if err != nil {
		log.Println(err)
	}

	//compare the user from the request, with the one we defined:

	if len(userData.UserArray) == 0 {
		w.WriteHeader(401) // Wrong password or username, Return 401.
		p.AccessToken = ""
		p.RefreshToken = ""
		p.Error = "Wrong Username and password"
		json.NewEncoder(w).Encode(p)
		return
	}

	for index, value := range userData.UserArray {
		fmt.Println(value)
		if u.Email == value.Email && u.Password == value.Password {
			u = value
			break
		}

		if index == len(userData.UserArray)-1 {
			w.WriteHeader(401) // Wrong password or username, Return 401.
			p.AccessToken = ""
			p.RefreshToken = ""
			p.Error = "Wrong Username and password"
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	token, err := CreateToken(uint64(u.Id))
	if err != nil {
		w.WriteHeader(500) // Return 500 Internal Server Error.
		p.AccessToken = ""
		p.RefreshToken = ""
		p.Error = "Error in creating token"
		json.NewEncoder(w).Encode(p)
		return
	}

	saveErr := CreateAuth(uint64(u.Id), token)
	if saveErr != nil {
		w.WriteHeader(500) // Return 500 Internal Server Error.
		p.AccessToken = ""
		p.RefreshToken = ""
		p.Error = "Error in creating token"
		json.NewEncoder(w).Encode(p)
		return
	}

	w.WriteHeader(200)

	p.AccessToken = token.AccessToken
	p.RefreshToken = token.RefreshToken
	p.UserObj = u
	p.Error = ""
	json.NewEncoder(w).Encode(p)

}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *AccessDetails) (uint64, error) {
	userid, err := client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resData := map[string]string{
		"error":   "",
		"success": "",
	}

	au, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(401) // Wrong password or username, Return 401.
		resData = map[string]string{
			"error":   err.Error(),
			"success": "",
		}
		json.NewEncoder(w).Encode(resData)
		json.NewEncoder(w).Encode(resData)
	}
	deleted, delErr := DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any thing goes wrong
		resData = map[string]string{
			"error":   delErr.Error(),
			"success": "",
		}
		json.NewEncoder(w).Encode(resData)
		return
	}

	w.WriteHeader(200)

	resData = map[string]string{
		"error":   "",
		"success": "success",
	}
	json.NewEncoder(w).Encode(resData)
}

func ShowUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	fmt.Println("show users start")

	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(401)
		resData := AllUsersPayload{
			UserArray: nil,
			Success:   "",
			Error:     "Failure",
		}
		json.NewEncoder(w).Encode(resData)
		return
	}

	userId, err := FetchAuth(tokenAuth)

	fmt.Println("Hellow is %d", userId)

	if err != nil {
		w.WriteHeader(401)
		resData := AllUsersPayload{
			UserArray: nil,
			Success:   "",
			Error:     "Failure",
		}
		json.NewEncoder(w).Encode(resData)
		return
	}

	fileExists()
	err, userData := readDataFromFile()

	if err != nil {
		log.Println(err)
	}

	//you can proceed to save the Todo to a database
	//but we will just return it to the caller here:

	fmt.Println("data is %+v", userData.UserArray)
	w.WriteHeader(200)
	resData := AllUsersPayload{
		UserArray: userData.UserArray,
		Success:   "success",
		Error:     "",
	}
	json.NewEncoder(w).Encode(resData)
}

func readDataFromFile() (err error, userData *UserData) {
	jsonFile, err := os.Open("/tmp/user-data.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return err, nil
	}
	fmt.Println("Successfully Opened users.json")

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err, nil
	}

	err = json.Unmarshal(byteValue, &userData)

	if err != nil {
		return err, nil
	}

	return nil, userData
}

func fileExists() bool {
	info, err := os.Stat("/tmp/user-data.json")
	if os.IsNotExist(err) {
		myfile, e := os.Create("/tmp/user-data.json")
		if e != nil {
			log.Fatal(e)
		}
		myfile.Close()
		return false
	}
	return !info.IsDir()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	fmt.Println("create users start")

	fileExists()
	err, userData := readDataFromFile()

	if err != nil {
		log.Println(err)
	}

	var u User
	_ = json.NewDecoder(r.Body).Decode(&u)

	if userData == nil || len(userData.UserArray) == 0 {
		u.Id = 1
		userData = &UserData{}
	} else {

		u.Id = len(userData.UserArray) + 1
	}

	userData.UserArray = append(userData.UserArray, u)

	jsonData, err := json.Marshal(userData)

	if err != nil {
		log.Println(err)
	}

	err = os.Truncate("/tmp/user-data.json", 200)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("/tmp/user-data.json", jsonData, 0777)

	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("bytes written successfully")

	//if err != nil {
	//	resData := map[string]string{
	//		"error":   "Error for User",
	//		"success": "",
	//	}
	//	json.NewEncoder(w).Encode(resData)
	//	return
	//}

	//you can proceed to save the Todo to a database
	//but we will just return it to the caller here:
	resData := map[string]string{
		"error":   "",
		"success": "User successfully Created",
	}
	json.NewEncoder(w).Encode(resData)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	//tokenAuth, err := ExtractTokenMetadata(r)
	//if err != nil {
	//	return
	//}

	//xuserId, err := FetchAuth(tokenAuth)

	//if err != nil {
	//	resData := map[string]string{
	//		"error":   "Error for User",
	//		"success": "",
	//	}
	//	json.NewEncoder(w).Encode(resData)
	//	return
	//}

	//you can proceed to save the Todo to a database
	//but we will just return it to the caller here:
	resData := map[string]string{
		"error":   "",
		"success": "",
	}
	json.NewEncoder(w).Encode(resData)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers",
		"Content-Type,access-control-allow-origin, access-control-allow-headers")
	fmt.Println("show users start")

	tokenAuth, err := ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(401)
		resData := AllUsersPayload{
			UserArray: nil,
			Success:   "",
			Error:     "Failure",
		}
		json.NewEncoder(w).Encode(resData)
		return
	}

	userId, err := FetchAuth(tokenAuth)

	var u User
	_ = json.NewDecoder(r.Body).Decode(&u)

	fmt.Println("Hellow is %d", userId)

	fileExists()
	err, userData := readDataFromFile()

	if err != nil {
		log.Println(err)
	}

	for index, value := range userData.UserArray {
		fmt.Println(value)
		if u.Id == value.Id {
			copy(userData.UserArray[index:], userData.UserArray[index+1:])
			userData.UserArray[len(userData.UserArray)-1] = User{}
			userData.UserArray = userData.UserArray[:len(userData.UserArray)-1]
			break
		}

		if index == len(userData.UserArray)-1 {
			w.WriteHeader(500)
			resData := map[string]string{
				"error":   "Error is deleteing User",
				"success": "",
			}
			json.NewEncoder(w).Encode(resData)
			return
		}
	}

	jsonData, err := json.Marshal(userData)

	if err != nil {
		log.Println(err)
	}

	err = os.Truncate("/tmp/user-data.json", 200)
	if err != nil {
		log.Fatal(err)

	}

	err = ioutil.WriteFile("/tmp/user-data.json", jsonData, 0777)

	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("bytes written successfully")

	resData := map[string]string{
		"error":   "",
		"success": "User succesfully deleted",
	}
	json.NewEncoder(w).Encode(resData)
}
