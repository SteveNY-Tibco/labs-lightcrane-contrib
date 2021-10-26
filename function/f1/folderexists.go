package f1

import (
	"os"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression/function"
)

func init() {
	function.Register(&fnFolderExists{})
}

type fnFolderExists struct {
}

func (fnFolderExists) Name() string {
	return "folderexists"
}

func (fnFolderExists) Sig() (paramTypes []data.Type, isVariadic bool) {
	return []data.Type{data.TypeString}, false
}

func (fnFolderExists) Eval(params ...interface{}) (interface{}, error) {
	log.Debug("FolderExists.eval] entering ..... ")
	defer log.Debug("FolderExists.eval] exit ..... ")

	log.Debug("FolderExists.eval] folder name = ", params[0])

	exist := true
	folderInfo, err := os.Stat("temp")
	if os.IsNotExist(err) {
		exist = false
	}

	log.Debug("FolderExists.eval] folder info = ", folderInfo)

	return exist, nil
}
