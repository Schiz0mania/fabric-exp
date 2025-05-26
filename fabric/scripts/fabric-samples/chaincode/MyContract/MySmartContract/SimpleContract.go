package SmartContract

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strings"
)

// SimpleContract  for handling writing and reading from the world state
type SimpleContract struct {
	contractapi.Contract
}

// 规定所有传输string类，如包含多个数据，皆为json格式
/*
todo:问题在于实验要达到什么地步 合约用go写，app交互可以用其他语言写
	1. 所有合约以首字母简写代表，S用于读写账本，A，V用于规定一些业务逻辑
	2. 对于S的读写账本部分，特别是写账本由四种行为数据，如下规定
	3. 流程，这里就涉及到app的交互了
	用户注册（需要补全），调用V中的IdCheck（）验证id，根据下方结构打包数据，调用S内的写
	用户加密数据，这里应该不需要定义，同上调用
	用户请求搜索，验证后调用S内的QueryCIndex（），获取Cindex然后进行关键词搜索，搜索完成后要求CS返回密文，进行解密，最后打包数据（这里成功或者失败都会记录在log里）
	用户请求修改属性或者访问结构更新，该交互调用A中的函数AltS（），打包数据
	（脚本负责业务交互，返回打包好的业务数据给S进行数据分类打包）
	4. 补全S中的读写操作
	5. 交互的脚本选择语言，并开始学习

*/

// Behavior 用户行为
type Behavior struct {
	Sid     string `json:"sid"` // 用户系统id
	UserCrd string // 用户信用
	ActTag  string `json:"actTag"`  // 行为标记，四种R,UP,SNF
	ActInfo string `json:"actInfo"` // 对应的行为信息
}

/*
四种行为的actInfo
Register:

	用户公钥pk
	用户身份信息（属性与访问控制结构）
	对应负责注册机构名称instCode

Upload：

	数据加密的一些参数，统称为ECP
	密文索引CIndex
	密文位置LOC

SNF：

	被请求者的Sid
	一些log即可

RA:

	同register
*/
type RInfo struct {
	Pk       string `json:"pk"`       // 用户公钥
	CrlInfo  string `json:"crlInfo"`  // 身份属性与控制结构
	InstCode string `json:"instCode"` // 负责的机构或者端口标识
}
type UpInfo struct {
	EncrptParam string `json:"encrptParam"` // 加密参数
	CIndex      string `json:"CIndex"`      // 密文索引
	Cloc        string `json:"Cloc"`        // 密文位置
}
type SNFInfo struct {
	DoSid string `json:"DoSid"` // DO的Sid
	Log   string `json:"log"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Behavior
}

func (s *SimpleContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	/*rinfo := []RInfo{
		{
			Pk:       "12-111",
			CrlInfo:  "patient",
			InstCode: "a21",
		},
		{
			Pk:       "13-222",
			CrlInfo:  "doctor",
			InstCode: "c31",
		},
		{
			Pk:       "14-333",
			CrlInfo:  "researcher",
			InstCode: "z12",
		},
	}
	auxstr := []string{
		"", "", "",
	}
	str, _ := json.Marshal(rinfo[0])
	auxstr[0] = string(str)
	str, _ = json.Marshal(rinfo[1])
	auxstr[1] = string(str)
	str, _ = json.Marshal(rinfo[2])
	auxstr[2] = string(str)

	uinfo := UpInfo{
		EncrptParam: "secret",
		CIndex:      "no encrypt",
		Cloc:        "12.21.43",
	}

	str, _ = json.Marshal(uinfo)
	newstr := string(str)

	bh := []Behavior{
		{
			Sid:     "111",
			UserCrd: "5",
			ActTag:  "R",
			ActInfo: auxstr[0],
		},
		{
			Sid:     "222",
			UserCrd: "5",
			ActTag:  "R",
			ActInfo: auxstr[1],
		},
		{
			Sid:     "333",
			UserCrd: "5",
			ActTag:  "R",
			ActInfo: auxstr[2],
		},
		{
			Sid:     "111",
			UserCrd: "5",
			ActTag:  "UP",
			ActInfo: newstr,
		},
	}

	for i, bhiter := range bh {
		bhAsByte, _ := json.Marshal(bhiter)
		err := ctx.GetStub().PutState("LOG"+strconv.Itoa(i), bhAsByte)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}*/

	return nil
}

func (s *SimpleContract) WriteLedger(ctx contractapi.TransactionContextInterface, lognum string, bh Behavior) error {
	bhAsBytes, _ := json.Marshal(bh)
	return ctx.GetStub().PutState(lognum, bhAsBytes)

}

func (s *SimpleContract) ReadAllLedger(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		bh := new(Behavior)
		_ = json.Unmarshal(queryResponse.Value, bh)

		queryResult := QueryResult{Key: queryResponse.Key, Record: bh}
		results = append(results, queryResult)
	}

	return results, nil
	// 读取world state账本数据进行输出
}

func (s *SimpleContract) ReadLedger(ctx contractapi.TransactionContextInterface, lognum string) (*Behavior, error) {

	bhAsBytes, err := ctx.GetStub().GetState(lognum)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if bhAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", err.Error())

	}
	bh := new(Behavior)
	_ = json.Unmarshal(bhAsBytes, bh)
	return bh, nil
}

func (s *SimpleContract) ReadLedger_Sid(ctx contractapi.TransactionContextInterface, sid string, actTag string) (*Behavior, error) {

	startKey := ""
	endKey := ""
	bhFinal := new(Behavior)
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		bh := new(Behavior)
		_ = json.Unmarshal(queryResponse.Value, bh)

		if strings.EqualFold(bh.Sid, sid) && strings.EqualFold(bh.ActTag, actTag) {
			bhFinal = bh
			// 只找最后一个

		} else {
			continue
		}

	}

	return bhFinal, err
	// 读取world state账本数据进行输出
}

func (s *SimpleContract) RegisterUser(ctx contractapi.TransactionContextInterface, sid string, pk string) error {
	// 用户提供sid进行验证,并产生访问策略数据
	flag := false
	var err error
	var crlinfo string
	flag, err = s.IdCheck(sid)

	if flag {
		crlinfo, err = s.GenIDS(sid)

	} else {
		return fmt.Errorf("register failed%s", err)
	}

	// 信用计算
	var usr_crd, _ = s.CreditCal(sid)

	// 数据打包
	rinfo := RInfo{
		Pk:       pk,
		CrlInfo:  crlinfo,
		InstCode: "a21",
	}
	actinfo, _ := json.Marshal(rinfo)

	bh := Behavior{
		Sid:     sid,
		UserCrd: usr_crd,
		ActTag:  "R",
		ActInfo: string(actinfo),
	}

	//写账本
	return s.WriteLedger(ctx, "LOG01", bh)
}

func (s *SimpleContract) Encryption(ctx contractapi.TransactionContextInterface, sid string, auxParam string) error {
	//用户验证

	var info RInfo
	// 根据sid寻找用户pk，crlinfo
	var bh = new(Behavior)

	var _ error
	bh, _ = s.ReadLedger_Sid(ctx, sid, "R")
	// 数据解包
	// string转结构体
	str := []byte(bh.ActInfo)
	_ = json.Unmarshal([]byte(str), &info)
	// 获取info中的crlinfo与pk

	// 计算cindex
	var Cindex = info.CrlInfo + " " + auxParam

	// 上传密文，等待cloc
	var cloc = "C12.R3.B4.11231"

	//打包数据
	var actbyte []byte
	var actinfo = UpInfo{
		EncrptParam: auxParam + info.CrlInfo,
		CIndex:      Cindex,
		Cloc:        cloc,
	}
	actbyte, _ = json.Marshal(actinfo)
	bh.ActTag = "UP"
	bh.ActInfo = string(actbyte)
	//写账本

	return s.WriteLedger(ctx, "LOG02", *bh)
}

func (s *SimpleContract) UpdateCIndex(ctx contractapi.TransactionContextInterface, sid string, newindex string) error {
	// 更新
	bh, err := s.ReadLedger_Sid(ctx, sid, "UP")
	var info UpInfo
	if err != nil {
		return err
	}
	if bh.ActTag == "UP" {
		// 所有行为都是结构体对象转为string类型

		// string转结构体
		str := []byte(bh.ActInfo)
		err = json.Unmarshal([]byte(str), &info)
		// 转换后更改属性
		info.CIndex = newindex
		// 再次转换为string存储
		actinfo, _ := json.Marshal(info)
		bh.ActInfo = string(actinfo)
		/*
			// 记录 以byte流数据记录入账本
			bhAsByte, _ := json.Marshal(bh)
			return ctx.GetStub().PutState(lognum, bhAsByte)
		*/

		return s.WriteLedger(ctx, "LOG05", *bh)

	} else {
		fmt.Errorf("log type wrong")

	}
	return nil
}

func (s *SimpleContract) UpdateCrl(ctx contractapi.TransactionContextInterface, sid string) error {
	// 更新
	bh, err := s.ReadLedger_Sid(ctx, sid, "R")
	var info RInfo
	if err != nil {
		return err
	}
	if bh.ActTag == "R" {
		// 所有行为都是结构体对象转为string类型

		// string转结构体
		str := []byte(bh.ActInfo)
		err = json.Unmarshal([]byte(str), &info)
		// 转换后更改属性
		info.CrlInfo, _ = s.AltS(bh.Sid)
		// 再次转换为string存储
		actinfo, _ := json.Marshal(info)
		bh.ActInfo = string(actinfo)
		/*// 记录 以byte流数据记录入账本
		bhAsByte, _ := json.Marshal(bh)
		return ctx.GetStub().PutState(lognum, bhAsByte)*/
		return s.WriteLedger(ctx, "LOG04", *bh)

	} else {
		fmt.Errorf("log type wrong")

	}
	return nil
}

func (s *SimpleContract) ItokenGen(ctx contractapi.TransactionContextInterface, sid string) (string, error) {

	flag := false
	flag, _ = s.IdCheck(sid)
	if flag { // 首先验证id
		bh, err := s.ReadLedger_Sid(ctx, sid, "R")
		var info RInfo
		if err != nil {
			return "itokengen failed", err
		}
		// 所有行为都是结构体对象转为string类型

		// string转结构体
		str := []byte(bh.ActInfo)
		err = json.Unmarshal([]byte(str), &info)

		if s.SCheck(info.CrlInfo) { // 检查属性
			// 返回itoken
			return " itoken generation succeed!", nil

		}

	}

	return "itokengen failed", nil
}
func (s *SimpleContract) Search(ctx contractapi.TransactionContextInterface, sid string) (string, error) {
	// 首先验证id：能力有限，这里应该传两个sid，检验前者，寻找后者
	flag := false
	flag, _ = s.IdCheck(sid)
	logtag := false
	if flag { // 首先验证id
		bh, err := s.ReadLedger_Sid(ctx, sid, "UP")
		var info UpInfo
		if err != nil {
			return "search failed", err
		}
		// 所有行为都是结构体对象转为string类型

		// string转结构体
		str := []byte(bh.ActInfo)
		err = json.Unmarshal([]byte(str), &info)

		//密文完整性验证
		logtag = s.RCheck(info.Cloc, info.EncrptParam)
		log := ""
		if logtag {
			log = "search succeed "
		} else {
			log = "search failed"
		}
		//先打包
		var actinfo = SNFInfo{
			sid,
			log,
		}

		str, _ = json.Marshal(actinfo)
		newstr := string(str)

		bh = &Behavior{
			Sid:     sid,
			UserCrd: "10",
			ActTag:  "SNF",
			ActInfo: newstr,
		}
		//写账本
		s.WriteLedger(ctx, "LOG03", *bh)

		// 返还密文或者报错
		if logtag {
			return "cipher text at " + info.Cloc + " intact and returned!", err
		} else {
			return "cipher text at " + info.Cloc + " corrupted!", err
		}

	}

	return "search failed", nil

	//search

	//返还ct

}
