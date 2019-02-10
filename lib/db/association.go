package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

// Association represents relationships between Resource, includes HasMany, HasOne and BelongsTo association.
type Association interface {
	onCreate(tran Transaction, resource Resource)
	onDelete(tran Transaction, resource Resource)
}

type HasManyAssociation interface {
	FetchAll(resources interface{}) error
	Attach(tran Transaction, resource Resource)
	Detach(tran Transaction, resource Resource)

	onCreate(tran Transaction, resource Resource)
	onDelete(tran Transaction, resource Resource)
}

type HasOneAssociation interface {
	Fetch(resource Resource) error

	onCreate(tran Transaction, resource Resource)
	onDelete(tran Transaction, resource Resource)
}

type BelongsToAssociation interface {
	Fetch(resource Resource) error

	onCreate(tran Transaction, resource Resource)
	onDelete(tran Transaction, resource Resource)
}

func HasManyThrough(key string) HasManyAssociation {
	return &hasManyThroughAssociation{
		relatedKey: key,
	}
}

func HasManyPrefix(prefix string) HasManyAssociation {
	return &hasManyPrefixAssociation{
		relatedPrefix: prefix,
	}
}

func HasOne(relatedKey string) HasOneAssociation {
	return &hasOneAssociation{
		relatedKey: relatedKey,
	}
}

func BelongsTo(relatedKey string, belongToKeys ...string) BelongsToAssociation {
	if len(belongToKeys) == 0 {
		return &belongsToAssociation{
			relatedKey: relatedKey,
		}
	} else if len(belongToKeys) == 1 {
		return &belongsToAssociation{
			relatedKey:  relatedKey,
			belongToKey: belongToKeys[0],
		}
	} else {
		panic("belongToKeys can only have one element")
	}
}

type hasManyThroughAssociation struct {
	relatedKey    string
	cascadeDelete bool
}

func (assoc *hasManyThroughAssociation) FetchAll(resources interface{}) error {
	resourceKeys := []string{}

	kv, err := FetchJSON(assoc.relatedKey, &resourceKeys)

	if err != nil {
		return err
	} else if kv == nil {
		return nil
	}

	resourceBytesList := [][]byte{}

	for _, resourceKey := range resourceKeys {
		resp, err := client.Get(context.Background(), resourceKey)

		if err != nil {
			return err
		}

		if resp.Count > 0 {
			resourceBytesList = append(resourceBytesList, resp.Kvs[0].Value)
		} else {
			fmt.Printf("Association.FetchList: `%s` is missing\n", resourceKey)
		}
	}

	resourcesBytes := bytes.Buffer{}
	resourcesBytes.Write([]byte("["))
	resourcesBytes.Write(bytes.Join(resourceBytesList, []byte(",")))
	resourcesBytes.Write([]byte("]"))

	return json.Unmarshal(resourcesBytes.Bytes(), resources)
}

func (assoc *hasManyThroughAssociation) onCreate(tran Transaction, resource Resource) {
}

func (assoc *hasManyThroughAssociation) onDelete(tran Transaction, resource Resource) {
	tran.DeleteKey(assoc.relatedKey)
}

func (assoc *hasManyThroughAssociation) Attach(tran Transaction, resource Resource) {
	tran.AddToStringSet(assoc.relatedKey, resource.ResourceKey())
}

func (assoc *hasManyThroughAssociation) Detach(tran Transaction, resource Resource) {
	tran.PullfromStringSet(assoc.relatedKey, resource.ResourceKey())
}

type hasManyPrefixAssociation struct {
	relatedPrefix string
	cascadeDelete bool
}

func (assoc *hasManyPrefixAssociation) FetchAll(resources interface{}) error {
	return FetchAllFrom(assoc.relatedPrefix, resources, strings.Count(assoc.relatedPrefix, "/"))
}

func (assoc *hasManyPrefixAssociation) onCreate(tran Transaction, resource Resource) {
}

func (assoc *hasManyPrefixAssociation) onDelete(tran Transaction, resource Resource) {
	tran.DeletePrefix(assoc.relatedPrefix)
}

func (assoc *hasManyPrefixAssociation) Attach(tran Transaction, resource Resource) {
	panic("can not attach to HasManyPrefix association")
}

func (assoc *hasManyPrefixAssociation) Detach(tran Transaction, resource Resource) {
	panic("can not detach from HasManyPrefix association")
}

type hasOneAssociation struct {
	relatedKey string
}

func (assoc *hasOneAssociation) Fetch(resource Resource) error {
	return FetchFrom(assoc.relatedKey, resource)
}

func (assoc *hasOneAssociation) onCreate(tran Transaction, resource Resource) {
}

func (assoc *hasOneAssociation) onDelete(tran Transaction, resource Resource) {
}

type belongsToAssociation struct {
	relatedKey  string
	belongToKey string
}

func (assoc *belongsToAssociation) Fetch(resource Resource) error {
	return FetchFrom(assoc.relatedKey, resource)
}

func (assoc *belongsToAssociation) onCreate(tran Transaction, resource Resource) {
	if assoc.belongToKey != "" {
		tran.AddToStringSet(assoc.belongToKey, resource.ResourceKey())
	}
}

func (assoc *belongsToAssociation) onDelete(tran Transaction, resource Resource) {
	if assoc.belongToKey != "" {
		tran.PullfromStringSet(assoc.belongToKey, resource.ResourceKey())
	}
}
