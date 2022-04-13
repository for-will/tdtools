package GameMsg

import "fmt"

func (x ReturnCode) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\u001b[38;5;202m\"%s\"\u001b[0m", x.String())), nil
}
