package main

import (
	"flag"
	"github.com/heartszhang/builder"

	"log"
)

var (
	FTP_HOST   = flag.String("host", "172.16.12.12", "internal ftp server host ip")
	FTP_USER   = flag.String("user", "admin", "ftp user name")
	FTP_PASSWD = flag.String("passwd", "admin@8", "ftp user password")
	FTP_DIR    = flag.String("ftp-dir", "/windgame", "remote target dir")
	DIR        = flag.String("dir", ".", "dir should be upload, no recursive")
)

func main() {
	flag.Parse()
	ftp, err := builder.NewFtpUploader(*FTP_HOST, *FTP_USER, *FTP_PASSWD)
	if err != nil {
		log.Fatal(err)
	}
	defer ftp.Close()
	err = ftp.UploadDir(*DIR, *FTP_DIR)
	if err != nil {
		log.Fatal(err)
	}
}
