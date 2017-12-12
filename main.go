package main

/*
#cgo pkg-config: vlc-plugin

#include <stdlib.h>

// VLC core API headers
#include <vlc_common.h>
#include <vlc_messages.h>
#include <vlc_plugin.h>
#include <vlc_interface.h>
#include <vlc_playlist.h>
#include <vlc_input.h>

input_thread_t* my_CurrentInput(vlc_object_t* obj);
void my_msg_Info(vlc_object_t* obj, char* str);
*/
import "C"

import (
    "fmt"
    "unsafe"
    "net/http"
    "io"
    "context"
    "time"
)
// unsafe.Pointer
func log_Info(ptr *C.vlc_object_t, format string, a ...interface{}) {
    msg := fmt.Sprintf(format, a...)
    cs := C.CString(msg)
    C.my_msg_Info(ptr, cs)
    C.free(unsafe.Pointer(cs))
}

func var_InheritString(ptr *C.vlc_object_t, arg string) string {
    name := C.CString(arg)
    who := C.var_InheritString(ptr, name);
    C.free(unsafe.Pointer(name))
    C.free(unsafe.Pointer(who))
    return C.GoString(who)
}

type Foo struct {
    name string
    intf *C.struct_vlc_object_t
    srv  *http.Server
}

var ctx = Foo{}

//export Open
func Open(ptr *C.vlc_object_t) C.int {
    // Try to read settings
    ctx.name = var_InheritString(ptr, "hello-who")

    ctx.intf = ptr

    // Greet!
    log_Info(ptr, "Go says: Hello %s!", ctx.name)

    // Call a goroutine
    go testGoroutine(ptr)

    log_Info(ptr, "Starting Server, %s!", ctx.name)

    ctx.srv = &http.Server{
        Addr:         ":8080",
        ReadTimeout:  time.Second,
        WriteTimeout: time.Second,
    }
    http.HandleFunc("/", rootHandler)
    go ctx.srv.ListenAndServe()

    return 0;
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    log_Info(ctx.intf, "Request got!")

    input := C.my_CurrentInput(ctx.intf)

    if (input == nil) {
        io.WriteString(w, "Nothing playing")
        return
    }

    defer C.vlc_object_release((*C.struct_vlc_object_t)(unsafe.Pointer(input)))
    item := C.input_GetItem(input)
    tmp_cstr := C.input_item_GetTitleFbName(item)
    title := C.GoString(tmp_cstr)
    io.WriteString(w, "Currently playing: ")
    io.WriteString(w, title)
    C.free(unsafe.Pointer(tmp_cstr))
}

func testGoroutine(ptr *C.vlc_object_t) {
    log_Info(ptr, "Hello from a Goroutine!")
}

//export Close
func Close(ptr *C.vlc_object_t) {
    if (ctx.srv == nil) {
        log_Info(ptr, "Oh noes!")
    }
    if err := ctx.srv.Shutdown(context.Background()); err != nil {
        log_Info(ptr, "failure to shutdown server!")
        panic(err) // failure/timeout shutting down the server gracefully
    }
    log_Info(ptr, "Go says: Bye %s!", ctx.name)
}

func main() { }
