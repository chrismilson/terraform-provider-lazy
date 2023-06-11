package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLazyStringResource_basic(t *testing.T) {
	dsn := "lazy_string.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = "explicit_value"
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsn, "id"),
					resource.TestCheckResourceAttrSet(dsn, "last_updated"),
					resource.TestCheckResourceAttr(dsn, "initially", "initial_value"),
					resource.TestCheckResourceAttr(dsn, "explicitly", "explicit_value"),
					resource.TestCheckResourceAttr(dsn, "result", "explicit_value"),
				),
			},
		},
	})
}

func TestAccLazyStringResource_InitialValue_basic(t *testing.T) {
	dsn := "lazy_string.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "initial_value"),
				),
			},
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = null
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "initial_value"),
				),
			},
		},
	})
}

func TestAccLazyStringResource_ExplicitValue_SetOnCreate(t *testing.T) {
	dsn := "lazy_string.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = "explicit_value"
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "explicit_value"),
				),
			},
		},
	})
}

func TestAccLazyStringResource_ExplicitValue_SetOnUpdate(t *testing.T) {
	dsn := "lazy_string.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = null
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "initial_value"),
				),
			},
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = "explicit_value"
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "explicit_value"),
				),
			},
		},
	})
}

func TestAccLazyStringResource_ExplicitValue_Unset_NoChanges(t *testing.T) {
	dsn := "lazy_string.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = "explicit_value"
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "explicit_value"),
				),
			},
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = null
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "explicit_value"),
					resource.TestCheckResourceAttr(dsn, "explicitly", "explicit_value"),
				),
			},
		},
	})
}

func TestAccLazyStringResource_Import_basic(t *testing.T) {
	dsn := "lazy_string.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = "explicit_value"
        }`,
				ImportState:        true,
				ResourceName:       dsn,
				ImportStateId:      "imported_value",
				ImportStatePersist: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "imported_value"),
					resource.TestCheckResourceAttr(dsn, "explicitly", "imported_value"),
					resource.TestCheckResourceAttr(dsn, "initially", "imported_value"),
				),
			},
		},
	})
}

func TestAccLazyStringResource_Import_Keep(t *testing.T) {
	dsn := "lazy_string.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = "explicit_value"
        }`,
				ImportState:        true,
				ResourceName:       dsn,
				ImportStateId:      "imported_value",
				ImportStatePersist: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "imported_value"),
					resource.TestCheckResourceAttr(dsn, "explicitly", "imported_value"),
					resource.TestCheckResourceAttr(dsn, "initially", "imported_value"),
				),
			},
			{
				Config: `
        resource "lazy_string" "test" {
          initially  = "initial_value"
          explicitly = null
        }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "result", "imported_value"),
					resource.TestCheckResourceAttr(dsn, "explicitly", "imported_value"),
					resource.TestCheckResourceAttr(dsn, "initially", "initial_value"),
				),
			},
		},
	})
}
