package compress

import "testing"

func TestName(t *testing.T) {
	//tarfile := "/data/work/Code/mincli/1594361096-yolov3.tar"
	tarfile := "/data/work/Code/mincli/mnt.tar"
	dest := "/data/work/Code/mincli/test"
	err := UnTar(tarfile, dest)
	t.Error(err)
}
