#!/bin/bash
set -e
cd /workspaces/retrocode 2>/dev/null || cd "$(git rev-parse --show-toplevel)"
echo "[1/5] tools:"; go version; cargo --version
echo "[2/5] building client (first time ~4 min, FLTK compiles)..."
cd client && cargo build 2>&1 | tail -3 && cargo test 2>&1 | tail -3 && cd ..
echo "[3/5] starting server..."
cd server && go build -o bin/retrocode-server ./cmd/retrocode-server
pkill retrocode-server 2>/dev/null; nohup ./bin/retrocode-server > server.log 2>&1 &
sleep 1; cat server.log; cd ..
echo "[4/5] smoke test over framed TCP..."
cat > /tmp/smoke.go <<'GOEOF'
package main
import ("encoding/binary";"fmt";"io";"net")
func main(){
 c,e:=net.Dial("tcp","127.0.0.1:7741"); if e!=nil{panic(e)}; defer c.Close()
 p:=[]byte(`{"m":"session/ping","params":{"t":42}}`)
 h:=make([]byte,4); binary.LittleEndian.PutUint32(h,uint32(len(p)))
 c.Write(append(h,p...))
 rh:=make([]byte,4); if _,e:=io.ReadFull(c,rh);e!=nil{panic(e)}
 ln:=binary.LittleEndian.Uint32(rh); b:=make([]byte,ln)
 if _,e:=io.ReadFull(c,b);e!=nil{panic(e)}
 fmt.Println("SMOKE-OK SERVER SAID:",string(b))
}
GOEOF
cd server && go run /tmp/smoke.go; cd ..
echo "[5/5] committing..."
git add -A; git commit -m "M0b: client+server built, framed TCP smoke test green" --allow-empty
git push
echo "==== ALL GREEN - copy everything above and paste it back to the assistant ===="
