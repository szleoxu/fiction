package common

import (
	"github.com/axgle/mahonia"
	"strings"
)

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func GbkToUtf8(src string)string{
	//GBK space is special char,so that i should replace it before convert
	src=strings.Replace(src,"\xC2\xA0"," ",-1)
	return ConvertToString(src, "gbk", "utf-8")
}
