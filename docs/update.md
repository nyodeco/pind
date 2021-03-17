# Update

* Run the following commands to update pind, all dependencies, and install it:

```bash
cd $GOPATH/src/github.com/nyodeco/pind
git pull && GO111MODULE=on go install -v . ./cmd/...
```
