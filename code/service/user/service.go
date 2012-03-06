package main

import (
   "github.com/Philio/GoMySQL"
   "log"
   "json"
   "github.com/hoisie/web.go"
)

type jsonRes struct {
   isSuc boolean
   errMsg string
   data 
}

func Login(ctx *web.Context,v string) string { 
   log.Printf("start login");
   
   prams := ctx.Request.Params;
   
   db := GetDbClient();
   if(db == nil){
      return "db conn err";
   }

   log.Printf("get conn ok");

   err := db.Query("select * from user where username = '" + prams["username"] + "' and password = password('" + prams["password"] + "')");
   if err != nil{
      log.Printf("%v\n",err);
      return "db queery error";
   }

   result,err := db.UseResult();
   if err != nil{
      log.Printf("%v\n",err);
      return "db use result err";
   }
   
   u := result.FetchMap();

   db.Close();

   log.Printf("query user ok");

   if(u == nil){
      b,_ := json.Marshal(jsonRes{isSuc:false},errMsg:"have no user found."});
	  return string(b);
   }

   b,_ := json.Marshal(jsonRes{isSuc:true},data:u});
   return string(b);
}

func Register(ctx *web.Context,v string) string{
   prams := ctx.Request.Params;

   db := GetDbClient();
   if(db == nil){
      return "db conn err";
   }

   err := db.Query("insert into user (username,password,nickname) values('" + prams["username"] + "',password('" + prams["password"] + "'),'" + prams["nickname"] + "')");
   if err != nil{
      log.Printf("%v\n",err);
      return "db insert error";
   }

   db.Close();

   return "ok";
}
 

func GetDbClient() *mysql.Client{
   //get client
   client, err := mysql.DialUnix(mysql.DEFAULT_SOCKET, "sit1", "sit1", "aimeedeer");
   if err != nil {
      log.Printf("%v\n", err);
      return nil;
   }
   
   return client;
}


func main() {
    web.Get("/user/login?(.*)", Login);
    web.Get("/user/register?(.*)", Register);
    web.Run("0.0.0.0:8080");
}