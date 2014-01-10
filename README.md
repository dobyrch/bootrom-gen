bootrom-gen
===========

A program for generating custom Game Boy boot ROMs


How it works
============

The program looks for two images in the current directory:

* A 48x8 pixel png called *logo.png*
* An 8x8 pixel png called *notice.png*

It modifies the standard Game Boy boot procedure to display your custom images
in place of the Nintendo logo.  Currently the patched ROM is outputted in
binary format, so it is recommended that you pipe the output to a file.
