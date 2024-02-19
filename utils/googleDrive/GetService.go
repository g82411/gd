package googleDrive

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"log"
	"os"
	"sync"
)

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	codeChan := make(chan string)
	go StartCallbackServer(context.Background(), 8080, codeChan)
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Open the following link in your browser then type the authorization code: \n%v\n", authURL)
	var authCode string

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case code := <-codeChan:
				StopCallbackServer(context.Background())
				authCode = code
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// tokenFromFile 嘗試從一個檔案中讀取 token。
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

// saveToken 將 token 保存到一個檔案中。
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetService(ctx context.Context) (*drive.Service, error) {
	// TODO: make clinet_secret.json as a parameter
	b, err := os.ReadFile("./client_secret.json") // Make sure to replace 'credentials.json' with your actual credentials file
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	// TODO: remove magic redirect URL
	config.RedirectURL = "http://localhost:8080/callback"
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	// 查看是否已經有保存的 token
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	client := config.Client(ctx, tok) // This function is similar to the one used in the Quickstart guide by Google
	return drive.NewService(ctx, option.WithHTTPClient(client))
}
