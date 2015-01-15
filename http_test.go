package builder

import "testing"

const (
	URL_LOGIN  = "http://192.168.1.241/"
	URL_UPLOAD = "http://192.168.1.241/game/upload/"
	USER       = "WH_mobilegame"
	PASSWORD   = "funshion"
)

func TestHttpUpload(t *testing.T) {
	client, err := NewFunshionUploader(USER, PASSWORD, URL_LOGIN)
	if err != nil {
		t.Fatal(err)
	}
	dest, err := client.Upload("./firstgame-test.touch.2", URL_UPLOAD)
	if err != nil {
		t.Fatal(err)
	}
	if dest != "http://neirong.funshion.com/download/funmgame/firstgame-test.touch.2" {
		t.Fatal(dest)
	}
}
