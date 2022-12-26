package utils

import (
	"fmt"
	"time"
)

func Spin(s *string, completed string, done chan bool) {
	ticker := time.NewTicker(100 * time.Millisecond)
	chars := []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	count := 0
	go func() {
		for {
			select {
			case <-done:
				fmt.Printf("\033[2K\r%s\n", completed)
				ticker.Stop()
				return
			case <-ticker.C:
				fmt.Printf("\033[2K\r%s %s", chars[count], *s)
				count++
				count = count % len(chars)
			}
		}
	}()

}
