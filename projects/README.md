# Project directory

This directory contains the source code for the different projects I make based on XPHTTPBridge and XPicoConnect. Each project is a separate directory with its own module file.

The `github.com/steveiliop56/xpicoconnect` module is replaced to point to the local version of the library which makes development easier. In your own projects, you may want to remove the `replace github.com/steveiliop56/xpicoconnect => ../..` line from your `go.mod` file.
