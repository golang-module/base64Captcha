// Copyright 2017 Eric Zhou. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package base64Captcha supports digits, numbers,alphabet, arithmetic, audio and digit-alphabet captcha.
// base64Captcha is used for fast development of RESTful APIs, web apps and backend services in Go. give a string identifier to the package and it returns with a base64-encoding-png-string
package base64Captcha

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var globalStore Store

//CaptchaInterface captcha interface for captcha engine to to write staff
type CaptchaInterface interface {
	//BinaryEncodeing covert to bytes
	BinaryEncodeing() []byte
	//WriteTo output captcha entity
	WriteTo(w io.Writer) (n int64, err error)
}


func Init(){
	var err error
	globalStore, err = NewGlobalStore("redis", storeConfig{expire:60, limitNumber:10240, extraConfig:"127.0.0.1:6379"})
	if err != nil {
		log.Fatalf("captcha Init error : %s\n", err.Error())
	}
}

//CaptchaWriteToBase64Encoding converts captcha to base64 encoding string.
//mimeType is one of "audio/wav" "image/png".
func CaptchaWriteToBase64Encoding(cap CaptchaInterface) string {
	binaryData := cap.BinaryEncodeing()
	var mimeType string
	if _, ok := cap.(*Audio); ok {
		mimeType = MimeTypeCaptchaAudio
	} else {
		mimeType = MimeTypeCaptchaImage
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(binaryData))
}

//CaptchaWriteToFile output captcha to file.
//fileExt is one of "png","wav"
func CaptchaWriteToFile(cap CaptchaInterface, outputDir, fileName, fileExt string) error {
	filePath := filepath.Join(outputDir, fileName+"."+fileExt)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("%s is invalid path.error:%v", filePath, err)
		return err
	}
	defer file.Close()
	_, err = cap.WriteTo(file)
	return err
}

//CaptchaItem captcha basic information.
type CaptchaItem struct {
	//Content captcha entity content.
	Content string
	//VerifyValue captcha verify value.
	VerifyValue string
	//ImageWidth image width pixel.
	ImageWidth int
	//ImageHeight image height pixel.
	ImageHeight int
}

// VerifyCaptcha by given id key, return boolean value.
// 验证图像验证码,返回boolean.
func VerifyCaptcha(identifier, verifyValue string) bool {
	if verifyValue == "" {
		return false
	}
	storeValue := globalStore.Get(identifier, false)
	if storeValue == "" {
		return false
	}
	result := strings.ToLower(storeValue) == strings.ToLower(verifyValue)
	if result {
		globalStore.Get(identifier, true)
	}
	return result
}

//GenerateCaptcha create captcha by config struct and id.
//idkey can be an empty string, base64 will create a unique id four you.
//if idKey is a empty string, the package will generate a random unique identifier for you.
//configuration struct should be one of those struct ConfigAudio, ConfigCharacter, ConfigDigit.
//
//Example Code
//	//config struct for digits
//	var configD = base64Captcha.ConfigDigit{
//		Height:     80,
//		Width:      240,
//		MaxSkew:    0.7,
//		DotCount:   80,
//		CaptchaLen: 5,
//	}
//	//config struct for audio
//	var configA = base64Captcha.ConfigAudio{
//		CaptchaLen: 6,
//		Language:   "zh",
//	}
//	//config struct for Character
//	var configC = base64Captcha.ConfigCharacter{
//		Height:             60,
//		Width:              240,
//		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
//		Mode:               base64Captcha.CaptchaModeNumber,
//		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
//		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
//		IsUseSimpleFont:    true,
//		IsShowHollowLine:   false,
//		IsShowNoiseDot:     false,
//		IsShowNoiseText:    false,
//		IsShowSlimeLine:    false,
//		IsShowSineLine:     false,
//		CaptchaLen:         6,
//	}
//	//create a audio captcha.
//	//GenerateCaptcha first parameter is empty string,so the package will generate a random uuid for you.
//	idKeyA,capA := base64Captcha.GenerateCaptcha("",configA)
//	//write to base64 string.
//	//GenerateCaptcha first parameter is empty string,so the package will generate a random uuid for you.
//	base64stringA := base64Captcha.CaptchaWriteToBase64Encoding(capA)
//	//create a characters captcha.
//	//GenerateCaptcha first parameter is empty string,so the package will generate a random uuid for you.
//	idKeyC,capC := base64Captcha.GenerateCaptcha("",configC)
//	//write to base64 string.
//	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(capC)
//	//create a digits captcha.
//	idKeyD,capD := base64Captcha.GenerateCaptcha("",configD)
//	//write to base64 string.
//	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
func GenerateCaptcha(idKey string, configuration interface{}) (id string, captchaInstance CaptchaInterface) {
	if idKey == "" {
		idKey = randomId()
	}
	id = idKey
	var verifyValue string
	switch config := configuration.(type) {
	case ConfigAudio:
		audio := EngineAudioCreate(idKey, config)
		verifyValue = audio.VerifyValue
		captchaInstance = audio

	case ConfigCharacter:
		char := EngineCharCreate(config)
		verifyValue = char.VerifyValue
		captchaInstance = char

	case ConfigDigit:
		dig := EngineDigitsCreate(idKey, config)
		verifyValue = dig.VerifyValue
		captchaInstance = dig

	default:
		log.Fatal("config type not supported", config)
	}

	globalStore.Set(idKey, verifyValue)

	return idKey, captchaInstance
}
