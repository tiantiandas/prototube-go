package prototube

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProducer(t *testing.T) {
	p := &Producer{}
	uuid, err := uuid.NewUUID()
	assert.Nil(t, err)
	uuidBytes, err := uuid.MarshalBinary()
	assert.Nil(t, err)
	bytes, err := p.encode(1234, uuidBytes, &PrototubeMessageHeader{})
	assert.Nil(t, err)
	assert.Equal(t, prototubeMessageHeader, bytes[:len(prototubeMessageHeader)])
}
