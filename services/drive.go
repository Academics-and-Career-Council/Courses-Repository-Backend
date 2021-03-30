package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/Mshivam2409/AnC-Courses/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// DriveClient : Helper struct for Google Drive Service
type DriveClient struct {
	*drive.Service
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "/secrets/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Panic("Unable to to parse Json file")
	}
	oldToken, err := json.Marshal(fiber.Map{"client_id": config.ClientID, "client_secret": config.ClientSecret, "refresh_token": tok.RefreshToken, "grant_type": "refresh_token"})
	if err != nil {
		log.Panic("Unable to to parse Json file")
	}
	r, err := http.Post("https://accounts.google.com/o/oauth2/token", "application/json", bytes.NewBuffer(oldToken))
	if err != nil {
		log.Panic("Unable to to parse Json file")
	}
	newToken := &oauth2.Token{}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(newToken)
	return config.Client(context.Background(), newToken)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// GetService : Connect to GoogleAPI and return the client
func GetService() *DriveClient {
	b, err := ioutil.ReadFile("/secrets/credentials.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	srv, err := drive.New(client)
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v", err)
	}
	return &DriveClient{srv}
}

// CreateFile : Creates a file in google drive and return the file id and filename
func (srv *DriveClient) CreateFile(folderID string, file *multipart.FileHeader) (models.File, error) {
	f, err := file.Open()
	if err != nil {
		log.Printf("Unable to upload file! %v", err)
		return models.File{}, err
	}
	driveFile, err := srv.Files.Create(&drive.File{Name: file.Filename, Parents: []string{folderID}}).Media(f).Do()
	if err != nil {
		log.Printf("Unable to upload file! %v", err)
		return models.File{}, err
	}
	return models.File{ID: driveFile.DriveId, Name: file.Filename}, err
}

// CreateFolder : Creates a folder in google drive
func (srv *DriveClient) CreateFolder(folderName string) (string, error) {
	folder, err := srv.Files.Create(&drive.File{MimeType: "application/vnd.google-apps.folder", Name: folderName}).Do()
	log.Printf("Folder with id %v created", folder.Id)
	if err != nil {
		log.Printf("Unable to create folder: %v", err)
	}
	return folder.Id, err
}

// DeleteFile : Deletes a file with the given ID
func (srv *DriveClient) DeleteFile(fileID string) error {
	err := srv.Files.Delete(fileID).Do()
	if err != nil {
		log.Printf("Unable to delete file! : %v"+"File ID:"+fileID, err)
	}
	return err
}

// GetFile : Downloads a file with the given ID.
func (srv *DriveClient) GetFile(fileID string) (io.Reader, error) {
	f, err := srv.Files.Get(fileID).Download()
	a := io.Reader(f.Body)
	if err != nil {
		log.Printf("Unable to delete file! : %v"+"File ID:"+fileID, err)
	}
	return a, err
}
