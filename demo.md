# Window 1

```
$ accesschecker
You are missing an argument.

Usage:
  accesschecker [accessgrant] [flags]

Flags:
  -h, --help                             help for accesschecker
  -o, --output-empty-passphrase-access   Output generated empty passphrase access (only use to download/remove files)

$ accesschecker 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP

You have 2 files uploaded without encryption in this project:

sj://asdf/DEVELOPING.md
sj://moby/README.md

moby@thinkpad:~/dev/storj/access-empty-passphrase-checker$ accesschecker 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP --output-empty-passphrase-access

You have 2 files uploaded without encryption in this project:

sj://asdf/DEVELOPING.md
sj://moby/README.md

==================================================

Generated empty passphrase grant:

1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP

WARNING: This access is capable of uploading and downloading unencrypted files to and from your project. We recommend using it only to download and subsequently remove files which are unencrypted.

```

# Window 2

```
$ uplink access import generatedempty 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP
Imported access "generatedempty" to "/home/moby/.config/storj/uplink/access.json"
$ uplink access use generatedempty
Switched default access to "generatedempty"
$ uplink rm sj://asdf/DEVELOPING.md
removed sj://asdf/DEVELOPING.md
moby@thinkpad:~/dev/storj/storj$ uplink rm sj://moby/README.md
removed sj://moby/README.md
```

# Window 1

```
$ accesschecker 1QiUipyf19MS5ZMC2y7knDt4jRtf53SgpuXnoX8mdfEqQkx2bUbXTdYhoAxZyioxaWhmjSJR2XHFcBuWJv85oPKSeJ7ZY6XQcUJdRQYrxiG1s6Gf2i6Xxrp1YyAbdG7kythpA6jxr5JBWnMffdioEha9DJWSeHrLUbLcAMTYRjZH1c2yXbg2uxNcAypjbs5dh4Gb4Ksyv4WN7jqQrRqx3eXqf873Z1zFCi8EaczhH4LhEenLgfmUbkL4NAkCRZ8EM5Zp89oP --output-empty-passphrase-access
You do not have any files uploaded without encryption in this project.
```
