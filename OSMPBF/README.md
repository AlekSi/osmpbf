`*.proto` files were downloaded from https://github.com/scrosby/OSM-binary/tree/master/src and slightly changed.

To eliminate continuous conversions from `[]byte` to `string`, this part

```protobuf
message StringTable {
   repeated bytes s = 1;
}
```

was changed to

```protobuf
message StringTable {
   repeated string s = 1;
}
```

This change is expected to be fully compatible with all PBF files.

References:
* [`osmpbf` commit](https://github.com/AlekSi/osmpbf/commit/e702813f2b9077cadbe3521a6a12785db6d7828c)
* [Upstream discussion](https://github.com/scrosby/OSM-binary/commit/8f0c4a8bc2fe7b6669e7550854deb497cc253c28#commitcomment-5804528)
