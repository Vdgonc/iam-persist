# iam-persist


## Example

```golang
package main

import (
	"fmt"
	"os"

	iampersist "github.com/vdgonc/iam-persist"
)

func main() {
	creds := iampersist.CreatePersistence(&iampersist.CreatePersistenceInput{
		AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Username: "test.iam",
	})

	fmt.Printf("access_key = %s\nsecret_key=%s", creds.AccessKey, creds.SecretKey)
}
```