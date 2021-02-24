package functions

import (
	"encoding/base64"
	//"fmt"
)

/* func main() {
	basicAuthB64("flipboard", "Dv3dSTXKhZA7")
} */

func basicAuthB64(u string, p string) (string, error) {
	msg := u + ":" + p
	authHeader := base64.StdEncoding.EncodeToString([]byte(msg))
	//fmt.Println("User Name: " + u)
	//fmt.Println("User Passwd: " + p)
	//fmt.Println("Basic " + authHeader)
	return "Basic " + authHeader, nil
}
