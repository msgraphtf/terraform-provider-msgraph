
# Terraform Provider for Microsoft Graph

The MSGraph Terraform Provider allows managing resources within Microsoft Graph.

> [!WARNING]
> This plugin is in early alpha! Functionality is incomplete, and may be broken. Expect bugs, and please report them!
> Releases may contain undocumented breaking changes, and no commitment to backwards compatibility is provided at this time.

## Usage Example

Please see the documentation in the [Terraform Registry](https://registry.terraform.io/providers/msgraphtf/msgraph/latest/docs) for the full usage and features.

## Known issues or missing features

- All resources and attributes are using a custom plan modifier which may not be suitable for all things.
- Not using write-only attributes where applicable.
- Lacks automated testing.
- Lacks extensive logging for troubleshooting.
- Error messages and handling can probably be improved.
- Some resources and data sources have had to be manually adjusted (`generate/augment`), such as by excluding properties. Removing/fixing this will simplify development.
- Documentation lacks examples.

## Requesting new endpoints
Not all Microsoft Graph endpoints are supported (yet). This is simply so that I don't get overwhelmed with issue reports or failed builds for the hundreds of endpoints that the Microsoft Graph API supports. So I am slowly adding new endpoints as the project matures.

That said, if you have a use case for a particular endpoint, then I am committed to adding new endpoints upon request. Please open an issue, requesting the endpoint you want, and I will add it as soon as I can. And please provide feedback and bug reports for it.
