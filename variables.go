package sina

const (
	encryptPem = `\
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCBpueN
weMbYdb+CMl8dUNv5g5THYLD9Z33cAMA4GNjmPYsbcNQ
LyO5QSlLNjpbCwopt7b5lFP8TGLUus4x0Ed6S4Wd9KmN
w6NLbszNEmppP9HXlT9sT4/ShL0CpVF4ofFS8O/g
XwCTJjYZJ0HvK3GBTSP2C9WlipTpWQ+9QJugewIDAQAB
-----END PUBLIC KEY-----
`

	privPem = `\
-----BEGIN RSA PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAO/
6rPCvyCC+IMalLzTy3cVBz/+wamCFNiq9qKEilEBDTttP7Rd/GAS51lsf
CrsISbg5td/w25+wulDfuMbjjlW9Afh0p7Jscmbo1skqIOIUPYfVQEL6
87B0EmJufMlljfu52b2efVAyWZF9QBG1vx/AJz1EVyfskMaYVqPiTesZ
AgMBAAECgYEAtVnkk0bjoArOTg/KquLWQRlJDFrPKP3CP25wHsU4
749t6kJuU5FSH1Ao81d0Dn9m5neGQCOOdRFi23cV9gdFKYMhwP
E6+nTAloxI3vb8K9NNMe0zcFksva9c9bUaMGH2p40szMoOpO6Tr
SHO9Hx4GJ6UfsUUqkFFlN76XprwE+ECQQD9rXwfbr9GKh9QMNv
nwo9xxyVl4kI88iq0X6G4qVXo1Tv6/DBDJNkX1mbXKFYL5NOW1wa
ZzR+Z/XcKWAmUT8J9AkEA8i0WT/ieNsF3IuFvrIYG4WUadbUqObcY
P4Y7Vt836zggRbu0qvYiqAv92Leruaq3ZN1khxp6gZKl/OJHXc5xzQJ
ACqr1AU1i9cxnrLOhS8m+xoYdaH9vUajNavBqmJ1mY3g0IYXhcbF
m/72gbYPgundQ/pLkUCt0HMGv89tn67i+8QJBALV6UgkVnsIbkkK
COyRGv2syT3S7kOv1J+eamGcOGSJcSdrXwZiHoArcCZrYcIhOxOW
B/m47ymfE1Dw/+QjzxlUCQCmnGFUO9zN862mKYjEkjDN65n1IUB
9Fmc1msHkIZAQaQknmxmCIOHC75u4W0PGRyVzq8KkxpNBq62IC
l7xmsPM=
-----END RSA PRIVATE KEY-----
`

	pubPem = `\
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDv0rdsn5FYPn0EjsCPqDyIsYRa
wNWGJDRHJBcdCldodjM5bpve+XYb4Rgm36F6iDjxDbEQbp/HhVPj0XgGlCRKpblu
yJJt8ga5qkqIhWoOd/Cma1fCtviMUep21hIlg1ZFcWKgHQoGoNX7xMT8/0bEslda
KdwxOlv3qGxWfqNV5QIDAQAB
-----END PUBLIC KEY-----
`

	MGS_GATE   = "https://testgate.pay.sina.com.cn/mgs/gateway.do"
	MAS_GATE   = "https://testgate.pay.sina.com.cn/mas/gateway.do"
	PARTNER_ID = "200004595271"
	NOTIFY_URL = "http://83.rmb.io/api/sina/notify"
	RETURN_URL = "http://83.rmb.io/user/recharges"
)

var MGS = []string{
	"create_activate_member", "set_real_name", "binding_verify", "unbinding_verify", "query_verify", "binding_bank_card",
	"binding_bank_card_advance", "unbinding_bank_card", "query_bank_card", "query_balance",
	"balance_freeze", "balance_unfreeze", "query_account_details", "query_middle_account",
}

type sinaResp struct {
	Response_time    string `json:"response_time"`
	Partner_id       string `json:"partner_id"`
	Input_charset    string `json:"_input_charset"`
	Sign             string `json:"sign"`
	Sign_type        string `json:"sign_type"`
	Sign_version     string `json:"sign_version"`
	Response_code    string `json:"response_code"`
	Response_message string `json:"response_message"`
	Memo             string `json:"memo"`
	Error_url        string `json:"error_url"`
}
