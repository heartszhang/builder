package builder

import "testing"

const (
	FTP_HOST   = "172.16.12.12"
	FTP_USER   = "admin"
	FTP_PASSWD = "admin@8"
)

func TestFtpUpload(t *testing.T) {
	ftp, err := NewFtpUploader(FTP_HOST, FTP_USER, FTP_PASSWD)
	if err != nil {
		t.Fatal(err)
	}
	defer ftp.Close()
	err = ftp.UpdateDir(".", "/windgame")
	if err != nil {
		t.Fatal(err)
	}
}
