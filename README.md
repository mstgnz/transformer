## Transformer
### Transformer for json, yaml and xml files

The main purpose of this project is to convert json, yaml and xml files into a golang object.

As you know, there are two ways to convert json, yaml and xml files into an object.

The first way is to create a struct and running unmarshal. But here we need to write the struct structure. For a large file this will be cumbersome.

The second way is to convert it to a map object. Unfortunately, this map object is unordered. Therefore, it does not provide the structure we want.

The third way is to write our own key-value structure and design it to be sequential.

This transformer gives us the 3rd path.

### What will this process bring us?

Our structure will hold the next, previous and parent objects along with the key/value. In this way, we will have the opportunity to move on the struct. We can also customize this struct as we want by adding the methods we want. By default, it will have search, remove, move and replace methods.

File conversion will be made easier using this package.

Now that our plan is like this, let's get started. I will start with json first, then xml and finally yaml.

json v1 will be xml v2 and yaml v3. Then we will continue in the form of bugfixes and improvements.

json v1 is under construction...