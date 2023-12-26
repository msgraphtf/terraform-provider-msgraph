# TODO

- Support validation for enums
- Support pagination for collection results
- User Required instead of Optional for attributes where possible

# Mapping OpenAPI Concepts to Terraform

## Open API schema object/property types -> Terraform Attribute types

| OpenAPI                    | Terraform             |
| -------------------------- | --------------------- |
| string                     | StringAttribute       |
| integer                    | Int64Attribute        |
| boolean                    | BoolAttribute         |
| object                     | SingleNestedAttribute |
| array (of primitive types) | ListAttribute         |
| array (of objects)         | ListNestedAttribute   |
| string (enum)              | StringAttribute (with validation, can only be one of a set of values |

