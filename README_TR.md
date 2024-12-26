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

Kütüphane kapsamlı test kapsamına sahiptir. Testleri çalıştırmak için:

```bash
# Tüm testleri çalıştır
go test ./...

# Kapsam raporuyla testleri çalıştır
go test -cover ./...

# Detaylı test çıktısıyla çalıştır
go test -v ./...
```

Mevcut test kapsamı: >90%

## Performans

Kütüphane şunlar için optimize edilmiştir:
- Bellek verimliliği
- CPU kullanımı
- Büyük dosya işleme
- Eşzamanlı işlemler

## Güvenlik

- XML varlık saldırılarını önlemek için girdi doğrulama
- Büyük dosyalar için bellek limiti kontrolleri
- Güvenli tip dönüşümleri
- Harici komut çalıştırma yok

## Katkıda Bulunma

Bu proje açık kaynaklıdır ve katkılara açıktır. Her türlü katkıda bulunmaktan veya geri bildirimde bulunmaktan çekinmeyin.


## Lisans

Bu proje Apache License, Version 2.0 altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.