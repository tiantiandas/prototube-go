package prototube

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	p := &Producer{}
	uuid, err := uuid.NewUUID()
	uuidBytes, err := uuid.MarshalBinary()
	bytes, err := p.encode(1234, uuidBytes, &PrototubeMessageHeader{})
	assert.Nil(t, err)
	assert.Equal(t, prototubeMessageHeader, bytes[:len(prototubeMessageHeader)])
}
