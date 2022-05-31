package rest

import "fmt"

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func setPort(_port int) string {
	return fmt.Sprintf(":%d", _port)
}
