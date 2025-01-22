package entity

type EntitySharded interface {
	Entity
	ShardKey() string
	SetShardKey(string)
}

type EntityShardedFactory func(shardKey string) EntitySharded

type BaseEntitySharded struct {
	BaseEntity
	shardKey string
}

func (b *BaseEntitySharded) ShardKey() string {
	return b.shardKey
}

func (b *BaseEntitySharded) SetShardKey(sk string) {
	b.shardKey = sk
}
