package builder

import (
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	ftp "code.google.com/p/ftp4go"
)

// ftp FTP.Cwd not work for vFtpd, I don't know why
// so, ftp4go' src is modified

type ftp_uploader struct {
	*ftp.FTP
}

type FtpUploadCloser interface {
	io.Closer
	UpdateDir(localfolder, remotefolder string) error
	UploadDir(localfolder, remotefolder string) error
}

func NewFtpUploader(host, user, password string) (FtpUploadCloser, error) {
	xtp := &ftp_uploader{ftp.NewFTP(0)}

	_, err := xtp.Connect(host, ftp.DefaultFtpPort, "")
	if err != nil {
		return nil, err
	}
	_, err = xtp.Login(user, password, "")
	if err != nil {
		xtp.Quit()
		return nil, err
	}
	return xtp, err
}

func (xtp *ftp_uploader) Close() error {
	_, err := xtp.Quit()
	return err
}

func (xtp *ftp_uploader) UpdateDir(localdir, remotedir string) (err error) {
	cwd, err := xtp.Pwd()
	defer xtp.Cwd(cwd)
	if _, err = xtp.Cwd(remotedir); err != nil {
		return err
	}
	remote_files, err := xtp.dir(remotedir)
	local_files, err := local_dir(localdir)
	// upload all local files except the remote one be same size with local
	for file, v := range local_files {
		if x, ok := remote_files[file]; !ok || x != v {
			if err = xtp.upload(localdir, file); err != nil {
				log.Println(file, err)
			}
			log.Println(file, v)
		}
	}
	return nil
}

func (xtp *ftp_uploader) UploadDir(localdir, remotedir string) (err error) {
	_, err = xtp.UploadDirTree(localdir, remotedir, 4, nil, dummy)
	return
}

//perm inode owner group filesize date filename
func (xtp *ftp_uploader) dir(dir string) (map[string]int64, error) {
	ls, err := xtp.Dir()
	if err != nil {
		return nil, err
	}
	v := make(map[string]int64)
	for _, line := range ls {
		fields := strings.Fields(line)
		if len(fields) == 9 { // total 9 fields
			//field 4 is filesize
			if sz, err := strconv.ParseInt(fields[4], 0, 64); err == nil {
				v[fields[8]] = sz
			}
		}
	}
	return v, nil
}

func local_dir(folder string) (map[string]int64, error) {
	fi, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	v := make(map[string]int64)
	for _, f := range fi {
		if !f.IsDir() {
			v[f.Name()] = f.Size()
		}
	}
	return v, nil
}

func (xtp *ftp_uploader) upload(dir, name string) error {
	return xtp.UploadFile(name, filepath.Join(dir, name), false, dummy)
}

func dummy(ci *ftp.CallbackInfo) {
	if ci.Eof {
		log.Println(ci.Filename, ci.Resourcename, ci.BytesTransmitted)
	}
}
