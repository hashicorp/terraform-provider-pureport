package tags

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func TagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

func FilterTags(tags map[string]interface{}) (out map[string]string) {

	if out == nil {
		out = map[string]string{}
	}

	for k, v := range tags {
		switch v.(type) {
		case string:
			out[k] = v.(string)
		default:
			// Do nothing
		}
	}

	return
}
