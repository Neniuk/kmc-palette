
# K-Means Clustering Palette

This project allows you to build a palette from an image (png, jpg, gif) through utilizing k-means clustering. By default it provides 8 colors (clusters), with an iteration of 2 (by increasing the number of iterations, the accuracy of clustering increases), additionally it adds color variants of the generated colors, to provide an extended palette of 16 colors by default (can be turned off using the "-scv" flag).  

Generated colors displayed with both the rgb values, hex color codes and a color block.


## Installation

Ready built binaries are available in the [bin folder](https://github.com/Neniuk/kmc-palette/tree/main/bin) and on the [releases page](https://github.com/Neniuk/kmc-palette/releases). Download the appropriate version for your OS.

## Usage/Examples

```bash
./kmc-palette [-i <number>] [-k <number>] [-scv] [-hhp] <filepath-to-image>
```

### Flags:  
-i, <number>, number of iterations for K-means clustering  
-k, <number>, number of clusters for K-means clustering  
-scv, skip adding color variants to the palette  
-hhp, hide horizontal palette  


