package ag

import (
	"encoding/json"
	"fmt"
)

// Entity is an interface used by ScanEntity. Any struct that has Vertex or
// Edge as its embedded field and implements EntitySaver can be an entity for
// vertex or edge.
type Entity interface {
	entityReader
	EntitySaver
}

type entityReader interface {
	readEntity(b []byte) (*entityData, error)
}

type entityData struct {
	core       interface{}
	properties []byte
}

// EntitySaver is an interface used by ScanEntity.
type EntitySaver interface {
	// SaveEntity assigns an entity from the database driver.
	//
	// valid is true if the entity is not NULL.
	//
	// core is VertexCore or EdgeCore that will be stored in the entity for
	// vertex or edge respectively. If valid is false, core will be nil.
	//
	// An error should be returned if the entity cannot be stored without
	// loss of information.
	SaveEntity(valid bool, core interface{}) error
}

// PropertiesSaver is an interface used by ScanEntity.
type PropertiesSaver interface {
	// By default, properties of an entity read by ScanEntity are stored in
	// the entity itself by calling json.Unmarshal over it. To modify this
	// default behavior, one may implement PropertiesSaver for the entity.
	//
	// The underlying array of b may be reused.
	//
	// An error should be returned if the properties cannot be stored
	// without loss of information.
	SaveProperties(b []byte) error
}

// ScanEntity reads an entity for vertex or edge from src and stores the result
// in the given entity.
//
// An error will be returned if the type of src is not []byte, or src is
// invalid for the given entity.
func ScanEntity(src interface{}, entity Entity) error {
	switch src := src.(type) {
	case []byte:
		if len(src) < 1 {
			return fmt.Errorf("invalid source for entity: %v", src)
		}
		d, err := entity.readEntity(src)
		if err != nil {
			return err
		}
		return saveEntityData(d, entity)
	case *entityData:
		return saveEntityData(src, entity)
	case nil:
		return entity.SaveEntity(false, nil)
	default:
		return fmt.Errorf("invalid source for entity: %T", src)
	}
}

func saveEntityData(d *entityData, entity Entity) error {
	if d == nil {
		panic("invalid entity data: nil")
	}

	err := entity.SaveEntity(true, d.core)
	if err != nil {
		return err
	}

	if p, ok := entity.(PropertiesSaver); ok {
		err = p.SaveProperties(d.properties)
	} else {
		err = json.Unmarshal(d.properties, entity)
	}
	return err
}
