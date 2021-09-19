/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"liveVideoToMp4/pkg"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "直播流转储工具",
	Short: "转储ts文件到mp4",
	Long: `转储ts文件到mp4.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		start,_ := cmd.Flags().GetString("start-time")
		duration, _ := cmd.Flags().GetInt("duration")
		m3u8Url, _ := cmd.Flags().GetString("m3u8-url")
		out, _ := cmd.Flags().GetString("out")
		temp, _ := cmd.Flags().GetString("temp")
		mergeTime, _ := cmd.Flags().GetInt("merge-time")
		pkg.StartDownloader(temp, out, mergeTime)
		pkg.StartTask(start, duration, m3u8Url)
		fmt.Println("root called")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// 配置参数
	rootCmd.Flags().StringP("start-time", "s", "08:00", "开始时间")
	rootCmd.Flags().IntP("duration", "d", 9, "持续时间小时")
	rootCmd.Flags().IntP("merge-time", "m", 600, "持续时间小时")
	rootCmd.Flags().StringP("m3u8-url", "u", "", "m3u8的链接地址")
	rootCmd.MarkFlagRequired("m3u8-url")
	rootCmd.Flags().StringP("out", "o", "./out", "mp4的保存路径")
	rootCmd.Flags().StringP("temp", "t", "./temp", "缓存文件路径")
}

