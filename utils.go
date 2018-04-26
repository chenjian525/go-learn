package sina

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("crypto/tls: found unknown private key type in PKCS#8 wrapping")
		}
	}

	return nil, errors.New("crypto/tls: failed to parse private key")
}

func installPubKey(path string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(path))
	if block == nil {
		fmt.Println("pem decode %s pubKey error", path)
		return nil, errors.New("install pubkey failed")
	}
	switch block.Type {
	case "PUBLIC KEY":
		res, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			fmt.Errorf("parse %s pubKey error.")
			return nil, err
		}
		return res.(*rsa.PublicKey), err
	default:
		fmt.Errorf("unsupport key type %s", path)
		return nil, nil
	}
}

func installPrivateKey(path string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(path))
	if block == nil {
		fmt.Printf("pem deocde %s error", path)
		return nil, errors.New("pem decode error")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		res, err := parsePrivateKey(block.Bytes)
		if err != nil {
			fmt.Printf("parse privkey privKey error.")
			fmt.Println()
			return nil, err
		}
		return res.(*rsa.PrivateKey), err
	default:
		return nil, nil
	}
}

func formateTime() string {
	t := time.Now()
	x, y, z := t.Clock()
	l, m, n := t.Date()
	a := fmt.Sprintf("%02d%02d%02d%02d%02d%02d", l, m, n, x, y, z)
	return a
}

func encrypt(s string, pubKey *rsa.PublicKey) (string, error) {
	b := []byte(s)
	encryptRes, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, b)
	if err != nil {
		fmt.Errorf("encrypt message error: %s", s)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptRes), nil
}

func encryptFields(m *map[string]interface{}, encryptKey *rsa.PublicKey, fields ...string) {
	for _, field := range fields {
		if val, ok := (*m)[field]; ok {
			encryptedStr, err := encrypt(val.(string), encryptKey)
			if err != nil {
				fmt.Errorf("%s", err)
			} else {
				(*m)[field] = encryptedStr
			}
		}
	}
}

func encryptPayMethod(m *map[string]interface{}, pubKey *rsa.PublicKey) {
	var pay_method string
	if payMethod, ok := (*m)["pay_method"]; ok {
		pay_method = payMethod.(string)
	} else {
		return
	}

	slice := strings.Split(pay_method, "^")
	if slice[0] == "quick_pay" {
		sub_slice := strings.Split(slice[2], ",")
		for _, i := range []int{1, 2, 6, 7} {
			encryptedVal, err := encrypt(sub_slice[i], pubKey)
			if err != nil {
				fmt.Errorf("%s", err)
			} else {
				sub_slice[i] = encryptedVal
			}
		}
		slice[2] = strings.Join(sub_slice, ",")
		(*m)["pay_method"] = strings.Join(slice, "^")
	}
}

func encryptSensitiveArgs(m *map[string]interface{}, encryptKey *rsa.PublicKey) {
	service := (*m)["service"]
	switch service {
	case "set_real_name":
		encryptFields(m, encryptKey, "real_name", "cert_no")
	case "binding_verify":
		encryptFields(m, encryptKey, "verify_entity")
	case "binding_bank_card":
		encryptFields(m, encryptKey, "bank_account_no", "account_name", "cert_no", "phone_no", "validity_period", "vertification_value")
	case "create_hosting_deposit":
		encryptPayMethod(m, encryptKey)
	}
}

func sha1Hash(s string) []byte {
	h := sha1.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func sign(s string, privKey *rsa.PrivateKey) (string, error) {
	d := sha1Hash(s)
	signRes, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA1, d)
	if err != nil {
		fmt.Errorf("rsa sign error: %s", s)
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(signRes), nil
	}
}

func signStr(m map[string]interface{}) string {
	keys := make([]string, 0)
	for k, _ := range m {
		if k != "sign" && k != "sign_type" && k != "sign_version" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var res string
	for _, k := range keys {
		v := fmt.Sprintf("%v", m[k])
		if v != "" {
			if len(res) > 0 {
				res += "&"
			}
			res += k + "=" + v
		}
	}
	return res
}

func signArgs(m map[string]interface{}, privKey *rsa.PrivateKey) string {
	signString := signStr(m)
	s, err := sign(signString, privKey)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		return s
	}
}

func buildArguments(m map[string]interface{}, privKey *rsa.PrivateKey, encryptKey *rsa.PublicKey) (map[string]interface{}, error) {
	m["version"] = "1.0"
	m["request_time"] = formateTime()
	m["partner_id"] = PARTNER_ID
	m["_input_charset"] = "UTF-8"
	m["notify_url"] = NOTIFY_URL
	m["return_url"] = RETURN_URL
	m["sign_type"] = "RSA"

	encryptSensitiveArgs(&m, encryptKey)
	m["sign"] = signArgs(m, privKey)
	return m, nil
}
