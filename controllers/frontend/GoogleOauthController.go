package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"fmt"
	"io/ioutil"

	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
	"golang.org/x/oauth2"
)

var endpoint = oauth2.Endpoint{
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://oauth2.googleapis.com/token",
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     "368713028639-232vujcid04q68dk0umdhjk60s84v774.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8081/oauth/GoogleCallBack",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     endpoint,
}

const oauthStateString = "random"

type GoogleOauthController struct {
	BaseController
}

func (c *GoogleOauthController) HandleGoogleLogin() {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(url, 302)
}

func (c *GoogleOauthController) HandleGoogleCallback() {
	state := c.GetString("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		c.Redirect("/", 302)
		return
	}
	code := c.GetString("code")
	// fmt.Println("code: %s", code)
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		c.Redirect("/", 302)
	}
	// fmt.Println("token: %s", token)

	/*
	   在 Google OAuth 2.0 官方文件有提到，在拿到回應之後，
	   最好把 access_token 和 refresh_token 都保存在一個長期有效的安全地方，
	   例如有人會存在 shared session。
	   refresh_token 這組 token 是當 access_token 過期時，
	   可以拿 refresh_token 再去交換一組有效的 access_token。
	*/
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	data, err := simplejson.NewJson(contents)
	if err != nil {
		panic(err.Error())
	}
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	user := models.User{
		Email:    data.Get("email").MustString(),
		Password: common.Md5("google_oauth" + data.Get("id").MustString()),
		LastIP:   ip,
	}
	models.DB.Create(&user)

	models.Cookie.Set(c.Ctx, "userinfo", user)
	c.Redirect("/", 302)
}
