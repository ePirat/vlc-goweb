# VLC GoWeb
VLC GoWeb is a simple VLC web interface module written in Go.
Currently it is just a simple and rough experiment, mainly to try out if writing a VLC module in Go could work.

## Build
To build the plugin, make sure you have the VLC 3.0 plugin headers. The include and lib flags are found with pkg-config,
so make sure it can find `vlc-plugin`.

If you have your VLC headers in a custom location, you can use the `PKG_CONFIG_PATH` env variable to point to the correct
location of the pkgconfig files for them.

Then you should be able to just compile it like that:

```
go build -buildmode=c-shared -o libgoweb_plugin.dylib
```

If you want a smaller library, turn off DWARF symbols (`-w`) and strip Go symbols (`-s`):

```
go build -ldflags="-s -w" -buildmode=c-shared -o libgoweb_plugin.dylib
```

## Installation

To install it, just drop the resulting library in your VLC plugin directory and reload the module cache.
