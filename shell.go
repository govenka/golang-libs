package shell 

import (
	"os/exec"
    //"os"
	"strings"
	"fmt"
	"unicode"
	"golang.org/x/crypto/ssh"
	"github.com/sparrc/go-ping"
    "github.com/sfreiberg/simplessh"
    "io/ioutil"
	//"log"
	//"os"
)

func init() {
   
}

func Execute(cmd string) string {

	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]
	
	out, err := exec.Command(head,parts...).Output()
	if err != nil {
		return fmt.Sprint(err)
	}
	return string(out)

}

func GetUUID() string {
	var ret string
	ret=strings.Replace(Execute("cat /sys/class/net/eth0/address"),":","",-1)
	ret=strings.Replace(ret,"\n","",-1)
	return ret
}

func WriteConfigInFile(file string,configs string) (error,bool) {

    var changed bool

    changed = false
    fileinfos, err := ioutil.ReadFile(file)
    if strings.Compare(string(fileinfos),configs) == 0 {
            fmt.Println("conf unchanged in file ",file)
            changed = false
    } else {
            fmt.Println("Rewrite conf in file ",file)
            ioutil.WriteFile(file,[]byte(configs),0644)
            changed = true
    }

    return err,changed
}

func connectToHost(host string, user string, password string) (*ssh.Client, *ssh.Session, error) {
	
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "blowfish-cbc", "cast128-cbc", "aes192-cbc", "aes256-cbc", "arcfour"},
		},
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func SendSSH(host string, user string, password string, command string) string {

	client, session, err := connectToHost(host, user, password)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	out, err := session.CombinedOutput(command)
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
	
}

func GetSshFile(host string, user string, password string, source string, dest string) error {
    
    client, err := simplessh.ConnectWithPassword(host,user,password)
    if err != nil {
        return err
    }
    defer client.Close()


    err=client.Download(source,dest)
    if err != nil {
        return err
    }

    return nil

}

func GetServerMemoryUsage() string {
 
	free := Execute("free");
	return free;
	/*
	$free = (string)trim($free);
	$free_arr = explode("\n", $free);
	$mem = explode(" ", $free_arr[1]);
	$mem = array_filter($mem);
	$mem = array_merge($mem);
	$memory_usage = $mem[2]/$mem[1]*100;
 
	return $memory_usage;
	*/
}

func GetServerCpuUsage() string {
 
	load := Execute("cat /proc/loadavg");
	return load;
 
}


func GetSystemMemInfo() string {       
    data := Execute("cat /proc/meminfo");
    return data;
}

func Ping(address string) string {
	Execute("sysctl -w net.ipv4.ping_group_range=\"0   2147483647\"")
	pinger, err := ping.NewPinger(address)
	if err != nil {
        return "ERROR"
	} else {
		
		pinger.Count = 1
		pinger.Run() 
		stats := pinger.Statistics().AvgRtt
		statsstring := strings.Replace(fmt.Sprint(stats),"ms","",-1)
		fmt.Printf("PINGER :%v\n",statsstring)
		return statsstring
	}
}

func SplitQuoted(s string) []string {
    var ret []string
    var curr = make([]rune, len(s))
    var cpos = 0
    var quoted = ' '
    var escaped = false
    sr := strings.NewReader(s)

    for {
        r, _, err := sr.ReadRune()
        if err != nil {
            // Append last
            if cpos != 0 {
                ret = append(ret, string(curr[0:cpos]))
                cpos = 0
            }
            break
        }
        switch r {
        case ' ':
            if quoted != ' ' {
                if escaped {
                    curr[cpos] = '\\'
                    cpos++
                    escaped = false
                }
                curr[cpos] = ' '
                cpos++
            } else if escaped {
                curr[cpos] = ' '
                cpos++
                escaped = false
            } else if cpos != 0 {
                ret = append(ret, string(curr[0:cpos]))
                cpos = 0
            }
        case '"', '\'':
            if escaped {
                curr[cpos] = r
                cpos++
                escaped = false
            } else if quoted == r {
                quoted = ' '
            } else if quoted == ' ' {

                quoted = r
            } else {
                curr[cpos] = r
                cpos++
            }
        case '\\':
            if escaped {
                curr[cpos] = '\\'
                cpos++
                escaped = false
            } else {
                escaped = true
            }
        default:
            if unicode.IsSpace(r) {
                
                if quoted != ' ' {
                    curr[cpos] = r
                    cpos++
                } else if cpos != 0 {
                    
                    ret = append(ret, string(curr[0:cpos]))
                    cpos = 0
                }
            } else {
                if escaped {
                    curr[cpos] = '\\'
                    cpos++
                    escaped = false
                }
                curr[cpos] = r
                cpos++
            }
        }
    }
    return ret
}
