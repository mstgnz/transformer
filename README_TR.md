# Transformer

Transformer, farklı veri formatları (JSON, XML, YAML) arasında dönüşüm yapmanızı sağlayan bir Go kütüphanesidir. Ortak bir veri yapısı kullanarak, bu formatlar arasında kayıpsız dönüşüm yapabilirsiniz.

[English Documentation (İngilizce Dokümantasyon)](README.md)

## Özellikler

- JSON, XML ve YAML formatları arasında dönüşüm
- Ortak veri yapısı ile tutarlı dönüşüm
- Kolay kullanım
- Tip güvenliği
- Özelleştirilebilir dönüşüm kuralları

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
yamlData := []byte(`name: John\nage: 30`)
node, _ := tyaml.DecodeYaml(yamlData)
jsonData, _ := tjson.NodeToJson(node)
```

## Paket Yapısı

- `node`: Temel veri yapısını ve operasyonlarını içerir
- `tjson`: JSON dönüşüm işlemlerini gerçekleştirir
- `txml`: XML dönüşüm işlemlerini gerçekleştirir
- `tyaml`: YAML dönüşüm işlemlerini gerçekleştirir
- `example`: Örnek kullanımları içerir

## Veri Tipleri

Node yapısı aşağıdaki veri tiplerini destekler:

- `TypeNull`: Boş değer
- `TypeObject`: Nesne (key-value çiftleri)
- `TypeArray`: Dizi
- `TypeString`: Metin
- `TypeNumber`: Sayı
- `TypeBoolean`: Mantıksal değer

## Katkıda Bulunma

1. Bu repository'yi fork edin
2. Yeni bir branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'feat: add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## Lisans

Bu proje Apache License, Version 2.0 altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.

## İletişim

Mesut GENEZ - [@mstgnz](https://github.com/mstgnz)

Proje Linki: [https://github.com/mstgnz/transformer](https://github.com/mstgnz/transformer)