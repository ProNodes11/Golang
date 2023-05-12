package wallet

import (
	"io/ioutil"
  "log"
	"fmt"
	"time"
  "net/http"
	"encoding/json"
	"strconv"
)
type AvailableBalance struct {
	Balances []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type DelegationsBalance struct {
	DelegationResponses []struct {
		Delegation struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Shares           string `json:"shares"`
		} `json:"delegation"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
	} `json:"delegation_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type UnboundingBalance struct {
	UnbondingResponses []struct {
		DelegatorAddress string `json:"delegator_address"`
		ValidatorAddress string `json:"validator_address"`
		Entries          []struct {
			CreationHeight string    `json:"creation_height"`
			CompletionTime time.Time `json:"completion_time"`
			InitialBalance string    `json:"initial_balance"`
			Balance        string    `json:"balance"`
		} `json:"entries"`
	} `json:"unbonding_responses"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type RewardBalance struct {
	Rewards []struct {
		ValidatorAddress string `json:"validator_address"`
		Reward           []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"reward"`
	} `json:"rewards"`
	Total []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"total"`
}

type WalletCheck struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}


func GetBalance(address string, api string) (ress string) {
	walletexist := CheckWallet(address, api)
	switch walletexist {
	case false:
		ress = "Wallet not exist"
	case true:
		tokenamount := GetAvailableBalance(address, api)
		tokendelegated :=GetDelegatedBalance(address, api)
		tokenunbounding := GetUnboundingBalance(address, api)
		tokenreward := GetRewardBalance(address, api)
		ress = fmt.Sprintf("🔹Wallet %s \n\n🔹Token amount: \n%s \n🔹Delegeted tokens: \n%s \n🔹Unbounding: \n%s \n🔹Rewards: \n%s", address, tokenamount, tokendelegated, tokenunbounding, tokenreward)
	}
	return
}


func GetAvailableBalance(address string, api string) (ress string) {
	resp, err := http.Get(api + "/cosmos/bank/v1beta1/balances/" + address)
   if err != nil {
      log.Fatalln(err)
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
	 var response AvailableBalance
	 if err := json.Unmarshal([]byte(body), &response); err != nil{
		 log.Printf("Ошибка")
	 }
	 // var fullamount float64
	 if len(response.Balances) == 0 {
		 ress = "0.0"
	 } else {
		 for i := range response.Balances{
			 if amount, err := strconv.ParseFloat(response.Balances[i].Amount, 64); err == nil {
        ress += fmt.Sprintf("\n%.6f %s ", amount / 1000000.0, response.Balances[i].Denom)
				}
		 	}

	 }
	return ress
}
func GetDelegatedBalance(address string, api string) (ress string) {
	resp, err := http.Get(api + "/cosmos/staking/v1beta1/delegations/" + address)
   if err != nil {
      log.Fatalln(err)
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
	 var response DelegationsBalance
	 if err := json.Unmarshal([]byte(body), &response); err != nil{
		 log.Printf("Ошибка")
	 }
	 var fullamount float64
	 if len(response.DelegationResponses) == 0 {
		 ress = "0.0"
		 } else {
			 for i := range response.DelegationResponses{
				if amount, err := strconv.ParseFloat(response.DelegationResponses[i].Balance.Amount, 64); err == nil {
				 fullamount += amount
				 }
			 }
			ress = fmt.Sprintf("\n%.6f %s", fullamount / 1000000.0, response.DelegationResponses[0].Balance.Denom)
	}
	return ress
}

func GetUnboundingBalance(address string, api string) (ress string) {
	resp, err := http.Get(api + "/cosmos/staking/v1beta1/delegators/" + address + "/unbonding_delegations")
   if err != nil {
      log.Fatalln(err)
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
	 var response UnboundingBalance
	 if err := json.Unmarshal([]byte(body), &response); err != nil{
		 log.Printf("Ошибка")
	 }
	 var fullamount float64
	 if len(response.UnbondingResponses) == 0 {
		 ress = "0.0"
		 } else {
			 for i := range response.UnbondingResponses{
				if amount, err := strconv.ParseFloat(response.UnbondingResponses[i].Entries[0].Balance, 64); err == nil {
				 fullamount += amount
				 }
			 }
			ress = fmt.Sprintf("\n%.6f", fullamount / 1000000.0)
	}
	return ress
}

func GetRewardBalance(address string, api string) (ress string) {
	resp, err := http.Get(api + "/cosmos/distribution/v1beta1/delegators/" + address + "/rewards")
   if err != nil {
      log.Fatalln(err)
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
	 var response RewardBalance
	 if err := json.Unmarshal([]byte(body), &response); err != nil{
		 log.Printf("Ошибка")
	 }
	 // var fullamount float64
	 if len(response.Total) == 0 {
		 ress = "0.0"
		 } else {
		  for i := range response.Total{
				if amount, err := strconv.ParseFloat(response.Total[i].Amount, 64); err == nil {
				 ress += fmt.Sprintf("\n%.6f %s", amount / 1000000.0, response.Total[i].Denom)
				 }

			 }
	}
	return ress
}

func CheckWallet(address string, api string) (ress bool) {
	resp, err := http.Get(api + "/cosmos/bank/v1beta1/balances/" + address)
   if err != nil {
      log.Fatalln(err)
   }
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Fatalln(err)
   }
	 var response WalletCheck
	 if err := json.Unmarshal([]byte(body), &response); err != nil{
		 log.Printf("Ошибка")
	 }
	 if response.Code == 3 {
		 ress = false
	 } else {
		 ress = true
	 }

	return ress
}
