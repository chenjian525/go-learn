package sina

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Install_key_string(t *testing.T) {
	//	a := map[string]interface{}{"service": "create_activate_member", "identity_id": "gomember2", "identity_type": "UID", "client_ip": "127.0.0.1"}
	//	a := map[string]interface{}{"service": "set_real_name", "identity_id": "gomember", "identity_type": "UID", "real_name": "陈科明", "cert_type": "IC", "cert_no": "513429198502117617", "client_ip": "127.0.0.1"}
	//	a := map[string]interface{}{"service": "query_middle_account"}
	a := map[string]interface{}{"service": "create_hosting_deposit", "identity_id": "83326", "identity_type": "UID", "out_trade_no": "go-trade-number11111", "summary": "go-trade-number11111", "account_type": "SAVING_POT", "amount": 555, "payer_ip": "127.0.0.1", "pay_method": "online_bank^555.00^SINAPAY,DEBIT,C"}
	si, _ := New(a)
	resp, err := si.Fetch()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	x := 3
	y, _ := json.Marshal(x)
	fmt.Println(len(y))
}
