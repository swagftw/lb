package backend

type Backend struct {
    IP               string
    TotalConnections int
    Alive            bool
}

func (b *Backend) IsAlive() bool {
    return b.Alive
}

func (b *Backend) SetAlive(alive bool) {
    b.Alive = alive
}

func (b *Backend) Ping() {

}
