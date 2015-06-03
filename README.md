# clihist

Move over gnuplot, there's a new sheriff in town.

## Installation

Using go get:

```
$ go get github.com/wernerandrew/clihist
```

To just make the binary, no funny stuff:

```
$ git clone git@github.com:wernerandrew/clihist.git
$ cd clihist
$ go build
```

## Usage

### Example

```
$ python -c 'import random; print "\n".join([str(random.normalvariate(0, 1)) for _ in xrange(100)])' | clihist

                                ####
                                ####
                                ####
                                ####
                                ####
                                ####
                                ####
                                ####
                                ####
                                ########
                            ############
                            ############
            ####            ############
            ####        ################    ########
            ####        ################    ########
            ########    ################    ########
            ############################    ########
            ########################################
            ########################################
            ########################################
        ############################################
        ############################################
####    ############################################        ####
```

Use `clihist -help` to get a full range of ocmmand line options.
