# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [[0.1.2](https://github.com/Endg4meZer0/lrcsnc/releases/tag/v0.1.2)] - 2026-01-03
### Added
- Ability to remove cached data for currently playing song by signaling USR2 at the lrcsnc process.
### Changed
- LrcLib provider now only sends one `search` request instead of multiple tries of `get` and `search` with different queries. Now, to get the best pick out of received data, there is an scoring algorithm inside lrcsnc itself that gives points for specific advantages (e.g. matching all artists or having exact duration match). Of course, that's not perfect, but it's actually better than LrcLib's `get` results sometimes, and on par otherwise.
- A fatal log should now allow a graceful shutdown to proceed (e.g. deleting the UNIX socket during shutdown if launched as server with UNIX specified as protocol).
- The default value for `cache.life-span` was changed from 168 (hours, so 1 week) to 720 (hours, so 1 month).
### Fixed
- The connection to LrcLib was insecure (using http) all this time (and I didn't notice T_T). That was changed (now it uses https properly).
- Fatal log was always duplicated due to a duplicated print inside the method.
- A proper error message for UNIX socket being busy during server's start-up that was introduced in v0.1.1 now works as intended.
### Removed
- Japanese romanization was removed due to being inconsistent and difficult to implement properly. Until I find a way to do so (preferrably not depending on outside binaries like `kakasi`), I will not try again. As such, the mention of `kakasi` required as optional dependency was removed as well.

## [[0.1.1](https://github.com/Endg4meZer0/lrcsnc/releases/tag/v0.1.1)] - 2025-11-02
### Added
- Client-server communication is realized: now lrcsnc can be launched in server-only, client-only, or standalone mode. Connections are made using UNIX sockets or TCP protocol. Standalone mode remains a default option and doesn't require any connection.

    Corresponding configuration and flag options were added: the `-s` (server flag), `-l` (listen at), `-p` (protocol to use), and a whole `[net]` section in config file for consistent connection configurations (more at wiki).
### Changed
- Due to client-server communication addition, the idea of TUI was completely scrapped, and so the remnants in configuration sections: the `[output]` section is now called `[client]`, and the `[output.piped]` section was removed completely as it is obsolete. The members of `[output.piped]` section were moved directly to `[client]` section.
- The plain and JSON variants were replaced by templates and formats: you can do any kind of output now by using the `template` config option and keys like `%text%`, `%title%` and others (check wiki for more).
- Store condition config option was changed from byte set to a more human-readable table.
    Was:
    ```
    [cache]
    ...
    store-condition = 100 # first for synced, second for plain, third for instrumental
    ```
    Now:
    ```
    [cache]
    ...

    [cache.store-condition]
    if-synced = true
    if-plain = false
    if-instrumental = true # was changed to true by default with this refactor
    ``` 

## [[0.1.0](https://github.com/Endg4meZer0/lrcsnc/releases/tag/v0.1.0)] - 2025-05-03
### Added
- Some simple unit tests like cache and romanization.
- Makefile for easier build and more control over linking stuff.
- JSON piped output.
### Changed
- **Everything is rewrote from scratch**:
- MPRIS support now works on signals instead of polling.
- Configuration format is now TOML.
- Japanese romanization now uses [kakasi](https://github.com/loretoparisi/kakasi) and is able to romanize kanji. Kakasi is installed as a separate dependency for it to work.
- Romanization is now able to handle multiple languages on one line.
- Flags usage is now handled by [go-flags](https://github.com/jessevdk/go-flags) instead of standard library.
### Removed
- Playerctl dependency is now completely abandoned and cut from the app in favor of direct MPRIS handling using [own library](https://github.com/Endg4meZer0/go-mpris).
- Terminal output in one line is removed.

## The changelog below describes [playerctl-lyrics](https://github.com/Endg4meZer0/playerctl-lyrics), the now archived project which lrcsnc is based on.

## [[0.2.1](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.2.1)] - 2024-08-29
### Added
- A command-line option `-o` to redirect the output to a set file.
- ~~A command-line option to display lyrics in one line~~ **is removed now**
- A configuration option to offset the lyrics by set seconds by @Endg4meZer0 in [#9](https://github.com/Endg4meZer0/playerctl-lyrics/pull/9)
### Changed
- More refactoring: `cmus` and other players that report position in integer seconds are now fully supported.
- Cache system is reverted back to JSON instead of LRC files to allow more additional data to be stored ([#10](https://github.com/Endg4meZer0/playerctl-lyrics/pull/10))
### Fixed
- Instrumental lyrics overlapped actual lyrics in some cases ([#11](https://github.com/Endg4meZer0/playerctl-lyrics/pull/11))

## [[0.2.0](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.2.0)] - 2024-08-24
### Changed
- A big concept rewrite happened to allow players like `cmus` that report position in integer seconds work on par with others.
- A rename of `doCacheLyrics` configuration option to `enabled`

## [[0.1.1](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.1.1)] - 2024-08-21
### Added
- A configuration option to control the format of repeated lyrics multiplier.
### Fixed
- Fixed a panic if there is no space after a timestamp.
- Fixed a panic when romanization of Japanese kanji failed and fell down to Chinese characters 

## [[0.1.0](https://github.com/Endg4meZer0/playerctl-lyrics/releases/tag/v0.1.0)] - 2024-08-15
### Added
- Initial unstable release of playerctl-lyrics.
- Display lyrics for currently playing song.
- Support for multiple music players using `playerctl`.
- Automatic lyric fetching from `lrclib`.
- Configuration file for custom settings.
- Romanization for several asian languages.
- Caching system to significantly reduce traffic usage.