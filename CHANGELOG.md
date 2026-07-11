# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- (placeholder)

### Changed
- (placeholder)

### Fixed
- (placeholder)

### Security
- (placeholder)

## [v1.1.0] - 2026-07-11

Fork release (Maestroff/azurecaf registry versioning).

### Added
- 14 new resource types for parity with azurerm v4.80:
  - `azurerm_windows_function_app` (slug `winfa`), `azurerm_windows_function_app_slot` (`winfas`), `azurerm_windows_web_app_slot` (`wwapps`)
  - `azurerm_managed_redis` (`amr`) — Azure Managed Redis
  - `azurerm_managed_devops_pool` (`mdp`) — Managed DevOps Pools
  - `azurerm_trusted_signing_account` (`tsa`) — Trusted Signing
  - `azurerm_eventgrid_partner_namespace` (`egpn`), `azurerm_eventgrid_partner_registration` (`egpr`), `azurerm_eventgrid_namespace_topic` (`egnt`)
  - `azurerm_mongo_cluster_firewall_rule` (`mongofw`), `azurerm_mongo_cluster_user` (`mongou`)
  - `azurerm_network_security_perimeter_profile` (`nspp`), `azurerm_network_security_perimeter_access_rule` (`nspar`)
  - `azurerm_cognitive_account_project` (`proj`) — AI Foundry project (new-style resource)
- Unit tests covering all new resource types (`new_resource_types_test.go`)
- Regenerated `docs/resource_types.md` and `docs/azure_resources.md`

## [v4.0.0] - 2026-03-06 - MAJOR RELEASE

### 🔴 BREAKING CHANGES
- **Default Slug Behavior Changed**: `use_legacy_slug` now defaults to `false` (was `true` in v3.x)
  - **Impact**: Resources will use modern Microsoft CAF slugs by default
  - `azurerm_mssql_database`: Uses `sqldb` slug (was `database`)
  - `azurerm_mssql_elasticpool`: Uses `sqlep` slug (was `elasticpool`)
  - **Action Required**: Add `use_legacy_slug = true` to existing resource configurations to maintain current naming

### Migration Path from v3.x to v4.0.0

#### For Existing Deployments (Maintain Current Names)
```terraform
# Add to ALL existing mssql_database and mssql_elasticpool resources
resource "azurecaf_name" "db" {
  resource_type     = "azurerm_mssql_database"
  name              = "mydb"
  use_legacy_slug   = true    # ← ADD THIS to maintain "database" slug
}
```

#### For New Deployments (Adopt Modern CAF Slugs)
```terraform
# Default behavior - no changes needed
resource "azurecaf_name" "db" {
  resource_type     = "azurerm_mssql_database"
  name              = "mydb"
  # use_legacy_slug = false (default) → uses "sqldb" slug
}
```

#### Gradual Migration Strategy
1. **Immediate**: Add `use_legacy_slug = true` to all existing resources in Terraform states
2. **Planning**: Run `terraform plan` to confirm no unexpected changes
3. **Selective Migration**: Remove `use_legacy_slug = true` from non-critical resources first
4. **Validation**: Test name changes in development/staging environments
5. **Production**: Schedule maintenance window for database renames if needed

### Alternative: Prevent Name Changes
If you want to freeze existing names permanently:
```terraform
resource "azurecaf_name" "db" {
  passthrough   = true
  name          = "existing-database-name"  # Exact current name
}
```

### What's Improved in v4.0.0
- Modern CAF slugs as default align with official Microsoft standards
- Consistent naming for new projects
- Legacy support still available via explicit opt-in
- Clearer migration path for future slug updates

### Rollback Strategy
If issues occur after upgrading:
1. Pin provider version to v3.1.0: `version = "~> 3.1.0"`
2. Or add `use_legacy_slug = true` to affected resources
3. Run `terraform plan` to verify rollback

## [v3.1.0] - 2026-03-06

### Added
- **Backward Compatibility Feature**: Added `use_legacy_slug` parameter to `azurecaf_name` resource and data source
  - Allows resources to maintain legacy slug naming from v1.x/v2.x versions
  - `azurerm_mssql_database`: Can use legacy `database` slug instead of `sqldb`
  - `azurerm_mssql_elasticpool`: Can use legacy `elasticpool` slug instead of `sqlep`
  - **Default**: `use_legacy_slug = true` (maintains backward compatibility)
  - Added `legacy_slug` field to `resourceDefinition.json` for affected resources

### Changed
- **Provider Migration Strategy**: Established deprecation path for legacy slugs
  - v3.x: `use_legacy_slug` defaults to `true` (current)
  - v4.0.0: `use_legacy_slug` will default to `false` (future major version)

### Migration Guidance
- **Existing Deployments**: No action required - legacy slugs remain active by default
- **New Deployments**: Set `use_legacy_slug = false` to adopt modern CAF slugs
- **Gradual Migration**: Can be controlled per-resource:
  ```terraform
  resource "azurecaf_name" "db" {
    resource_type     = "azurerm_mssql_database"
    use_legacy_slug   = true   # Maintains "database" slug
  }
  ```
- **Preparing for v4.0.0**: Explicitly set `use_legacy_slug = true` on existing resources to ensure stability across version upgrades

### Documentation
- See migration guidance in this changelog under `v4.0.0` → `Migration Path from v3.x to v4.0.0`
- For detailed v3.0.0 to v4.0.0 migration planning: Check CHANGELOG v4.0.0 notes

## [v3.0.0] - 2025-10-30

### Fixed
- **Microsoft CAF Slug Compliance**: Updated resource slugs to match official Microsoft Cloud Adoption Framework recommendations
  - `azurerm_mssql_database`: Changed slug from `database` to `sqldb` (official Microsoft CAF abbreviation)
  - `azurerm_mssql_elasticpool`: Changed slug from `elasticpool` to `sqlep` (official Microsoft CAF abbreviation)
  - Added official resource provider namespaces to resource definitions for better compliance tracking
  - All changes verified against https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations

### Changed
- **AI Agent Documentation**: Added comprehensive `.github/copilot-instructions.md` with project-specific patterns, architecture, and workflows
  - Documented core components and data flows (resourceDefinition.json → gen.go → models_generated.go)
  - Added detailed development workflows with Makefile targets and testing hierarchy
  - Included resource definition schema and slug validation process
  - Documented critical conventions: data source vs resource usage, regex escaping, test organization

### Impact Assessment
- **Potential Breaking Behavior**: The slug changes for `azurerm_mssql_database` (`database` → `sqldb`) and `azurerm_mssql_elasticpool` (`elasticpool` → `sqlep`) modify the generated name when your configuration uses `use_slug = true` (default) in `azurecaf_name` or `azurecaf_naming_convention`. This can trigger a resource replacement (destroy/create) in Terraform due to the name change.
- **Existing Resources**: Previously deployed resources keep their current name as long as you don't reprovision or change inputs that regenerate it. A diff appears only if the name is recalculated and differs from the existing one.
- **Naming Convention**: New names use official Microsoft CAF abbreviations; improves alignment but may alter prefixes.
- **State Impact**: `terraform plan` will show a change for the `name` attribute and mark recreation for resources whose name derives directly from `azurecaf_name.*.result`.
- **Migration Guidance (Options)**:
  1. Preserve existing names: Set `passthrough = true` and `name = "<current_name>"` to validate without regenerating.
  2. Remove slug influence: Use `use_slug = false` to avoid the slug change and keep the stable main portion.
  3. Pin provider version: Lock the previous version temporarily if you need time (`version = "~> 2.0.3"`).
  4. Controlled migration: Create new resources with CAF-compliant names and migrate data/manual steps before destroying old ones.
  5. Prevent immediate destroy: Add `lifecycle { prevent_destroy = true }` short-term while executing migration steps.
  6. Progressive rename: Apply first in non-prod; validate impact on dependencies (backups, alerts, connection strings) before production.
  7. Manual import (if needed): If adopting a new name without destruction, create resource with new name, migrate, then remove the old after verification.
- **Recommended Approach**: For future CAF consistency, adopt the new slugs and schedule a maintenance window for databases / elastic pools if the name changes.
- **Automation Scripts**: Update any external script (CI/CD, monitoring) that derives names via slug.
- **Rollback Strategy**: If issues arise, revert to the previous provider version and use passthrough to freeze names; then retry with a staged transition strategy.


## [v2.0.3] - 2025-01-27

### Fixed
- **GoReleaser Configuration**: Fixed terraform registry manifest file handling in release process
- **Release Assets**: Corrected missing manifest.json asset generation for Terraform Registry compatibility

## [v2.0.0]

### Fixed
- **CAF Prefix Compliance**: Fixed 'general' and 'general_safe' resources to use valid CAF slugs ('gen' and 'gensafe') and updated their official slug fields to pass CAF prefix validation tests.
- **CRITICAL**: Achieved 100% CAF Compliance with zero duplicates
  - **Load Balancer**: Updated `azurerm_lb` slug from `lb` to `lbe` (external load balancer default)
  - **Logic App Standard**: Updated `azurerm_logic_app_standard` slug from `logic` to `logicstd` to differentiate from workflow
  - **SQL Server Legacy**: Updated `azurerm_sql_server` slug from `sql` to `sqlsrv` to differentiate from `azurerm_mssql_server`
  - **Perfect Optimization**: Eliminated all remaining duplicate slugs (0 duplicates from 2 groups)
  - **CAF Validation Scripts**: Updated validation logic to properly handle differentiated resources
  - **Workflow Integration**: Enhanced CAF correction scripts to maintain compliance automatically
- **CAF Abbreviations Compliance**: Updated resource definitions to follow official Microsoft Cloud Adoption Framework abbreviation recommendations
  - Fixed Front Door (classic) to use official CAF abbreviation `afd` instead of `fd`
  - Updated all Azure Automation resources to use official CAF abbreviation `aa` instead of `automati`
  - Fixed App Configuration resources to use official CAF abbreviation `appcs` instead of `appconfi`
  - Updated Active Directory resources to use official CAF abbreviation `aad` instead of `activedi`
  - Added missing official slug mappings for automation resources that lacked them
  - Comprehensive test coverage added to ensure ongoing CAF compliance

### Added
- **MAJOR**: Authoritative Microsoft CAF Compliance System with real-time corrections
  - Complete 154-resource official CAF mapping from Microsoft Learn documentation
  - Authoritative workflow that enforces official CAF abbreviations with absolute priority
  - Intelligent duplicate resolution algorithm with meaningful slug generation
  - Automatic backup creation before applying corrections
  - Real-time CAF compliance validation and correction pipeline
  - Enhanced conflict resolution with semantic slug naming (webapp, ch, gw, cd, func, etc.)
  - Comprehensive reporting system with detailed compliance metrics
- **Enhanced Azure Resource Sync Tool v4.0**: Complete workflow automation
  - Integrated CAF validation, duplicate detection, and optimization
  - Streamlined workflow commands: sync, test, docs, analyze, status, auto-add, cleanup
  - Hashicorp Registry deep linking with improved resource extraction
  - Automated missing resource detection and prioritized addition
  - Comprehensive documentation generation with Azure-specific categorization
- **CAF Compliance Automation**: Production-ready validation and correction system
  - Automated application of 19+ CAF corrections per sync cycle
  - Reduced duplicate slugs from 10 groups to 1 (90% improvement)
  - Increased unique slugs from 1,170 to 1,184 (14 additional unique resources)
  - Idempotent operations ensuring system stability across multiple executions
- **CAF Abbreviations Test Suite**: Comprehensive test coverage for CAF compliance
  - `TestResourceAbbreviationsCAFCompliance`: Tests specific resources for correct abbreviations
  - `TestCAFOfficialAbbreviations`: Validates all automation resources use correct CAF abbreviations
  - `TestAppConfigurationCAFAbbreviation`: Validates App Configuration resources use correct abbreviations

### Changed
- **BREAKING**: Enhanced sync workflow to apply permanent CAF corrections
  - `enhanced_sync_official_resources_caf.sh` now applies and saves corrections directly
  - Eliminated audit-only mode in favor of authoritative correction system
  - Official CAF abbreviations now take absolute priority over existing slugs
  - Improved conflict resolution generates semantic suffixes instead of truncated names
- **Workflow Optimization**: Streamlined and consolidated automation scripts
  - Integrated Python optimization logic directly into bash workflow
  - Removed redundant automation scripts for simplified maintenance
  - Enhanced JQ-based processing for better performance and reliability
  - Improved error handling and rollback capabilities

### Fixed
- **Go Version Alignment**: Resolved conflicting Go version declarations in go.mod
  - Changed from conflicting `go 1.23.0` and `toolchain go1.24.4` to unified `go 1.24`
  - Eliminates version mismatch errors during builds
  - Ensures consistent Go toolchain usage across all environments
  - Impact: Medium - Fixes build reliability and development environment consistency
- **Linting Issues**: Fixed non-constant format string errors in logging and error handling
  - Fixed `fmt.Errorf` call in `resource_name.go` to use proper format string
  - Fixed `log.Printf` call in `resource_naming_convention.go` to use proper format string
  - Resolves Go vet warnings and ensures build passes all checks
  - Impact: Low - Improves code quality and eliminates build warnings
- **CAF Compliance**: Addressed systematic naming convention inconsistencies
  - Fixed Hashicorp resource extraction to prevent false azurerm_com_ artifacts
  - Improved regex filtering to exclude URLs and non-resource patterns
  - Enhanced resource validation to ensure Azure Resource Manager compatibility
  - Impact: High - Ensures reliable resource synchronization and CAF compliance

## [v1.2.30]

### Fixed
- **CI/CD Pipeline**: Fixed GoReleaser failure due to git tag mismatch and dirty state
  - Removed problematic auto-commit step that was creating commits during release process
  - Fixed generated file timestamp stability to prevent dirty git state in CI
  - Added `fetch-depth: 0` to GitHub Actions checkout for full git history
  - Stabilized `models_generated.go` timestamp format to be environment-independent
  - Resolves GoReleaser errors: "git tag was not made against commit" and "git is in a dirty state"
  - Impact: High - Fixes release automation and ensures reliable tag-based releases
- **GoReleaser Configuration**: Updated GoReleaser configuration to v2 format
  - Added `version: 2` to support GoReleaser v2.x
  - Changed `changelog.skip: true` to `changelog.disable: true`
  - Removed deprecated `archives.format` property to use automatic format selection
  - Fixes release pipeline compatibility with goreleaser-action@v6
- **GitHub Workflow**: Fixed workflow step ordering and improved GPG key import
  - Moved "Set up Go" step before "Install tfproviderlint" to resolve dependency issues
  - Enhanced GPG key import with additional configuration options
  - Added `continue-on-error: true` for GPG import to handle missing secrets gracefully
  - Improved Git signing configuration with proper trust levels
- **README Display**: Fixed GitHub repository homepage README display issue
  - Converted README.md line endings from Windows-style (CRLF) to Unix-style (LF)
  - Renamed .github/README.md to .github/README-workflows.md to prevent GitHub display conflict
  - Resolves issue where GitHub was showing workflows documentation instead of main project README
  - Ensures proper display of comprehensive project documentation on repository homepage
- **Code Generation**: Removed timestamp from generated `models_generated.go` file
  - Eliminated dynamic timestamp that was causing git dirty state during CI/CD
  - Removed `GeneratedTime` field from template data structure
  - Updated template to exclude timestamp comment from generated code
  - Impact: High - Prevents GoReleaser "git is in a dirty state" errors during releases
  - Resolves: CI builds no longer modify tracked files during generation process

### Security
- **CRITICAL**: Fixed security vulnerabilities in Go dependencies:
  - Updated `golang.org/x/net` from v0.23.0 to v0.38.0 to resolve:
    - GO-2025-3595: Cross-site Scripting vulnerability in html package
    - GO-2025-3503: HTTP Proxy bypass using IPv6 Zone IDs
    - GO-2024-3333: Non-linear parsing vulnerability in html package
  - Updated `golang.org/x/crypto` from v0.21.0 to v0.36.0
  - Updated `golang.org/x/sys` from v0.18.0 to v0.31.0
  - Updated `golang.org/x/text` from v0.14.0 to v0.23.0
- Updated Go toolchain from 1.20 to 1.23.0 with Go 1.24.4 for enhanced security
- **SECURITY**: Fixed loose POSIX file permissions in E2E testing framework:
  - Changed directory permissions from 0755 to 0750 (removed world access)
  - Changed executable file permissions from 0755 to 0750 (removed world access)
  - Affected files: `e2e/framework/e2e_test.go`, `e2e/framework/framework.go`

### Added
- **E2E Testing Infrastructure**: Complete end-to-end testing framework
  - Comprehensive test suite covering all provider functionality
  - Import functionality testing with real Terraform state management
  - Data source validation with cross-platform compatibility
  - Naming convention testing across multiple resource types
  - Multi-resource type testing for complex scenarios
- **CI/CD Integration**: Full GitHub Actions integration for automated testing
  - Quick E2E tests on every push (~10-15 seconds)
  - Full E2E tests on pull requests (~25-30 seconds) 
  - Manual workflow dispatch with selective test execution
  - Smart triggering based on file changes
- **Local CI Simulation**: Act integration for local CI environment testing
  - Complete workflow validation before pushing to GitHub
  - Docker-based CI environment simulation
  - Cross-platform testing (macOS M-series compatibility)
  - Comprehensive testing scripts for development workflow
- **Testing Scripts**: Production-ready testing automation
  - `scripts/complete-e2e-validation.sh` - Full validation pipeline
  - `scripts/quick-ci-test.sh` - Quick CI environment validation
  - `scripts/test-ci-with-act.sh` - Interactive CI simulation
  - `scripts/validate-ci-e2e.sh` - Enhanced local + CI testing
- **Documentation**: Complete testing and CI/CD documentation
  - `E2E_IMPLEMENTATION_SUMMARY.md` - Implementation overview
  - `ACT_TESTING_GUIDE.md` - Local CI testing guide
  - `CI_E2E_INTEGRATION.md` - CI/CD integration documentation
  - `e2e/README.md` - E2E testing framework documentation
- **GitHub Copilot Integration**: Enhanced development workflow automation
  - `copilot-setup-steps.yml` - GitHub Actions workflow for Copilot environment setup
  - Automated Go and Terraform environment configuration for Copilot sessions
  - Streamlined development environment preparation with proper versioning
- **MAJOR**: Comprehensive end-to-end (E2E) testing framework for real-world validation
  - Complete E2E test suite covering provider build → Terraform usage → Azure integration
  - Mock Azure RM provider integration for testing without actual Azure API calls
  - Automated provider compilation and local installation testing
  - Azure resource naming compliance validation for all supported resource types
  - Support for all naming conventions (cafclassic, cafrandom, random, passthrough)
  - Edge case testing including length limits, special characters, and error conditions
  - Integration testing with azurerm provider using mock scenarios
  - Test runner CLI with flexible execution options and debugging support
  - Comprehensive documentation and troubleshooting guides
- New Makefile targets for E2E testing:
  - `test_e2e` - Complete E2E test suite
  - `test_e2e_quick` - Fast E2E tests for CI/CD
  - `test_e2e_integration` - AzureRM integration tests
  - `test_e2e_naming` - Naming convention validation
  - `test_e2e_edge_cases` - Edge case scenarios
  - `test_e2e_verbose` - Verbose output for debugging
  - `test_complete_with_e2e` - Complete testing including E2E
- Official Azure Cloud Adoption Framework documentation mapping for 55 resources
- New nested `official` object structure containing Azure CAF documentation attributes
- Comprehensive official resource provider namespace mappings
- GitHub Copilot Agent firewall configuration for improved CI/CD testing
- Enhanced resource validation and testing framework
- Comprehensive CI testing pipeline with resource validation, matrix testing, and coverage analysis
- Advanced Makefile targets for comprehensive testing (`test_ci`, `test_ci_fast`, `test_ci_complete`)
- Shared testing utilities to reduce code duplication (SonarQube compliance)
- Refactored naming convention tests to use centralized test helpers

### Changed
- **BREAKING**: Consolidated `resourceDefinition.json` and `resourceDefinition_out_of_docs.json` into single unified file
- **BREAKING**: Refactored JSON structure to nest official Azure CAF attributes under `official` object
- Updated resource definitions to include proper Azure CAF documentation mapping for key resources:
  - API Management service instance (`apim`) - Microsoft.ApiManagement/service
  - AKS cluster (`aks`) - Microsoft.ContainerService/managedClusters
  - Container apps (`ca`) - Microsoft.App/containerApps
  - Application gateway (`agw`) - Microsoft.ApplicationGateway/applicationGateways
  - Virtual network (`vnet`) - Microsoft.Network/virtualNetworks
  - Storage account (`st`) - Microsoft.Storage/storageAccounts
  - And 49 additional resources with official mappings
- Simplified resource definition structure for non-official resources (only `resource` field in `official` object)
- Enhanced code generation logic to handle nested official attributes
- Updated documentation and contribution guidelines to reflect new structure

### Fixed
- DNS blocking issues with `checkpoint-api.hashicorp.com` during integration tests
- Resource provider namespace accuracy for officially documented Azure resources
- Resource generation and validation processes for unified file structure

### Removed
- `resourceDefinition_out_of_docs.json` file (consolidated into main file)
- Legacy flat structure for official documentation attributes

## Migration Guide

### For Contributors
- Use the new nested `official` object structure when adding or modifying resources
- Resources in official Azure CAF documentation should include `slug`, `resource`, and `resource_provider_namespace` in the `official` object
- Resources not in official documentation should only include the `resource` field in the `official` object

### For Consumers
- The root-level `slug` field remains unchanged for backward compatibility
- New official documentation data is available through the `official` object
- No breaking changes to existing provider functionality

## Statistics
- **Total Resources**: 395 (previously 364 + 31 across two files)
- **Official Azure CAF Mappings**: 55 resources with complete official documentation data
- **Non-Official Resources**: 340 resources with simplified official structure
- **Files Consolidated**: 2 → 1 resource definition file

---

*This changelog consolidates major structural changes made to the terraform-provider-azurecaf resource definitions and documentation mapping. Future releases will continue to document changes in this format for semantic versioning purposes.*