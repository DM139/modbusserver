package main


import "log"
type People struct{
}
//定义一个类型
func (p *People)PrePing()  {
	log.Println("pre ping")
}
func (p *People)Ping()  {
	log.Println("ping")
}
//定义另一个类型，继承上面的并覆盖某些方法
type GoodGuys struct{
	*People
	i int
}
func (g *GoodGuys)Ping()  {
	log.Println("pong")
}
func main(){
	p := &People{}
	i := 1
	g:=&GoodGuys{p, i}
	g.PrePing()
	g.Ping()
}
