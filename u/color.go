package u

import (
	"fmt"
	"strings"
)

const ResetCode string = "\033[0m"

var Colors = []*Color{
	&Color{"green", "\033[32m"},
	&Color{"yellow", "\033[33m"},
	&Color{"blue", "\033[34m"},
	&Color{"purple", "\033[35m"},
	&Color{"cyan", "\033[36m"},
	&Color{"white", "\033[37m"},
	&Color{"red", "\033[31m"},
}

type Color struct {
	Name string
	Code string
}

func GetColorByName(name string) (*Color, error) {
	lower := strings.ToLower(name)
	for _, color := range Colors {
		if lower == color.Name {
			return color, nil
		}
	}
	return nil, fmt.Errorf("color with name %s not found", name)
}

func ColorText(text, colorCode string) string {
	return colorCode + text + ResetCode
}

func ColorTextByName(text, colorName string) (string, error) {
	color, err := GetColorByName(colorName)
	if err != nil {
		return "", err
	}
	return ColorText(text, color.Code), nil
}
