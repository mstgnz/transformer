# Transformer

Transformer, farklı veri formatları (JSON, XML, YAML) arasında dönüşüm yapmanızı sağlayan bir Go kütüphanesidir. Bir belge tek bir ortak agaca cozulur ve butun yazicilar bu agaci okur; boylece anahtar sirasi, deger tipleri ve tekrarli element'ler yolculuktan sag cikar. Her formatin ifade edemedigi seyler [Format-Spesifik Davranislar](#format-spesifik-davranislar) altinda listelidir.

[English Documentation (İngilizce Dokümantasyon)](README.md)

## Özellikler

- JSON, XML ve YAML formatları arasında dönüşüm
- Anahtar sırası, değer tipleri ve tekrarlı elementler her dönüşümden sağ çıkar
- Belgede anahtar arama: her eşleşmenin yolu ve derecesi ile birlikte
- Bul değiştir: eşleşen her şeyi topluca ya da yolun adreslediği tek geçişi
- Herhangi bir node veya değere yol ile erişim, array indeksleri dahil
- Sayılar hiçbir zaman float64 üzerinden geçmez, uzun kimlikler her hanesini korur
- Yüksek test kapsamı
- Tek bağımlılık: `gopkg.in/yaml.v3`

## Gereksinimler

- Go 1.23 veya üzeri
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
yamlData, err := tyaml.NodeToYaml(node) // []byte döner
if err != nil {
    log.Fatal(err)
}
```

### Formatlar Arası Dönüşüm (tek satir API)

```go
import transformer "github.com/mstgnz/transformer"

jsonData := []byte(`{"name": "John", "age": 30}`)

// JSON -> XML
xmlData, err := transformer.JsonToXml(jsonData)

// JSON -> YAML
yamlData, err := transformer.JsonToYaml(jsonData)

// XML -> JSON
xmlInput := []byte(`<root><name>John</name></root>`)
jsonOut, err := transformer.XmlToJson(xmlInput)

// XML -> YAML
yamlOut, err := transformer.XmlToYaml(xmlInput)

// YAML -> JSON
yamlInput := []byte("name: John\nage: 30\n")
jsonOut, err = transformer.YamlToJson(yamlInput)

// YAML -> XML
xmlOut, err := transformer.YamlToXml(yamlInput)
```

### Formatlar Arası Dönüşüm (Node uzerinden)

Degerleri adimlar arasinda incelemek veya degistirmek icin Node API'si kullanilabilir:

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

## Belge Icinde Arama

Node modelinin varlik sebebinin yarisi donusum. Digeri arama: buyuk bir belgede bir anahtar yirmi yerde gecebilir, hangisinin kastedildigini de yol soyler.

```go
matches, _ := transformer.SearchJson(data, "name")
for _, m := range matches {
    fmt.Println(m)
}
// root.meta.name (depth 2) = catalog
// root.groups[0].members[0].name (depth 5) = alice
// root.groups[0].members[0].profile.name (depth 6) = Alice Doe
// root.deep[0][0].items[0].entries[0].detail.name (depth 9) = deepest
```

`SearchJson`, `SearchXml` ve `SearchYaml` decode ve aramayi tek cagriada yapar. Elinizde zaten bir agac varsa karsiligi `n.Search(key)`.

Her `Match` uc sey tasir:

- `Node` — eslesen node
- `Path` — ona nasil ulasilacagi, aranan node'un kendi anahtariyla baslayarak: `root.groups[0].members[1].profile.name`
- `Depth` — bu yolun kac adim oldugu. Dogrudan cocuk 1. derecededir, array indeksi de kendi basina bir adimdir.

### Bir anahtar hangi derecelerde

```go
root.Depths("name")   // [1 2 3 4 5 6 7 8 9 10]
```

`Depths` bir anahtarin gectigi dereceleri artan sirada ve tekrarsiz verir; tek bir seviyede gecen anahtarla belgenin tamamina yayilan anahtar boylece ayirt edilir.

### Degeri geri okumak

Bir eslesmeden gelen yol `GetNodeByPath` ile geri okunur; `GetValueByPath` ise listenin icindeki, kendi node'u olmayan duz degerlere de ulasir:

```go
root.GetNodeByPath("root.groups[0].members[1].profile.name")   // *Node
root.GetValueByPath("root.matrix[1][2]")                       // *Value, true
root.GetValueByPath("root.cube[1][0][1]")                      // *Value, true
```

Icinde nokta ya da kose parantez gecen bir anahtar bu yolla adreslenemez; eslesmenin tasidigi `Node` kullanilir.

### Baska olcutlerle arama

```go
// belgedeki her liste
lists := root.SearchFunc(func(n *node.Node) bool { return n.Type() == node.TypeArray })

// bu anahtar, listelerin ici dahil, herhangi bir yerde var mi
root.Exists("milestones")   // true
```

Array'lerin ve array icindeki array'lerin icindeki anahtarlara da digerleri gibi ulasilir. Listenin icindeki duz bir degerin kendi node'u yoktur, dolayisiyla arama sonucu olamaz; listesi aranir ve degerler oradan okunur.


## Belgeyi Degistirme

Bir seyi bulmak isin genelde ilk yarisi. `EditJson`, `EditXml` ve `EditYaml` belgeyi cozer, agaci size verir ve sonucu geldigi formatta geri yazar:

```go
out, _ := transformer.EditJson(data, func(n *node.Node) {
    n.ReplaceText("staging.example.com", "example.com")
    n.Rename("labels", "tags")
    n.Remove("debug")
})
```

### Eslesen her seyi degistirmek

| | |
| --- | --- |
| `Replace(key, value)` | o anahtara sahip her node'a degerin bir kopyasini verir |
| `ReplaceFunc(predicate, replace)` | kabul edilen her node'u `replace`'e verir |
| `Rename(oldKey, newKey)` | `oldKey` okunan her anahtari yeniden adlandirir |
| `Remove(key)` | o anahtara sahip her obje uyesini siler |
| `ReplaceText(old, new)` | listeler dahil her metin degerinin icindeki metni degistirir |

Her biri kac node degistirdigini dondurur; hicbir seyle eslesmeyen bir degistirme sessiz kalmaz.

Her hedef degerin kendi kopyasini alir, dolayisiyla sonuclardan birini degistirmek digerlerini degistirmez. `ReplaceText` yalnizca metne dokunur: ayni karakterleri iceren bir anahtar, sayi veya boolean oldugu gibi kalir; anahtarlari degistiren `Rename`'dir.

### Tam olarak birini degistirmek

Aramadan gelen bir yol tek bir node'u adresler; hepsini degil kastettiginiz tek gecisi degistirmenin yolu budur:

```go
matches := root.Search("name")
root.SetByPath(matches[3].Path, &node.Value{Type: node.TypeString, Worth: "renamed"})
```

| | |
| --- | --- |
| `SetByPath(path, value)` | yolun adresledigi degeri degistirir; son segment olmayan bir anahtari adlandiriyorsa uyeyi ekler |
| `RemoveByPath(path)` | bir obje uyesini ya da listeden bir elemani cikarir |

Ikisi de yolun bir yere varip varmadigini bildirir. `SetByPath` liste elemanlarina da, liste icindeki listelere de ulasir (`root.matrix[1][2]`) ve bir liste elemanina verilen objeyi decoder'larin sardigi gibi sarar.

Agac butun bunlar boyunca saglam kalir: herhangi birinden sonra `Validate` gecmeye devam eder ve belge uc formattan herhangi birinde geri yazilabilir.


## Format-Spesifik Davranislar

### Anahtar Sirasi

Anahtarlar, kaynak belgede yazildiklari sirayi her yonde korur. Ne decode ne de encode tarafi Go map uzerinden gecer; map iterasyonu her calismada rastgelelestigi icin belge her donusumde farkli bir sirayla cikardi.

```go
out, _ := transformer.JsonToXml([]byte(`{"zebra":1,"apple":2}`))
// her calismada <root><zebra>1</zebra><apple>2</apple></root>
```

### XML Tekrarli Element'ler

Ayni ismi paylasan kardes element'ler tek bir array node'unda toplanir; array node'u da her deger icin bir element olarak geri yazilir. XML'in liste sozdizimi olmadigi icin bir listenin XML'den sag cikmasini saglayan sey budur.

```go
out, _ := transformer.XmlToJson([]byte(`<root><item>1</item><item>2</item></root>`))
// {"item":[1,2]}
```

### XML Attribute ve Metin

Attribute'lar `@ad` seklinde cocuk node olarak saklanir. Hem attribute hem metin tasiyan bir element, metnini `#text` cocugunda tutar; boylece ikisi birbirinin yerini almaz. Yalnizca metin tasiyan bir element duz string, sayi veya boolean olarak kalir.

```go
out, _ := transformer.XmlToJson([]byte(`<root id="5"><name lang="en">John</name></root>`))
// {"@id":"5","name":{"@lang":"en","#text":"John"}}
```

Ayni yapi geri donusturuldugunde orijinal XML elde edilir.

### XML Tip Isaretleri

JSON ve YAML'in yazabildigi birkac degerin XML'de kendine ait bir sekli yok. Bos bir etiket null da olabilir, bos string de, bos obje de; tek bir `<port>` element'i de tek bir deger olabilir, tek elemanli bir liste de. Sekil kaybolacak oldugunda yazici element'i kutuphanenin kendi namespace'inde isaretler, okuyucu da geri okur:

```xml
<n t:type="null"/>                          <!-- null -->
<s t:type="string"/>                        <!-- bos string -->
<a t:type="array"/>                         <!-- bos liste -->
<a t:type="array"><item>80</item></a>       <!-- tek elemanli liste -->
```

Isaret yalnizca gerektigi yerde yazilir, namespace de yalnizca belge bir isaret kullaniyorsa root element'te bildirilir. Iki ve daha fazla elemanli duz listeler yine duz tekrarli element olarak gider; siradan belgeler isaretler eklenmeden onceki ciktinin aynisini verir.

```go
out, _ := transformer.JsonToXml([]byte(`{"ports":[80,443]}`))
// <root><ports>80</ports><ports>443</ports></root>   isarete gerek yok

out, _ = transformer.JsonToXml([]byte(`{"ports":[80]}`))
// <root xmlns:t="..."><ports t:type="array"><item>80</item></ports></root>
```

### XML Element Adlari

XML adlari serbest metin degildir: bir element `first name` ya da `1st` olarak adlandirilamaz. Gecerli bir ad olmayan anahtar cikista kacislanir, giriste geri okunur; boylece anahtar korunur ve belge gecerli kalir. Zaten gecerli olan adlar, ASCII disi olanlar dahil, oldugu gibi yazilir.

```go
out, _ := transformer.JsonToXml([]byte(`{"first name":"John","şehir":"İstanbul"}`))
// <root><first_x0020_name>John</first_x0020_name><şehir>İstanbul</şehir></root>
```

Iki nokta ancak belge o prefix'i bildirmisse namespace prefix'i sayilir. `xmlns:soap` bildirimi yanindaki `soap:Envelope` prefix'li ad olarak kalir; icinde sadece iki nokta gecen bir anahtarin iki noktasi kacislanir, cunku onu prefix olarak yazmak hic bildirilmemis bir namespace'e atifta bulunurdu.

### XML Namespace'leri

Namespace prefix'leri uçtan uca korunur. `soap:Envelope` gibi bir element, Node'da `"soap:Envelope"` anahtariyla saklanir; namespace bildirimleri `@xmlns:soap` gibi attribute'lar olarak tutulur. Namespace'li bir XML belgesinin round-trip'i ayni ciktiy uretir.

```xml
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body/>
</soap:Envelope>
```

Bu XML decode edildiginde key'i `"soap:Envelope"`, attribute'u `@xmlns:soap` olan bir Node olusur ve geri ayni XML'e encode edilir.

### JSON Buyuk Tamsayilar

JSON sayilari `json.Decoder.UseNumber()` ile decode edilir ve belgenin yazdigi metin olarak tutulur, hicbir zaman `float64` olarak degil. Bu yuzden 19 haneli bir kimlik 19 hanesiyle geri gelir, `1e10` da geldigi gosterimle geri doner.

```go
// 9007199254740993 tam olarak korunur, float64 yuvarlamasi olmaz
node, _ := tjson.DecodeJson([]byte(`{"id": 9007199254740993}`))
out, _ := tjson.NodeToJson(node)
// out == {"id":9007199254740993}
```

### YAML Anchor ve Alias'lar

YAML anchor'lari (`&anchor`) ve alias'lari (`*alias`), veri Node yapisina ulasmadan once YAML parser tarafindan otomatik olarak cozumlenir. Cozumlenen degerler dogru sekilde korunur; ancak anchor/alias sozdizimi yeniden encode sirasinda saklanmaz.

```yaml
defaults: &defaults
  timeout: 30

production:
  <<: *defaults # birlestirilir: decode sonrasi production.timeout == 30
  host: prod.example.com
```

### Yuvalama Derinligi Siniri

Tum decoder'lar maksimum 100 seviye yuvalama derinligini zorunlu kilar ve asildigi durumda hata doner. Bu, kotü niyetli girdilerde stack overflow'u onler.

### Round-Trip'te Korunmayanlar

- **Karisik icerikte metin konumu.** Cocuk element'lerin arasina serpistirilmis metin korunur, ancak geri yazilirken cocuk element'lerin arasina degil oncesine yazilir.
- **YAML anchor sozdizimi.** Anchor ve alias'lar giriste cozumlenir, merge key (`<<`) yerinde genisletilir. Degerler dogrudur, sozdizimi geri yazilmaz.
- **XML yorum ve isleme talimatlari.** Decode sirasinda atilir.

Bu kutuphanenin yazmadigi XML icin bilinmesi gereken bir sey var: isaretsiz bir `<a/>` bos obje olarak okunur, cunku isaret yokken onun null mi, bos string mi, bos obje mi kastedildigini soyleyecek bir sey yoktur.

Bunlarin disindaki her sey korunur. JSON, XML ve YAML'in her biri Node modelinden degismeden gecer; baska bir formata cevrilip geri getirilen belge, basladigi belge olarak geri gelir.

## Paket Yapısı

- `node`: Temel veri yapısını ve operasyonlarını içerir
  - Hiyerarşik veriyi temsil eden Node yapısı
  - Değer tipleri ve tip dönüşüm operasyonları
  - Ağaç gezinme ve manipülasyon fonksiyonları
  - Arama: `Search`, `SearchFunc`, `Depths`, `Exists`, `GetNodeByPath`, `GetValueByPath`
  - Değiştirme: `Replace`, `ReplaceFunc`, `ReplaceText`, `Rename`, `Remove`, `SetByPath`, `RemoveByPath`
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
