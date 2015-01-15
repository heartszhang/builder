package main

import "flag"
import (
	"log"
	"github.com/heartszhang/builder"
)

var (
	URL_LOGIN  = flag.String("login-url", "http://192.168.1.241/", "funshion cdn portal")
	URL_UPLOAD = flag.String("upload-url", "http://192.168.1.241/game/upload/", "funshion upload url")
	USER       = flag.String("user", "WH_mobilegame", "funshion cdn uploaders name")
	PASSWORD   = flag.String("passwd", "funshion", "password")
)

func main() {
	flag.Parse()
	client, err := builder.NewFunshionUploader(*USER, *PASSWORD, *URL_LOGIN)
	if err != nil {
		log.Fatal(err)
	}
	dest, err := client.Upload("./firstgame-test.touch.2", *URL_UPLOAD)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dest)
}
