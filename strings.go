package utils

import "strings"

func trimSpace(in string) string {
	in = strings.TrimSpace(in)
	return in
}

func trimNewLine(in string) string {
	// for windows 换行符
	for {
		out := strings.TrimPrefix(in, "\r\n")
		if out == in {
			break
		} else {
			in = out
		}
	}
	for {
		out := strings.TrimSuffix(in, "\r\n")
		if out == in {
			break
		} else {
			in = out
		}
	}

	// for linux 换行符
	for {
		out := strings.TrimPrefix(in, "\n")
		if out == in {
			break
		} else {
			in = out
		}
	}
	for {
		out := strings.TrimSuffix(in, "\n")
		if out == in {
			break
		} else {
			in = out
		}
	}

	return in
}

// TrimSpaceNewLine 去掉首尾所有的空格和换行
func TrimSpaceNewLine(in string) string {
	changed := false
	for {
		out := trimSpace(in)
		if out != in {
			changed = true
		}
		out = trimNewLine(in)
		if out != in {
			changed = true
		}

		if changed == false {
			return out
		}
	}
}
