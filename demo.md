# Window 1

## Help text

```
$ accesschecker
You are missing an argument.

Usage:
  accesschecker [accessgrant] [flags]

Flags:
  -h, --help                             help for accesschecker
  -o, --output-empty-passphrase-access   Output generated empty passphrase access (only use to download/remove files)
```

## Basic usage to check if any files are uploaded with an empty passphrase in this project

The access you provide can be a new unrestricted access grant created from the Satellite UI, with any passphrase:

```
$ accesschecker 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP

You have 2 files uploaded without encryption in this project:

sj://asdf/DEVELOPING.md
sj://moby/README.md
```

## Usage with --output-empty-passphrase-access flag to get an access that can be used with Uplink to download/remove files:

```
$ accesschecker 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP --output-empty-passphrase-access

You have 2 files uploaded without encryption in this project:

sj://asdf/DEVELOPING.md
sj://moby/README.md

==================================================

Generated empty passphrase grant:

1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP

WARNING: This access is capable of uploading and downloading unencrypted files to and from your project. We recommend using it only to download and subsequently remove files which are unencrypted.

```

# Window 2

## Import empty passphrase access outputted by tool:
```
$ uplink access import generatedempty 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP
Imported access "generatedempty" to "/home/moby/.config/storj/uplink/access.json"
$ uplink access use generatedempty
Switched default access to "generatedempty"
```

## Download unencrypted files:
```
$ uplink cp sj://moby/README.md ./
download sj://moby/README.md to README.md
2.79 KiB / 2.79 KiB [-----------------------------------------------------------------------] 100.00% ? p/s
$ uplink cp sj://asdf/DEVELOPING.md ./
download sj://asdf/DEVELOPING.md to DEVELOPING.md
15.69 KiB / 15.69 KiB [---------------------------------------------------------------------] 100.00% ? p/s
```

## Delete unencrypted files from Storj DCS:
```
$ uplink rm sj://asdf/DEVELOPING.md
removed sj://asdf/DEVELOPING.md
$ uplink rm sj://moby/README.md
removed sj://moby/README.md
```

## Switch to an access with a passphrase, and delete empty passphrase access:

```
$ uplink access use <someotheraccess> 
Switched default access to "<someotheraccess>"
$ uplink access remove generatedempty
Removed access "generatedempty" from "/home/moby/.config/storj/uplink/access.json"
```


## Re-upload files with an access containing a passphrase:

```
$ uplink cp README.md sj://moby
upload README.md to sj://moby/README.md
2.79 KiB / 2.79 KiB [-----------------------------------------------------------------------] 100.00% ? p/s
$ uplink cp DEVELOPING.md sj://asdf
upload DEVELOPING.md to sj://asdf/DEVELOPING.md
15.69 KiB / 15.69 KiB [---------------------------------------------------------------------] 100.00% ? p/s
```

# Window 1

## Check that no files were uploaded with an empty encryption passphrase:

```
$ accesschecker 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP --output-empty-passphrase-access
You do not have any files uploaded without encryption in this project.
```
