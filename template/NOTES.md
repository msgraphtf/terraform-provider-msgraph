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

