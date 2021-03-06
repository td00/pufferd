/*
 Copyright 2016 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package httphandlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/config"
	"github.com/pufferpanel/pufferd/logging"
	"fmt"
)

func OAuth2Handler(gin *gin.Context) {
	authHeader := gin.Request.Header.Get("Authorization")
	var authToken string;
	if authHeader == "" {
		authToken = gin.Query("accessToken")
		if authToken == "" {
			gin.AbortWithStatus(401)
			return
		}
	} else {
		authArr := strings.SplitN(authHeader, " ", 2)
		if len(authArr) < 2 || authArr[0] != "Bearer" {
			gin.AbortWithStatus(400)
			return
		}
		authToken = authArr[1];
	}
	validateToken(authToken, gin)
}

func validateToken(accessToken string, gin *gin.Context) {
	authUrl := config.Get("infoserver")
	token := config.Get("authtoken")
	client := &http.Client{}
	data := url.Values{}
	data.Set("token", accessToken)
	request, _ := http.NewRequest("POST", authUrl, bytes.NewBufferString(data.Encode()))
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(request)
	if err != nil {
		logging.Error("Error talking to auth server", err)
		errMsg := make(map[string]string)
		errMsg["error"] = err.Error();
		gin.JSON(500, errMsg);
		return
	}
	if response.StatusCode != 200 {
		logging.Error("Unexpected response code from auth server", response.StatusCode)
		errMsg := make(map[string]string)
		errMsg["error"] = fmt.Sprintf("Receieved response %i", response.StatusCode);
		gin.JSON(500, errMsg);
		return
	}
	var respArr map[string]interface{}
	json.NewDecoder(response.Body).Decode(&respArr)
	if respArr["error"] != nil {
		logging.Error("Error parsing response from auth server", err)
		errMsg := make(map[string]string)
		errMsg["error"] = "Failed to parse auth server response";
		gin.JSON(500, errMsg);
		return
	}
	if respArr["active"].(bool) == false {
		gin.AbortWithStatus(401)
		return
	}
	gin.Set("server_id", respArr["server_id"].(string))
	gin.Set("scopes", strings.Split(respArr["scope"].(string), " "))
}
