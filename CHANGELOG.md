# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- TBD

### Changed

- Replaced channel-based sequence generation with a mutex-protected counter to improve speed.
- Updated time generation logic within the Next() method.
- Refactored NewGenerator and Close methods to better support test isolation, particularly for the singleton generator instance.
- Updated BenchmarkNext to accurately measure the performance of the Next() method after changes.
- Fixed issues in tests that arose from the refactoring of the generator logic.

### Deprecated

- TBD

### Reomved

- TBD

### Fixed

- TBD

### Security

- TBD

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
