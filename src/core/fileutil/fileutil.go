package fileutil

import (
	"fmt"

	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	//EnvGoExeName 执行文件名称系统变量
	EnvGoExeName = "EnvGoExeName"
)

var (
	//DefaultFileMode 默认的文件权限 0640
	DefaultFileMode os.FileMode = 0640
	//DefaultFolderMode 默认的文件夹权限 0750
	DefaultFolderMode os.FileMode = 0750

	//DefaultFileFlag linux下需加上O_WRONLY或是O_RDWR
	DefaultFileFlag int = os.O_APPEND | os.O_CREATE | os.O_WRONLY
)
var (
	DefaultDirs []string
	//ExeName 执行文件名称
	ExeName string
)

func init() {
	wd, _ := os.Getwd()
	fmt.Println("execute path:", wd)
	exeFile := filepath.Clean(os.Args[0])
	parent, exeName := filepath.Split(exeFile)
	names := strings.Split(exeName, ".")
	exeName = names[0]
	//1：命令执行所在目录
	wdDir := TransPath(wd)
	DefaultDirs = append(DefaultDirs, wdDir)
	exeDir := TransPath(parent)
	if wdDir != exeDir { //2：可执行文件所在目录
		DefaultDirs = append(DefaultDirs, exeDir)
	}
	exeLastDir := TransPath(filepath.Join(exeDir, ".."))
	if wdDir != exeLastDir { //3：可执行文件上一级目录
		DefaultDirs = append(DefaultDirs, exeLastDir)
	}
	fmt.Println("default file base dirs:", DefaultDirs)
	SetOSEnv(EnvGoExeName, exeName)
	ExeName = exeName
}

//SetOSEnv 设置环境变量
func SetOSEnv(option, value string) {
	if old, ok := os.LookupEnv(option); ok {
		os.Setenv(option, value)
		fmt.Printf("update sys env [%s=%s] ==>[%s=%s]\n", option, old, option, value)
	} else {
		os.Setenv(option, value)
		fmt.Printf("set sys env [%s=%s]\n", option, value)
	}
}

//InitPathDirs 初始化工程文件根目录 retset：是否重置默认的目录列表 。 rootPaths：新增的目录列表，如果没有设置指定根目录并未重置默认的目录的话，文件搜索规则为: 1：命令执行所在目录、2：可执行文件所在目录、3：可执行文件上一级目录...新增的路径列表，否则直接按照给的文件名称查找文件。
func InitPathDirs(retset bool, rootPaths ...string) {
	if retset {
		DefaultDirs = DefaultDirs[:0]
	}
	for _, rootPath := range rootPaths {
		if len(rootPath) > 0 {
			rootPath = TransPath(rootPath)
			rootPath, _ = filepath.Abs(rootPath)
			rootPath = TransPath(rootPath)
		}
		if rootPath != "" {
			found := false
			for _, defaultDir := range DefaultDirs {
				if defaultDir == rootPath {
					found = true
					break
				}
			}
			if !found {
				DefaultDirs = append(DefaultDirs, rootPath)
			}
		}
	}
}

//FindFile 查找文件，根据初始化的文件目录顺序查找文件
func FindFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	for _, dir := range DefaultDirs {
		fpath := PathJoin(dir, name)
		if FileExists(fpath) {
			f, err := os.OpenFile(fpath, flag, perm)
			if err != nil {
				return nil, fmt.Errorf("open file err,path:%s,err:%v", fpath, err)
			}
			return f, nil
		}
	}
	fpath := TransPath(name)
	if FileExists(fpath) {
		f, err := os.OpenFile(fpath, flag, perm)
		if err != nil {
			return nil, fmt.Errorf("open file err,path:%s,err:%v", fpath, err)
		}
		return f, nil
	}
	return nil, fmt.Errorf("file no found,file:%s,dirs:%v", fpath, DefaultDirs)
}

//FindFullFilePath 查找相对目录文件的全路径文件 根据初始化的文件目录顺序查找文件（查文件不是查目录）
func FindFullFilePath(name string) (string, error) {
	for _, dir := range DefaultDirs {
		fpath := PathJoin(dir, name)
		if FileExists(fpath) {
			return fpath, nil
		}
	}
	fpath := TransPath(name)
	if FileExists(fpath) {
		return fpath, nil
	}
	return "", fmt.Errorf("file no found,file:%s,dirs:%v", fpath, DefaultDirs)
}

//FindFullPathPath 查找相对文件目录的全路径目录 根据初始化的文件目录顺序查找文件（查目录不是查文件）
func FindFullPathPath(name string) (string, error) {
	for _, dir := range DefaultDirs {
		fpath := PathJoin(dir, name)
		if DirExists(fpath) {
			return fpath, nil
		}
	}
	fpath := TransPath(name)
	if DirExists(fpath) {
		return fpath, nil
	}
	return "", fmt.Errorf("file no found,path:%s,dirs:%v", fpath, DefaultDirs)
}

//OpenFile 打开文件，如果目录文件不存在则创建一个文件
func OpenFile(pathfile string, fileflag int, filemode os.FileMode) (wc *os.File, err error) {
	pathfile = strings.Replace(filepath.Clean(pathfile), "\\", "/", -1)
	dir := path.Dir(pathfile)
	if _, err = os.Stat(dir); err != nil && !os.IsExist(err) {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.MkdirAll(dir, DefaultFolderMode); err != nil {
			return nil, err
		}
		if _, err = os.Stat(dir); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(pathfile, fileflag, filemode)
}

//DirExists 指定目录是否存在
func DirExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}
	return true
}

//FileExists 指定文件是否存在
func FileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

//ChangeFileExt 更改文件后缀名 eg:filename=`../a/test/aa.txt` newExt=`.csv` -->return=`../a/test/aa.csv`
func ChangeFileExt(filename, newExt string) string {
	filename = strings.Replace(filepath.Clean(filename), "\\", "/", -1)
	file := path.Base(filename)
	file = strings.TrimSuffix(file, path.Ext(file)) + newExt
	dir := path.Dir(filename)
	return PathJoin(dir, file)
}

//PathJoin 路径合并 并将 “\\” 转换成 “/”
func PathJoin(dir, filename string) string {
	return strings.Replace(filepath.Join(filepath.Clean(dir), filename), "\\", "/", -1)
}

//TransPath 路径连接符转换 将路径 “\\” 转换成 “/”
func TransPath(path string) string {
	return strings.Replace(filepath.Clean(path), "\\", "/", -1)
}
