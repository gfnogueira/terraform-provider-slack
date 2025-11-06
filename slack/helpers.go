package slack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// convertSchemaSetToStringSlice safely converts a *schema.Set to []string
func convertSchemaSetToStringSlice(set *schema.Set) []string {
	result := make([]string, 0, set.Len())
	for _, v := range set.List() {
		result = append(result, v.(string))
	}
	return result
}

// convertStringSliceToInterface converts []string to []interface{} for schema.Set
func convertStringSliceToInterface(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
