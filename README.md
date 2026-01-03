# <p align="center">lrcsnc</p>
Gets the currently playing song's synced lyrics and displays them in sync with song's actual position.

lrcsnc was primarily designed for bars like [Waybar](https://github.com/Alexays/Waybar), but grew into something that can be used basically anywhere (check below!)

https://github.com/user-attachments/assets/467a5cc0-28cd-4c61-8bf7-f1f533ad2a95

<sub>^ - lrcsnc launched as a server and [lrcsnc-lbl-client](https://github.com/Endg4meZer0/lrcsnc-lbl-client) used for client side.</sub>

https://github.com/user-attachments/assets/1bc93e59-385f-41cb-a23e-49298e5887b0

<sub>^ - a basic example in Waybar using the simple client that's included in the base package</sub>

## Features

- Precise synchronizing to any* player that supports MPRIS
- Can be tailor-fit into a lot of things; the UNIX way, as they say.
- Client-server communication, allowing for different types of clients to exist simultaneously
- A decent level of customization and configuration using TOML
- Barebones romanization for Japanese, Chinese and Korean languages
- ...and more!

<sub>* - player should be precise itself. There are examples of players that don't handle timings well, or have problems with their MPRIS implementations. Check [compatibility wiki page](https://github.com/Endg4meZer0/lrcsnc/wiki/Compatibility-with-players) for more.</sub>

## Install
lrcsnc is available at AUR!
```
yay -S lrcsnc
```

Also you can build it from source (see below)

## Build
```
git clone https://github.com/Endg4meZer0/lrcsnc.git
cd lrcsnc
make # or `sudo make all` for automatic install
```
Make sure to have go v1.23 or above.

## Usage
```
lrcsnc [OPTION]
```
Get more info on on available options with `lrcsnc -h`.

## TODO
- [ ] Check [compatibility](https://github.com/Endg4meZer0/lrcsnc/wiki/Compatibility-with-players) with different players
- [ ] More lyrics providers (maybe local files too?)
- [ ] More configuration options?
- [ ] There is definitely always more!

## Need help or want to contribute?
You can always make an issue for either a bug or a feature suggestment! If your question is more general, consider opening a discussion.

## Your favorite song's lyrics were not found?
Consider adding them! Currently lrcsnc uses only *[LrcLib](https://lrclib.net)*, which is a great open-source lyrics provider service that has its own easy-to-use [app](https://github.com/tranxuanthang/lrcget) to download or upload lyrics. Once the lyrics are uploaded, lrcsnc should be able to pick them up on the next play of the song (that is if the cached version is not available though - check the docs for how to clear the cache). Also, other ways to get lyrics will be implemented later.
