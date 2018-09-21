/*
 * Copyright (c) 2014-2018 Bitnine, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
