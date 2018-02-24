package api

import (
	"encoding/json"
	"fmt"
	"github.com/SpaTmole/gustly/app"
	"github.com/SpaTmole/gustly/app/models"
	"github.com/SpaTmole/gustly/app/utils"
	"github.com/revel/revel"
	"io/ioutil"
	"time"
	// "os"
	// "net/http"
	// "strings"
)

type StoriesController struct {
	APIv1
}

type GetStoryForm struct {
	ID          uint
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	UserID      uint  `json:owner`
	NextStoryID *uint `json:next_story_id`
	Content     string
	Type        string
	UsersViewed []uint `json:users_viewed`
}

func (c StoriesController) AddStory(upload []byte) revel.Result {
	notExpiredDate := time.Now().Add(time.Hour * -24)
	user, _ := c.Args["user"].(models.User)
	fmt.Println("File upload:  ", len(upload), " bytes.")
	contentType := utils.GuessContentType(upload)
	if contentType == "" {
		return c.RenderJSON(map[string]interface{}{"errors": map[string]interface{}{"upload": "Unsupported file format."}})
	}
	name := fmt.Sprintf("media/%s_%s", utils.MakeUniqueKey(64), c.Params.Files["upload"][0].Filename)
	err := ioutil.WriteFile(name, upload, 0644)
	if err != nil {
		revel.ERROR.Println(err)
		return c.RenderJSON(map[string]interface{}{"errors": []string{"Failed to upload story."}})
	}
	lastStory := &models.Story{}
	err = app.DB.Table("stories").Where(
		"user_id = ? AND created_at > ?",
		user.ID,
		notExpiredDate,
	).Order("created_at desc").First(lastStory).Error
	if err != nil {
		revel.WARN.Println("No stories found for user: ", user)
		lastStory = nil
	}
	story := &models.Story{
		Type:    contentType,
		Content: name,
	}
	err = app.DB.Model(&user).Association("Stories").Append(story).Error
	if err != nil {
		panic(err)
	}
	if lastStory != nil {
		lastStory.NextStoryID = &story.ID
		err = app.DB.Save(lastStory).Error
		if err != nil {
			panic(err)
		}
	}
	return c.RenderJSON(map[string]interface{}{"story": story})
}

func (c StoriesController) GetStory(user_id uint) revel.Result {
	notExpiredDate := time.Now().Add(time.Hour * -24)
	user, _ := c.Args["user"].(models.User)
	storyId := c.Params.Query.Get("story_id")
	query := app.DB.Table("stories").Order("created_at").Where("user_id = ? AND created_at > ?", user_id, notExpiredDate)
	if storyId != "" {
		query = query.Where("id = ?", storyId)
	}
	story := &models.Story{}
	err := query.First(story).Error
	if err != nil {
		story = nil
	}
	if story != nil {
		err = app.DB.Model(&user).Association("Views").Append(story).Error
	}
	var user_ids []uint
	form := GetStoryForm{}
	err = app.DB.Table("views").Where("story_id = ?", story.ID).Pluck("user_id", &user_ids).Error
	if err != nil {
		user_ids = []uint{}
		revel.ERROR.Println(err)
	}
	marshaled, _ := json.Marshal(story)
	json.Unmarshal(marshaled, &form)
	form.UsersViewed = user_ids

	return c.RenderJSON(map[string]interface{}{"result": form})
}
