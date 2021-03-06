package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"
)

type Dog struct {
	Age int `json:",omitempty"`
}

func TestHandleExecShell(t *testing.T) {
	age := 1
	d := Dog{
		Age: age,
	}

	b, _ := json.Marshal(d)
	fmt.Println(string(b))
	var dd Dog
	err := json.Unmarshal(b, &dd)
	if err != nil {
		fmt.Println("dd", err)
	}

	ddd := new(Dog)
	err = json.Unmarshal(b, ddd)
	if err != nil {
		fmt.Println("ddd", err)
	}
}

func TestGet(t *testing.T) {
	var f float64 = 1.23
	fmt.Printf("%0.f", f)
}

func TestResponseRecorder(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))

	// Output:
	// 200
	// text/html; charset=utf-8
	// <html><body>Hello World!</body></html>
}

func TestPipe(t *testing.T) {
	c1 := exec.Command("ls")
	c2 := exec.Command("wc", "-l")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	io.Copy(os.Stdout, &b2)
}

func TestEnv(t *testing.T) {
	os.Setenv("FOO", "1")
	fmt.Println("FOO:", os.Getenv("FOO"))
}

func TestJson(t *testing.T) {
	type FruitBasket struct {
		Name    string
		Fruit   []string
		Id      int64 `json:"ref"` // 声明对应的json key
		Created time.Time
	}

	jsonData := []byte(`
    {
        "Name": "Standard",
        "Fruit": [
             "Apple",
            "Banana",
            "Orange"
        ],
        "ref": 999,
        "Created": "2018-04-09T23:00:00Z"
    }`)

	var basket FruitBasket
	err := json.Unmarshal(jsonData, &basket)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(basket.Name, basket.Fruit, basket.Id)

	data, err := json.Marshal(jsonData)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("jsonData:", string(jsonData))
	fmt.Println("jsonMarshal:", string(data))
}
