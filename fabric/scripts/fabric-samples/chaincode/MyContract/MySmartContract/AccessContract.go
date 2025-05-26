package SmartContract

import (
	_ "errors"
	_ "fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AccessControlContract  Contract:generate and alter user access control scheme
type AccessControlContract struct {
	contractapi.Contract
}

// GenIDS generate user inital control and id scheme
func (s *SimpleContract) GenIDS(sid string) (crlinfo string, err error) {
	//通过上层验证，生成身份属性与访问结构

	// 生成属性和访问结构
	return sid + " lsss matrix", nil

}

// AltS aid user update access control plan
func (s *SimpleContract) AltS(sid string) (crlinfo string, err error) {

	// 生成属性和访问结构
	return sid + " new lsss matrix", nil
}
