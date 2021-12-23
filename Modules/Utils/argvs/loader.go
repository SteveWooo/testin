package argvs

import (
	os "os"
)

type Argv struct {
	Items map[string]string
}

func (a *Argv) Load() {
	a.Items = map[string]string{}
	for index, value := range os.Args {
		// 启动文件命令不需要理会
		if index == 0 {
			continue
		}

		if len(value) > 2 && value[0:2] == "--" && os.Args[index+1] != "" {
			key := value[2:]
			val := os.Args[index+1]
			a.Items[key] = val
			continue
		}

		if len(value) > 1 && value[0:1] == "-" && os.Args[index+1] != "" {
			key := value[1:]
			val := os.Args[index+1]
			a.Items[key] = val
			continue
		}
	}
}

func (a *Argv) Get(key string) string {
	return a.Items[key]
}
