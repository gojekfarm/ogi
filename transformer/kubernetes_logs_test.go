package ogitransformer

import (
	"log"
	"regexp"
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"

	ogiproducer "github.com/gojektech/ogi/producer"
)

func TestKubernetesKafkaLogTransformSuccess(t *testing.T) {
	kubernetesKafkaLogMessage := `{
		"message": "what-message",
		"stream": "somestream",
		"log": "this is log",
		"docker": {
			"container_id": "dckr-123"
		},
		"kubernetes": {
			"container_name": "svc-123",
			"namespace_name": "backend",
			"pod_name": "svc-123-abc",
			"pod_id": "12345",
			"labels": {
				"app": "mysvc",
				"stack": "rails"
			},
			"host": "1.2.3.4",
			"master_url": "some-url"
		}
	}`

	kubernetesKafkaLogMessageWithKey := `{
		"message": "what-message",
		"stream": "somestream",
		"log": "this is log",
		"docker": {
			"container_id": "dckr-123"
		},
		"kubernetes": {
			"container_name": "svc-123",
			"namespace_name": "backend",
			"pod_name": "svc-123-abc",
			"pod_id": "12345",
			"labels": {
				"app": "mysvc",
				"stack": "rails"
			},
			"host": "1.2.3.4",
			"master_url": "some-url"
		},
		"message_key": "svc-123-abc"
	}`

	reg, err := regexp.Compile("[^a-zA-Z0-9{}:\"]+")
	if err != nil {
		log.Fatal(err)
	}

	var guard *monkey.PatchGuard
	var guardB bool
	var guardMessage, guardKey, guardTopic string
	guard = monkey.Patch(ogiproducer.Produce, func(topic string, message []byte, message_key string) {
		guard.Unpatch()
		defer guard.Restore()
		guardB = true
		guardMessage, guardKey, guardTopic = string(message), message_key, topic
		return
	})

	kkl := KubernetesKafkaLog{}
	err = kkl.Transform([]byte(kubernetesKafkaLogMessage))
	assert.Nil(t, err)

	actual := reg.ReplaceAllString(string(guardMessage), "")
	expected := reg.ReplaceAllString(kubernetesKafkaLogMessageWithKey, "")
	assert.Equal(t, expected, actual)

	assert.Equal(t, "svc-123-abc", guardKey)
	assert.Equal(t, "mysvc", guardTopic)
}

func TestKubernetesKafkaLogTransformUnmarshallError(t *testing.T) {
	kubernetesKafkaLogMessageBad := `{
		"message": "what-message",
		"some": "thing"
	}`

	kkl := KubernetesKafkaLog{}
	err := kkl.Transform([]byte(kubernetesKafkaLogMessageBad))
	assert.NotNil(t, err)
}

func TestKubernetesKafkaLogTransformNoTopicLabel(t *testing.T) {
	kubernetesKafkaLogMessageWithoutTopic := `{
		"message": "what-message",
		"stream": "somestream",
		"log": "this is log",
		"docker": {
			"container_id": "dckr-123"
		},
		"kubernetes": {
			"container_name": "svc-123",
			"namespace_name": "backend",
			"pod_name": "svc-123-abc",
			"pod_id": "12345",
			"labels": {
				"stack": "rails"
			},
			"host": "1.2.3.4",
			"master_url": "some-url"
		}
	}`

	kkl := KubernetesKafkaLog{}
	err := kkl.Transform([]byte(kubernetesKafkaLogMessageWithoutTopic))
	assert.NotNil(t, err)
}
