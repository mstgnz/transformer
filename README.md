## Transformer
#### Transformer for json, yaml and xml files

The main purpose of this project is to consolidate what I know, apply, learn new things, understand file systems and at the end of the day serve humanity :)

This project aims to convert json, yaml and xml files to each other.

During the conversion process, the file format will be checked first. If the file is in an error-free format, an error-free conversion takes place.

Users will be able to do this conversion process in two ways with command line or api.

I have assigned the adopted path for conversion between files as follows. Rather than converting the given file directly, converting it to an object and giving it a chance to make certain changes to that object and then produce the output in the desired format.

This way the project will not resist development.

The map that comes embedded in the Go programming language is unfortunately unordered. so I'm looking for a way to decode sequentially.

- Check File Format - DONE
- Decode File - DONE (but the map is unordered so UNDONE) - do refactor
- Create recursive function and write to node - DONE
- Create and print the desired file - UNDONE  
- .......