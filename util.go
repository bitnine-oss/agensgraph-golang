package ag

import "fmt"

func readJSONObject(b []byte) ([]byte, error) {
	if b[0] != byte('{') {
		return nil, fmt.Errorf("invalid JSON object: %s", b)
	}
	depth := 1

	for pos, length := 1, len(b); pos < length; pos++ {
		switch b[pos] {
		case byte('"'):
			// Leave pos unchanged if unmatched " were found.
			escape := false
		InString:
			for i := pos + 1; i < length; i++ {
				switch b[i] {
				case byte('\\'):
					escape = !escape
				case byte('"'):
					if escape {
						escape = false
					} else {
						pos = i
						break InString
					}
				default:
					escape = false
				}
			}
		case byte('{'):
			depth++
		case byte('}'):
			depth--
			if depth == 0 {
				return b[:pos+1], nil
			}
		}
	}

	return nil, fmt.Errorf("invalid JSON object: %s", b)
}
