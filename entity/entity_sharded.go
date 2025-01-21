package entity

type EntitySharded interface {
	Entity
	ShardKey() string
}

type EntityShardedFactory func(shardKey string) EntitySharded

type BaseEntitySharded struct {
	shardKey string
}

func (e *BaseEntitySharded) ID() string { return "" }

func (e *BaseEntitySharded) TABLE() string { return "" }

func (e *BaseEntitySharded) NAME() string { return "" }

func (e *BaseEntitySharded) KEY() string { return "" }

func (e *BaseEntitySharded) ShardKey() string { return e.shardKey }
