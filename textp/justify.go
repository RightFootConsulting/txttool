package textp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func returnLeadingSpace(line string) string {
	output := make([]rune, 0)
	for _, rune := range line {
		if unicode.IsSpace(rune) {
			if rune == '\t' {
				output = append(output, ' ', ' ', ' ')
			}
			output = append(output, rune)
		} else {
			break
		}
	}
	return string(output)
}

func isEmpty(str string) bool {
	temp := strings.TrimSpace(str)
	return temp == ""
}

func JustifyText(input string, output string, linelength int) error {
	if input == output {
		err := fmt.Errorf("input and output files cannot be the same, this is unsafe")
		return err
	}
	inStat, err := os.Stat(input)
	if err != nil {
		return err
	}
	ofMode := inStat.Mode()

	inputFile, err := os.OpenFile(input, os.O_RDONLY, os.ModeExclusive)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, ofMode.Perm())
	if err != nil {
		return err
	}
	defer outputFile.Close()

	inputReader := bufio.NewReader(inputFile)
	lineScanner := bufio.NewScanner(inputReader)
	lineScanner.Split(bufio.ScanLines)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	var outputLine int = 1
	var inputLine int = 0
	var paraBegin bool = true
	var linePrefix string = ""
	var wordPrefix string = ""
	var currentLineLength int = 0
	for lineScanner.Scan() {
		nextLine := lineScanner.Text()
		inputLine++
		if isEmpty(nextLine) {
			// End the paragraph
			//log.Printf("Ending paragraph at input line %d\n", inputLine)
			_, err := writer.WriteString("\n\n")
			if err != nil {
				return err
			}
			outputLine += 2
			currentLineLength = 0
			wordPrefix = ""
			linePrefix = ""
			paraBegin = true
			continue
		}
		leadingSpace := returnLeadingSpace(nextLine)
		if leadingSpace != linePrefix {
			//log.Printf("Line prefix changed from '%s' to '%s' at input line %d", leadingSpace, linePrefix, inputLine)
			linePrefix = leadingSpace
			wordPrefix = ""
			currentLineLength = 0
			n, err := writer.WriteString(fmt.Sprintf("\n%s", linePrefix))
			if err != nil {
				return err
			}
			outputLine++
			currentLineLength += n
		}
		if false && paraBegin {
			n, err := writer.WriteString(linePrefix)
			if err != nil {
				return err
			}
			//log.Printf("New line prefix: '%s' at output line %d", linePrefix, outputLine)
			currentLineLength += n
			paraBegin = false
		}

		wordScanner := bufio.NewScanner(strings.NewReader(nextLine))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			nextWord := wordScanner.Text()
			toWrite := len(wordPrefix) + len(nextWord)
			tentativeNewLength := currentLineLength + toWrite
			if tentativeNewLength > linelength {
				// Do line break in current paragraph respecting indentation
				//log.Printf("At output line %d, current length is %d, next word is '%s%s' with length %d, resulting in tentative length of %d. Doing linebreak.\n",
				//	outputLine, currentLineLength, wordPrefix, nextWord, toWrite, tentativeNewLength)
				wordPrefix = ""
				currentLineLength = 0
				n, err := writer.WriteString(fmt.Sprintf("\n%s", linePrefix))
				if err != nil {
					return err
				}
				outputLine++
				currentLineLength += n
			}
			n, err := writer.WriteString(fmt.Sprintf("%s%s", wordPrefix, nextWord))
			if err != nil {
				return err
			}
			currentLineLength += n
			wordPrefix = " "
			writer.Flush()
		}
		writer.Flush()
	}
	writer.WriteString("\n")
	outputLine++
	wordPrefix = ""
	writer.Flush()
	return nil
}
