package pureport

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestCheckResourceConnectionIdChanged(start *string, end *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if start == end {
			return fmt.Errorf("ID was not updated so connection was not recreated: id=%s", *start)
		}
		return nil
	}
}
