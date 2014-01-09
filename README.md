bootrom-gen
===========

A program for generating custom Game Boy boot ROMs


How it works
============

The program looks for two images in the current directory:

* A 48x8 pixel png called *logo.png*
* An 8x8 pixel png called *notice.png*

It will output a list of comma separated bytes which can be inserted into the Game Boy boot ROM to show your custom image in place of the Nintendo logo.
