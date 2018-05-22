// Copyright 2018, all rights reserved.
//-------------------------------------

/*
file包提供文件操作工具类
 */
package file

import (
	"os"
)

//判断文件/文件夹是否存在
//在获取到os.Stat的err后，不可以使用os.IsExist判断
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
