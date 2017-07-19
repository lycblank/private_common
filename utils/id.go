package utils

import (
	"github.com/sony/sonyflake"
	"github.com/sony/sonyflake/awsutil"
)

var traceSF *sonyflake.Sonyflake

func InitUniqueID() {
	var st sonyflake.Settings
	st.MachineID = awsutil.AmazonEC2MachineID
	traceSF = sonyflake.NewSonyflake(st)
	if traceSF == nil {
		panic("init unique id panic")
	}
}

func GetTranceID() (uint64, error) {
	return traceSF.NextID()
}
