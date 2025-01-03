# Transformer

Transformer, farklı veri formatları (JSON, XML, YAML) arasında dönüşüm yapmanızı sağlayan bir Go kütüphanesidir. Ortak bir veri yapısı kullanarak, bu formatlar arasında kayıpsız dönüşüm yapabilirsiniz.

[English Documentation (İngilizce Dokümantasyon)](README.md)

## Özellikler

- JSON, XML ve YAML formatları arasında dönüşüm
- Ortak veri yapısı ile tutarlı dönüşüm
- Kolay kullanım
- Tip güvenliği
- Özelleştirilebilir dönüşüm kuralları
- Yüksek test kapsamı
- İş parçacığı güvenli operasyonlar
- Minimum bağımlılıklar

## Gereksinimler

- Go 1.16 veya üzeri
- Bağımlılıklar:
  - YAML işlemleri için `gopkg.in/yaml.v3`

## Kurulum

```bash
go get github.com/mstgnz/transformer
```

## Kullanım

### JSON Dönüşümleri

```go
import "github.com/mstgnz/transformer/tjson"

// JSON dosyasını okuma
data, err := tjson.ReadJson("data.json")
if err != nil {
    log.Fatal(err)
}

// JSON formatını doğrulama
if !tjson.IsJson(data) {
    log.Fatal("Geçersiz JSON formatı")
}

// JSON'ı Node yapısına dönüştürme
node, err := tjson.DecodeJson(data)
if err != nil {
    log.Fatal(err)
}

// Node yapısını JSON'a dönüştürme
jsonData, err := tjson.NodeToJson(node)
if err != nil {
    log.Fatal(err)
}
```

### XML Dönüşümleri

```go
import "github.com/mstgnz/transformer/txml"

// XML dosyasını okuma
data, err := txml.ReadXml("data.xml")
if err != nil {
    log.Fatal(err)
}

// XML formatını doğrulama
if !txml.IsXml(data) {
    log.Fatal("Geçersiz XML formatı")
}

// XML'i Node yapısına dönüştürme
node, err := txml.DecodeXml(data)
if err != nil {
    log.Fatal(err)
}

// Node yapısını XML'e dönüştürme
xmlData, err := txml.NodeToXml(node)
if err != nil {
    log.Fatal(err)
}
```

### YAML Dönüşümleri

```go
import "github.com/mstgnz/transformer/tyaml"

// YAML dosyasını okuma
data, err := tyaml.ReadYaml("data.yaml")
if err != nil {
    log.Fatal(err)
}

// YAML formatını doğrulama
if !tyaml.IsYaml(data) {
    log.Fatal("Geçersiz YAML formatı")
}

// YAML'ı Node yapısına dönüştürme
node, err := tyaml.DecodeYaml(data)
if err != nil {
    log.Fatal(err)
}

// Node yapısını YAML'a dönüştürme
yamlData, err := tyaml.NodeToYaml(node)
if err != nil {
    log.Fatal(err)
}
```

### Formatlar Arası Dönüşüm

```go
// JSON -> XML dönüşümü
jsonData := []byte(`{"name": "John", "age": 30}`)
node, _ := tjson.DecodeJson(jsonData)
xmlData, _ := txml.NodeToXml(node)

// XML -> YAML dönüşümü
xmlData := []byte(`<root><name>John</name><age>30</age></root>`)
node, _ := txml.DecodeXml(xmlData)
yamlData, _ := tyaml.NodeToYaml(node)

// YAML -> JSON dönüşümü
yamlData := []byte("name: John\nage: 30")
node, _ := tyaml.DecodeYaml(yamlData)
jsonData, _ := tjson.NodeToJson(node)
```

## Paket Yapısı

- `node`: Temel veri yapısını ve operasyonlarını içerir
  - Hiyerarşik veriyi temsil eden Node yapısı
  - Değer tipleri ve tip dönüşüm operasyonları
  - Ağaç gezinme ve manipülasyon fonksiyonları
- `tjson`: JSON dönüşüm işlemlerini gerçekleştirir
  - JSON kodlama/kod çözme
  - JSON doğrulama
  - JSON dosya işlemleri
- `txml`: XML dönüşüm işlemlerini gerçekleştirir
  - XML kodlama/kod çözme
  - XML doğrulama
  - XML dosya işlemleri
  - XML öznitelik işleme
- `tyaml`: YAML dönüşüm işlemlerini gerçekleştirir
  - YAML kodlama/kod çözme
  - YAML doğrulama
  - YAML dosya işlemleri
- `example`: Örnek kullanımları içerir
  - Temel dönüşüm örnekleri
  - Karmaşık veri yapısı örnekleri
  - Hata işleme örnekleri

## Veri Tipleri

Node yapısı aşağıdaki veri tiplerini destekler:

- `TypeNull`: Boş değer
- `TypeObject`: Nesne (anahtar-değer çiftleri)
  - İç içe nesneleri destekler
  - Anahtar sıralamasını korur
  - Döngüsel referansları işler
- `TypeArray`: Dizi
  - Karışık tipleri destekler
  - Sıralamayı korur
- `TypeString`: Metin
- `TypeNumber`: Sayı (tam sayılar ve ondalıklı sayılar)
- `TypeBoolean`: Mantıksal değer

## Hata İşleme

Kütüphane, çeşitli senaryolar için detaylı hata bilgisi sağlar:

- Dosya işlemi hataları
- Format doğrulama hataları
- Dönüşüm hataları
- Tip uyuşmazlığı hataları
- Yapı doğrulama hataları

Örnek hata işleme:

```go
if err := validateAndConvert(); err != nil {
    switch e := err.(type) {
    case *FormatError:
        log.Printf("Geçersiz format: %v", e)
    case *ConversionError:
        log.Printf("Dönüşüm başarısız: %v", e)
    default:
        log.Printf("Beklenmeyen hata: %v", e)
    }
}
```

## Test

Kütüphane kapsamlı test kapsamına sahiptir. Testleri çalıştırmak için aşağıdaki make komutlarını kullanabilirsiniz:

### Genel Test Komutları
```bash
# Tüm testleri çalıştır
make test

# Tüm testleri detaylı çıktı ile çalıştır
make test-verbose

# Testleri kapsam raporu ile çalıştır
make test-cover

# HTML kapsam raporu oluştur
make test-coverage-report
```

### Paket Özel Testler
```bash
# JSON testlerini çalıştır
make test-json

# XML testlerini çalıştır
make test-xml

# YAML testlerini çalıştır
make test-yaml

# Node testlerini çalıştır
make test-node

# Benchmark testlerini çalıştır
make test-bench
```

### Paket Özel Kapsam Raporları
```bash
# JSON testlerini kapsam raporu ile çalıştır
make test-json-cover

# XML testlerini kapsam raporu ile çalıştır
make test-xml-cover

# YAML testlerini kapsam raporu ile çalıştır
make test-yaml-cover

# Node testlerini kapsam raporu ile çalıştır
make test-node-cover
```

Mevcut test kapsamı: >90%

## Performans

Kütüphane aşağıdaki konularda optimize edilmiştir:
- Bellek verimliliği
- CPU kullanımı
- Büyük dosya işleme
- Eşzamanlı işlemler

### Benchmark Sonuçları

```bash
goos: darwin
goarch: arm64
cpu: Apple M1
BenchmarkJSONMarshal-8           4416622               261.0 ns/op           192 B/op          2 allocs/op
BenchmarkXMLMarshal-8             975189              1230 ns/op            4704 B/op         10 allocs/op
BenchmarkYAMLMarshal-8            213493              5284 ns/op           16728 B/op         47 allocs/op
BenchmarkJSONUnmarshal-8         1000000              1742 ns/op             272 B/op          9 allocs/op
BenchmarkXMLUnmarshal-8           370683              3104 ns/op            2328 B/op         56 allocs/op
BenchmarkYAMLUnmarshal-8          142972              8640 ns/op           10128 B/op        108 allocs/op
BenchmarkLargeJSONMarshal-8        66734             17580 ns/op           10953 B/op          2 allocs/op
BenchmarkLargeXMLMarshal-8         12298             97192 ns/op           33456 B/op         15 allocs/op
BenchmarkLargeYAMLMarshal-8         2500            466568 ns/op         1581555 B/op       3149 allocs/op
```

#### Analiz
- **JSON** hem marshal hem de unmarshal işlemlerinde en iyi performansı göstermektedir
  - Marshal: ~261 ns/op ve sadece 2 bellek tahsisi
  - Unmarshal: ~1.7 µs/op ve 9 bellek tahsisi
- **XML** JSON'a göre daha yavaş performans göstermektedir
  - Marshal: ~1.2 µs/op ve 10 bellek tahsisi
  - Unmarshal: ~3.1 µs/op ve 56 bellek tahsisi
- **YAML** en yüksek kaynak kullanımına sahiptir
  - Marshal: ~5.2 µs/op ve 47 bellek tahsisi
  - Unmarshal: ~8.6 µs/op ve 108 bellek tahsisi
- Büyük veri işlemleri için:
  - JSON minimal bellek tahsisi ile verimliliğini korumaktadır
  - XML orta düzeyde performans düşüşü göstermektedir
  - YAML hem süre hem de bellek kullanımında önemli artış göstermektedir

## Güvenlik

- XML varlık saldırılarını önlemek için girdi doğrulama
- Büyük dosyalar için bellek limiti kontrolleri
- Güvenli tip dönüşümleri
- Harici komut çalıştırma yok

## Katkıda Bulunma

Bu proje açık kaynaklıdır ve katkılara açıktır. Her türlü katkıda bulunmaktan veya geri bildirimde bulunmaktan çekinmeyin.


## Lisans

Bu proje Apache License, Version 2.0 altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.