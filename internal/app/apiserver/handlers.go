package apiserver

import (
	"dialogue/internal/models"
	"dialogue/internal/validator"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type storyForm struct {
	Title   string `schema:"title"`
	Content string `schema:"content"`
	Options string `schema:"options"`
	Privacy bool   `schema:"privacy"`
	validator.Validator
}

type userForm struct {
	Nickname string `schema:"nickname"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
	validator.Validator
}

type userLoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
	validator.Validator
}

type accountPasswordUpdateForm struct {
	CurrentPassword         string `schema:"currentPassword"`
	NewPassword             string `schema:"newPassword"`
	NewPasswordConfirmation string `schema:"newPasswordConfirmation"`
	validator.Validator
}

// redirectHome redirects default query to the home page.
func (s *server) redirectHomePage(c *gin.Context) {
	if c.Request.URL.Path != "/" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusFound, "/home")
}

// homePage renders home page with public and private stories related to a user.
func (s *server) homePage(c *gin.Context) {

	//Get user ID from context.
	userID := s.getID(c)

	//Pass related data and stories and render the page.
	data := s.newTemplateData(c)
	data.DataDialogues.DialoguesToDisplay, _ = s.services.Story().DisplayStories(userID)
	s.render(c, http.StatusOK, "home.html", data)
}

// emptyFBView renders an empty form where user can create a starting point of a new story.
func (s *server) startingBlockForm(c *gin.Context) {
	data := s.newTemplateData(c)
	s.render(c, http.StatusOK, "createFB.html", data)
}

// createFB method parse form, get the values from it and create first block of the story.
// The first block has it's unique ID, and it is equal to the ID of the story itself.
func (s *server) startingBlockCreate(c *gin.Context) {

	//Get values from the form and store them into form variable.
	var storyForm storyForm
	s.parse(c, &storyForm)

	//Basic validations checks.
	storyForm.validateStoryForm()
	if !storyForm.Valid() {
		data := s.newTemplateData(c)
		data.StoryForm = storyForm
		s.render(c, http.StatusUnprocessableEntity, "createFB.html", data)
		return
	}

	//Parse options from the form and store them into the slice of strings.
	optionsSlice := strings.Split(storyForm.Options, "\r\n")

	//Get user ID from context and put gathered data into DB, then get the ID of fresh created first block of the story.
	userID := s.getID(c)
	newStoryID := s.services.Story().Create(userID, storyForm.Title, storyForm.Content, optionsSlice, storyForm.Privacy)

	s.setFlash(c, "First step is done, and the story have been created!")
	path := strconv.Itoa(int(newStoryID))
	c.Redirect(http.StatusFound, path)
}

// createdFBView renders view of fresh created story with nessessary data.
func (s *server) startingBlockRender(c *gin.Context) {

	//Get the ID of fresh story.
	storyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	//Get the data related to the story with ID and pass it to the view.
	data := s.newTemplateData(c)
	data.DataDialogues = s.services.Story().StartingBlockData(storyID)
	s.render(c, http.StatusOK, "renderFB.html", data)
}

// editFBView allows to edit first block of the story.
func (s *server) startingBlockEditionForm(c *gin.Context) {

	//get the ID of the story and data of the first block.
	storyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	// Render the form for editing with existing data.
	data := s.newTemplateData(c)
	data.DataDialogues = s.services.Story().StartingBlockData(storyID)
	s.render(c, http.StatusOK, "editFB.html", data)
}

// editFB passes edited data to the data base.
func (s *server) startingBlockEdit(c *gin.Context) {

	//Parse edited data for the first block of a story.
	var storyForm storyForm
	s.parse(c, &storyForm)
	optionsSlice := strings.Split(storyForm.Options, "\r\n")

	//Get ID of the story and update it's data with a new one.
	storyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userID := s.getID(c)
	s.services.Story().Edit("starting_blocks", storyID, userID, storyForm.Title, storyForm.Content, optionsSlice)
	path := "/stories/startingblocks/" + strconv.Itoa(storyID)
	c.Redirect(http.StatusFound, path)
}

// deleteFB deletes the whole story and all blocks related to it.
func (s *server) deleteWholeStory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	s.services.Story().DeleteWholeStory(id)
	c.Redirect(http.StatusFound, "/home")
}

// createdBView renders existing block of a story.
func (s *server) blockRender(c *gin.Context) {

	//Get ID of a block.
	blockID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	//Retrieve data from database and render the block.
	data := s.newTemplateData(c)
	data.DataDialogues = s.services.Story().BlockData(blockID)
	s.render(c, http.StatusOK, "renderB.html", data)
}

// editBView allows to edit block of the story.
func (s *server) blockEditionForm(c *gin.Context) {

	//Get ID of a block.
	blockID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	// Render the form for editing with existing data.
	data := s.newTemplateData(c)
	data.DataDialogues = s.services.Story().BlockData(blockID)
	s.render(c, http.StatusOK, "editB.html", data)
}

// editFB passes edited data to the data base.
func (s *server) blockEdit(c *gin.Context) {

	//Parse form and store it.
	var blockForm storyForm
	s.parse(c, &blockForm)

	//Parse options from the form and store them into the slice of strings.
	optionsSlice := strings.Split(blockForm.Options, "\r\n")

	//Get ID of the editing block and update it's data with a new one.
	blockID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	userID := s.getID(c)
	s.services.Story().Edit("common_blocks", blockID, userID, blockForm.Title, blockForm.Content, optionsSlice)
	path := "/stories/blocks/" + strconv.Itoa(blockID)
	c.Redirect(http.StatusFound, path)
}

// deleteB deletes a block and other blocks if they are not related to other blocks.
func (s *server) deleteBlock(c *gin.Context) {
	blockID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	s.services.Story().DeleteOneBlock(blockID)
	c.Redirect(http.StatusFound, "/home")
}

// userSignupView renders the page for signing user up.
func (s *server) userSignupForm(c *gin.Context) {
	data := s.newTemplateData(c)
	s.render(c, http.StatusOK, "signup.html", data)
}

// userSignup signs up a new user with provided data.
func (s *server) userSignup(c *gin.Context) {

	//Parse provided form.
	var userForm userForm
	s.parse(c, &userForm)

	if userForm.FieldErrors == nil {
		userForm.FieldErrors = make(map[string]string)
	}

	//Basic validation checks.
	userForm.validateUserForm()
	if !userForm.Valid() {
		data := s.newTemplateData(c)
		data.UserForm = userForm
		s.render(c, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	u := &models.User{
		Nickname: userForm.Nickname,
		Email:    userForm.Email,
		Password: userForm.Password,
	}

	//Save new user into the data base.
	err := s.store.User().Create(u)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			userForm.AddFieldError("email", "Email address is already in use")
			data := s.newTemplateData(c)
			data.UserForm = userForm
			s.render(c, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			s.serverError(c, err)
		}
		return
	}

	s.setFlash(c, "You successfully signed up. Please, log in for more content.")
	c.Redirect(http.StatusFound, "/user/login")
}

// userLoginView renders the page for logging user in.
func (s *server) userLoginForm(c *gin.Context) {
	data := s.newTemplateData(c)
	s.render(c, http.StatusOK, "login.html", data)
}

// userLogin logs user in.
func (s *server) userLogin(c *gin.Context) {

	//Parse provided form.
	var userForm userLoginForm
	s.parse(c, &userForm)

	//Basic validations check.
	userForm.validateUserForm()
	if !userForm.Valid() {
		data := s.newTemplateData(c)
		data.UserLoginForm = userForm
		s.render(c, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	//Authenticate user and log him in if no errors.
	user, err := s.services.User().Authenticate(userForm.Email, userForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			userForm.AddNonFieldError("Email or password is incorrect")
			data := s.newTemplateData(c)
			data.UserLoginForm = userForm
			s.render(c, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			s.serverError(c, err)
		}
		return
	}
	s.setFlash(c, "You logged in with a geat success.")

	//Generate new session ID for the user and save related data via cookies.
	sessionID := generateSessionID()
	c.SetCookie("session_id", sessionID, 3600, "/", "", false, true)
	sessionData := map[string]any{
		"userID":       user.ID,
		"createdAt":    time.Now().Unix(),
		"lastActiveAt": time.Now().Unix(),
	}
	s.redisClient.Set(c, sessionID, serialize(sessionData), 12*time.Hour)
	s.redisClient.Set(c, sessionID+":userID", user.ID, 12*time.Hour)

	c.Redirect(http.StatusFound, "/home")
}

// userLogoutPost logouts the user and destroy current session.
func (s *server) userLogout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err == nil {
		s.redisClient.Del(c, sessionID)
	}
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("flash_message", "You logged out with a great success", 5, "/", "", false, true)

	c.Redirect(http.StatusFound, "/home")
}

// accountView renders a page with data related to the user (nickname and other).
func (s *server) userAccountRender(c *gin.Context) {

	//Get user ID and then other data related to the user.
	userID := s.getID(c)
	user, err := s.store.User().Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			c.Redirect(http.StatusFound, "/user/login")
		} else {
			s.serverError(c, err)
		}
		return
	}

	//Renders the page with all related data.
	data := s.newTemplateData(c)
	data.UserData = user
	s.render(c, http.StatusOK, "account.html", data)
}

// passwordUpdateView renders a page where user can change the password for the account.
func (s *server) passwordEditionForm(c *gin.Context) {
	data := s.newTemplateData(c)
	data.PasswordForm = accountPasswordUpdateForm{}
	s.render(c, http.StatusOK, "password.html", data)
}

// passwordUpdate updates user's password with new provided information.
func (s *server) passwordUpdate(c *gin.Context) {

	//Parse the form provided by the user.

	var passwordForm accountPasswordUpdateForm
	s.parse(c, &passwordForm)

	//Basic validations check.
	passwordForm.validatePasswordUpdateForm()
	if !passwordForm.Valid() {
		data := s.newTemplateData(c)
		data.PasswordForm = passwordForm
		s.render(c, http.StatusUnprocessableEntity, "password.html", data)
		return
	}

	//Get ID of a user.
	userID := s.getID(c)

	//Update password with new information.
	err := s.store.User().PasswordUpdate(userID, passwordForm.CurrentPassword, passwordForm.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			passwordForm.AddFieldError("currentPassword", "Current password is incorrect")
			data := s.newTemplateData(c)
			data.PasswordForm = passwordForm
			s.render(c, http.StatusUnprocessableEntity, "password.tmpl", data)
		} else {
			s.serverError(c, err)
		}
		return
	}
	s.setFlash(c, "Password successfully updated.")
	c.Redirect(http.StatusFound, "/account/view")
}

// about contains basic idea of the site.
func (s *server) about(c *gin.Context) {
	data := s.newTemplateData(c)
	s.render(c, http.StatusOK, "about.html", data)
}
