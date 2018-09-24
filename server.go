package main

import (
  "net/http"
  "fmt"
  "strconv"
  "os"
  "crypto/rand"
  "errors"
  "io"
  "math/big"
  "encoding/binary"
  //"g//ithub.com/VividCortex/godaemon"
 )

func vPrime(w http.ResponseWriter, r *http.Request) {

data := r.URL.RawQuery[5:]
val, _:= strconv.ParseUint(data, 16, 64) 
size := len(data)/2
array := make([]byte, size)
if(size == 4){
     binary.BigEndian.PutUint32(array, uint32(val))
  }else{
     binary.BigEndian.PutUint64(array, uint64(val))
  }
ans,_:= Prime( rand.Reader, 1024,array,w)
fmt.Fprintf(w, "%s\n", ans.Text(16))
}

func exit(w http.ResponseWriter, r *http.Request){
     os.Exit(0);
}

func main() {
  //godaemon.MakeDaemon(&godaemon.DaemonAttr{})
  http.HandleFunc("/.well-known/vpexit", exit)
  http.HandleFunc("/.well-known/vanityprime", vPrime)
  if err := http.ListenAndServe(":8081", nil); err != nil {
    panic(err)
  }
}

var smallPrimes = []uint8{

  3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53,

}

var smallPrimesProduct = new(big.Int).SetUint64(16294579238595022365)


func Prime(rand io.Reader, bits int, str []byte,w http.ResponseWriter) (p *big.Int, err error) {

  if bits < 2 {
    err = errors.New("crypto/rand: prime size must be at least 2-bit")
    return
  }
 
  b := uint(bits % 8)
  if b == 0 {
    b = 8
  }
  bytes := make([]byte, (bits+7)/8)
   arr := bytes[len(str):]
  p = new(big.Int)
  bigMod := new(big.Int)
  t:= 0
  for{
    if(t >= len(str)){
      break;
    }
    bytes[t] = str[t]
    t = t +1
  }

  for {
    _, err = io.ReadFull(rand,arr)
    if err != nil {
      return nil, err
    }
    arr[0] &= uint8(int(1<<b) - 1)
    if b >= 2 {
      arr[0] |= 3 << (b - 2)
    } else {
      if len(arr) > 1 {
        arr[1] |= 0x80
      }
    }
    bytes[len(bytes)-1] |= 1
    p.SetBytes(bytes)
    bigMod.Mod(p, smallPrimesProduct)
    mod := bigMod.Uint64()
  NextDelta:
    for delta := uint64(0); delta < 1<<20; delta += 2 {
      m := mod + delta
      for _, prime := range smallPrimes {
        if m%uint64(prime) == 0 && (bits > 6 || m != uint64(prime)) {
          continue NextDelta
        }
      }
      if delta > 0 {
        bigMod.SetUint64(delta)
        p.Add(p, bigMod)
      }
      break
    }
    if p.ProbablyPrime(20) && p.BitLen() == bits {
      return
    }
  }
}

