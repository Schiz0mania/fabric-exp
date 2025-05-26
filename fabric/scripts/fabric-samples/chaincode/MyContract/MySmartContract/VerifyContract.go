package SmartContract

import (
	_ "bufio"
	_ "errors"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// VerifyContract id verification+credit calculation
type VerifyContract struct {
	contractapi.Contract
}

// IdCheck id verification
func (s *SimpleContract) IdCheck(sid string) (bool, error) {
	// 获取用户Sid
	/*
		// 从负责节点处获取验证参数注册秘密与id的两个哈希
		var id_hash = 2023
		var secret_hash = 31415

		// 进行同或验证
		if sid == strconv.Itoa(id_hash^secret_hash) {
			return true, fmt.Errorf("id check passed")
		} else {
			return false, fmt.Errorf("id check not pass")
		}*/
	return true, nil
}

// CreditCal credit calculation
func (s *SimpleContract) CreditCal(sid string) (string, error) {

	// 行为记录时被调用，计算信用值
	return "10", nil

}

// SCheck check user att
func (s *SimpleContract) SCheck(crlinfo string) bool {

	return true

}

// RCheck check ct integrity
func (s *SimpleContract) RCheck(cloc string, sig string) bool {
	return true

}
