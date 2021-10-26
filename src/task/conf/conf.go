package conf

import "os"

var Conf *conf

func init() {
	Conf = newLocalConf()

	_ = os.Mkdir(Conf.LogPath, 0755)
}
