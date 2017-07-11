package chlib

import (
	"fmt"

	"github.com/kfeofantov/chkit-v2/helpers"
)

type UserInfo struct {
	Token     string `mapconv:"token"`
	Namespace string `mapconv:"namespace"`
}

const userBucket = "user"

func init() {
	defaultInfo := UserInfo{Token: "", Namespace: "default"}
	initializers[userBucket] = helpers.StructToMap(defaultInfo)
}

func GetUserInfo() (info UserInfo, err error) {
	m, err := readFromBucket(userBucket)
	if err != nil {
		return info, fmt.Errorf("user bucket get: %s", err)
	}
	err = helpers.FillStruct(&info, m)
	if err != nil {
		return info, fmt.Errorf("user data fill: %s", err)
	}
	return info, nil
}

func UpdateUserInfo(info UserInfo) error {
	return pushToBucket(userBucket, helpers.StructToMap(info))
}

func UserLogin(login, password string) error {
	info, err := GetUserInfo()
	if err != nil {
		return err
	}
	client, err := NewClient(helpers.CurrentClientVersion, helpers.UuidV4())
	if err != nil {
		return err
	}
	token, err := client.Login(login, password)
	if err != nil {
		return err
	}
	info.Token = token
	err = UpdateUserInfo(info)
	if err != nil {
		return err
	}
	return nil
}

func UserLogout() error {
	info, err := GetUserInfo()
	if err != nil {
		return err
	}
	info.Token = ""
	err = UpdateUserInfo(info)
	return err
}
