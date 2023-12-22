# mp3renamer

I had accidentally deleted some of my music files and directories from my mp3 folder.
Luckily I had a full backup on Google Drive, but the downloaded directories were zipped
and all special characters were replaced with an underscore in file and directory names
in the downloaded zip file.

This simple Go app traverses the given path and merges directories which have the same
name, but have an underscore in the name instead of a special character. It also removes
duplicate files based on this rule.

Usage: ./mp3renamer <path to directory>
