package users

import (
	"github.com/aftaab60/bookstore_oauth-go/oauth"
	"github.com/aftaab60/bookstore_users-api/domain/users"
	"github.com/aftaab60/bookstore_users-api/services"
	"github.com/aftaab60/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/*
//this controller can be replaced with below type. Instead of ioutil, unmarshalling, we can use directly shouldBind
func Create(c *gin.Context) {
	var user users.User
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)
	result, saveErr := services.Create(user)
	if saveErr != nil {
		c.String(http.StatusInternalServerError, "Error in creating user")
		return
	}
	c.JSON(http.StatusCreated, result)
}
 */

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.UserService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status, err)
		return
	}
	userId, userErr := strconv.ParseInt(c.Param("userId"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		c.JSON(err.Status, err)
		return
	}
	user, getErr := services.UserService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("userId"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		c.JSON(err.Status, err)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Id = userId

	isPatch := false
	if c.Request.Method == http.MethodPatch {
		isPatch = true
	}
	result, updateErr := services.UserService.UpdateUser(isPatch, user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
	}
	c.JSON(http.StatusOK, result)
}

func Delete(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("userId"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		c.JSON(err.Status, err)
		return
	}

	if err := services.UserService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	userStatus, ok := c.GetQuery("status")
	if !ok || userStatus=="" {
		err := errors.NewBadRequestError("invalid user status")
		c.JSON(err.Status, err)
		return
	}

	users, err := services.UserService.SearchUser(userStatus)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, users.Marshall(c.GetHeader("X-Public")=="true"))
}

func Login(c *gin.Context) {
	var loginRequest users.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	user, err := services.UserService.LoginUser(loginRequest)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public")=="true"))
}
