package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2"

	"log"
	"os"
	"io/ioutil"
	"crypto/sha1"
	"encoding/hex"
	"time"
)


type LinkData struct {
	Link     string `json:"link"`
	Username string `json:"username"`
	ID       string `json:"id"`
	UpVotes  int    `json:"upVotes`
	DownVotes int    `json:"downVotes`
	comment string  `json:comment`
}


func main() {
    app := fiber.New()
    
	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://socialbuying--thefoschini.myvtex.com, https://gofiber.net",
		AllowHeaders:  "Origin, Content-Type, Accept",
	}))

	app.Post("/generateLink", func(c *fiber.Ctx) error {
		data := new(LinkData)

		if err := c.BodyParser(data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON payload")
		}


		uniqueLink := generateUniqueLink(data.Link, data.Username)
		if err := storeLinkData(uniqueLink, data  ); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to store link data")
		}

		return c.JSON(fiber.Map{
			"link": uniqueLink,
		})
	})
    
    app.Listen(":3000")
}


func generateUniqueLink(link string, username string) string {

		data := link + username + time.Now().Format(time.RFC3339Nano)
		hash := sha1.Sum([]byte(data))
		hashString := hex.EncodeToString(hash[:])
		return hashString

}




func storeLinkData(link string, data *LinkData) error {

	file, err := os.OpenFile("link_data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Failed to open JSON file:", err)
		return err
	}
	defer file.Close()


	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Failed to read JSON file:", err)
		return err
	}

	var linkDataArray []LinkData
	if len(contents) > 0 {
		err = json.Unmarshal(contents, &linkDataArray)
		if err != nil {
			log.Println("Failed to decode link data:", err)
			return err
		}
	}


	data.ID = link


	linkDataArray = append(linkDataArray, *data)


	jsonData, err := json.Marshal(linkDataArray)
	if err != nil {
		log.Println("Failed to encode link data:", err)
		return err
	}


	err = file.Truncate(0)
	if err != nil {
		log.Println("Failed to truncate JSON file:", err)
		return err
	}


	_, err = file.WriteAt(jsonData, 0)
	if err != nil {
		log.Println("Failed to write link data to file:", err)
		return err
	}

	return nil
}




