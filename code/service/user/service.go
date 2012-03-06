package main

import (
   "github.com/Philio/GoMySQL"
   "log"
   "json"
   "github.com/hoisie/web.go"
   "rand"
   "strconv"
)

func Login(ctx *web.Context,v string) string { 
   log.Printf("start login");

   if getSession(ctx,"user") != nil{
      b,_ := json.Marshal(getSession(ctx,"user"));
      return string(b);
   }
   
   prams := ctx.Request.Params;
   
   db := GetDbClient();
   if db == nil {
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

   if u == nil || len(u) == 0 {
      log.Printf("eerrrr");
      b,_ := json.Marshal(map[string]string{"ecode":"no user","emsg":"has no user found."});
      return string(b);
   }

   b,_ := json.Marshal(u);
   putSession(ctx,"user",u);
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

var sessionMap map[string]map[string]interface{};

func putSession(ctx *web.Context,key string,value interface{}){
   sessionid, _ := ctx.GetSecureCookie("sessionid");
   if sessionid == ""{
      sessionid = strconv.Itoa64(rand.Int63())
      ctx.SetSecureCookie("sessionid", sessionid, 3600)
      log.Printf("create session:session id=%v",sessionid);
   }
   if sessionMap == nil{
      sessionMap = map[string]map[string]interface{}{sessionid:nil};
   }
   if sessionMap[sessionid] == nil{
      sessionMap[sessionid] = map[string] interface{}{key:value};
   }else{
      sessionMap[sessionid][key] = value;
   }
}

func getSession(ctx *web.Context,key string) interface{}{
   sessionid, _ := ctx.GetSecureCookie("sessionid");
   return sessionMap[sessionid][key];
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
    web.Config.CookieSecret = "66d337519aa14ac4ac150f8569e2b719";
    web.Get("/service/user/login?(.*)", Login);
    web.Get("/service/user/register?(.*)", Register);
    web.Run("0.0.0.0:8080");
}