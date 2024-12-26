package benchmark

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"gopkg.in/yaml.v3"
)

// Defining root element for XML
type TestData struct {
	XMLName xml.Name `xml:"data"`
	Name    string   `json:"name" xml:"name" yaml:"name"`
	Age     int      `json:"age" xml:"age" yaml:"age"`
	Hobbies []string `json:"hobbies" xml:"hobby" yaml:"hobbies"`
}

var testData = TestData{
	Name:    "Test User",
	Age:     30,
	Hobbies: []string{"okuma", "yazma", "kodlama"},
}

// XML structure for large data set
type LargeTestData struct {
	XMLName xml.Name `xml:"data"`
	Items   []Item   `json:"items" xml:"items>item" yaml:"items"`
}

type Item struct {
	ID          int      `json:"id" xml:"id" yaml:"id"`
	Name        string   `json:"name" xml:"name" yaml:"name"`
	Description string   `json:"description" xml:"description" yaml:"description"`
	Tags        []string `json:"tags" xml:"tag" yaml:"tags"`
}

var largeTestData = LargeTestData{
	Items: make([]Item, 100),
}

var (
	jsonData []byte
	xmlData  []byte
	yamlData []byte
)

func init() {
	var err error
	jsonData, err = json.Marshal(testData)
	if err != nil {
		panic(err)
	}
	xmlData, err = xml.Marshal(testData)
	if err != nil {
		panic(err)
	}
	yamlData, err = yaml.Marshal(testData)
	if err != nil {
		panic(err)
	}

	// Büyük veri setini doldur
	for i := 0; i < len(largeTestData.Items); i++ {
		largeTestData.Items[i] = Item{
			ID:          i,
			Name:        "Item " + string(rune(i)),
			Description: "Bu bir test açıklamasıdır " + string(rune(i)),
			Tags:        []string{"tag1", "tag2", "tag3"},
		}
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkXMLMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := xml.Marshal(testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkYAMLMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := yaml.Marshal(testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	var result struct {
		Name    string   `json:"name"`
		Age     int      `json:"age"`
		Hobbies []string `json:"hobbies"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(jsonData, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkXMLUnmarshal(b *testing.B) {
	var result struct {
		Name    string   `xml:"name"`
		Age     int      `xml:"age"`
		Hobbies []string `xml:"hobbies"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := xml.Unmarshal(xmlData, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkYAMLUnmarshal(b *testing.B) {
	var result struct {
		Name    string   `yaml:"name"`
		Age     int      `yaml:"age"`
		Hobbies []string `yaml:"hobbies"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := yaml.Unmarshal(yamlData, &result); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(largeTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeXMLMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := xml.Marshal(largeTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeYAMLMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := yaml.Marshal(largeTestData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
