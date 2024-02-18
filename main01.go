package main

import (
    "fmt"
    "time"
    "strconv"
	"net/http"
    "strings"
	"sync"
    //"os"
)

type Operation struct {
    Value       float64
    OperRune    rune
}
type Info struct {
    CountStr    string
    IsOpered    bool
    Opers       []Operation
    TimePlus    time.Duration
    TimeMinus   time.Duration
    TimeMult    time.Duration
    TimeDivis   time.Duration
    ID          int
    ProcessTime time.Duration
    Value       float64
    Error       error
}

var TimePlus, TimeMinus, TimeMult, TimeDivis time.Duration
var Information []Info
var ch []int
var mu sync.Mutex

func MainHandler(w http.ResponseWriter, r *http.Request) {
    st := r.URL.Query().Get("nm")
    st = strings.ReplaceAll(st, string(rune(32)), "+")
    fmt.Fprintln(w, st)
    
    id := len(Information)
    Information = append(Information, Info{st, false, nil, TimePlus, TimeMinus, TimeMult, TimeDivis, id, 0, 0, fmt.Errorf("200")})
    //value, timeProcess, err := 
    fmt.Fprintln(w, id, len(Information))
    //go CountProcess(&Information[id])
    mu.Lock()
    ch = append(ch, id)
    mu.Unlock()
 }
func TimeHandler(w http.ResponseWriter, r *http.Request) {
    sleepDurPlus, errPlus := time.ParseDuration(r.URL.Query().Get("timePlus"))    //time.Millisecond * 1100
    sleepDurMinus, errMinus := time.ParseDuration(r.URL.Query().Get("timeMinus")) 
    sleepDurMult, errMult := time.ParseDuration(r.URL.Query().Get("timeMult")) 
    sleepDurDivis, errDivis := time.ParseDuration(r.URL.Query().Get("timeDivis")) 
    if errPlus != nil {
        fmt.Fprintln(w, errPlus)
    } else {
        TimePlus = sleepDurPlus
    }
    if errMinus != nil {
        fmt.Fprintln(w, errMinus)
    } else {
        TimeMinus = sleepDurMinus
    }
    if errDivis != nil {
        fmt.Fprintln(w, errDivis)
    } else {
        TimeDivis = sleepDurDivis
    }
    if errMult != nil {
        fmt.Fprintln(w, errMult)
    } else {
        TimeMult = sleepDurMult
    }
    fmt.Fprintln(w, "Время сложения:", TimePlus)
    fmt.Fprintln(w, "Время вычитания:", TimeMinus)
    fmt.Fprintln(w, "Время деления:", TimeDivis)
    fmt.Fprintln(w, "Время умножения:", TimeMult)
 }
func InfoHandler(w http.ResponseWriter, r *http.Request) {
    id,err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id >= len(Information) {
        fmt.Fprintln(w, fmt.Errorf("Неправильный ID"))
    } else {
        fmt.Fprintf(w, "ID: %d\nValue: %f\nCode: %s\nProcessTime: %s\nString: %s", id, Information[id].Value, Information[id].Error, Information[id].ProcessTime,Information[id].CountStr)
    }
 }
func DataHandler(w http.ResponseWriter, r *http.Request) {
    for i := range Information {
        fmt.Fprintf(w, "ID: %d Value: %f Code: %s ProcessTime: %s String: %s\n", Information[i].ID, Information[i].Value, Information[i].Error, Information[i].ProcessTime,Information[i].CountStr)
    }  
    mu.Lock()
    fmt.Fprintln(w,ch)
    mu.Unlock()
 }

func Meine(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//start := time.Now()

		next.ServeHTTP(w, r)

		//duration := time.Since(start)
	})
 }

func main() {
    TimePlus = time.Second
    TimeMinus = time.Second
    TimeDivis = time.Second
    TimeMult = time.Second
    /*file_path := "C:\\Users\\aleks\\codeing\\go\\code\\ЯндексЛицей\\2\\финальная_задача\\data.txt"
    _, err := os.ReadFile(file_path)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("ok")
    }*/
    Information = make([]Info, 0)
    ch = make([]int, 0)

	mux := http.NewServeMux()

	main_ := http.HandlerFunc(MainHandler)
    time_ := http.HandlerFunc(TimeHandler)
    info_ := http.HandlerFunc(InfoHandler)
    data_ := http.HandlerFunc(DataHandler)
    
	mux.Handle("/", Meine(main_))
    mux.Handle("/times/", Meine(time_))
    mux.Handle("/get/", Meine(info_))
    mux.Handle("/data/", Meine(data_))

    go func(){
        for {
            mu.Lock()
            if len(ch) > 0 {
                ch = ch[1:]
                go CountProcess(&Information[ch[0]])
            }
            mu.Unlock()
            time.Sleep(time.Millisecond * 200)
        }
    }()
    fmt.Println("start ok")
	http.ListenAndServe(":1229", mux)
 }

func CountProcess(inf *Info) (float64, time.Duration, error){
    ls := inf.Opers
    runes := []rune(inf.CountStr)
    if !inf.IsOpered {
      ls = make([]Operation, 0)
      for i := 0; i < len(runes); i+=1 {
        char := runes[i]
        a := float64(0)
        b := false
        numberFromRune, err := strconv.Atoi(string(char))
        for (err == nil) {
            b = true
            a = a * 10 + float64(numberFromRune)
            i += 1
            if i >= len(runes) {
                break
            }
            char = runes[i]
            numberFromRune, err = strconv.Atoi(string(char))
        }
        if char == '.' || char == ',' {
            if b == false {
                inf.Error = fmt.Errorf("400")
                return 0,0,fmt.Errorf("400")
            }
            numAfterPoint := 1
            i+=1
            char = runes[i]
            numberFromRune, err = strconv.Atoi(string(char))
            for (err == nil) {
                b = true
                a += float64(numberFromRune) / float64(nInDegree(10, numAfterPoint))
                i += 1
                if i >= len(runes) {
                    break
                }
                numAfterPoint += 1
                char = runes[i]
                numberFromRune, err = strconv.Atoi(string(char))
            }
        }
        if b == true {
            i-=1
            ls = append(ls, Operation{a, 'n'})
        } else {
            if isRuneCorrect(char) == false {
                inf.Error = fmt.Errorf("400")
                return 0,0,fmt.Errorf("400")
            }
            ls = append(ls, Operation{0,char})
        }

      }
      inf.IsOpered = true
      inf.Opers = ls
    }
    t0 := time.Now()
    if ls[0].OperRune != 'n' || ls[len(ls) - 1].OperRune !='n' {
        return 0,0,fmt.Errorf("400")
    } else {
        for i := 1; i < len(ls) - 1; i += 1 {
            x := ls[i].OperRune
            if x == '*' || x == '/' {
                if ls[i-1].OperRune != 'n' || ls[i+1].OperRune != 'n' {
                    inf.Error = fmt.Errorf("400")
                    return 0,0,fmt.Errorf("400")
                }
                res, err := doCount(ls[i-1].Value, ls[i+1].Value, x, inf)
                if err!=nil {
                    inf.Error = err
                    return 0,0,err
                }
                ls[i - 1] = Operation{res, 'n'}
                ls = remove(ls, i, i + 1)
                i-=2
            }
        }
        for i := 1; i < len(ls) - 1; i += 1 {
            x := ls[i].OperRune
            if x == '-' || x == '+' {
                if ls[i-1].OperRune != 'n' || ls[i+1].OperRune != 'n' {
                    return 0,0,fmt.Errorf("400")
                }
                res, err := doCount(ls[i-1].Value, ls[i+1].Value, x, inf)
                if err!=nil {
                    inf.Error = err
                    return 0,0,err
                }
                ls[i - 1] = Operation{res, 'n'}
                ls = remove(ls, i, i + 1)
                i-=2
            }
        }
    }
    inf.ProcessTime = time.Now().Sub(t0)
    inf.Value = ls[0].Value
    //fmt.Println(ls)
    return inf.Value, inf.ProcessTime, nil
}

func doCount(n1, n2 float64, operation rune, inf *Info) (float64, error) {
    x := 0.0
    var dur time.Duration
    if operation == '+' {
        x = n1 + n2
        dur = inf.TimePlus
    } else if operation == '-' {
        x = n1 - n2
        dur = inf.TimeMinus
    } else if operation == '*' {
        x = n1 * n2
        dur = inf.TimeMult
    } else if operation == '/' {
        if n2 == 0 {
            return 0, fmt.Errorf("400")
        } else {
            x = n1 / n2
            dur = inf.TimeMinus
        }
    } 
    time.Sleep(dur)
    return x, nil
 }
func remove(slice []Operation, s, s2 int) []Operation {
    return append(slice[:s], slice[s2+1:]...)
 }
func isRuneCorrect(r rune) bool {
    return (r == '*' || r == '/' || r == '-' || r == '+')
 } 
func nInDegree(a, b int) int {
    p := 1
    for i := 0; i < b; i+= 1{
        p*=a
    }
    return p 
 }





/*
    go run "C:\Users\aleks\codeing\go\code\ЯндексЛицей\2\финальная_задача\main01.go"
*/