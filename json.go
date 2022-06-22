package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func isJSON(doc []byte) bool {
	return json.Unmarshal(doc, new(map[string]any)) == nil
}

// jsonDecode
func jsonDecode(doc []byte) (*node, error) {
	var (
		knot   *node
		parent *node
	)
	dec := json.NewDecoder(strings.NewReader(string(doc)))
	var (
		key         string
		types       string
		typeVal     string
		objStart    bool
		arrStart    bool
		arrObjStart bool
	)
	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return knot, errors.Wrap(err, "no more")
		}
		// Get type and value of t object in each loop
		types = reflect.TypeOf(t).String()
		typeVal = fmt.Sprintf("%v", t)
		// If the type of the object is json.Delim
		if types == "json.Delim" {
			// Ff no node has been created yet, don't enter here, skip json start -> {
			if !knot.Exists() {
				continue
			}
			// If the value of object t is json object or array
			switch typeVal {
			case "{": // set open object - {
				/*
					burada iki tür opsiyon var.
					1- key boş değilse; bu bir objedir. mevcut düğümün nextine yeni bir düğüm eklenecek ve geriye yeni eklenen düğüm dönecektir.
					2- key boş ise; bir array objesi kesin vardır ve array içerisinde düğüm oluşturulacak ve geriye yeni eklenen düğüm dönecektir.
				*/
				if len(key) > 0 {
					knot = knot.AddToNext(knot, parent, key, &node{})
					parent = knot
				} else {
					// bu aslında olmazsa olmazdır çünkü bir obje sadece ve sadece array içersinde keysiz başlar.
					knot = knot.AddObjToArr(knot)
				}
				if arrStart {
					arrObjStart = true
					arrStart = false
				}
				objStart = true
				key = ""
			case "[": // set open array - [
				/*
					burada iki tür opsiyon var
					1- key boş değilse ve objStart true ise; nodun ilk değeri set edilecektir.
					2- key boş değilse ve objStart false ise; nodun nextine yeni bir node oluşturulacaktır.
					3- key boş ise; bu bir nested array nesnesidir. direk olarak mevcut düğümün arrayine eklenecektir.
				*/
				if key != "" {
					// eğer objStart true ise nodun ilk değeri set ediliyor.
					if objStart {
						knot = knot.AddToValue(knot, parent, key, []any{})
					} else {
						// eğer objStart false ise mevcut nodun nextine yeni bir node oluşturuluyor.
						knot = knot.AddToNext(knot, parent, key, []any{})
					}
					parent = knot
					key = ""
				} else {
					// eğer key yok ise iç içe arraydir.
					// TODO nested array için indis tutulacak ve array kapatılana kadar bu arraye append edilecektir.
					//knot = knot.AddToArr(knot, typeVal)
				}
				objStart = false
				arrStart = true
			case "]": // set close array
				arrStart = false
			case "}": // set close object and set parent node
				if arrObjStart {
					arrStart = true
				}
				parent = nil
				if knot.parent != nil {
					knot = knot.parent
					parent = knot.parent
				}
			default: // shouldn't go here
				fmt.Println("default not set -> ", t)
			}
		} else {
			// döngü nesnesi bir json.Delim değil ise key ve value alanları set edilecektir.
			// json nesnesi bir key value değer çifti olduğu için önce key set edilecek daha sonra value set edilecektir.
			if len(key) == 0 {
				// eğer bir array objesi açık ise bu key değeri esasen array nesnesidir.
				// if the array is not empty
				if arrStart && !objStart {
					knot = knot.AddToArr(knot, typeVal)
				} else {
					key = typeVal
				}
			} else {
				// eğer objStart true ise nodun ilk değeri set ediliyor.
				if objStart {
					knot = knot.SetToValue(knot, key, typeVal)
				} else {
					knot = knot.AddToNext(knot, parent, key, typeVal)
				}
				// burada objStart ve key değerleri sıfırlanıyor.
				objStart = false
				key = ""
			}
		}
	}
}

// array içinde obje ve array işi karıştırıyor buna bi çözüm düşüneceğiz.
