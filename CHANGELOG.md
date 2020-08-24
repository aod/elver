# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.4] - 2020-08-24

## [0.4.3] - 2020-08-16
### Fixed
 - Fix missing AoC year header in first line of output

## [0.4.2] - 2020-07-21
### Added
- Besides the `AOC_SESSION` environment variable, the session ID can now
  _also_ be stored in a file `aoc_session` inside the elver config dir. (The
  environment variable does take precedence over the file.)
    - Windows: `%AppData%\elver\`
    - MacOS: `/Library/Application Support/elver/`
    - Linux: `$HOME/.config/elver/`

## [0.3.2] - 2020-06-03
### Added
- `-y` flag to specify which year to run
- `-d` flag to specify which day to run

## [0.2.2] - 2020-05-21
### Added
- New logo and badges in README.md

## [0.2.1] - 2020-05-10
### Added
- Benchmarking of the latest solution behind the `-b` flag

## [0.1.1] - 2020-05-03
### Fixed
- Refactored procedure that runs the latest solution for better maintainability
- Refactored procedure which fetches and caches the input for better maintainability

## [0.1.0] - 2020-05-01
### Added
- Automatically run the latest solution
- Automatically download and cache input
- Display the time it took to run a solution
