package prepare

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
)

func Do(inPath string, outPath string) error {
	inBytes, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}

	preprocessedInStr := preprocess(string(inBytes))

	md := goldmark.New(
		goldmark.WithExtensions(NewCoditraGoldmark()),
	)

	var outBuf bytes.Buffer
	if err := md.Convert([]byte(preprocessedInStr), &outBuf); err != nil {
		return fmt.Errorf("converting file: %v", err)
	}

	outStr := fmt.Sprintf("<html>\n<body>\n%s</body>\n</html>", outBuf.String())
	err = os.WriteFile(outPath, []byte(outStr), 0644)
	if err != nil {
		return fmt.Errorf("writing output file: %v", err)
	}
	fmt.Printf("Wrote to '%s'", outPath)

	return nil
}

func preprocess(s string) string {
	s = strings.ReplaceAll(s, "...", "…")
	s = strings.ReplaceAll(s, "..", "‥")
	s = strings.ReplaceAll(s, "\u00a0", " ") // Non-breaking space
	s = strings.ReplaceAll(s, "\u00ad", "")  // Soft hyphen
	return s
}
