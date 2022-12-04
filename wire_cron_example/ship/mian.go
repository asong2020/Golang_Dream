//+build wireinject

package main

import (
	"github.com/google/wire"
)

type Ship struct {
	pulp *Pulp
}
func NewShip(pulp *Pulp) *Ship {
	return &Ship{
		pulp: pulp,
	}
}
type Pulp struct {
	count int
}

func NewPulp() *Pulp {
	return &Pulp{
	}
}

func (c *Pulp)set(count int)  {
	c.count = count
}

func (c *Pulp)get() int {
	return c.count
}

func InitShip() *Ship {
	wire.Build(
		NewPulp,
		NewShip,
		)
	return &Ship{}
}

func main(){

}

