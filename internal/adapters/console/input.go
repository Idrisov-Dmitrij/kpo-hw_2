package console

import (
	"fmt"
	"strconv"
	"strings"
)

func (c *Console) readLine() (string, bool) {
	if !c.scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(c.scanner.Text()), true
}

func (c *Console) readString(prompt string) (string, bool) {
	fmt.Fprintf(c.out, "%s: ", prompt)
	return c.readLine()
}

func (c *Console) readID(prompt string) (string, bool) {
	value, ok := c.readString(prompt)
	return strings.ToUpper(value), ok
}

func (c *Console) readFloat(prompt string) (float64, bool) {
	for {
		fmt.Fprintf(c.out, "%s: ", prompt)
		line, ok := c.readLine()
		if !ok {
			return 0, false
		}

		value, err := strconv.ParseFloat(strings.ReplaceAll(line, ",", "."), 64)
		if err == nil {
			return value, true
		}
		fmt.Fprintln(c.out, "Ошибка: введите число.")
	}
}

func (c *Console) readInt(prompt string) (int, bool) {
	for {
		fmt.Fprintf(c.out, "%s: ", prompt)
		line, ok := c.readLine()
		if !ok {
			return 0, false
		}

		value, err := strconv.Atoi(line)
		if err == nil {
			return value, true
		}
		fmt.Fprintln(c.out, "Ошибка: введите целое число.")
	}
}
