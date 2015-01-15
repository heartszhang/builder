// builder project builder.go
package builder

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
)

type FunshionUploader interface {
	Upload(filepath, upurl string) (string, error)
}

type fs_uploader struct {
	http.Client
}

// formnames are magic
// username
// password
func NewFunshionUploader(user, password, login string) (FunshionUploader, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	v := &fs_uploader{http.Client{Jar: jar}}
	form := url.Values{"username": {user}, "password": {password}}
	resp, err := v.PostForm(login, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, os.ErrInvalid
	}
	return v, nil
}

// http://neirong.funshion.com/download/funmgame/firstgame-test.touch
// i dont known what does 'channel => csol' mean, but it's necessory
// file formname must be 'video'
func (client *fs_uploader) Upload(filepath, upurl string) (dest string, err error) {
	body_start := &bytes.Buffer{}
	mwriter := multipart.NewWriter(body_start)
	mwriter.WriteField("channel", "csol")
	if _, err = mwriter.CreateFormFile("video", filepath); err != nil {
		return
	}

	fi, err := os.Stat(filepath)
	if err != nil {
		return
	}
	fb, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer fb.Close()

	boundary := mwriter.Boundary()
	body_end := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	reqreader := io.MultiReader(body_start, fb, body_end)
	req, err := http.NewRequest("POST", upurl, reqreader)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_start.Len()) + int64(body_end.Len())
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = os.ErrInvalid
		return
	}

	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	re := regexp.MustCompile("http://[a-zA-Z0-9_%.\\-/]+")
	dest = re.FindString(string(cont))
	return
}
