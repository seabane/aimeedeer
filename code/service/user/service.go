package main

import (
   "github.com/Philio/GoMySQL"
   "log"
   "json"
   "github.com/hoisie/web.go"
   "rand"
   "strconv"
   "fmt"
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

func AddThing(ctx *web.Context,v string) string{
	//check login status
	if getSession(ctx,"user") == nil{
		b,_ := json.Marshal(map[string]string{"ecode":"no login","emsg":"please login at first."});
		return string(b);
	}
	sessionUser := getSession(ctx,"user");
	user := sessionUser.(mysql.Map);
	prams := ctx.Request.Params;
	
	db := GetDbClient();
	if(db == nil){
      return "db conn err";
	}
	err := db.Query("insert into things (username,content) values('" + user["username"].(string) + "','" + prams["content"] + "')");
    if err != nil{
       log.Printf("%v\n,sql:%v",err,"insert into things (username,content) values('" + user["username"].(string) + "','" + prams["content"] + "'");
       return "db insert things error";
    }
	db.Close();
	
    return "ok";
}


func DelThing(ctx *web.Context,v string) string{
	if getSession(ctx,"user") == nil{
		b,_ := json.Marshal(map[string]string{"ecode":"no login","emsg":"please login at first."});
		return string(b);
	}
	
	prams := ctx.Request.Params;
	
	db := GetDbClient();
	if(db == nil){
      return "db conn err";
	}
	
	err := db.Query("delete from things where id = " +  prams["id"]);
    if err != nil{
       log.Printf("%v\n",err);
       return "db delete things error";
    }
	db.Close();
	
    return "ok";
}

type things struct{
     id string
     username string
     time_create string
     content string

}

func QueryThing(ctx *web.Context,v string) string{
	//check login status
	if getSession(ctx,"user") == nil{
	   log.Printf("%v\n",getSession(ctx,"user"));
	   b,_ := json.Marshal(map[string]string{"ecode":"no login","emsg":"please login at first."});
	   return string(b);
	}
	
	sessionUser := getSession(ctx,"user");
	user := sessionUser.(mysql.Map);
	
	db := GetDbClient();
	if(db == nil){
      return "db conn err";
	}
	stmt,err := db.Prepare("select * from things where username = ? order by time_create");
    if err != nil{
       log.Printf("%v\n",err);
       return "db select things error";
    }
	
    stmt.BindParams(user["username"].(string));
    stmt.Execute();

    thing := things{};
    stmt.BindResult(&thing.id,&thing.username,&thing.time_create,&thing.content)

    
    list := make([]map[string]string,1,100);
    i := 0;
    for {  
       eof, _ := stmt.Fetch()  
       if eof {  
          break  
       }  
       list[i] = map[string]string{"id":thing.id,"username":thing.username,"time_create":thing.time_create,"content":thing.content};
       i++;
    }
    db.Close();
    log.Printf("i=%v,len(list=%v,list=%v",i,len(list),list);
    b,_ := json.Marshal(list);
    return string(b);
    
}


var sessionMap map[string]map[string]interface{};

func putSession(ctx *web.Context,key string,value interface{}){
   sessionid, _ := ctx.GetSecureCookie("sessionid");
   if sessionid == ""{
      sessionid = strconv.Itoa64(rand.Int63())
      //ctx.SetSecureCookie("sessionid", sessionid, 3600)
      cookie := fmt.Sprintf("%s=%s;path=/;", key, sessionid)
      ctx.SetHeader("Set-Cookie", cookie, false)
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
   //sessionid, err := ctx.GetSecureCookie("sessionid");
   for _, cookie := range ctx.Request.Cookie {
        if cookie.Name != key {
            continue
        }

        return sessionMap[cookie.Value][key];
   }
   
   return nil;
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
    web.Get("/user/login?(.*)", Login);
    web.Get("/user/register?(.*)", Register);
    web.Get("/things/query?(.*)",QueryThing);
    web.Get("/things/del?(.*)",DelThing);
    web.Get("/things/add?(.*)",AddThing);
    web.Run("0.0.0.0:8080");
}