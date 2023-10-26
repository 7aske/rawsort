## Description

A simple CLI tool to organize photos.

## Installation

```bash
sudo make install
```

## Usage

```bash
rawsort -s $HOME/Pictures/import -d $HOME/Pictures/export -f "%K/%L/%D/%D_%t_%K_%L%e"
```

## Options

```
Filename format options:
    %D - Date
    %t - Time
    %y - Year
    %m - Month (mm)
    %d - Day
    %K - Make
    %L - Model
    %e - Extension
```

## Contact

Nikola Tasić – nik@7aske.com

Distributed under the GPL v2 license. See [LICENSE](./LICENSE) for more information.

[7aske.com](https://7aske.com)

[github.com/7aske](https://github.com/7aske)

## Contributing

1. Fork it (<https://github.com/7aske/rawsort/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request