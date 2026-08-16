package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/misonikomipan/homebox-cli/internal/client"
)

// entityBasePath is the v0.26 "entities" endpoint. Homebox v0.26 merged the
// previous /v1/items and /v1/locations resources into a single entity resource.
const entityBasePath = "/v1/entities"

// locationEntityTypeID returns the entity type marked isLocation=true (a
// "location" type) for the current group. Locations must be created with an
// entity type whose isLocation flag is true.
func locationEntityTypeID(c *client.Client) (string, error) {
	data, err := c.Get("/v1/entity-types", nil)
	if err != nil {
		return "", err
	}
	var types []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		IsLocation bool   `json:"isLocation"`
	}
	if err := json.Unmarshal(data, &types); err != nil {
		return "", err
	}
	for _, t := range types {
		if t.IsLocation {
			return t.ID, nil
		}
	}
	return "", fmt.Errorf("no location entity type found for this group")
}

// fetchEntity decodes the JSON response of GET /v1/entities/{id}.
func fetchEntity(c *client.Client, id string) (map[string]any, error) {
	data, err := c.Get(entityBasePath+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// updatePayload converts a GET /v1/entities/{id} response into a payload the
// PUT /v1/entities/{id} endpoint accepts. Nested objects (parent, entityType)
// and the tags edge are flattened into parentId / entityTypeId / tagIds, and
// read-only fields are removed. Unknown extra keys are ignored by the server.
func updatePayload(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}

	if p, ok := out["parent"].(map[string]any); ok {
		if id, ok := p["id"].(string); ok {
			out["parentId"] = id
		}
	} else {
		out["parentId"] = ""
	}
	delete(out, "parent")

	if et, ok := out["entityType"].(map[string]any); ok {
		if id, ok := et["id"].(string); ok {
			out["entityTypeId"] = id
		}
	} else {
		out["entityTypeId"] = ""
	}
	delete(out, "entityType")

	if tags, ok := out["tags"].([]any); ok {
		ids := make([]string, 0, len(tags))
		for _, tg := range tags {
			if t, ok := tg.(map[string]any); ok {
				if id, ok := t["id"].(string); ok {
					ids = append(ids, id)
				}
			}
		}
		out["tagIds"] = ids
	}
	delete(out, "tags")

	for _, k := range []string{"createdAt", "updatedAt", "imageId", "thumbnailId", "children", "totalPrice", "itemCount", "attachments"} {
		delete(out, k)
	}
	return out
}

// putEntity performs PUT /v1/entities/{id} with an EntityUpdate payload.
func putEntity(c *client.Client, id string, payload map[string]any) ([]byte, error) {
	return c.Put(entityBasePath+"/"+id, payload)
}

// parentName returns the name of an entity's parent, or "".
func parentName(m map[string]any) string {
	if p, ok := m["parent"].(map[string]any); ok {
		if n, ok := p["name"].(string); ok {
			return n
		}
	}
	return ""
}

// stringField returns a string field from a map, or "".
func stringField(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
