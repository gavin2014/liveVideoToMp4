package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/oopsguy/m3u8/parse"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func StartTask (start string, duration int, m3u8Url string) error {
	startArr := strings.Split(start, ":")
	if len(startArr)<2 {
		return errors.New("参数不足")
	}
	startHour, _ := strconv.Atoi(startArr[0])
	startMinute, _ := strconv.Atoi(startArr[1])

	bigResultMap := make(map[string]bool, 0)
	urlArr := strings.Split(m3u8Url, "/")
	urlPre := strings.Join(urlArr[0:len(urlArr)-1], "/")
	for {
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMinute,0,0, now.Location())
		end := start.Add(time.Hour * time.Duration(duration))
		//fmt.Println("now" + now.String())
		//fmt.Println("start" + start.String())
		//fmt.Println("end" + end.String())
		//fmt.Println(now.After(start) && now.Before(end))

		if now.After(start) && now.Before(end) {
			fmt.Println(time.Now().String()+": 执行任务开始")
			// 获取 m3u8
			re, err := parse.FromURL(m3u8Url)
			if err!=nil {
				fmt.Println("m3u8 下载 err:"+err.Error())
				time.Sleep(time.Second*30)
			}else{
				for _, segmentItem := range re.M3u8.Segments {

					if _,ok := bigResultMap[segmentItem.URI]; !ok {
						bigResultMap[segmentItem.URI] = true
						TsFiles <- urlPre+"/"+segmentItem.URI
					}
				}
				time.Sleep(time.Second*time.Duration(int64(re.M3u8.TargetDuration)))
			}
		}else{
			bigResultMap = make(map[string]bool, 0)
			time.Sleep(time.Second*60)
		}

	}
	return nil
}


var TsFiles chan string

func StartDownloader(tempPath, out string, mergeTime int) {
	_, err := os.Stat(out)
	if os.IsNotExist(err) {
		os.Mkdir(out, 0777)
	}
	// ts 目录是否存在
	_, err = os.Stat(tempPath)
	if os.IsNotExist(err) {
		os.Mkdir(tempPath, 0777)
	}
	fmt.Println("启动downloer")
	TsFiles = make(chan string)
	tempFiles := make(TempTsFiles, 0)
	// 获取缓存文件
	fileInfoList, err := ioutil.ReadDir(tempPath)
	if err != nil {
		fmt.Println("缓存文件读取失败：" + err.Error())
	}
	for _, file := range fileInfoList {
		if !file.IsDir() {
			tempFiles = append(tempFiles, TsFile(tempPath+"/"+file.Name()))
			// fmt.Println(file.Name()) //打印当前文件或目录下的文件或目录名
		}
	}
	timeOut := time.After(time.Second*time.Duration(mergeTime))
	go func(){
		for  {
			select {
				case tsFile :=  <- TsFiles :
					fmt.Println("下载" + tsFile)
					// continue
					// 下载文件
					resp, err := http.Get(tsFile)
					if err != nil {
						if err != nil {
							fmt.Println("下载文件 err:"+err.Error())
						}
					}
					fileBytes,err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("读取文件 err:"+err.Error())
					}
					tsUrlArr := strings.Split(tsFile, "/")
					err = os.WriteFile(filepath.Join(tempPath, tsUrlArr[len(tsUrlArr)-1]), fileBytes, 0755)
					if err != nil {
						fmt.Println("读取文件 err:"+err.Error())
					}
					tempFiles = append(tempFiles, TsFile(filepath.Join(tempPath, tsUrlArr[len(tsUrlArr)-1])))
				case <-timeOut:
					if tempFiles.Len()>0 {
						// 合并
						fmt.Println("合并文件")
						sort.Sort(tempFiles)
						// 合并完格式转换
						MakeMp4(MergeTsFiles(tempFiles, out), out)
						// 清理
						tempFiles = make(TempTsFiles, 0)
					}
					timeOut = time.After(time.Second*time.Duration(mergeTime))
			}

		}
	}()
}

func MergeTsFiles(tsFiles TempTsFiles, out string) string{
	// Create a TS file for merging, all segment files will be written to this file.
	mergeTSFilename := time.Now().Unix()
	mFilePath := filepath.Join(out, strconv.FormatInt(mergeTSFilename, 10)+".ts")
	mFile, err := os.Create(mFilePath)
	if err != nil {
		fmt.Println("create main TS file failed：%s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer mFile.Close()
	writer := bufio.NewWriter(mFile)
	for _, tsFile := range tsFiles {
		bytes, err := ioutil.ReadFile(string(tsFile))
		_, err = writer.Write(bytes)
		if err != nil {
			continue
		}
	}
	_ = writer.Flush()
	for _, tsFile := range tsFiles {
		// Remove `ts` folder
		_ = os.Remove(string(tsFile))
	}
	return mFilePath
}

// ffmpeg -i "concat:segment-1-v1-a1.ts|segment-2-v1-a1.ts" -acodec copy -vcodec copy -absf aac_adtstoasc output.mp4
func MakeMp4(str string, out string) {
	//binary, lookErr := exec.LookPath("ffmpeg")
	//if lookErr != nil {
	//	panic(lookErr)
	//}
	no := time.Now().Format("2006年01月02日15时04分05")+".mp4"
	args := []string{
		"-i",
		fmt.Sprintf("concat:%s", str),
		"-acodec",
		"copy",
		"-vcodec",
		"copy",
		"-absf",
		"aac_adtstoasc",
		filepath.Join(out, no),
	}
	cmd := exec.Command("./ffmpeg.exe", args...)
	r, err := cmd.CombinedOutput()
	if err!=nil {
		fmt.Println("生成mp4错误"+err.Error())
		fmt.Println(string(r))
	}else {
		fmt.Println("删除"+str)
		os.Remove(str)
	}

}


type TsFile string

type TempTsFiles []TsFile
// 获取此 slice 的长度
func (p TempTsFiles) Len() int { return len(p) }
// 根据元素的年龄降序排序 （此处按照自己的业务逻辑写）
func (p TempTsFiles) Less(i, j int) bool {
	nameIArr := strings.Split(strings.TrimSuffix(strings.ToLower(string(p[i])), ".ts"), "-")
	nameI,_ := strconv.Atoi(nameIArr[1])
	nameJArr := strings.Split(strings.TrimSuffix(strings.ToLower(string(p[j])), ".ts"), "-")
	nameJ,_ := strconv.Atoi(nameJArr[1])
	return  nameI < nameJ
}
// 交换数据
func (p TempTsFiles) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
