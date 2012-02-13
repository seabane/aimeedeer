package main

import (
    "github.com/hoisie/web.go"
)

func loginProcess(val string) string { 
	pramas := strings.Split(val, "&");
	
	username := "";
	password := "";
	
	for i := 0; i < len(pramas); i++ {
		prama := pramas[i];
		keyValue := strings.Split(prama,"=");
		//取用户名
		if keyValue[0] == "um"{
			username = keyValue[1];
		}
		
		//取密码
		if keyValue[0] == "pw"{
			password = keyValue[1];
		}
	}
	
	return "username:" + username + " and password:" + password;

} 

func main() {
    web.Get("/user/login/(.*)", loginProcess);
    web.Run("0.0.0.0:8080");
}