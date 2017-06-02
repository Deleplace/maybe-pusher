// Package maybe-pusher offers HTTP/2 server push capability if current go net/http package has it,
// and doesn't fail to compile if current go version (less than 1.8) doesn't have it.
package pusher

import (
	"log"
	"net/http"
	"reflect"
)

// PushCapable returns whether HTTP/2 server push capability is available through htpp.Pusher.
func PushCapable(w http.ResponseWriter) bool {
	_, ok := getPushFunc(w)
	return ok
}

// Push initiates an HTTP/2 server push if available through http.Pusher, and returns true.
// Otherwise, it returns false.
func Push(w http.ResponseWriter, target string) bool {
	push, ok := getPushFunc(w)
	if !ok {
		// Not available
		return false
	}

	args := make([]reflect.Value, 2)
	args[0] = reflect.ValueOf(target)
	log.Println("args[0] is", args[0])
	optionsStructType := push.Type().In(1).Elem()
	log.Println("optionsStructType is", optionsStructType)
	args[1] = reflect.New(optionsStructType)
	log.Println("args[1] is", args[1])
	// args[1] = reflect.ValueOf(nil)
	push.Call(args)

	return true
}

// Push initiates an HTTP/2 server push if available through http.Pusher, and returns true.
// Otherwise, it returns false.
func PushWithOptions(w http.ResponseWriter, target string, method string, header http.Header) bool {
	push, ok := getPushFunc(w)
	if !ok {
		// Not available
		return false
	}

	//
	// TODO
	//
	_ = push

	return true
}

func getPushFunc(w http.ResponseWriter) (pushFunc reflect.Value, ok bool) {
	if w == nil {
		return pushFunc, false
	}
	// wType := reflect.TypeOf(w)
	wValue := reflect.ValueOf(w)
	pushFunc = wValue.MethodByName("Push")
	if pushFunc.IsNil() {
		// Doesn't implement Pusher
		return pushFunc, false
	}

	// If w.Push exists but doesn't match http.Pusher, then just return nil (don't crash).
	if pushFunc.Type().Kind() != reflect.Func {
		log.Println("Warning: Push exists, but is not a method")
		return pushFunc, false
	}
	if pushFunc.Type().NumIn() != 2 {
		log.Println("Warning: method Push exists, but has", pushFunc.Type().NumIn(), "arguments instead of 2")
		return pushFunc, false
	}
	if pushFunc.Type().NumOut() != 1 {
		log.Println("Warning: method Push exists, but has", pushFunc.Type().NumOut(), "return values instead of 1")
		return pushFunc, false
	}
	if pushFunc.Type().In(0).Kind() != reflect.String {
		log.Println("Warning: expected string as Push first argument, found", pushFunc.Type().In(0).Kind())
		return pushFunc, false
	}

	log.Println("Method type is", pushFunc.Type())
	return pushFunc, true
}
