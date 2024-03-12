package main

type Handler func (a int)

func xc(pa int, handler Handler) {
  handler(pa)
}

func main(){
  xc(123, func(a int){
	  print (a) //=> 123
  })
}
