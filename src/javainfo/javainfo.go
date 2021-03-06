package javainfo

import (
	"bitbucket.org/reckhou/DoomAnalysis/src/dbinfo"
	"bitbucket.org/reckhou/DoomAnalysis/src/file"
	"crypto/md5"
	"encoding/hex"
	//"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

// 上传的JAVA文件行
var key_arr_java [6]string = [...]string{"JAVA", "UUID", "device", "version", "product_name", "file"}

type JavaFileInfo struct {
	info_      map[string]string
	file_name_ string
	project    string
	lianyun    string
}

func (info *JavaFileInfo) SetProjectInfo(pro string, lianyun string) {
	info.project = pro
	info.lianyun = lianyun
}

func (info *JavaFileInfo) GenJavaInfo(s string) {
	info.info_ = make(map[string]string)
	line_num := 0
	start_index := 0
	context_start_index := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if start_index > 0 {
				start_index++
			}
			context_start_index = start_index + len(key_arr_java[line_num]) + 1

			if line_num == 5 {
				info.info_[key_arr_java[line_num]] = s[start_index:]
			} else {
				info.info_[key_arr_java[line_num]] = s[context_start_index:i]
			}

			start_index = i
			line_num++
			if line_num >= len(key_arr_java) {
				break
			}
		}
	}

	if info.info_["file"] == "" {
		info.info_["file"] = s[start_index+len(key_arr_java[5])+2:]
	}

	path := "./" + info.project + "/dump/" + info.info_["version"]
	file.CreateDir(path)
	info.file_name_ = info.info_["UUID"] + ".txt"
	file.WriteFile(path+"/"+info.file_name_, []byte(info.info_["file"]), os.O_TRUNC)
}

func (info *JavaFileInfo) GenJavaDBInfo() {

	s := file.ReadFile("./" + info.project + "/dump/" + info.info_["version"] + "/" + info.file_name_)
	context := string(s)

	start_pos := 0
	info_key := ""
	key_index := 0
	for i := 0; i < len(context); i++ {
		if context[i] == '\n' {
			if i-start_pos > 1 {
				temp_str := string(context[start_pos:i])

				re := regexp.MustCompile("\\([^\\(\\)]*\\.java:[0-9]{0,10}\\)")
				matched := re.FindString(temp_str)
				if matched != "" {
					//fmt.Println("matched : ", matched)
					key_index++
					info_key += matched
					if key_index >= 3 {
						break
					}
				}
			}
			start_pos = i + 1

		}
	}

	h := md5.New()
	h.Write([]byte(info_key))
	result_str := hex.EncodeToString(h.Sum(nil))

	mysql_c, db_err := dbinfo.Init()
	if db_err == nil {
		mysql_c.AddInfo(info.project, info.info_["version"], result_str, context, info.info_["UUID"], info.lianyun)
		mysql_c.AddDeviceInfo(info.project, info.info_["version"], info_key, info.info_["device"], info.lianyun, info.info_["UUID"])
	}
}

func (info *JavaFileInfo) GenTar(mode string) {
	// info.info_["UUID"]
	cmd := exec.Command("/bin/sh", "gen_tar.sh", info.info_["version"], info.project, info.info_["UUID"], mode)
	_, err := cmd.Output()
	if err != nil {
		log.Println("GenSym err:" + err.Error())
	}
}
