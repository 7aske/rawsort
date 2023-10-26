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

Photos will recursively be searched for in the source directory and moved to the destination directory with the provided format:

```
$HOME/Pictures/import/IMG_001.jpg
```

will be moved to

```
$HOME/Pictures/export/Canon/5D Mark IV/2019-01-01/2019-01-01_12:34:56_Canon_5D Mark IV.jpg
```

where the format is

```
Canon - %K
5D Mark IV - %L
2019-01-01 - %D
2019-01-01_12:34:56_Canon_5D Mark IV.jpg - %D_%t_%K_%L%e
```

## Options

Filename format options:

```
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
