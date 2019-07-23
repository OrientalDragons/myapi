package myapi

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//DoCmd 执行cmd命令
func DoCmd(inward string) ([]byte, error) {
	cmd := exec.Command("cmd")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in //绑定输入
	var out bytes.Buffer
	cmd.Stdout = &out //绑定正常输出
	cmd.Stderr = &out //绑定错误输出
	inwardByte, _ := EncodeGBK([]byte(inward))
	go func() {
		in.Write(inwardByte)
		// in.WriteString(inward) //写入你的命令，可以有多行，"\n"表示回车
	}()
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	err2 := cmd.Wait()
	if err2 != nil {
		log.Printf("Command finished with error: %v", err2)
	}
	op, e := DecodeGBK(out.Bytes())
	if e != nil {
		return nil, e
	}
	return op, nil
}

//DecodeGBK 将GBK(本机的cmd编码)转化为utf-8方便输出
func DecodeGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//EncodeGBK 将utf-8转为GBK供cmd输入
func EncodeGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//MakeMd5 md5
func MakeMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
