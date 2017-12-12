#include <stdlib.h>

// VLC core API headers
#include <vlc_common.h>
#include <vlc_plugin.h>
#include <vlc_interface.h>
#include <vlc_playlist.h>

// Module string
#define MODULE_STRING "goweb"

extern int Open(vlc_object_t*);
extern void Close(vlc_object_t*);

// Module descriptor
vlc_module_begin()
    set_shortname("goweb")
    set_description("Go webinterface")
    set_capability("interface", 0)
    set_callbacks(Open, Close)
    set_category(CAT_INTERFACE)
    add_string("hello-who", "world", "Target", "Whom to say hello to.", false)
vlc_module_end()

void my_msg_Info(vlc_object_t* obj, char* str)
{
    msg_Info(obj, "%s", str);
}

input_thread_t* my_CurrentInput(vlc_object_t* obj)
{
    return pl_CurrentInput((intf_thread_t*)obj);
}
