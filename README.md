# Catppuccin CLI (WIP)

A work-in-progress CLI for Catppuccin themes.

It allows you to:
- Install themes with one command
- Uninstall themes with one command
- Update themes with one command

### `.ctprc`
The current `.ctprc` spec, which is intended to be kept inside each Catppuccin repo moving forward.
```yaml
app_name: # the app name
path_name: # the name by which you invoke this program from the command line
install_location:
  unix: # for unix systems, where the configs are installed
  windows: # same thing for windows
install_flavours:
  all: #instructions to install all flavours(if possible)
    default:                   # the default mode, because themes can have "variants"
      - themes/catppuccin.toml # files to install
    additional: # optional additional variants
      no-italics: # an example variant
        - themes/catppuccin-no-italics.toml # files to install
  # same thing for latte, frappe, blah blah blah
  to: themes/ # where to install the files/directories to
one_flavour: false # if you can only install one flavour at a time
modes:
  - no-italics # the modes you specified in install flavours
```
That was lazy. I'll make a full spec tomorrow, like I said. If ya want a good(and updated) example, check out [my helix .ctprc](https://github.com/catppuccin/helix-new/blob/main/.ctprc).
