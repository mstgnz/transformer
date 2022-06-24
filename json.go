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
		key      string
		value    string
		objStart bool
		arrStart bool
		objCount int
		arrCount int
	)
	for {
		t, err := dec.Token()
		if err == io.EOF || err != nil {
			return knot, errors.Wrap(err, "no more")
		}
		// Get type and value of t object in each loop
		value = fmt.Sprintf("%v", t)
		// If the type of the object is json.Delim
		if reflect.TypeOf(t).String() == "json.Delim" {
			// If no node has been created yet, don't enter here, skip json start -> {
			if !knot.Exists() {
				continue
			}
			// If the value of object t is json object or array
			switch value {
			case "{": // set open object - {
				/*
					burada iki tür opsiyon var.
					1- key boş değilse; bu bir objedir. mevcut düğümün nextine yeni bir düğüm eklenecek ve geriye yeni eklenen düğüm dönecektir.
					2- key boş ise; bir array objesi kesin vardır ve array içerisinde düğüm oluşturulacak ve geriye yeni eklenen düğüm dönecektir.
				*/
				if arrStart {
					// bu aslında olmazsa olmazdır çünkü bir obje sadece ve sadece array içersinde keysiz başlar.
					knot = knot.AddObjToArr(knot)
					arrStart = true
				} else {
					knot = knot.AddToNext(knot, parent, key, &node{})
					parent = knot
				}
				objStart = true
				arrStart = false
				objCount++
				key = ""
			case "[": // set open array - [
				/*
					burada iki tür opsiyon var
					1- key boş değilse ve objStart true ise; nodun ilk değeri set edilecektir.
					2- key boş değilse ve objStart false ise; nodun nextine yeni bir node oluşturulacaktır.
					3- key boş ise; bu bir nested array nesnesidir. direk olarak mevcut düğümün arrayine eklenecektir.
				*/
				if len(key) > 0 {
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
					//knot = knot.AddToArr(knot, value)
				}
				arrStart = true
				objStart = false
				arrCount++
			case "]": // set close array
				/*if objCount > 0 {
					knot = knot.prev
				}*/
				arrCount--
				arrStart = false
			case "}": // set close object and set parent node
				if arrCount > 0 {
					arrStart = true
				}
				objCount--
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
				if arrStart {
					knot = knot.AddToArr(knot, value)
				} else {
					key = value
				}
			} else {
				// eğer objStart true ise nodun ilk değeri set ediliyor.
				if objStart {
					knot = knot.AddToValue(knot, parent, key, value)
				} else {
					knot = knot.AddToNext(knot, parent, key, value)
				}
				// burada objStart ve key değerleri sıfırlanıyor.
				objStart = false
				key = ""
			}
		}
	}
}
