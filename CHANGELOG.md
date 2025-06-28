# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- TBD

### Changed

- TBD

### Deprecated

- TBD

### Removed

- TBD

### Fixed

- TBD

### Security

- TBD

## [1.0.0] - 2025-06-28

### Changed

- Refactored Snowflake ID generator to use a mutex and an `int64` for sequence management instead of a channel, improving performance and concurrency control.
- Enhanced `Close()` method to properly reset the singleton generator's state, enabling clean re-initialization for testing scenarios.
- Improved `NewGenerator` validation to ensure the node ID is non-negative.
- Implemented clock-backwards detection and sequence overflow handling within the `Next()` method for increased robustness.
- Updated benchmark tests to reflect the new generator implementation.

## [0.0.4] - 2024-03-04

### Added

- Added Close function.

## [0.0.3] - 2024-02-10

### Fixed

- Fixed the Next function exceed condition.

## [0.0.2] - 2024-02-10

### Fixed

- Fixed typo on test.

## [0.0.1] - 2024-02-10

### Added

- Initialize the snowflake sequence generator.
