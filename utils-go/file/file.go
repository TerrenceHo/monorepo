package file

import (
	"os"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
)

func CreateFileIfNotExists(fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return stackerrors.Wrapf(err, "failed to create file %s", fileName)
		}
		defer file.Close()
	}
	return nil
}
