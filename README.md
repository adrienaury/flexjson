# Ordered JSON

Parse semi-structured JSON (i.e. data that can't be mapped to a known structure) in a streaming manner, and decide what types will be used to parse objects and arrays.

Example

```go
package main

import (
	"os"

	"github.com/adrienaury/flexjson"
	"github.com/goccy/go-json"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type (
	Object = *orderedmap.OrderedMap[string, any]
	Array  = flexjson.Array
)

func main() {
	dec := flexjson.NewFlexDecoder(
		json.NewDecoder(os.Stdin),
		func() Object { return orderedmap.New[string, any]() },
		func(obj Object, key string, value any) Object {
			obj.Set(key, value)
			return obj
		},
		flexjson.StandardArrayMaker(),
		flexjson.StandardArrayAdder(),
	)

	dec.Decode()
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
