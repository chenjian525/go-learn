package sina

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	//	"strings"
	//	"os"
	"encoding/json"
)

type Sina struct {
	privKey    *rsa.PrivateKey
	pubKey     *rsa.PublicKey
	encryptKey *rsa.PublicKey
	url        string
	arguments  map[string]interface{}
}

func New(m map[string]interface{}) (*Sina, error) {
	privKey, _ := installPrivateKey(privPem)
	pubKey, _ := installPubKey(pubPem)
	encryptKey, _ := installPubKey(encryptPem)

	reqUrl := MAS_GATE
	if service, ok := m["service"]; ok {
		for _, val := range MGS {
			if val == service {
				reqUrl = MGS_GATE
			}
		}
	} else {
		return nil, errors.New("service variable must be needed")
	}

	arguments, err := buildArguments(m, privKey, encryptKey)
	if err != nil {
		return nil, errors.New("build arguments failed")
	}
	return &Sina{privKey, pubKey, encryptKey, reqUrl, arguments}, nil
}

func (s Sina) Fetch() (map[string]string, error) {
	//	data := url.Values{}
	//	fmt.Printf("%+v", s.arguments)
	//	fmt.Println()
	//	for k, v := range s.arguments {
	//		data.Add(k, fmt.Sprintf("%s", v))
	//	}

	//	resp, err := http.PostForm(s.url, data)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}

	//	tmp := make([]string, 0)
	//	for k, v := range s.arguments {
	//		fmt.Println(k, v)
	//		tmp = append(tmp, k+"="+fmt.Sprintf("%s", v))
	//	}
	//	ss := strings.Join(tmp, "&")
	//	req, _ := http.NewRequest("POST", s.url, bytes.NewBuffer([]byte(ss)))
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//	client := http.Client{}
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	for k, v := range s.arguments {
		s.arguments[k] = quote(fmt.Sprintf("%v", v))
	}

	tmp := []string{}
	for k, _ := range s.arguments {
		v := fmt.Sprintf("%v", s.arguments[k])
		if v != "" {
			v = url.QueryEscape(v)
			tmp = append(tmp, k+"="+v)
		}
	}
	a := strings.Join(tmp, "&")
	body := []byte(a)
	resp, err := http.Post(s.url, "application/x-www-form-urlencoded;charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("network error, no resp")
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("read resp error")
	}
	escaptedStr, _ := url.QueryUnescape(string(bodyBytes))
	bodyBytes = []byte(escaptedStr)
	resp_ := map[string]string{}
	json.Unmarshal(bodyBytes, &resp_)
	if len(resp_) == 0 {
		return map[string]string{"result": escaptedStr}, nil
	}

	if s.checkSign(resp_) {
		return resp_, nil
	} else {
		return nil, errors.New("wrong sign")
	}
}

func (s Sina) checkSign(m map[string]string) bool {
	m_ := make(map[string]interface{})
	for k, v := range m {
		m_[k] = v
	}
	signBytes := sha1Hash(signStr(m_))
	sig, _ := base64.StdEncoding.DecodeString(m["sign"])
	err := rsa.VerifyPKCS1v15(s.pubKey, crypto.SHA1, signBytes, sig)
	if err == nil {
		return true
	} else {
		return false
	}
}
