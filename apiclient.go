package apiclient 

import (
    "github.com/ddliu/go-httpclient"
    "encoding/json"
    "fmt"
    "github.com/firewut/go-json-map"
    //"net/url"
    "net/http"
    "time"
)

func init() {
   
}


func GetJson(url string, target interface{}) error {
    var myClient = &http.Client{Timeout: 10 * time.Second}
    r, err := myClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}


func SendGet(champs []string, url string,params map[string]string) (  map[string]string,error) {
    var x map[string]interface{}
    var resultat map[string]string
    resultat = make(map[string]string)


    httpclient.Defaults(httpclient.Map {
        httpclient.OPT_USERAGENT: "theagent",
        "Accept-Language": "fr-fr",
    })

    res, err := httpclient.Get(url, params)
    
    if res!=nil && err==nil {
        bodyString,_ := res.ToString()
        json.Unmarshal([]byte(bodyString), &x)
    
        for _,champ :=range champs {
            property, _ := gjm.GetProperty(x, champ)

            resultat[champ]=fmt.Sprint(property)
        }
    }
    return resultat,err
}

func SendDelete(url string) (error) {
    httpclient.Defaults(httpclient.Map {
        httpclient.OPT_USERAGENT: "theagent",
        "Accept-Language": "fr-fr",
    })

    _, err := httpclient.Post(url, nil)
    //bodyString,_ := res.ToString()
    //json.Unmarshal([]byte(bodyString), &x)
    
    return err
}

func SendPost(adresse string,params map[string]string) (map[string]string,error) {
  var x map[string]interface{}


    httpclient.Defaults(httpclient.Map {
        
        "opt_useragent": "theagent",
        "opt_timeout": 10,
        "Accept-Encoding": "gzip,deflate,sdch",
        "Accept-Language": "en-us",
        "Content-Type": "application/json"})


    
   // adresse=url.QueryEscape(adresse)
    
    var ret map[string]string
    ret = make(map[string]string)
    adresse = fmt.Sprint(adresse)
    fmt.Printf("%s%s",adresse,"fin")

    res, err := httpclient.Post(adresse, params)
    
    if res!=nil && err==nil{
    
        
        bodyString,err := res.ToString()
        json.Unmarshal([]byte(bodyString), &x)
    
        for _, record := range x {
            rec, _ := record.(map[string]interface{})
            for key,value := range rec {
                ret[fmt.Sprint(key)]=fmt.Sprint(value)
                //fmt.Println("VAR MAP:",)
            }
            return ret,err
        }
    } else {
        
        fmt.Println(err,ret)
    }

    return ret,err

    
}   
