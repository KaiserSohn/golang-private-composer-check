package html_tools

import (
	"PrivateComposerChecker/composer"
	"PrivateComposerChecker/http_tools"
	"github.com/pejman-hkh/gdp/gdp"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const HtmlPath = "/service/rest/repository/browse/composer-proxy-packagist/"

type AvailableVersion struct {
	Name    string
	Version string
}

func ParseVersions(auth composer.Auth, lib composer.Lib) []AvailableVersion {
	parsedUrl := http_tools.Parse(lib.Dist["url"])
	url := parsedUrl.Scheme + "://" + parsedUrl.Host + HtmlPath + lib.Name

	document := gdp.Default(http_tools.GetHtml(auth, url))

	var versionsMap = make(map[int]AvailableVersion)
	var versionsIntMap []int

	document.Find("table").Eq(0).Find("tr").Each(func(i int, tag *gdp.Tag) {
		innerHtml := tag.Find("a").Eq(0).Html()

		versionInt := convertVersionToInteger(innerHtml)

		if versionInt == 0 {
			return
		}

		versionsMap[versionInt] = AvailableVersion{
			Name:    lib.Name,
			Version: innerHtml,
		}
		versionsIntMap = append(versionsIntMap, versionInt)
	})

	nearestVersionsInt := nearestHeap(versionsIntMap, convertVersionToInteger(lib.Version), 4)

	var result []AvailableVersion
	for _, nearestVersionInt := range nearestVersionsInt {
		result = append(result, versionsMap[nearestVersionInt])
	}

	return result
}

func convertVersionToInteger(requestedVersion string) int {
	reg, _ := regexp.Compile("[^\\d.]")
	requestedVersion = reg.ReplaceAllString(requestedVersion, "")

	parts := strings.Split(requestedVersion, ".")

	result := 0
	for i, value := range parts {
		intValue, _ := strconv.Atoi(value)
		if i == 0 {
			result += intValue * 10000000
		} else if i == 1 {
			result += intValue * 100
		} else {
			result += intValue * 10
		}
	}
	return result
}

// Нейросетевое хрючево
func absInt(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func nearestHeap(nums []int, target int, size int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	type item struct {
		value int
		dist  int
	}

	heap := make([]item, 0, size)

	for _, v := range nums {
		dist := absInt(v, target)
		heap = append(heap, item{value: v, dist: dist})

		if len(heap) > size {
			maxIdx := 0
			for i := 1; i < len(heap); i++ {
				if heap[i].dist > heap[maxIdx].dist ||
					(heap[i].dist == heap[maxIdx].dist && heap[i].value > heap[maxIdx].value) {
					maxIdx = i
				}
			}
			heap[maxIdx] = heap[len(heap)-1]
			heap = heap[:len(heap)-1]
		}
	}

	sort.Slice(heap, func(i, j int) bool {
		if heap[i].dist == heap[j].dist {
			return heap[i].value < heap[j].value
		}
		return heap[i].dist < heap[j].dist
	})

	result := make([]int, len(heap))
	for i, it := range heap {
		result[i] = it.value
	}
	return result
}
