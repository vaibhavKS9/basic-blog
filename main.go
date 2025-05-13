package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// BASIC PURPOSE OF THE PROJECT TO LEARN CRUD and creating RESTful API and Learn about common HTTP methods like GET, POST, PUT, PATCH, DELETE

type blogpost struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Email           string    `json:"email"`
	Content         string    `json:"content"`
	Category        string    `json:"category"`
	Date            time.Time `json:"date"` // RFC3339 by default
	AffiliatedLinks string    `json:"affiliatedLinks"`
	Tags            []string  `json:"tags"`
}

var mpPosts = make(map[int64]blogpost)
var nextID int64 = 1

// CRUD OPERATION

func createPost(c *gin.Context) {
	var newPost blogpost
	if err := c.BindJSON(&newPost); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newPost.ID = nextID
	nextID++
	mpPosts[newPost.ID] = newPost
	c.JSON(http.StatusCreated, newPost)
}

func getPost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	post, ok := mpPosts[id]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, post) // fixed: "newPost" → "post", "StatusOk" → "StatusOK"

}

func deletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	_, ok := mpPosts[id]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	delete(mpPosts, id)

	c.JSON(http.StatusOK,
		gin.H{"message": "Post with id = " + strconv.FormatInt(id, 10) + " successfully deleted!"}) // message formatting corrected

}

func updatePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var updatedPost blogpost
	if err := c.BindJSON(&updatedPost); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, ok := mpPosts[id]; !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	updatedPost.ID = id
	mpPosts[id] = updatedPost
	c.JSON(http.StatusOK,
		gin.H{"message": "Post with id = " + strconv.FormatInt(id, 10) + " successfully updated!", "post": updatedPost})

}

func getAllPosts(c *gin.Context) {
	var posts []blogpost
	for _, post := range mpPosts {
		posts = append(posts, post)
	}
	c.IndentedJSON(http.StatusOK, posts)
}

func main() {
	router := gin.Default()

	router.POST("/makeposts", createPost)
	router.GET("/getpost/:id", getPost)
	router.DELETE("/deletepost/:id", deletePost)
	router.PUT("/updatepost/:id", updatePost)
	router.GET("/getallposts", getAllPosts)

	router.Run(":8080")
}
