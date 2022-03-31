package requester

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PickHD/pickablog/config"
	"github.com/PickHD/pickablog/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type (
	//IOAuthGoogle interface for requester oauth provider google
	IOAuthGoogle interface {
		GetUserInfo(code string) (model.GoogleOauthResponse,error)
	}

	//OAuthGoogle struct for requester oauth provider google
	OAuthGoogle struct {
		Context context.Context
		Config *config.Configuration
		Logger *logrus.Logger
		GConfig *oauth2.Config
		HTTPClient *http.Client
	}
)

// GetUserInfo retrieved user data from google services
func (og *OAuthGoogle) GetUserInfo(code string) (model.GoogleOauthResponse,error) {
	token,err := og.GConfig.Exchange(og.Context,code)
	if err != nil {
		og.Logger.Error(fmt.Errorf("OAuthGoogle.GetUserInfo Exchange ERROR : %v MSG : %s",err,err.Error()))
		return model.GoogleOauthResponse{},err
	}

	resp,err := og.HTTPClient.Get(og.Config.Const.OauthGoogleAPIURL + token.AccessToken)
	if err != nil {
		og.Logger.Error(fmt.Errorf("OAuthGoogle.GetUserInfo HttpClient.Get ERROR : %v MSG : %s",err,err.Error()))
		return model.GoogleOauthResponse{},err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		og.Logger.Error(fmt.Errorf("OAuthGoogle.GetUserInfo ioutil.ReadAll ERROR : %v MSG : %s",err,err.Error()))
		return model.GoogleOauthResponse{}, err
	}

	var getUser model.GoogleOauthResponse

	err = json.Unmarshal(data,&getUser)
	if err != nil {
		og.Logger.Error(fmt.Errorf("OAuthGoogle.GetUserInfo json.Unmarshal ERROR : %v MSG : %s",err,err.Error()))
		return model.GoogleOauthResponse{}, err
	}

	return getUser, nil
}